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
    
    // Fetch category details on mount
    async function fetchCategoryDetails() {
        try {
            const response = await categoryHelpers.getCategory(Number(categoryId));
            if (response.success && response.data) {
                category = {...response.data as Category};
                categoryName = category.name;
                categoryDescription = category.description;
            } else {
                error = response.error || "Failed to load category details";
            }
        } catch (err) {
            error = err instanceof Error ? err.message : "An unexpected error occurred";
        } finally {
            isLoading = false;
        }
    }
    
    onMount(fetchCategoryDetails);
</script>

<div class="container mx-auto py-3 px-4">
    {#if isLoading}
        <div class="flex justify-center items-center h-48">
            <div class="text-gray-500 animate-pulse">Loading category details...</div>
        </div>
    {:else if error}
        <Card.Root>
            <Card.Header>
                <Card.Title class="text-red-600">Error Loading Category</Card.Title>
            </Card.Header>
            <Card.Content>
                <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
                    {error}
                </div>
                <div class="mt-4 flex justify-center">
                    <Button onclick={() => goto('/categories')} variant="outline">
                        Return to Categories
                    </Button>
                </div>
            </Card.Content>
        </Card.Root>
    {:else if category}
        <div class="flex flex-row">
            <SquareArrowLeft class="m-1 cursor-pointer" onclick={() => goto('/categories')}/>
            <h1 class="text-xl font-bold ml-2">{category.name}</h1>
        </div>
        <Card.Root>
            <Card.Header class="flex flex-row items-center justify-between space-y-0">
                <div>
                    <Card.Title class="text-lg font-semibold text-gray-700">Details</Card.Title>
                </div>
                <Button variant="outline" size="icon" onclick={() => goto(`/categories/${category?.id}/edit`)}>
                    <Pencil class="h-2 w-2"/>
                </Button>    
            </Card.Header>
            
            <Card.Content class="pt-1">
                <div class="">
                    <div>
                        <div class="space-y-2">
                            <div>
                                <span class="text-sm text-gray-500">Category Name</span>
                                <p class="font-medium">{category.name}</p>
                            </div>
                            <div>
                                <span class="text-sm text-gray-500">Description</span>
                                <p class="font-medium">{category.description || 'No description provided'}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </Card.Content>
        </Card.Root>
        <!-- Table of related items could go here -->
    {:else}
        <Card.Root>
            <Card.Header>
                <Card.Title>Category Not Found</Card.Title>
            </Card.Header>
            <Card.Content>
                <div class="text-center text-gray-500">
                    The category may have been deleted or the ID is invalid.
                    <div class="mt-4 flex justify-center">
                        <Button onclick={() => goto('/categories')} variant="outline">
                            Return to Categories
                        </Button>
                    </div>
                </div>
            </Card.Content>
        </Card.Root>
    {/if}
</div>