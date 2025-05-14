<script lang="ts">
    import { onMount } from "svelte";
    import { Button } from "$lib/components/ui/button/index";
    import * as Table from "$lib/components/ui/table/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Search } from "lucide-svelte";
    import { brands, isLoading, error, brandHelpers } from "$lib/stores/brands";
    import type { Brand } from "$lib/stores/brands";
    import { goto } from "$app/navigation";
  
    // Local state
    let searchQuery = $state<string>("");
    
    // Fetch brands on component mount
    onMount(async () => {
        const response = await brandHelpers.fetchBrands();
        if (!response.success) {
            // Handle error - maybe show a notification
            console.error("Failed to fetch brands:", response.error);
        }
    });
</script>
  
<div class="container mx-auto py-3 px-4">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-2">
        <h2 class="text-xl font-semibold">Brands</h2>
        <div class="flex items-center gap-3 w-full sm:w-auto">
            <div class="relative w-[400px]">
                <Search class="absolute left-3 top-3 h-4 w-4 text-gray-500" />
                <Input
                    id="search"
                    type="search"
                    placeholder="Search brands..."
                    class="pl-10"
                    bind:value={searchQuery}
                />
            </div>
        </div>
        
        <div class="">
            <Button variant="outline" onclick={() => goto(`/brands/create`)}>
                Create
            </Button>
        </div>
    </div>
    
    {#if $isLoading && $brands.length === 0}
        <div class="flex justify-center items-center h-48">
            <div class="text-gray-500">Loading brands...</div>
        </div>
    {:else if $error}
        <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
            Error: {$error}
        </div>
    {:else if !$brands || $brands.length === 0}
        <div>
            No Brands Found
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
                    {#each $brands as brand (brand.id)}
                    <Table.Row class="hover:bg-gray-50 cursor-pointer" data-id={brand.id} onclick={() => goto(`/brands/${brand.id}`, { state: { brand } })}>
                        <Table.Cell>{brand.id}</Table.Cell>
                        <Table.Cell>{brand.name}</Table.Cell>
                        <Table.Cell>{brand.description}</Table.Cell>
                        <Table.Cell>{brand.is_active ? "Yes" : "No"}</Table.Cell>
                    </Table.Row>
                    {/each}
                </Table.Body>
            </Table.Root>
        </div>
    {/if}
</div>