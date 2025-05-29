<script lang="ts">
    import { onMount, tick } from "svelte";
    import { page } from "$app/state";
    import { Button } from "$lib/components/ui/button/index";
    import * as Card from "$lib/components/ui/card/index";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import * as Tooltip from "$lib/components/ui/tooltip/index"
    import { Input } from "$lib/components/ui/input/index";
    import { Label } from "$lib/components/ui/label/index";
    import { Pencil, SquareArrowLeft } from "lucide-svelte";
    import { categoryHelpers } from "$lib/stores/category";
    import type { Category } from "$lib/stores/category";
    import { goto } from "$app/navigation";
    
    // Get category ID from route params
    const categoryId = $derived(page.params.id);
    
    // Local state
    let category: Category | null = $state(null);
    let isLoading = $state<boolean>(true);
    let error: string | null = $state(null);
    let categoryName = $state<string>("");
    let categoryDescription = $state<string>("");
    let isUpdating = $state<boolean>(false);
    let updateResult = $state<string>("");
    let isDialogOpen = $state<boolean>(false);
    
    // Handle form submission
    async function handleSubmitEditCategory(event: SubmitEvent) {
        event.preventDefault();
        isUpdating = true;
        updateResult = "";
        
        try {
            const response = await categoryHelpers.updateCategory(Number(categoryId), categoryName, categoryDescription);
            if (response.success) {
                updateResult = "Brand updated successfully";
            } else {
                updateResult = response.error || "Failed to update brand";
            }
        } catch (err) {
            updateResult = err instanceof Error ? err.message : "An unexpected error occurred";
        } finally {
            isUpdating = false;
        }
    }
</script>

<div class="container mx-auto py-3 px-40">
    <form onsubmit={handleSubmitEditCategory} class="space-y-4">
        <Card.Root> 
            <Card.Header> 
                <Card.Title class="text-xl font-bold">Edit {categoryName}</Card.Title>
            </Card.Header>
            <Card.Content> 
                    <div class="space-y-2">
                        <Label for="name" class="text-sm font-medium">Name</Label>
                        <Input 
                            id="name"
                            type="text" 
                            bind:value={categoryName}
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
                            bind:value={categoryDescription}
                            placeholder="Enter category description"
                            required
                            class="w-full"
                        />
                    </div>
                    
                    {#if updateResult}
                        <div class={updateResult.includes("successfully") 
                            ? "text-green-600 text-sm" 
                            : "text-red-600 text-sm"}>
                            {updateResult}
                        </div>
                    {/if}
            </Card.Content>
            <Card.Footer class="flex justify-end gap-2">
                <Button variant="outline" onclick = {() => goto (`/brands`)}> Cancel</Button>
                <Button variant="outline" type="submit" class="" disabled={isUpdating} >
                    {isUpdating ? 'Updating...' : 'Save Changes'}
                </Button>
            </Card.Footer>
        </Card.Root>
    </form>  
</div>