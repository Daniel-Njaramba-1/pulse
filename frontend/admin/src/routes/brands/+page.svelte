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
  
<div class="container mx-auto py-8 px-4">
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
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
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
        <div class="bg-white shadow-md rounded-lg overflow-hidden">
            <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                    <tr>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Id</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                    </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                    {#each $brands as brand (brand.id)}
                    <tr class="hover:bg-gray-50 cursor-pointer" data-id={brand.id} onclick={() => goto(`/brands/${brand.id}`, { state: { brand } })}>
                        <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                            {brand.id}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                            {brand.name}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                            {brand.description}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap">
                            <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full {brand.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
                                {brand.is_active ? 'Active' : 'Inactive'}
                            </span>
                        </td>
                    </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>