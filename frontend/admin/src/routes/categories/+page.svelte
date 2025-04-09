<script lang="ts">
    import { onMount } from "svelte";
    import { Button } from "$lib/components/ui/button/index";
    import * as Table from "$lib/components/ui/table/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Search, } from "lucide-svelte";
    import { categories, isLoading, error, categoryHelpers } from "$lib/stores/category";
    import type { Category } from "$lib/stores/category";
    import { goto } from "$app/navigation";

    // Local state
    let searchQuery = $state<string>("");

    // Fetch categories on component mount
    onMount(async () => {
        const response = await categoryHelpers.fetchCategories();
        if (!response.success) {
            console.error("Failed to fetch categories:", response.error);
        }
    });
</script>

<div class="container mx-auto py-3 px-4">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-2">
        <h2 class="text-xl font-semibold">Categories</h2>
        <div class="flex items-center gap-3 w-full sm:w-auto">
            <div class="relative w-[400px]">
                <Search class="absolute left-3 top-3 h-4 w-4 text-gray-500" />
                <Input
                    id="search"
                    type="search"
                    placeholder="Search categories..."
                    class="pl-10"
                    bind:value={searchQuery}
                />
            </div>
        </div>

        <div>
            <Button variant="outline" onclick={() => goto(`/categories/create`)}>
                Create
            </Button>
        </div>
    </div>

    {#if $isLoading && $categories.length === 0}
        <div class="flex justify-center items-center h-48">
            <div class="text-gray-500">Loading categories...</div>
        </div>
    {:else if $error}
        <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
            Error: {$error}
        </div>
    {:else if !$categories || $categories.length === 0}
        <div>
            No Categories Found
        </div>
    {:else}
        <div class="border border-gray-200 rounded-lg overflow-hidden">
            <Table.Root>
                <Table.Header>
                    <Table.Row class="bg-gray-100">
                        <Table.Head>ID</Table.Head>
                        <Table.Head>Name</Table.Head>
                        <Table.Head>Description</Table.Head>
                        <Table.Head>Active</Table.Head>
                    </Table.Row>
                </Table.Header>
                <Table.Body>
                    {#each $categories as category (category.id)}
                    <Table.Row class="hover:bg-gray-50 cursor-pointer" data-id={category.id} onclick={() => goto(`/categories/${category.id}`, { state: { category } })}>
                        <Table.Cell>{category.id}</Table.Cell>
                        <Table.Cell>{category.name}</Table.Cell>
                        <Table.Cell>{category.description}</Table.Cell>
                        <Table.Cell>{category.is_active ? "Yes" : "No"}</Table.Cell>
                    </Table.Row>
                    {/each}
                </Table.Body>
            </Table.Root>
        </div>
    {/if}
</div>