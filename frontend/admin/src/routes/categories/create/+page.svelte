<script lang="ts">
    import { Button } from "$lib/components/ui/button/index";
    import * as Card from "$lib/components/ui/card/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Label } from "$lib/components/ui/label/index";
    import { categories, isLoading, error, categoryHelpers } from "$lib/stores/category";
    import type { Category } from "$lib/stores/category";
    import { goto } from "$app/navigation";

    let newCategoryName = '';
    let newCategoryDescription = '';
    let newCategoryResult = '';
    let newCategoryLoading = false;

    // Handle new category submission
    async function handleSubmitNewCategory(event: SubmitEvent): Promise<void> {
        event.preventDefault();
        newCategoryLoading = true;
        newCategoryResult = "";

        try {
            const result = await categoryHelpers.createCategory(newCategoryName, newCategoryDescription);
            if (result.success) {
                newCategoryResult = "Category created successfully";
                newCategoryName = "";
                newCategoryDescription = "";
            } else {
                newCategoryResult = result.error || "Unknown error occurred";
            }
        } catch (error) {
            newCategoryResult = "Error occurred while creating new category";
            console.error(error);
        } finally {
            newCategoryLoading = false;
        }
    }
</script>

<div class="container mx-auto py-3 px-40"> 
    <form onsubmit={handleSubmitNewCategory} class="space-y-4">
        <Card.Root>
            <Card.Header> 
                <Card.Title class="text-xl font-bold">Add New Category</Card.Title>
            </Card.Header> 
            <Card.Content> 
                <div class="space-y-2">
                    <Label for="name" class="text-sm font-medium">Name</Label>
                    <Input
                        id="name"
                        type="text"
                        bind:value={newCategoryName}
                        placeholder="Enter category name"
                        required
                        class="w-full"
                    />
                </div>
                <div class="space-y-2">
                    <Label for="description" class="text-sm font-medium">Description</Label>
                    <Input
                        id="description"
                        type="text"
                        bind:value={newCategoryDescription}
                        placeholder="Enter category description"
                        required
                        class="w-full"
                    />
                </div>
            
                {#if newCategoryResult}
                    <div class={newCategoryResult.includes("successfully")
                        ? "text-green-600 text-sm"
                        : "text-red-600 text-sm"}>
                        {newCategoryResult}
                    </div>
                {/if}
            </Card.Content>
            <Card.Footer class="flex justify-end gap-2"> 
                <Button variant="outline" onclick = {() => goto (`/categories`)}> Cancel</Button>
                <Button variant="outline" type="submit" class=""  disabled={newCategoryLoading} >
                    {newCategoryLoading ? 'Creating...' : 'Create Category'}
                </Button>
            </Card.Footer>
        </Card.Root>
    </form>
</div>
