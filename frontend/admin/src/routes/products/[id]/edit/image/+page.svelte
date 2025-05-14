<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { toast } from "svelte-sonner";
    import { Button } from "$lib/components/ui/button/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Label } from "$lib/components/ui/label/index";
    import * as Card from "$lib/components/ui/card/index";
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
    let image_path = $state<string>("");
    let imageFile = $state<File | null>(null);
    let imagePreview = $state<string | null>(null);

    // Form processing states
    let isUpdatingImage = $state<boolean>(false);

    // Base URL for product images
    const imageBaseUrl = "http://localhost:8080/assets/products/";

    async function fetchProductDetails() {
        isProductLoading = true;
        productError = null;

        try {
            const response = await productHelpers.getProduct(Number(productId));
            if (response.success && response.data) {
                productDetail = { ...response.data as ProductDetail };
                image_path = productDetail.image_path || "";
                imagePreview = image_path ? `${imageBaseUrl}${image_path}` : null;
            } else {
                productError = response.error || "Failed to load product details";
            }
        } catch (err) {
            productError = err instanceof Error ? err.message : "An unexpected error occurred";
        } finally {
            isProductLoading = false;
        }
    }

    function handleImageChange(event: Event) {
        const target = event.target as HTMLInputElement;
        if (target.files && target.files[0]) {
            imageFile = target.files[0];
            imagePreview = URL.createObjectURL(target.files[0]);
        }
    }

    async function handleUpdateImage() {
        if (!imageFile) {
            toast.error("Validation Error", {
                description: "Please select an image to upload."
            });
            return;
        }

        isUpdatingImage = true;

        try {
            const result = await productHelpers.updateProductImage(
                Number(productId),
                imageFile
            );

            if (result.success) {
                toast.success("Image Updated", {
                    description: "Product image has been updated successfully."
                });
                await fetchProductDetails(); // Refresh data
            } else {
                toast.error("Image Update Failed", {
                    description: result.error || "Unknown error occurred."
                });
            }
        } catch (error) {
            toast.error("System Error", {
                description: "An unexpected error occurred while updating the image."
            });
            console.error(error);
        } finally {
            isUpdatingImage = false;
        }
    }

    onMount(async () => {
        await fetchProductDetails();
    });
</script>

<!-- Image Edit Card -->
<div class="container mx-auto py-2 px-40">
    <form onsubmit={handleUpdateImage} class="space-y-8"> 
        <Card.Root class="p-6 space-y-6">
            <Card.Header>
                <Card.Title>Update Product Image</Card.Title>
                <Card.Description>
                    Upload a new image for this product.
                </Card.Description>
            </Card.Header>

            <div class="space-y-6">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6 items-center">
                    <div class="flex flex-col gap-3">
                        <div class="border rounded-md p-3">
                            <Label for="image" class="flex flex-col items-center gap-2 cursor-pointer py-4">
                                <span class="text-sm text-muted-foreground">Click to upload or Drag and Drop</span>
                                <span class="text-sm text-muted-foreground">SVG, PNG, JPG, GIF (2MB max)</span>
                                <Input id="image" type="file" accept="image/*" class="hidden" onchange={handleImageChange} />
                            </Label>
                        </div>
                        <p class="text-sm text-gray-500">
                            {imageFile ? `Selected file: ${imageFile.name}` : 
                            image_path ? `Current image: ${image_path}` : 
                            "No image currently set"}
                        </p>
                    </div>
                    <div class="flex justify-center">
                        {#if imagePreview}
                            <div class="relative w-48 h-48 border rounded-md overflow-hidden bg-gray-50">
                                <img
                                    src={imagePreview}
                                    alt="Product Preview"
                                    class="object-contain w-full h-full"
                                />
                            </div>
                        {:else}
                            <div class="w-48 h-48 border rounded-md flex items-center justify-center bg-gray-50">
                                <span class="text-sm text-muted-foreground">No image preview</span>
                            </div>
                        {/if}
                    </div>
                </div>
            </div>

            <Card.Footer class="flex justify-end space-x-4">
                <Button variant="outline" onclick={() => { imageFile = null; imagePreview = null; }}>
                    Cancel
                </Button>
                <Button type="submit" disabled={isUpdatingImage} class="bg-primary">
                    {isUpdatingImage ? 'Uploading...' : 'Update Image'}
                </Button>
            </Card.Footer>
        </Card.Root>
    </form>
</div>
