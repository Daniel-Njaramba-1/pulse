<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { toast } from "svelte-sonner";
    import * as Card from "$lib/components/ui/card/index";
    import { Button } from "$lib/components/ui/button/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Label } from "$lib/components/ui/label/index";
    import { DollarSign } from "lucide-svelte";
    import { productHelpers } from "$lib/stores/product";
    import type { ProductDetail } from "$lib/stores/product";

    const productId = $derived(page.params.id);

    // Main product data state
    let productDetail: ProductDetail | null = null;
    let isProductLoading = true;
    let productError: string | null = null;

    // Edit fields
    let editBasePrice = $state<number>(0);

    // Form processing states
    let isUpdatingPrice = $state<boolean>(false);

    async function fetchProductDetails() {
        isProductLoading = true;
        productError = null;

        try {
            const response = await productHelpers.getProduct(Number(productId));
            if (response.success && response.data) {
                productDetail = { ...response.data as ProductDetail };
                editBasePrice = productDetail.base_price || 0;
            } else {
                productError = response.error || "Failed to load product details";
            }
        } catch (err) {
            productError = err instanceof Error ? err.message : "An unexpected error occurred";
        } finally {
            isProductLoading = false;
        }
    }

    async function handleUpdatePrice() {
        if (editBasePrice < 0) {
            toast.error("Validation Error", {
                description: "Price cannot be negative."
            });
            return;
        }

        isUpdatingPrice = true;

        try {
            const result = await productHelpers.updateBasePrice(
                Number(productId),
                editBasePrice
            );

            if (result.success) {
                toast.success("Price Updated", {
                    description: `Base price set to $${editBasePrice.toFixed(2)}.`
                });
                await fetchProductDetails(); // Refresh data
            } else {
                toast.error("Price Update Failed", {
                    description: result.error || "Unknown error occurred."
                });
            }
        } catch (error) {
            toast.error("System Error", {
                description: "An unexpected error occurred while updating price."
            });
            console.error(error);
        } finally {
            isUpdatingPrice = false;
        }
    }

    onMount(async () => {
        await fetchProductDetails();
    });
</script>

<div class="container mx-auto py-2 px-40">
    <form onsubmit={handleUpdatePrice} class="space-y-8">
        <Card.Root class="p-6 space-y-6">
            <Card.Header>
                <Card.Title>Update Product Price</Card.Title>
                <Card.Description>
                    Change the base price for this product.
                </Card.Description>
            </Card.Header>

            <div class="space-y-4">
                <div class="space-y-2">
                    <Label for="editBasePrice">Base Price ($)</Label>
                    <div class="relative">
                        <DollarSign class="absolute left-3 top-1/2 transform -translate-y-1/2" />
                        <Input
                            id="editBasePrice"
                            type="number"
                            class="pl-10"
                            placeholder="Enter base price"
                            bind:value={editBasePrice}
                            min="0"
                        />
                    </div>
                </div>
            </div>

            <Card.Footer class="flex justify-end space-x-4">
                <Button variant="outline" onclick={() => fetchProductDetails()}>
                    Cancel
                </Button>
                <Button type="submit" disabled={isUpdatingPrice} class="bg-primary">
                    {isUpdatingPrice ? 'Updating...' : 'Update Price'}
                </Button>
            </Card.Footer>
        </Card.Root>
    </form>
</div>