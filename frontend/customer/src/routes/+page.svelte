<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { products, isLoading, error, productHelpers } from "$lib/stores/product";
    import type { Product } from "$lib/stores/product";
    import { connectToSSE, disconnectFromSSE, connectionStatus } from "$lib/stores/pricing";
    import { Search, Star, ShoppingCart } from "lucide-svelte";
    import { Button } from "$lib/components/ui/button/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Badge } from "$lib/components/ui/badge/index";

    // Local state
    let searchQuery = $state<string>("");

    // Base URL for product images
    const imageBaseUrl = "http://localhost:8080/assets/products/";


    onMount(() => {
        (async () => {
            const productsResponse = await productHelpers.fetchProducts();
            console.log("Products fetched:", productsResponse);
            if (!productsResponse.success) {
                console.error("Failed to fetch products:", productsResponse.error);
            }
        })();
    
        connectToSSE();
    
        return () => {
            disconnectFromSSE();
        };
    });

    // Helper function to check if product is in stock
    function isInStock(product: Product): boolean {
        return product.stock_quantity > 0;
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


    {#if $isLoading && $products.length === 0}
        <div class="flex justify-center items-center h-48">
            <div class="text-gray-500">Loading products...</div>
        </div>
    {:else if $error}
        <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
            Error: {$error}
        </div>
    {:else}
        {#each $products as product (product.id)}
        <a href={`/${productHelpers.createProductSlug(product)}`} class="block">
            <div class="flex flex-col sm:flex-row gap-4 p-4 mb-5 border rounded-lg hover:shadow-md transition-shadow">
                <!-- Product image -->
                <div class="flex-shrink-0">
                    <img
                        src="{imageBaseUrl}{product.image_path}"
                        alt={product.name}
                        class="w-28 h-28 rounded-md object-cover bg-gray-100"
                    />
                </div>

                <!-- Product information -->
                <div class="flex-grow space-y-2">
                    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
                        <h3 class="text-xl font-semibold">{product.name}</h3>
                        <div class="text-lg font-bold text-blue-700">
                            {product.adjusted_price}
                        </div>
                    </div>

                    <p class="text-gray-600">{product.description}</p>
                    
                    <div class="flex items-center text-sm text-gray-500">
                        <span class="flex items-center">
                            <Star class="w-4 h-4 text-yellow-400 fill-yellow-400" />
                            <span class="ml-1">{product.average_rating}</span>
                        </span>
                        <span class="mx-2">•</span>
                        <span>{product.review_count} reviews</span>
                        <span class="mx-2">•</span>
                        <span>{product.brand_name}</span>
                    </div>

                    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 pt-2">
                        {#if isInStock(product)}
                            <Badge variant="outline" class="bg-green-50 text-green-700 border-green-200">
                                In Stock ({product.stock_quantity})
                            </Badge>
                        {:else}
                            <Badge variant="outline" class="bg-red-50 text-red-700 border-red-200">
                                Out of Stock
                            </Badge>
                        {/if}
                    </div>
                </div>
            </div>
        </a>
        {/each}
    {/if}
</div>