<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { products, isLoading, error, productHelpers } from "$lib/stores/product";
    import type { Product } from "$lib/stores/product";
    import { connectToSSE, disconnectFromSSE, connectionStatus, productPrices } from "$lib/stores/pricing";
    import { Search, Star, ShoppingCart } from "lucide-svelte";
    import { Button } from "$lib/components/ui/button/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Badge } from "$lib/components/ui/badge/index";

    // Local state
    let searchQuery = $state<string>("");

    // Base URL for product images
    const imageBaseUrl = "http://localhost:8080/assets/products/";
    let currentPrices = $state<Map<number, number>>(new Map());

    onMount(() => {
        (async () => {
            const productsResponse = await productHelpers.fetchProducts();
            console.log("Products fetched:", productsResponse);
            if (!productsResponse.success) {
                console.error("Failed to fetch products:", productsResponse.error);
            }
        })();
    
        connectToSSE();

        // Subscribe to price updates
        const unsubscribeFromPrices = productPrices.subscribe(prices => {
            currentPrices = new Map(prices);
        });
    
        return () => {
            disconnectFromSSE();
            unsubscribeFromPrices();
        };
    });

    // Helper function to check if product is in stock
    function isInStock(product: Product): boolean {
        return product.stock_quantity > 0;
    }

    // Get current price for a product (either updated price or original)
    function getCurrentPrice(product: Product): string {
        const updatedPrice = currentPrices.get(product.id);
        if (updatedPrice !== undefined) {
            return `$${updatedPrice.toFixed(2)}`;
        }
        return `$${product.adjusted_price.toFixed(2)}`;
    }

    // Check if price has been updated via SSE
    function isPriceUpdated(product: Product): boolean {
        return currentPrices.has(product.id);
    }

</script>

<div class="container mx-auto px-10 py-5">
    <h2 class="text-xl font-semibold">Products</h2>
    
    <div class="mb-6">
        <div class="relative w-full md:w-96">
            <Search class="absolute left-3 top-3 h-4 w-4 text-gray-500" />
            <Input
                id="search"
                type="search"
                placeholder="Search products..."
                class="pl-10"
                bind:value={searchQuery}
            />
        </div>
    </div>

    <!-- Connection status indicator -->
    <div class="mb-4">
        {#if $connectionStatus === 'connected'}
            <Badge variant="outline" class="bg-green-50 text-green-700 border-green-200">
                🟢 Live price updates active
            </Badge>
        {:else if $connectionStatus === 'connecting'}
            <Badge variant="outline" class="bg-yellow-50 text-yellow-700 border-yellow-200">
                🟡 Connecting to price updates...
            </Badge>
        {:else if $connectionStatus === 'error'}
            <Badge variant="outline" class="bg-red-50 text-red-700 border-red-200">
                🔴 Connection error - Will retry automatically
            </Badge>
        {:else}
            <Badge variant="outline" class="bg-gray-50 text-gray-700 border-gray-200">
                ⚫ Price updates disconnected
            </Badge>
        {/if}
    </div>


    {#if $isLoading && $products.length === 0}
        <div class="flex justify-center items-center h-48">
            <div class="text-center">
                <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                <div class="text-gray-500 text-lg">Loading products...</div>
            </div>
        </div>
    {:else if $error}
        <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
            Error: {$error}
        </div>
    {:else}
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {#each $products as product (product.id)}
            <a href={`/${productHelpers.createProductSlug(product)}`} class="block">
                <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden transition-all duration-200 hover:shadow-lg hover:border-gray-300 hover:-translate-y-1">
                    <!-- Product Image -->
                    <div class="aspect-square bg-gray-100 overflow-hidden">
                        <img
                            src="{imageBaseUrl}{product.image_path}"
                            alt={product.name}
                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                            loading="lazy"
                        />
                    </div>

                    <!-- Product Info -->
                    <div class="p-4 space-y-3">
                        <!-- Product Name -->
                        <h3 class="font-semibold text-gray-900 group-hover:text-blue-600 transition-colors line-clamp-2 leading-tight">
                            {product.name}
                        </h3>

                        <!-- Brand -->
                        <p class="text-sm text-gray-500 font-medium">
                            {product.brand_name}
                        </p>

                        <!-- Description -->
                        <p class="text-sm text-gray-600 line-clamp-2 leading-relaxed">
                            {product.description}
                        </p>

                        <!-- Rating and Reviews -->
                        <div class="flex items-center text-sm text-gray-500">
                            <div class="flex items-center">
                                <Star class="w-4 h-4 text-yellow-400 fill-yellow-400" />
                                <span class="ml-1 font-medium">{product.average_rating}</span>
                            </div>
                            <span class="mx-2">•</span>
                            <span>{product.review_count} reviews</span>
                        </div>

                        <!-- Stock Status -->
                        <div class="flex justify-between items-center">
                            {#if isInStock(product)}
                                <Badge variant="outline" class="bg-green-50 text-green-700 border-green-200 text-xs">
                                    In Stock ({product.stock_quantity})
                                </Badge>
                            {:else}
                                <Badge variant="outline" class="bg-red-50 text-red-700 border-red-200 text-xs">
                                    Out of Stock
                                </Badge>
                            {/if}
                        </div>

                        <!-- Price -->
                        <div class="flex items-center justify-between pt-2">
                            <div class="flex items-center space-x-2">
                                <span class="text-xl font-bold text-gray-900 {isPriceUpdated(product) ? 'text-blue-600' : ''}">
                                    {getCurrentPrice(product)}
                                </span>
                            </div>
                            
                            {#if isInStock(product)}
                                <Button 
                                    size="sm" 
                                    class="opacity-0 group-hover:opacity-100 transition-opacity"
                                    onclick={(e) => {
                                        e.preventDefault();
                                        // Add to cart logic here
                                        console.log('Add to cart:', product.name);
                                    }}
                                >
                                    <ShoppingCart class="w-4 h-4" />
                                </Button>
                            {/if}
                        </div>
                    </div>
                </div>
            </a>
            {/each}
        </div>
    {/if}
</div>

<style>
    .line-clamp-2 {
        display: -webkit-box;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }
</style>