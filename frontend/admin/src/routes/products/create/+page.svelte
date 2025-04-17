<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import {Plus, Funnel, Upload, SquareArrowLeft} from "lucide-svelte";
    import * as Card from "$lib/components/ui/card/index";
    import * as Select from "$lib/components/ui/select/index";
    import { Button } from "$lib/components/ui/button/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Label } from "$lib/components/ui/label/index";
    import { Textarea } from "$lib/components/ui/textarea/index"
    import { brands, brandHelpers } from "$lib/stores/brands";
    import type { Brand } from "$lib/stores/brands";
    import { categories, categoryHelpers } from "$lib/stores/category";
    import type { Category } from "$lib/stores/category";
    import { products, isLoading, error, productHelpers } from "$lib/stores/product";
    import type { Product } from "$lib/stores/product";


    // Local state
    let newProductName = $state<string>("");
    let newProductDescription = $state<string>("");
    let newProductResult = $state("");
    let newProductLoading = $state<boolean>(false);
    let selectedBrandId = $state<number | null>(null);
    let selectedCategoryId = $state<number | null>(null);
    let imageFile: File | null = $state<File | null>(null);
    let imagePreview: string | null = $state<string | null>(null);

    // Fetch initial data
    onMount(async () => {
        //fetch products
        const productsResponse = await productHelpers.fetchProducts();
        if (!productsResponse.success) {
            console.error("Failed to fetch products:", productsResponse.error);
        }

        // Fetch brands for dropdown
        const brandResponse = await brandHelpers.fetchBrands();
        if (!brandResponse.success) {
            console.error("Failed to fetch brands:", brandResponse.error);
        }
        
        // Fetch categories for dropdown
        const categoryResponse = await categoryHelpers.fetchCategories();
        if (!categoryResponse.success) {
            console.error("Failed to fetch categories:", categoryResponse.error);
        }
    });

    
    // Brand and category trigger content
    let brandTriggerContent = $derived($brands.find((f) => f.id === selectedBrandId)?.name ?? "Select Brand");
    let categoryTriggerContent = $derived($categories.find((f) => f.id === selectedCategoryId)?.name ?? "Select category");


    // Handle image selection
    function handleImageChange(event: Event) {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files.length > 0) {
            imageFile = input.files[0];
            const reader = new FileReader();
            reader.onload = (e) => {
                imagePreview = e.target?.result as string;
            };
            reader.readAsDataURL(imageFile);
        } else {
            imageFile = null;
            imagePreview = null;
        }
    }

    // Handle new product submission
    async function handleSubmitNewProduct(event: SubmitEvent): Promise<void> {
        event.preventDefault();
        newProductLoading = true;
        newProductResult = "";

        if (!selectedBrandId || !selectedCategoryId) {
            newProductResult = "Please select a brand and category.";
            newProductLoading = false;
            return;
        }

        try {
            const result = await productHelpers.createProduct(
                selectedBrandId,
                selectedCategoryId,
                newProductName,
                newProductDescription,
                imageFile
            );
            if (result.success) {
                newProductResult = "Product created successfully";
                newProductName = "";
                newProductDescription = "";
                imageFile = null;
                imagePreview = null;
                goto("/products")
            } else {
                newProductResult = result.error || "Unknown error occurred";
            }
        } catch (error) {
            newProductResult = "Error occurred while creating new product";
            console.error(error);
        } finally {
            newProductLoading = false;
        }
    }
</script>

<div class="container mx-auto py-2 px-40">    
    <form onsubmit={handleSubmitNewProduct} class="space-y-8">
        <Card.Root>
            <Card.Header> 
                <Card.Title class="text-xl font-bold">Add New Product</Card.Title>
            </Card.Header>
            <Card.Content class="space-y-8">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div class="space-y-2">
                        <Label for="brand">Brand</Label>
                        <Select.Root type="single" name="brandId" onValueChange={(value) => selectedBrandId = Number(value)}>
                            <Select.Trigger class="">{brandTriggerContent}</Select.Trigger>
                            <Select.Content>
                                <Select.Group>
                                    <Select.GroupHeading>Brands</Select.GroupHeading>
                                    {#each $brands as brand (brand.id)}
                                        <Select.Item value={Number(brand.id).toString()} label={brand.name} />
                                    {/each}
                                </Select.Group>
                            </Select.Content>
                        </Select.Root>
                    </div>
                    <div class="space-y-2">
                        <Label for="category">Category</Label>
                        <Select.Root type="single" name="categoryId" onValueChange={(value) => selectedCategoryId = Number(value)}>
                            <Select.Trigger class="">{categoryTriggerContent}</Select.Trigger>
                            <Select.Content>
                                <Select.Group>
                                    <Select.GroupHeading>Categories</Select.GroupHeading>
                                    {#each $categories as category (category.id)}
                                        <Select.Item value={Number(category.id).toString()} label={category.name} />
                                    {/each}
                                </Select.Group>
                            </Select.Content>
                        </Select.Root>
                    </div>
                </div>

                <div class="space-y-2">
                    <Label for="name">Product Name</Label>
                    <Input id="name" class="" placeholder="" bind:value={newProductName}/>
                </div>

                <div class="space-y-2">
                    <Label for="description">Product Description</Label>
                    <Input id="description" class="" placeholder="" bind:value={newProductDescription}/>
                </div>

                <div class="space-y-4">
                    <Label for="image">Product Image</Label>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 items-center">
                        <div class="flex flex-col gap-2">
                            <div class="border rounded-md p-2">
                                <Label for="image" class="flex flex-col items-center gap-2 cursor-pointer py-4">
                                    <Upload class="h-8 w-8 text-muted-foreground" />
                                    <span class="text-sm text-muted-foreground">Click to upload or Drag and Drop</span>
                                    <span class="text-sm text-muted-foreground">SVG, PNG, JPG, GIF (2MB max)</span>
                                    <Input id="image" type="file" accept="image/*" class="hidden" onchange={handleImageChange}/>
                                </Label>
                            </div>
                        </div>
                        <div class="flex justify-center">
                            {#if imagePreview}
                                <div class="relative w-40 h-40 border rounded-md overflow-hidden">
                                    <img
                                        src={imagePreview}
                                        alt="Product Preview"
                                        class="object-cover w-full h-full"
                                    />
                                </div>
                            {:else}
                                <div class="w-40 h-40 border rounded-md flex items-center justify-center bg-muted">
                                    <span class="text-sm text-muted-foreground">No image</span>
                                </div>
                            {/if}
                        </div>
                    </div>
                </div>
            </Card.Content>
            <Card.Footer class="flex justify-end gap-2">
                <Button variant="outline" onclick = {() => goto (`/products`)}> Cancel</Button>
                <Button variant="outline" type="submit" class="" disabled={newProductLoading} >
                    {newProductLoading ? 'Creating...' : 'Create Product'}
                </Button>
            </Card.Footer>
        </Card.Root>
    </form>
</div>
