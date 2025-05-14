<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { toast } from "svelte-sonner";
    import { Button } from "$lib/components/ui/button/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Label } from "$lib/components/ui/label/index";
    import * as Card from "$lib/components/ui/card/index";
    import { products, productHelpers } from "$lib/stores/product";
    import type { ProductDetail } from "$lib/stores/product";

    const productId = $derived(page.params.id);

    // Main product data state
    let productDetail: ProductDetail | null = $state(null);
    let isProductLoading = $state<boolean>(true);
    let productError = $state<string | null>(null);

    // Product display fields
    let stock_quantity = $state<number>(0);
    let restockQuantity = $state<number>(0);

    // Form processing states
    let isRestocking = $state<boolean>(false);

    async function fetchProductDetails() {
        isProductLoading = true;
        productError = null;

        try {
            const response = await productHelpers.getProduct(Number(productId));
            if (response.success && response.data) {
                productDetail = { ...response.data as ProductDetail };
                stock_quantity = productDetail.stock_quantity || 0;
            } else {
                productError = response.error || "Failed to load product details";
            }
        } catch (err) {
            productError = err instanceof Error ? err.message : "An unexpected error occurred";
        } finally {
            isProductLoading = false;
        }
    }

    async function handleRestock() {
        if (restockQuantity <= 0) {
            toast.error("Validation Error", {
                description: "Restock quantity must be greater than zero."
            });
            return;
        }

        isRestocking = true;

        try {
            const result = await productHelpers.restock(
                Number(productId),
                restockQuantity
            );

            if (result.success) {
                toast.success("Inventory Updated", {
                    description: `Added ${restockQuantity} units to inventory.`
                });
                await fetchProductDetails(); // Refresh data
                restockQuantity = 0; // Reset input
            } else {
                toast.error("Restock Failed", {
                    description: result.error || "Unknown error occurred."
                });
            }
        } catch (error) {
            toast.error("System Error", {
                description: "An unexpected error occurred while restocking."
            });
            console.error(error);
        } finally {
            isRestocking = false;
        }
    }

    onMount(async () => {
        await fetchProductDetails();
    });
</script>

<div class="container mx-auto py-2 px-40">
    <form onsubmit={handleRestock} class="space-y-8">
        <Card.Root class="p-6 space-y-6">
            <Card.Header>
                <Card.Title>Manage Inventory</Card.Title>
                <Card.Description>
                    Update inventory levels and configure stock thresholds.
                </Card.Description>
            </Card.Header>

            <div class="space-y-4">
                <div class="space-y-2">
                    <Label for="restockQuantity">Restock Quantity</Label>
                    <Input
                        id="restockQuantity"
                        type="number"
                        class="w-full"
                        placeholder="Enter quantity to add"
                        bind:value={restockQuantity}
                        min="1"
                    />
                    <p class="text-xs text-gray-500">Current stock: {stock_quantity}</p>
                </div>
            </div>

            <Card.Footer class="flex justify-end space-x-3">
                <Button variant="outline" onclick={() => (restockQuantity = 0)}>
                    Cancel
                </Button>
                <Button type="submit" disabled={isRestocking} size="sm">
                    {isRestocking ? "Adding..." : "Add Stock"}
                </Button>
            </Card.Footer>
        </Card.Root>
    </form>
</div>