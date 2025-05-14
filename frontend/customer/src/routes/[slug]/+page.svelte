<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { goto } from "$app/navigation";

    import { isLoading, error, productHelpers } from "$lib/stores/product";
    import { cartHelpers } from "$lib/stores/cart";
    import type { Product } from "$lib/stores/product";
    import { connectToSSE, disconnectFromSSE, connectionStatus } from "$lib/stores/pricing";
    
    import { Check, Star, ShoppingCart } from "lucide-svelte";
    import { Badge } from "$lib/components/ui/badge/index";
    import { Button } from "$lib/components/ui/button/index";
    import { Input } from "$lib/components/ui/input/index";
    import { toast } from "svelte-sonner";
    import * as Card from "$lib/components/ui/card/index";
    
    const imageBaseUrl = "http://localhost:8080/assets/products/";
    let slug = $derived(page.params.slug);
    let productId = $state<number | null>();
    let isProductLoading = $state<boolean>(true);
    let productError: string | null = $state(null);
    let product: Product | null = $state(null);

    let quantity = $state<number>(1);
    let isAddingToCart = $state<boolean>(false);
    let addToCartError = $state<string | null>(null);
    let addToCartSuccess = $state<boolean>(false);
    
    onMount(() => {
        (async () => {
            try {
                productId = productHelpers.extractIdFromSlug(slug);
                if (productId === null) {
                    console.log("Invalid product slug format: ", slug)
                    productError = "Invalid product URL.";
                    isProductLoading = false;
                    return
                }
    
                const response = await productHelpers.getProduct(productId);
                
                if (!response.success) {
                    console.error("Failed to fetch product:", response.error);
                    productError = response.error || "Failed to load product details.";
                } else {
                    product = { ...response.data as Product };
                    console.log("Product details:", product);
                }
            }   catch (err) {
                console.error("Error loading product: ", err)
            } finally {
                isProductLoading = false;
            }
        })();
    
        connectToSSE();
    
        return () => {
            disconnectFromSSE();
        }
    });

    function isInStock(product: Product): boolean {
        return product.stock_quantity > 0;
    }

    async function handleAddToCart() {
        if (!product || !productId) return;
        
        // Reset states
        isAddingToCart = true;
        addToCartError = null;
        addToCartSuccess = false;

        try {
            const response = await cartHelpers.addToCart(productId, quantity);
            
            if (response.success) {
                addToCartSuccess = true;
                toast.success("Product added to Cart successfully")
                setTimeout(() => {
                    addToCartSuccess = false;
                }, 3000);
            } else {
                addToCartError = response.error || "Failed to add to cart";
                toast.error("Add Product to cart failed")
            }
        } catch (err) {
            addToCartError = "An unexpected error occurred";
            console.error("Error adding to cart:", err);
            toast.error("Add Product to cart failed")
        } finally {
            isAddingToCart = false;
        }
    }
</script>

<div class="container mx-auto px-10 py-5">
    <!-- Connection status indicator -->
    <div class="mb-4">
        {#if $connectionStatus === 'connected'}
            <Badge variant="outline" class="bg-green-50 text-green-700 border-green-200">
                Connected - Receiving price updates
            </Badge>
        {:else if $connectionStatus === 'connecting'}
            <Badge variant="outline" class="bg-yellow-50 text-yellow-700 border-yellow-200">
                Connecting to price updates...
            </Badge>
        {:else if $connectionStatus === 'error'}
            <Badge variant="outline" class="bg-red-50 text-red-700 border-red-200">
                Connection error - Will retry automatically
            </Badge>
        {:else}
            <Badge variant="outline" class="bg-gray-50 text-gray-700 border-gray-200">
                Disconnected from price updates
            </Badge>
        {/if}
    </div>

    {#if isProductLoading }
        <div class="flex justify-center items-center h-64">
            <div class="text-gray-500 animate-pulse text-lg">Loading product details...</div>
        </div>
    {:else if productError}
        <Card.Root class="w-full max-w-2xl mx-auto">
            <Card.Header>
                <Card.Title class="text-red-600">Error Loading Product</Card.Title>
            </Card.Header>
            <Card.Content>
                <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
                    {error}
                </div>
                <div class="mt-4 flex justify-center">
                    <Button onclick={() => goto('/')} variant="outline">
                        Return to Home
                    </Button>
                </div>
            </Card.Content>
        </Card.Root>
    {:else if product}
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8"> 
            <div class="lg:col-span-1">
                <Card.Root> 
                    <Card.Content>
                        <img
                            src="{imageBaseUrl}{product.image_path}"
                            alt={product.name}
                            class="w-30 h-30 rounded-md object-cover bg-gray-100"
                        />
                    </Card.Content>
                </Card.Root>
            </div>
            <div class="lg:col-span-2"> 
                <Card.Root>
                    <Card.Header>
                        <Card.Title class="text-lg font-semibold text-gray-700">{product.name}</Card.Title>
                    </Card.Header>
                    <Card.Content>
                        <div class="text-gray-500 mb-4">
                            {product.description}
                        </div>
                        <div class="flex items-center text-sm text-gray-500">
                            <span class="flex items-center">
                                <Star class="h-4 w-4 text-yellow-500 mr-1" />
                                {product.average_rating} / 5
                            </span>
                            <span class="mx-2">•</span>
                            <span>{product.review_count} reviews</span>
                            <span class="mx-2">•</span>
                            <span>{product.brand_name}</span>
                        </div>
                    </Card.Content>
                </Card.Root>
            </div>

            <div class="mt-6"> 
                <div class="text-lg font-bold text-blue-700">
                    {product.adjusted_price}
                </div>
                
                <div>
                    <label for="quantity" class="sr-only">Quantity</label>
                    <Input
                        id="quantity"
                        type="number"
                        bind:value={quantity}
                        min={1}
                        max={product.stock_quantity}
                        class="w-20 text-center"
                    />

                    <div class="flex-1">
                        {#if isInStock(product)}
                            <Badge variant="outline" class="bg-green-50 text-green-700 border-green-200">
                                In Stock: {product.stock_quantity} available
                            </Badge>
                        {:else}
                            <Badge variant="outline" class="bg-red-50 text-red-700 border-red-200">
                                Out of Stock
                            </Badge>
                        {/if}
                    </div>
                        
                    <div class="mt-4">
                        <Button 
                            onclick={handleAddToCart} 
                            disabled={!isInStock(product) || isAddingToCart || quantity < 1}
                            class="w-full md:w-auto"
                        >
                            {#if isAddingToCart}
                                Adding...
                            {:else if addToCartSuccess}
                                <Check class="h-4 w-4 mr-2" />
                                Added to Cart
                            {:else}
                                <ShoppingCart class="h-4 w-4 mr-2" />
                                Add to Cart
                            {/if}
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    {/if}
</div>