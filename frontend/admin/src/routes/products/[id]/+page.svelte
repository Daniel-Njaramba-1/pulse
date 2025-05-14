<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import * as Card from "$lib/components/ui/card/index";
    import { toast } from "svelte-sonner";
    import { Button } from "$lib/components/ui/button/index";

    import { 
        Pencil, 
        SquareArrowLeft 
    } from "lucide-svelte";

    import { 
        products, 
        productHelpers,
    } from "$lib/stores/product";

    import type { ProductDetail } from "$lib/stores/product";

    const productId = $derived(page.params.id);

    // Main product data state
    let productDetail: ProductDetail | null = $state(null);
    let isProductLoading = $state<boolean>(true);
    let productError = $state<string | null>(null);

    // Product display fields
    let name = $state<string>(""); 
    let description = $state<string>("");
    let is_active = $state<boolean>(false);
    let image_path = $state<string>("");
    let brand = $state<string>("");
    let category = $state<string>("");
    let average_rating = $state<number>(0);
    let review_count = $state<number>(0);
    let wishlist_count = $state<number>(0);
    let base_price = $state<number>(0);
    let adjusted_price = $state<number>(0);
    let stock_quantity = $state<number>(0);
    let stock_threshold = $state<number>(0);

    // Base URL for product images
    const imageBaseUrl = "http://localhost:8080/assets/products/";

    // Product image url
    let imageUrl = $state<string>("");

        async function fetchProductDetails() {
        isProductLoading = true;
        productError = null;
        
        try {
            const response = await productHelpers.getProduct(Number(productId));
            if (response.success && response.data) {
                productDetail = {...response.data as ProductDetail};
                
                // Set display values
                name = productDetail.name;
                description = productDetail.description || "";
                is_active = productDetail.is_active;
                image_path = productDetail.image_path || "";
                brand = productDetail.brand_name;
                category = productDetail.category_name;
                average_rating = productDetail.average_rating || 0;
                review_count = productDetail.review_count || 0;
                wishlist_count = productDetail.wishlist_count || 0;
                base_price = productDetail.base_price || 0;
                adjusted_price = productDetail.adjusted_price || 0;
                stock_quantity = productDetail.stock_quantity || 0;
                stock_threshold = productDetail.stock_threshold || 0;

                imageUrl = image_path ? imageBaseUrl + image_path : "";
                
                console.log("Product details fetched successfully:", productDetail);
                console.log("Image URL:", imageUrl);
            } else {
                productError = response.error || "Failed to load product details";
            }
        } catch (err) {
            productError = err instanceof Error ? err.message : "An unexpected error occurred";
        } finally {
            isProductLoading = false;
        }
    }
    
    async function handleActivateDeactivate() {
        try {
            let result;
            
            if (is_active) {
                result = await productHelpers.deactivateProduct(Number(productId));
            } else {
                result = await productHelpers.reactivateProduct(Number(productId));
            }
            
            if (result.success) {
                toast.success(is_active ? "Product Deactivated" : "Product Activated", {
                    description: is_active ? 
                        "Product has been deactivated and is no longer visible to customers." :
                        "Product has been activated and is now visible to customers."
                });
                await fetchProductDetails(); // Refresh data
            } else {
                toast.error("Status Update Failed", {
                    description: result.error || "Unknown error occurred."
                });
            }
        } catch (error) {
            toast.error("System Error", {
                description: "An unexpected error occurred while updating status."
            });
            console.error(error);
        }
    }

    onMount(async () => {
        await fetchProductDetails();
    });
</script>

<div class="container mx-auto py-6 px-4 max-w-6xl">
    <!-- Back button -->
    <div class="mb-6">
        <Button variant="outline" class="flex items-center gap-2" onclick={() => goto('/products')}>
            <SquareArrowLeft size={18} />
            <span>Back to Products</span>
        </Button>
    </div>

    {#if isProductLoading}
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
                    {productError}
                </div>
                <div class="mt-4 flex justify-center">
                    <Button onclick={() => goto('/products')} variant="outline">
                        Return to Products
                    </Button>
                </div>
            </Card.Content>
        </Card.Root>
    {:else if productDetail}
        <!-- Product Header -->
        <div class="mb-8">
            <div class="flex items-center justify-between">
                <h1 class="text-3xl font-bold text-gray-800 mb-2">{name}</h1>
                <Button 
                    variant="outline" 
                     
                    onclick={handleActivateDeactivate}
                    class={`px-3 py-1 rounded-full text-sm font-medium ${
                        is_active ? "border-red-200 text-red-800 hover:bg-red-50" : 
                        "border-green-200 text-green-800 hover:bg-green-50"
                    }`}
                >
                    {is_active ? "Deactivate" : "Activate"}
                </Button>
            </div>
            <div class="flex items-center gap-3">
                <span class={`px-3 py-1 rounded-full text-sm font-medium ${is_active ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-800"}`}>
                    {is_active ? "Active" : "Inactive"}
                </span>
                <span class="text-gray-500">ID: {productId}</span>
            </div>
        </div>

        <!-- Main Content Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <!-- Left Column: Image -->
            <div class="lg:col-span-1">
                <Card.Root>
                    <Card.Header class="flex flex-row items-center justify-between p-4 pb-0">
                        <Card.Title class="text-lg">Product Image</Card.Title>
                        <Button 
                            variant="outline" 
                            class="flex items-center gap-2" 
                            onclick={() => goto(`/products/${productId}/edit/image`)}>
                            <Pencil size={18} />
                            <span>Edit Image</span>
                        </Button>
                    </Card.Header>
                    <Card.Content class="p-4">
                        {#if image_path}
                            <img 
                                src={imageUrl}
                                alt="{name}" 
                                class="w-full aspect-square rounded-md object-cover shadow-sm" 
                            />
                        {:else}
                            <div class="w-full aspect-square bg-gray-200 rounded-md flex items-center justify-center">
                                <span class="text-gray-500">No image available</span>
                            </div>
                        {/if}
                    </Card.Content>
                </Card.Root>
            </div>
            
            <!-- Right Column: Product Info -->
            <div class="lg:col-span-2">
                <!-- Product Details Card -->
                <Card.Root class="mb-6">
                    <Card.Header class="flex flex-row items-center justify-between pb-1">
                        <Card.Title class="text-xl font-semibold">Product Details</Card.Title>
                        <Button 
                            variant="outline" 
                            class="flex items-center gap-2" 
                            onclick={() => goto(`/products/${productId}/edit/details`)}
                        >
                            <Pencil size={18} />
                            <span>Edit Product</span>
                        </Button>
                    </Card.Header>
                    <Card.Content>
                        <div class="space-y-4">
                            <div>
                                <h3 class="text-sm font-medium text-gray-500 mb-1">Description</h3>
                                <p class="text-gray-800">{description || "No description available"}</p>
                            </div>
                            <div class="grid grid-cols-2 gap-4">
                                <div>
                                    <h3 class="text-sm font-medium text-gray-500 mb-1">Brand</h3>
                                    <p class="text-gray-800 font-medium">{brand || "Unknown"}</p>
                                </div>
                                <div>
                                    <h3 class="text-sm font-medium text-gray-500 mb-1">Category</h3>
                                    <p class="text-gray-800 font-medium">{category || "Uncategorized"}</p>
                                </div>
                            </div>
                        </div>
                    </Card.Content>
                </Card.Root>

                <!-- Product Metrics Card -->
                <Card.Root class="mb-6">
                    <Card.Header class="pb-1">
                        <Card.Title class="text-xl font-semibold">Product Metrics</Card.Title>
                    </Card.Header>
                    <Card.Content>
                        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
                            <div class="bg-gray-50 p-4 rounded-md">
                                <h3 class="text-sm font-medium text-gray-500 mb-1">Average Rating</h3>
                                <p class="text-2xl font-bold text-gray-800">{average_rating?.toFixed(1) || "N/A"}</p>
                                <p class="text-sm text-gray-500">{review_count || 0} reviews</p>
                            </div>
                            <div class="bg-gray-50 p-4 rounded-md">
                                <h3 class="text-sm font-medium text-gray-500 mb-1">Wishlist</h3>
                                <p class="text-2xl font-bold text-gray-800">{wishlist_count || 0}</p>
                                <p class="text-sm text-gray-500">users</p>
                            </div>
                            <div class="bg-gray-50 p-4 rounded-md">
                                <h3 class="text-sm font-medium text-gray-500 mb-1">Reviews</h3>
                                <p class="text-2xl font-bold text-gray-800">{review_count || 0}</p>
                                <p class="text-sm text-gray-500">total</p>
                            </div>
                        </div>
                    </Card.Content>
                </Card.Root>
            </div>
        </div>

        <!-- Bottom Cards Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8">
            <!-- Inventory Card -->
            <Card.Root>
                <Card.Header class="flex flex-row items-center justify-between pb-3">
                    <Card.Title class="text-xl font-semibold">Inventory</Card.Title>
                    <Button 
                        variant="outline" 
                        class="flex items-center gap-2" 
                        onclick={() => goto(`/products/${productId}/inventory`)}
                    >
                        <Pencil size={18} />
                        <span>Edit Inventory</span>
                    </Button>
                </Card.Header>
                <Card.Content>
                    <div class="grid grid-cols-2 gap-6">
                        <div>
                            <h3 class="text-sm font-medium text-gray-500 mb-1">Stock Quantity</h3>
                            <div class="flex items-end gap-2">
                                <p class="text-2xl font-bold text-gray-800">{stock_quantity || 0}</p>
                                <p class="text-sm text-gray-500 mb-1">units</p>
                            </div>
                        </div>
                        <div>
                            <h3 class="text-sm font-medium text-gray-500 mb-1">Stock Threshold</h3>
                            <div class="flex items-end gap-2">
                                <p class="text-2xl font-bold text-gray-800">{stock_threshold || 0}</p>
                                <p class="text-sm text-gray-500 mb-1">units</p>
                            </div>
                        </div>
                    </div>
                    
                    <div class="mt-4">
                        <div class="h-2 bg-gray-200 rounded-full overflow-hidden">
                            {#if stock_threshold > 0}
                                <div class="h-full {stock_quantity > stock_threshold * 2 ? 'bg-green-500' : stock_quantity >= stock_threshold ? 'bg-yellow-500' : 'bg-red-500'}" 
                                    style="width: {Math.min(stock_quantity / (stock_threshold * 3) * 100, 100)}%"></div>
                            {:else}
                                <div class="h-full bg-gray-400" style="width: 100%"></div>
                            {/if}
                        </div>
                        <div class="mt-2 text-sm">
                            {#if stock_quantity <= 0}
                                <span class="text-red-600 font-medium">Out of stock</span>
                            {:else if stock_threshold > 0 && stock_quantity < stock_threshold}
                                <span class="text-red-600 font-medium">Low stock</span>
                            {:else if stock_threshold > 0 && stock_quantity < stock_threshold * 2}
                                <span class="text-yellow-600 font-medium">Medium stock</span>
                            {:else}
                                <span class="text-green-600 font-medium">Well stocked</span>
                            {/if}
                        </div>
                    </div>
                </Card.Content>
            </Card.Root>

            <!-- Pricing Card -->
            <Card.Root>
                <Card.Header class="flex flex-row items-center justify-between pb-3">
                    <Card.Title class="text-xl font-semibold">Pricing</Card.Title>
                    <Button 
                        variant="outline" 
                        class="flex items-center gap-2" 
                        onclick={() => goto(`/products/${productId}/pricing`)}
                    >
                        <Pencil size={18} />
                        <span>Edit Pricing</span>
                    </Button>
                </Card.Header>
                <Card.Content>
                    <div class="grid grid-cols-2 gap-6">
                        <div>
                            <h3 class="text-sm font-medium text-gray-500 mb-1">Base Price</h3>
                            <div class="flex items-end gap-1">
                                <p class="text-2xl font-bold text-gray-800">${base_price?.toFixed(2) || '0.00'}</p>
                            </div>
                        </div>
                        <div>
                            <h3 class="text-sm font-medium text-gray-500 mb-1">Adjusted Price</h3>
                            <div class="flex items-end gap-1">
                                <p class="text-2xl font-bold {adjusted_price < base_price ? 'text-green-600' : 'text-gray-800'}">
                                    ${adjusted_price?.toFixed(2) || '0.00'}
                                </p>
                            </div>
                        </div>
                    </div>
                    
                    {#if adjusted_price !== base_price && base_price > 0}
                        <div class="mt-4 p-3 bg-gray-50 rounded-md">
                            <div class="flex justify-between items-center">
                                <span class="text-sm font-medium text-gray-600">Delta</span>
                                <span class="text-sm font-medium 
                                    {adjusted_price < base_price ? 'text-green-600' : 'text-red-600'}">
                                    {adjusted_price < base_price ? '-' : '+'}
                                    {Math.abs(((base_price - adjusted_price) / base_price * 100)).toFixed(1)}%
                                </span>
                            </div>
                        </div>
                    {/if}
                </Card.Content>
            </Card.Root>
        </div>
    {:else}
        <div class="text-center text-gray-500 py-16">
            <p class="text-lg">Product not found or an error occurred.</p>
            <Button class="mt-4" variant="outline" onclick={() => goto('/products')}>
                Return to Products
            </Button>
        </div>
    {/if}
</div>
