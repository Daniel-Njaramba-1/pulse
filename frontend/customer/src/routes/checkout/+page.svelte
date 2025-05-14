<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { checkoutHelpers } from "$lib/stores/checkout";
    import type { Order, OrderItem } from "$lib/stores/checkout";
    
    import { ShoppingBag, ArrowLeft, RefreshCw, CreditCard } from "lucide-svelte";
    import { Badge } from "$lib/components/ui/badge/index";
    import { Button } from "$lib/components/ui/button/index";
    import { toast } from "svelte-sonner";
    import * as Card from "$lib/components/ui/card/index";
    import * as Table from "$lib/components/ui/table/index";
    
    const imageBaseUrl = "http://localhost:8080/assets/products/";
    
    let isLoading = $state<boolean>(true);
    let order = $state<Order | null>(null);
    let orderItems = $state<OrderItem[]>([]);
    let orderTotalItems = $state<number>(0);
    let orderTotalPrice = $state<number>(0);
    let orderError = $state<string | null>(null);
    
    let isProcessingPayment = $state<boolean>(false);
    
    onMount(() => {
        loadOrder();
    });
    
    async function loadOrder() {
        isLoading = true;
        orderError = null;
        
        try {
            const response = await checkoutHelpers.fetchOrderWithItems();
            
            if (!response.success) {
                orderError = response.error || "Failed to load order";
                console.error("Failed to load order:", orderError);
            } else if (response.data) {
                order = response.data;
                orderItems = response.data.items || [];
                orderTotalPrice = response.data.total_price;
                console.log("Order loaded:", response.data);
            }
        } catch (err) {
            console.error("Error loading order:", err);
            orderError = "An unexpected error occurred";
        } finally {
            isLoading = false;
        }
    }
    
    async function handleProcessPayment() {
        if (isProcessingPayment) return;
        
        isProcessingPayment = true;
        
        try {
            const response = await checkoutHelpers.payment();
            
            if (response.success) {
                toast.success("Payment successful!");
                // Navigate to a success page or order confirmation
                goto('/cart');
            } else {
                toast.error(response.error || "Payment failed");
            }
        } catch (err) {
            console.error("Error processing payment:", err);
            toast.error("Payment failed");
        } finally {
            isProcessingPayment = false;
        }
    }
    
    function formatPrice(price: number): string {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD'
        }).format(price);
    }
    
    function formatDate(dateString: string): string {
        return new Date(dateString).toLocaleString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }
</script>

<div class="container mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold mb-6">Your Order</h1>
    
    {#if isLoading}
        <div class="flex justify-center items-center h-64">
            <div class="text-gray-500 animate-pulse text-lg">Loading your order...</div>
        </div>
    {:else if orderError}
        <Card.Root class="w-full max-w-2xl mx-auto">
            <Card.Header>
                <Card.Title class="text-red-600">Error Loading Order</Card.Title>
            </Card.Header>
            <Card.Content>
                <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
                    {orderError}
                </div>
                <div class="mt-4 flex justify-center">
                    <Button onclick={() => loadOrder()} variant="outline">
                        <RefreshCw class="h-4 w-4 mr-2" />
                        Try Again
                    </Button>
                </div>
            </Card.Content>
        </Card.Root>
    {:else if !order || orderItems.length === 0}
        <Card.Root class="w-full max-w-2xl mx-auto">
            <Card.Header>
                <Card.Title>No Active Order</Card.Title>
            </Card.Header>
            <Card.Content>
                <div class="text-center py-8">
                    <ShoppingBag class="h-16 w-16 mx-auto text-gray-400 mb-4" />
                    <p class="text-gray-600 mb-6">You don't have an active order. Check out your cart to create one.</p>
                    <Button onclick={() => goto('/cart')}>
                        <ArrowLeft class="h-4 w-4 mr-2" />
                        Go to Cart
                    </Button>
                </div>
            </Card.Content>
        </Card.Root>
    {:else}
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div class="lg:col-span-2">
                <Card.Root>
                    <Card.Header class="flex flex-row items-center justify-between">
                        <Card.Title>Order #{order.id}</Card.Title>
                        <Badge variant={order.status === 'pending' ? 'outline' : 'default'}>
                            {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                        </Badge>
                    </Card.Header>
                    <Card.Content>
                        <div class="mb-4 text-sm text-gray-500">
                            <div>Created: {formatDate(order.created_at)}</div>
                        </div>
                        <Table.Root>
                            <Table.Header>
                                <Table.Row>
                                    <Table.Head class="w-[100px]">Product</Table.Head>
                                    <Table.Head>Details</Table.Head>
                                    <Table.Head>Quantity</Table.Head>
                                    <Table.Head>Price</Table.Head>
                                </Table.Row>
                            </Table.Header>
                            <Table.Body>
                                {#each orderItems as item (item.id)}
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
                                                Unit Price: {formatPrice(item.price)}
                                            </div>
                                        </Table.Cell>
                                        <Table.Cell>
                                            <div class="flex items-center">
                                                <span class="w-8 text-center">{item.quantity}</span>
                                            </div>
                                        </Table.Cell>
                                        <Table.Cell>
                                            {formatPrice(item.quantity * item.price)}
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
                                <span>Items ({orderTotalItems}):</span>
                                <span>{formatPrice(orderTotalPrice)}</span>
                            </div>
                            <div class="flex justify-between">
                                <span>Shipping:</span>
                                <span>Free</span>
                            </div>
                            <div class="flex justify-between">
                                <span>Tax:</span>
                                <span>Included</span>
                            </div>
                            <div class="border-t pt-4 flex justify-between font-bold">
                                <span>Total:</span>
                                <span>{formatPrice(orderTotalPrice)}</span>
                            </div>
                            
                            {#if order.status === 'pending'}
                                <Button 
                                    onclick={handleProcessPayment} 
                                    class="w-full mt-4"
                                    disabled={isProcessingPayment}
                                >
                                    <CreditCard class="h-4 w-4 mr-2" />
                                    {isProcessingPayment ? 'Processing...' : 'Pay Now'}
                                </Button>
                            {:else}
                                <div class="bg-green-50 border border-green-200 p-4 rounded-lg text-green-700 text-center">
                                    This order has been {order.status}
                                </div>
                            {/if}
                            
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