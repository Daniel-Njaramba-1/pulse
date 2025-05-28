<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { goto } from "$app/navigation";

    import { isLoading, error, productHelpers } from "$lib/stores/product";
    import { cartHelpers } from "$lib/stores/cart";
    import type { Product, Review } from "$lib/stores/product";
    import { connectToSSE, disconnectFromSSE, connectionStatus } from "$lib/stores/pricing";
    
    import { Check, Star, ShoppingCart, Heart, Minus, Plus, MessageCircle, ThumbsUp, Send } from "lucide-svelte";
    import { Badge } from "$lib/components/ui/badge/index";
    import { Button } from "$lib/components/ui/button/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Textarea } from "$lib/components/ui/textarea/index";
    import { toast } from "svelte-sonner";
    import * as Card from "$lib/components/ui/card/index";
    
    const imageBaseUrl = "http://localhost:8080/assets/products/";
    let slug = $derived(page.params.slug);
    let productId = $state<number | null>();
    let isProductLoading = $state<boolean>(true);
    let productError: string | null = $state(null);
    let product: Product | null = $state(null);

    // Cart related states
    let quantity = $state<number>(1);
    let isAddingToCart = $state<boolean>(false);
    let addToCartError = $state<string | null>(null);
    let addToCartSuccess = $state<boolean>(false);
    
    // Wishlist states
    let isInWishlist = $state<boolean>(false);
    let isWishlistLoading = $state<boolean>(false);
    
    // Review states
    let showReviewForm = $state<boolean>(false);
    let newReview = $state({ rating: 0, comment: "" });
    let isSubmittingReview = $state<boolean>(false);
    let reviews = $state<Review[]>([]);
    let isReviewsLoading = $state<boolean>(false);
    let canReview = $state<boolean>(false); // Track if user can review (has purchased)
    
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
                    
                    // Load wishlist status and reviews
                    await Promise.all([
                        loadWishlistStatus(),
                        loadReviews(),
                        checkPurchaseStatus()
                    ]);
                }
            } catch (err) {
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

    async function loadWishlistStatus() {
        if (!productId) return;
        try {
            const response = await productHelpers.checkProductInWishlist(productId);
            if (response.success) {
                isInWishlist = response.inWishlist || false;
            } else {
                console.error("Error loading wishlist status:", response.error);
            }
        } catch (err) {
            console.error("Error loading wishlist status:", err);
        }
    }

    async function loadReviews() {
        if (!productId) return;
        isReviewsLoading = true;
        try {
            const response = await productHelpers.fetchProductReviews(productId);
            if (response.success) {
                reviews = response.data || [];
            } else {
                console.error("Error loading reviews:", response.error);
                toast.error("Failed to load reviews");
            }
        } catch (err) {
            console.error("Error loading reviews:", err);
            toast.error("Failed to load reviews");
        } finally {
            isReviewsLoading = false;
        }
    }

    async function checkPurchaseStatus() {
        if (!productId) return;
        try {
            const response = await productHelpers.verifyPurchase(productId);
            if (response.success) {
                canReview = response.purchased || false;
            } else {
                console.error("Error verifying purchase:", response.error);
            }
        } catch (err) {
            console.error("Error verifying purchase:", err);
        }
    }

    function isInStock(product: Product): boolean {
        return product.stock_quantity > 0;
    }

    function incrementQuantity() {
        if (product && quantity < product.stock_quantity) {
            quantity++;
        }
    }

    function decrementQuantity() {
        if (quantity > 1) {
            quantity--;
        }
    }

    async function handleAddToCart() {
        if (!product || !productId) return;
        
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

    async function toggleWishlist() {
        if (!productId) return;
        
        isWishlistLoading = true;
        try {
            let response;
            if (isInWishlist) {
                response = await productHelpers.removeFromWishlist(productId);
            } else {
                response = await productHelpers.addToWishlist(productId);
            }
            
            if (response.success) {
                isInWishlist = !isInWishlist;
                toast.success(response.message || (isInWishlist ? "Added to wishlist" : "Removed from wishlist"));
            } else {
                toast.error(response.error || "Failed to update wishlist");
            }
        } catch (err) {
            console.error("Error toggling wishlist:", err);
            toast.error("Failed to update wishlist");
        } finally {
            isWishlistLoading = false;
        }
    }

    async function submitReview() {
        if (!productId || !newReview.comment.trim()) return;
        
        isSubmittingReview = true;
        try {
            const response = await productHelpers.reviewProduct(
                productId,
                newReview.rating,
                newReview.comment
            );
            
            if (response.success) {
                toast.success(response.message || "Review submitted successfully");
                newReview = { rating: 0, comment: "" };
                showReviewForm = false;
                
                // Reload reviews to show the new one
                await loadReviews();
                
                // Reload product to update average rating and review count
                if (productId) {
                    const productResponse = await productHelpers.getProduct(productId);
                    if (productResponse.success) {
                        product = { ...productResponse.data as Product };
                    }
                }
            } else {
                toast.error(response.error || "Failed to submit review");
            }
        } catch (err) {
            console.error("Error submitting review:", err);
            toast.error("Failed to submit review");
        } finally {
            isSubmittingReview = false;
        }
    }

    function renderStars(rating: number, size: string = "h-4 w-4") {
        return Array.from({ length: 5 }, (_, i) => i < rating);
    }

    function formatDate(dateString: string): string {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        });
    }
</script>

<div class="container mx-auto px-6 py-8 max-w-7xl">
    <!-- Connection status indicator -->
    <div class="mb-6">
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

    {#if isProductLoading}
        <div class="flex justify-center items-center h-96">
            <div class="text-gray-500 animate-pulse text-xl">Loading product details...</div>
        </div>
    {:else if productError}
        <Card.Root class="w-full max-w-2xl mx-auto">
            <Card.Header>
                <Card.Title class="text-red-600">Error Loading Product</Card.Title>
            </Card.Header>
            <Card.Content class="space-y-4">
                <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
                    {productError}
                </div>
                <div class="flex justify-center">
                    <Button onclick={() => goto('/')} variant="outline">
                        Return to Home
                    </Button>
                </div>
            </Card.Content>
        </Card.Root>
    {:else if product}
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 mb-12">
            <!-- Product Image -->
            <div class="space-y-4">
                <Card.Root>
                    <Card.Content class="p-6">
                        <img
                            src="{imageBaseUrl}{product.image_path}"
                            alt={product.name}
                            class="w-full h-96 rounded-lg object-cover bg-gray-100"
                        />
                    </Card.Content>
                </Card.Root>
            </div>

            <!-- Product Details -->
            <div class="space-y-6">
                <div>
                    <h1 class="text-3xl font-bold text-gray-900 mb-2">{product.name}</h1>
                    <div class="flex items-center gap-4 text-sm text-gray-600 mb-4">
                        <div class="flex items-center">
                            {#each renderStars(Math.floor(product.average_rating)) as filled}
                                <Star class="h-4 w-4 {filled ? 'text-yellow-500 fill-current' : 'text-gray-300'}" />
                            {/each}
                            <span class="ml-2">{product.average_rating} / 5</span>
                        </div>
                        <span>•</span>
                        <span>{product.review_count} reviews</span>
                        <span>•</span>
                        <span class="font-medium">{product.brand_name}</span>
                    </div>
                </div>

                <div class="prose prose-gray">
                    <p class="text-gray-700 leading-relaxed">{product.description}</p>
                </div>

                <!-- Price -->
                <div class="bg-blue-50 p-4 rounded-lg">
                    <div class="text-3xl font-bold text-blue-700 mb-2">
                        ${product.adjusted_price}
                    </div>
                    <div class="flex items-center gap-3">
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
                </div>

                <!-- Quantity and Actions -->
                <div class="space-y-4">
                    <div class="flex items-center gap-4">
                        <label for="" class="text-sm font-medium text-gray-700">Quantity:</label>
                        <div class="flex items-center border border-gray-300 rounded-md">
                            <Button
                                variant="ghost"
                                size="sm"
                                onclick={decrementQuantity}
                                disabled={quantity <= 1}
                                class="h-10 w-10 p-0"
                            >
                                <Minus class="h-4 w-4" />
                            </Button>
                            <span class="w-12 text-center text-sm font-medium">{quantity}</span>
                            <Button
                                variant="ghost"
                                size="sm"
                                onclick={incrementQuantity}
                                disabled={!product || quantity >= product.stock_quantity}
                                class="h-10 w-10 p-0"
                            >
                                <Plus class="h-4 w-4" />
                            </Button>
                        </div>
                    </div>

                    <div class="flex gap-3">
                        <Button 
                            onclick={handleAddToCart} 
                            disabled={!isInStock(product) || isAddingToCart || quantity < 1}
                            class="flex-1"
                            size="lg"
                        >
                            {#if isAddingToCart}
                                Adding...
                            {:else if addToCartSuccess}
                                <Check class="h-5 w-5 mr-2" />
                                Added to Cart
                            {:else}
                                <ShoppingCart class="h-5 w-5 mr-2" />
                                Add to Cart
                            {/if}
                        </Button>

                        <Button
                            variant="outline"
                            onclick={toggleWishlist}
                            disabled={isWishlistLoading}
                            size="lg"
                            class="px-4"
                        >
                            <Heart class="h-5 w-5 {isInWishlist ? 'fill-current text-red-500' : ''}" />
                        </Button>
                    </div>

                    {#if addToCartError}
                        <div class="text-red-600 text-sm bg-red-50 p-3 rounded-md">
                            {addToCartError}
                        </div>
                    {/if}
                </div>
            </div>
        </div>

        <!-- Reviews Section -->
        <div class="border-t pt-12">
            <div class="flex items-center justify-between mb-8">
                <h2 class="text-2xl font-bold text-gray-900">Customer Reviews</h2>
                {#if canReview}
                    <Button onclick={() => showReviewForm = !showReviewForm} variant="outline">
                        <MessageCircle class="h-4 w-4 mr-2" />
                        Write a Review
                    </Button>
                {:else}
                    <div class="text-sm text-gray-500">
                        Purchase this product to write a review
                    </div>
                {/if}
            </div>

            <!-- Review Form -->
            {#if showReviewForm}
                <Card.Root class="mb-8">
                    <Card.Header>
                        <Card.Title>Write Your Review</Card.Title>
                    </Card.Header>
                    <Card.Content class="space-y-4">
                        <div>
                            <label for="" class="block text-sm font-medium mb-2">Rating</label>
                            <div class="flex gap-1">
                                {#each Array.from({ length: 5 }, (_, i) => i + 1) as rating}
                                    <button
                                        onclick={() => newReview.rating = rating}
                                        class="p-1 hover:scale-110 transition-transform"
                                    >
                                        <Star class="h-6 w-6 {rating <= newReview.rating ? 'text-yellow-500 fill-current' : 'text-gray-300'}" />
                                    </button>
                                {/each}
                            </div>
                        </div>

                        <div>
                            <label for="review-comment" class="block text-sm font-medium mb-2">Your Review</label>
                            <Textarea
                                id="review-comment"
                                bind:value={newReview.comment}
                                placeholder="Share your experience with this product..."
                                rows={4}
                                class="w-full"
                            />
                        </div>

                        <div class="flex gap-3">
                            <Button
                                onclick={submitReview}
                                disabled={isSubmittingReview || !newReview.comment.trim()}
                            >
                                {#if isSubmittingReview}
                                    Submitting...
                                {:else}
                                    <Send class="h-4 w-4 mr-2" />
                                    Submit Review
                                {/if}
                            </Button>
                            <Button variant="outline" onclick={() => showReviewForm = false}>
                                Cancel
                            </Button>
                        </div>
                    </Card.Content>
                </Card.Root>
            {/if}

            <!-- Reviews List -->
            {#if isReviewsLoading}
                <div class="flex justify-center items-center py-12">
                    <div class="text-gray-500 animate-pulse">Loading reviews...</div>
                </div>
            {:else}
                <div class="space-y-6">
                    {#each reviews as review (review.id)}
                        <Card.Root>
                            <Card.Content class="p-6">
                                <div class="flex items-start justify-between mb-4">
                                    <div>
                                        <div class="flex items-center gap-2 mb-2">
                                            <span class="font-medium text-gray-900">{review.customer_name}</span>
                                            <span class="text-gray-500 text-sm">•</span>
                                            <span class="text-gray-500 text-sm">{review.created_at}</span>
                                        </div>
                                        <div class="flex items-center">
                                            {#each renderStars(review.rating) as filled}
                                                <Star class="h-4 w-4 {filled ? 'text-yellow-500 fill-current' : 'text-gray-300'}" />
                                            {/each}
                                        </div>
                                    </div>
                                </div>
                                
                                <p class="text-gray-700 mb-4 leading-relaxed">{review.review_text}</p>
                            </Card.Content>
                        </Card.Root>
                    {/each}

                    {#if reviews.length === 0}
                        <div class="text-center py-12 text-gray-500">
                            <MessageCircle class="h-12 w-12 mx-auto mb-4 text-gray-300" />
                            <p class="text-lg">No reviews yet</p>
                            <p class="text-sm">Be the first to share your experience!</p>
                        </div>
                    {/if}
                </div>
            {/if}
        </div>
    {/if}
</div>