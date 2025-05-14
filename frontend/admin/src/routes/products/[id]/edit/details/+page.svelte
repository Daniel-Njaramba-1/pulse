<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { toast } from "svelte-sonner";
    import * as Card from "$lib/components/ui/card/index";
    import * as Select from "$lib/components/ui/select/index.js";
    import { Button } from "$lib/components/ui/button/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Textarea } from "$lib/components/ui/textarea/index";
    import { Label } from "$lib/components/ui/label/index";
    import { Switch } from "$lib/components/ui/switch/index";
    import { products, productHelpers } from "$lib/stores/product";
    import { brands } from "$lib/stores/brands";
    import { categories } from "$lib/stores/category";
    import type { ProductDetail } from "$lib/stores/product";
    import { goto } from "$app/navigation";
    import { SquareArrowLeft } from "lucide-svelte";

    const productId = $derived(page.params.id);

    let productDetail: ProductDetail | null = null;
    let isProductLoading = true;
    let productError: string | null = null;

    let editName = $state<string>("");
    let editDescription = $state<string>("");
    let editIsActive = $state<boolean>(false);
    let selectedBrandId = $state<number | null>(null);
    let selectedCategoryId = $state<number | null>(null);

    let isSaving = $state<boolean>(false);

    async function fetchProductDetails() {
        isProductLoading = true;
        productError = null;

        try {
            const response = await productHelpers.getProduct(Number(productId));
            if (response.success && response.data) {
                productDetail = response.data as ProductDetail;
                setEditValues();
            } else {
                productError = response.error || "Failed to load product details";
            }
        } catch (err) {
            productError = err instanceof Error ? err.message : "An unexpected error occurred";
        } finally {
            isProductLoading = false;
        }
    }

    function setEditValues() {
        if (!productDetail) return;
        editName = productDetail.name;
        editDescription = productDetail.description || "";
        editIsActive = productDetail.is_active;
        selectedBrandId = productDetail.brand_id;
        selectedCategoryId = productDetail.category_id;
    }

    async function handleSaveBasicInfo() {
        if (!selectedBrandId || !selectedCategoryId) {
            toast.error("Validation Error", {
                description: "Please select both brand and category."
            });
            return;
        }

        isSaving = true;

        try {
            const result = await productHelpers.updateProductDetails(
                Number(productId),
                selectedBrandId,
                selectedCategoryId,
                editName,
                editDescription,
                editIsActive
            );

            if (result.success) {
                toast.success("Product Updated", {
                    description: "Basic information has been saved successfully."
                });
                goto(`/products/${productId}`);
            } else {
                toast.error("Update Failed", {
                    description: result.error || "Unknown error occurred."
                });
            }
        } catch (error) {
            toast.error("System Error", {
                description: "An unexpected error occurred while saving."
            });
        } finally {
            isSaving = false;
        }
    }

    onMount(fetchProductDetails);
</script>
<div class="container mx-auto py-2 px-40">
    <form onsubmit={handleSaveBasicInfo} class="space-y-8"> 
        <Card.Root class="p-6 border rounded-md">
            <Card.Header>
                <Card.Title>Edit Product Information</Card.Title>
                <Card.Description>
                    Make changes to the product's basic information.
                </Card.Description>
            </Card.Header>
            <Card.Content>
                <div class="space-y-6 py-4">
                    <div class="space-y-2">
                        <Label for="editName">Product Name</Label>
                        <Input 
                            id="editName" 
                            class="w-full" 
                            placeholder="Enter product name" 
                            bind:value={editName} 
                        />
                    </div>

                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div class="space-y-2">
                            <Label for="brand">Brand</Label>
                            <Select.Root type="single" value={selectedBrandId?.toString()} onValueChange={(value) => selectedBrandId = Number(value)}>                           >
                                <Select.Trigger class="w-full">
                                    {$brands.find((b) => b.id === selectedBrandId)?.name || "Select Brand"}
                                </Select.Trigger>
                                <Select.Content>
                                    <Select.Group>
                                        <Select.GroupHeading>Brands</Select.GroupHeading>
                                        {#each $brands as brand (brand.id)}
                                            <Select.Item value={brand.id.toString()} label={brand.name} />
                                        {/each}
                                    </Select.Group>
                                </Select.Content>
                            </Select.Root>
                        </div>
                        <div class="space-y-2">
                            <Label for="category">Category</Label>
                            <Select.Root type="single" value={selectedCategoryId?.toString()} onValueChange={(value) => selectedCategoryId = Number(value)}>
                                <Select.Trigger class="w-full">
                                    {$categories.find((c) => c.id === selectedCategoryId)?.name || "Select Category"}
                                </Select.Trigger>
                                <Select.Content>
                                    <Select.Group>
                                        <Select.GroupHeading>Categories</Select.GroupHeading>
                                        {#each $categories as category (category.id)}
                                            <Select.Item value={category.id.toString()} label={category.name} />
                                        {/each}
                                    </Select.Group>
                                </Select.Content>
                            </Select.Root>
                        </div>
                    </div>

                    <div class="space-y-2">
                        <Label for="editDescription">Description</Label>
                        <Textarea 
                            id="editDescription" 
                            class="min-h-[100px]" 
                            placeholder="Enter product description" 
                            bind:value={editDescription} 
                        />
                    </div>

                    <div class="flex items-center space-x-2">
                        <Switch 
                            id="editIsActive" 
                            checked={editIsActive} 
                            onchange={() => editIsActive = !editIsActive} 
                        />
                        <Label for="editIsActive">Product {editIsActive ? 'Active' : 'Inactive'}</Label>
                    </div>
                </div>

                <Card.Footer>
                    <Button variant="outline" onclick={() => goto(`/products/${productId}`)}>
                        <SquareArrowLeft size={18} />
                        <span>Cancel</span>
                    </Button>
                    <Button type="submit" disabled={isSaving} class="bg-primary">
                        {isSaving ? 'Saving...' : 'Save Changes'}
                    </Button>
                </Card.Footer>
            </Card.Content>
        </Card.Root>
    </form>
</div>
