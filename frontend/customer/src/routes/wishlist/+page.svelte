<script lang="ts">
    import { onMount } from "svelte";
    import { wishlistHelpers, type WishlistItem } from "$lib/stores/wishlist";
    import { Badge } from "$lib/components/ui/badge/index";
    import { Button } from "$lib/components/ui/button/index";
    import { toast } from "svelte-sonner";

    let wishlistError: string | null = null;
    let wishlistItems: WishlistItem[] = [];
    let isLoading = true;

    onMount(async () => {
        isLoading = true;
        wishlistError = null;
        try {
            const response = await wishlistHelpers.fetchWishlist();
            if (!response.success) {
                wishlistError = response.error || "Failed to load wishlist";
            } else if (response.data) {
                wishlistItems = response.data.items || [];
            }
        } catch (err) {
            wishlistError = "An unexpected error occurred";
        }
        isLoading = false;
    });

    async function removeItem(productId: number) {
        const res = await wishlistHelpers.removeFromWishlist(productId);
        if (res.success) {
            wishlistItems = wishlistItems.filter(item => item.product_id !== productId);
            toast.success("Removed from wishlist");
        } else {
            toast.error(res.error || "Failed to remove");
        }
    }
</script>

<div class="container mx-auto px-10 py-5">
    <h2 class="text-xl font-semibold mb-4">My Wishlist</h2>

    {#if isLoading}
        <div>Loading wishlist...</div>
    {:else if wishlistError}
        <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
            Error: {wishlistError}
        </div>
    {:else if wishlistItems.length === 0}
        <div class="text-gray-500">Your wishlist is empty.</div>
    {:else}
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {#each wishlistItems as item (item.product_id)}
                <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                    <div class="aspect-square bg-gray-100 overflow-hidden">
                        <img
                            src={`http://localhost:8080/assets/products/${item.product_image_path}`}
                            alt={item.product_name}
                            class="w-full h-full object-cover"
                            loading="lazy"
                        />
                    </div>
                    <div class="p-4 space-y-2">
                        <h3 class="font-semibold text-gray-900">{item.product_name}</h3>
                        <div class="flex items-center justify-between">
                            <span class="text-lg font-bold text-gray-900">${item.product_adjusted_price.toFixed(2)}</span>
                            <Badge class={item.product_stock_quantity > 0 ? "bg-green-50 text-green-700" : "bg-red-50 text-red-700"}>
                                {item.product_stock_quantity > 0 ? "In Stock" : "Out of Stock"}
                            </Badge>
                        </div>
                        <Button size="sm" variant="destructive" onclick={() => removeItem(item.product_id)}>
                            Remove
                        </Button>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>
