<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { cartHelpers } from "$lib/stores/cart";
    import type { CartItem } from "$lib/stores/cart";
    
    import { Trash2, ShoppingBag, ArrowLeft, RefreshCw, Minus, Plus } from "lucide-svelte";
    import { Badge } from "$lib/components/ui/badge/index";
    import { Button } from "$lib/components/ui/button/index";
    import { toast } from "svelte-sonner";
    import * as Card from "$lib/components/ui/card/index";
    import * as Table from "$lib/components/ui/table/index";
    import { checkoutHelpers } from "$lib/stores/checkout";
    
    const imageBaseUrl = "http://localhost:8080/assets/products/";
    
    let isLoading = $state<boolean>(true);
    let cartItems = $state<CartItem[]>([]);
    let cartTotalItems = $state<number>(0);
    let cartTotalPrice = $state<number>(0);
    let cartError = $state<string | null>(null);
    
    let isRemoving = $state<Record<number, boolean>>({});
    let isClearing = $state<boolean>(false);
    let isGeneratingOrder = $state<boolean>(false);
    
    onMount(() => {
        loadCart();
    });
    
    async function loadCart() {
        isLoading = true;
        cartError = null;
        
        try {
            const response = await cartHelpers.fetchCartWithItems();
            
            if (!response.success) {
                cartError = response.error || "Failed to load cart";
                console.error("Failed to load cart:", cartError);
            } else if (response.data) {
                cartItems = response.data.items || [];
                cartTotalItems = response.data.total_items;
                cartTotalPrice = response.data.total_price;
                console.log("Cart loaded:", response.data);
            }
        } catch (err) {
            console.error("Error loading cart:", err);
            cartError = "An unexpected error occurred";
        } finally {
            isLoading = false;
        }
    }
    
    async function handleRemoveItem(itemId: number) {
        if (isRemoving[itemId]) return;
        
        isRemoving = { ...isRemoving, [itemId]: true };
        
        try {
            const response = await cartHelpers.removeFromCart(itemId);
            
            if (response.success) {
                toast.success("Item removed from cart");
                await loadCart(); // Reload cart after removal
            } else {
                toast.error(response.error || "Failed to remove item");
            }
        } catch (err) {
            console.error("Error removing item:", err);
            toast.error("Failed to remove item");
        } finally {
            isRemoving = { ...isRemoving, [itemId]: false };
        }
    }
    
    async function handleClearCart() {
        if (isClearing) return;
        
        isClearing = true;
        
        try {
            const response = await cartHelpers.clearCart();
            
            if (response.success) {
                toast.success("Cart cleared");
                await loadCart(); // Reload cart after clearing
            } else {
                toast.error(response.error || "Failed to clear cart");
            }
        } catch (err) {
            console.error("Error clearing cart:", err);
            toast.error("Failed to clear cart");
        } finally {
            isClearing = false;
        }
    }
    
    function formatPrice(price: number): string {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD'
        }).format(price);
    }
    
    async function goToCheckout() {
        if (isGeneratingOrder) return;
        isGeneratingOrder = true
        
        try {
            const response = await checkoutHelpers.generateOrder();
            goto('/checkout')
        } catch (err) {
            console.error("Error generating cart: ", err);
            toast.error("Failed generating cart")
        } finally {
            isGeneratingOrder= false;
        }
    }
</script>

<div class="container mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold mb-6">Your Shopping Cart</h1>
    
    {#if isLoading}
        <div class="flex justify-center items-center h-64">
            <div class="text-gray-500 animate-pulse text-lg">Loading your cart...</div>
        </div>
    {:else if cartError}
        <Card.Root class="w-full max-w-2xl mx-auto">
            <Card.Header>
                <Card.Title class="text-red-600">Error Loading Cart</Card.Title>
            </Card.Header>
            <Card.Content>
                <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
                    {cartError}
                </div>
                <div class="mt-4 flex justify-center">
                    <Button onclick={() => loadCart()} variant="outline">
                        <RefreshCw class="h-4 w-4 mr-2" />
                        Try Again
                    </Button>
                </div>
            </Card.Content>
        </Card.Root>
    {:else if cartItems.length === 0}
        <Card.Root class="w-full max-w-2xl mx-auto">
            <Card.Header>
                <Card.Title>Your Cart is Empty</Card.Title>
            </Card.Header>
            <Card.Content>
                <div class="text-center py-8">
                    <ShoppingBag class="h-16 w-16 mx-auto text-gray-400 mb-4" />
                    <p class="text-gray-600 mb-6">Your shopping cart is empty. Add some products to get started!</p>
                    <Button onclick={() => goto('/')}>
                        <ArrowLeft class="h-4 w-4 mr-2" />
                        Continue Shopping
                    </Button>
                </div>
            </Card.Content>
        </Card.Root>
    {:else}
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div class="lg:col-span-2">
                <Card.Root>
                    <Card.Header class="flex flex-row items-center justify-between">
                        <Card.Title>Cart Items ({cartTotalItems})</Card.Title>
                        <Button 
                            onclick={handleClearCart} 
                            variant="outline" 
                            class="text-red-600 border-red-200 hover:bg-red-50"
                            disabled={isClearing}
                        >
                            <Trash2 class="h-4 w-4 mr-2" />
                            {isClearing ? 'Clearing...' : 'Clear Cart'}
                        </Button>
                    </Card.Header>
                    <Card.Content>
                        <Table.Root>
                            <Table.Header>
                                <Table.Row>
                                    <Table.Head class="w-[100px]">Product</Table.Head>
                                    <Table.Head>Details</Table.Head>
                                    <Table.Head>Quantity</Table.Head>
                                    <Table.Head>Price</Table.Head>
                                    <Table.Head class="w-[80px]"></Table.Head>
                                </Table.Row>
                            </Table.Header>
                            <Table.Body>
                                {#each cartItems as item (item.id)}
                                    <Table.Row>
                                        <Table.Cell>
                                            <img 
                                                src="{imageBaseUrl}{item.product_image_path}" 
                                                alt={item.product_name}
                                                class="w-20 h-20 object-cover rounded-md"
                                            />
                                        </Table.Cell>
                                        <Table.Cell>
                                            <div class="font-medium">{item.product_name}</div>
                                            <div class="text-sm text-gray-500">
                                                Unit Price: {formatPrice(item.product_adjusted_price)}
                                            </div>
                                        </Table.Cell>
                                        <Table.Cell>
                                            <div class="flex items-center space-x-2">
                                                <span class="w-8 text-center">{item.quantity}</span>
                                            </div>
                                        </Table.Cell>
                                        <Table.Cell>
                                            {formatPrice(item.quantity * item.product_adjusted_price)}
                                        </Table.Cell>
                                        <Table.Cell>
                                            <Button 
                                                onclick={() => handleRemoveItem(item.id)}
                                                variant="ghost" 
                                                size="icon"
                                                class="text-red-500 hover:text-red-700 hover:bg-red-50"
                                                disabled={isRemoving[item.id]}
                                            >
                                                <Trash2 class="h-4 w-4" />
                                            </Button>
                                        </Table.Cell>
                                    </Table.Row>
                                {/each}
                            </Table.Body>
                        </Table.Root>
                    </Card.Content>
                </Card.Root>
            </div>
            
            <div class="lg:col-span-1">
                <Card.Root>
                    <Card.Header>
                        <Card.Title>Order Summary</Card.Title>
                    </Card.Header>
                    <Card.Content>
                        <div class="space-y-4">
                            <div class="flex justify-between">
                                <span>Items ({cartTotalItems}):</span>
                                <span>{formatPrice(cartTotalPrice)}</span>
                            </div>
                            <div class="flex justify-between">
                                <span>Shipping:</span>
                                <span>Calculated at checkout</span>
                            </div>
                            <div class="border-t pt-4 flex justify-between font-bold">
                                <span>Total:</span>
                                <span>{formatPrice(cartTotalPrice)}</span>
                            </div>
                            
                            <Button 
                                onclick={goToCheckout} 
                                class="w-full mt-4"
                                disabled={cartItems.length === 0}
                            >
                                Proceed to Checkout
                            </Button>
                            
                            <Button 
                                onclick={() => goto('/')} 
                                variant="outline" 
                                class="w-full mt-2"
                            >
                                <ArrowLeft class="h-4 w-4 mr-2" />
                                Continue Shopping
                            </Button>
                        </div>
                    </Card.Content>
                </Card.Root>
            </div>
        </div>
    {/if}
</div>