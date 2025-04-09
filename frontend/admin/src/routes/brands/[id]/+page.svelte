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
    
    // Fetch brand details on mount
    async function fetchBrandDetails() {
        try {
            const response = await brandHelpers.getBrand(Number(brandId));
            if (response.success && response.data) {
                brand = {...response.data as Brand};
                brandName = brand.name;
                brandDescription = brand.description;
            } else {
                error = response.error || "Failed to load brand details";
            }
        } catch (err) {
            error = err instanceof Error ? err.message : "An unexpected error occurred";
        } finally {
            isLoading = false;
        }
    }
    
    onMount(fetchBrandDetails);
    
    // Handle form submission
    async function handleSubmitEditBrand(event: SubmitEvent) {
        event.preventDefault();
        isUpdating = true;
        updateResult = "";
        
        try {
            const response = await brandHelpers.updateBrand(Number(brandId), brandName, brandDescription);
            if (response.success) {
                updateResult = "Brand updated successfully";
                
                // Refetch the brand details to ensure we have the most up-to-date data
                await fetchBrandDetails();
                
                // Close dialog after successful update
                isDialogOpen = false;
                await tick();
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

<div class="container mx-auto py-3 px-4">
    {#if isLoading}
        <div class="flex justify-center items-center h-48">
            <div class="text-gray-500 animate-pulse">Loading brand details...</div>
        </div>
    {:else if error}
        <Card.Root>
            <Card.Header>
                <Card.Title class="text-red-600">Error Loading Brand</Card.Title>
            </Card.Header>
            <Card.Content>
                <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
                    {error}
                </div>
                <div class="mt-4 flex justify-center">
                    <Button onclick={() => goto('/brands')} variant="outline">
                        Return to Brands
                    </Button>
                </div>
            </Card.Content>
        </Card.Root>
    {:else if brand}
        <div class="flex flex-row">
            <SquareArrowLeft class="m-1 cursor-pointer" onclick={() => goto('/brands')}/>
            <h1 class="text-xl font-bold ml-2">{brand.name}</h1>
        </div>
        <Card.Root>
            <Card.Header class="flex flex-row items-center justify-between space-y-0">
                <div>
                    <Card.Title class="text-lg font-semibold text-gray-700">Details</Card.Title>
                </div>
                
                <Dialog.Root bind:open={isDialogOpen}> 
                    <Tooltip.Provider> 
                        <Tooltip.Root> 
                            <Tooltip.Trigger> 
                                <Dialog.Trigger>
                                    <Button variant="outline" size="icon">
                                        <Pencil class="h-2 w-2"/>
                                    </Button>
                                </Dialog.Trigger>
                            </Tooltip.Trigger>
                            <Tooltip.Content> 
                                Edit Brand
                            </Tooltip.Content>
                        </Tooltip.Root>
                    </Tooltip.Provider>
                    
                    <Dialog.Content interactOutsideBehavior="ignore"> 
                        <Dialog.Header> 
                            <Dialog.Title>Edit Brand</Dialog.Title>
                            <Dialog.Description>
                                Update the details of this brand in your inventory system.
                            </Dialog.Description>
                        </Dialog.Header>
                        
                        
                    </Dialog.Content>
                </Dialog.Root>
            </Card.Header>
            
            <Card.Content class="pt-1">
                <div class="">
                    <div>
                        <div class="space-y-2">
                            <div>
                                <span class="text-sm text-gray-500">Brand Name</span>
                                <p class="font-medium">{brand.name}</p>
                            </div>
                            <div>
                                <span class="text-sm text-gray-500">Description</span>
                                <p class="font-medium">{brand.description || 'No description provided'}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </Card.Content>
        </Card.Root>

        <!-- Table of Products -->
    {:else}
        <Card.Root>
            <Card.Header>
                <Card.Title>Brand Not Found</Card.Title>
            </Card.Header>
            <Card.Content>
                <div class="text-center text-gray-500">
                    The brand may have been deleted or the ID is invalid.
                    <div class="mt-4 flex justify-center">
                        <Button onclick={() => goto('/brands')} variant="outline">
                            Return to Brands
                        </Button>
                    </div>
                </div>
            </Card.Content>
        </Card.Root>
    {/if}
</div>