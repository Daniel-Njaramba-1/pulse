<script lang="ts">
    import { Button } from "$lib/components/ui/button/index";
    import * as Card from "$lib/components/ui/card/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Label } from "$lib/components/ui/label/index";
    import { brands, isLoading, error, brandHelpers } from "$lib/stores/brands";
    import type { Brand } from "$lib/stores/brands";
    import { goto } from "$app/navigation";
  
    // Local state
    let newBrandName = $state<string>("");
    let newBrandDescription = $state<string>("");
    let newBrandResult = $state("");
    let newBrandLoading = $state<boolean>(false);

    // Handle new brand submission
    async function handleSubmitNewBrand(event: SubmitEvent): Promise<void> {
        event.preventDefault();
        newBrandLoading = true;
        newBrandResult = "";

        try {
            const result = await brandHelpers.createBrand(newBrandName, newBrandDescription);
            if (result.success) {
                newBrandResult = "Brand created successfully";
                newBrandName = "";
                newBrandDescription = "";
                goto ('/brands')
            } else {
                newBrandResult = result.error || "Unknown error occurred";
            }
        } catch (error) {
            newBrandResult = "Error occurred while creating new brand";
            console.error(error);
        } finally {
            newBrandLoading = false;
        }
    }
</script>

<div class="container mx-auto py-3 px-40">
    <form onsubmit={handleSubmitNewBrand} class="space-y-4">
        <Card.Root> 
            <Card.Header> 
                <Card.Title class="text-xl font-bold">Add New Brand</Card.Title>
            </Card.Header>
            <Card.Content> 
                <div class="space-y-2">
                    <Label for="name" class="text-sm font-medium">Name</Label>
                    <Input 
                        id="name"
                        type="text" 
                        bind:value={newBrandName}
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
                        bind:value={newBrandDescription}
                        placeholder="Enter brand description"
                        required
                        class="w-full"
                    />
                </div>
                
                {#if newBrandResult}
                    <div class={newBrandResult.includes("successfully") 
                        ? "text-green-600 text-sm" 
                        : "text-red-600 text-sm"}>
                        {newBrandResult}
                    </div>
                {/if}
            </Card.Content>
            <Card.Footer class="flex justify-end gap-2"> 
                <Button variant="outline" onclick = {() => goto (`/brands`)}> Cancel</Button>
                <Button variant="outline" type="submit" class=""  disabled={newBrandLoading} >
                    {newBrandLoading ? 'Creating...' : 'Create Brand'}
                </Button>
            </Card.Footer>
        </Card.Root>
    </form>
</div>

