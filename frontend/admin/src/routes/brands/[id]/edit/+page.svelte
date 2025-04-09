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
    import { brandHelpers } from "$lib/stores/brands";
    import type { Brand } from "$lib/stores/brands";
    import { goto } from "$app/navigation";
    
    // Get brand ID from route params
    const brandId = $derived(page.params.id);
    
    // Local state
    let brand: Brand | null = $state(null);
    let isLoading = $state<boolean>(true);
    let error: string | null = $state(null);
    let brandName = $state<string>("");
    let brandDescription = $state<string>("");
    let isUpdating = $state<boolean>(false);
    let updateResult = $state<string>("");
    let isDialogOpen = $state<boolean>(false);
    
    // Handle form submission
    async function handleSubmitEditBrand(event: SubmitEvent) {
        event.preventDefault();
        isUpdating = true;
        updateResult = "";
        
        try {
            const response = await brandHelpers.updateBrand(Number(brandId), brandName, brandDescription);
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
    <form onsubmit={handleSubmitEditBrand} class="space-y-4">
        <Card.Root> 
            <Card.Header> 
                <Card.Title class="text-xl font-bold">Edit {brandName}</Card.Title>
            </Card.Header>
            <Card.Content> 
                    <div class="space-y-2">
                        <Label for="name" class="text-sm font-medium">Name</Label>
                        <Input 
                            id="name"
                            type="text" 
                            bind:value={brandName}
                            placeholder="Enter brand name"
                            required
                            class="w-full"
                        />
                    </div>
                    <div class="space-y-2">
                        <Label for="description" class="text-sm font-medium">Description</Label>
                        <Input 
                            id="description"
                            type="text" 
                            bind:value={brandDescription}
                            placeholder="Enter brand description"
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