<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import * as Card from "$lib/components/ui/card/index";
    import { Button } from "$lib/components/ui/button/index";
    import * as Table from "$lib/components/ui/table/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Search } from "lucide-svelte";
    import { products, isLoading, error, productHelpers } from "$lib/stores/product";
    import type { Product } from "$lib/stores/product";
    
  
    // Local state
    let searchQuery = $state<string>("");
    
    // Fetch initial data
    onMount(async () => {
        //fetch products
        const productsResponse = await productHelpers.fetchProducts();
        if (!productsResponse.success) {
            console.error("Failed to fetch products:", productsResponse.error);
        }
    });
</script>
  
<div class="container mx-auto py-3 px-4">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-2">
        <h2 class="text-xl font-semibold">Products</h2>
        <div class="flex items-center gap-3 w-full sm:w-auto">
            <div class="relative w-[400px]">
                <Search class="absolute left-3 top-3 h-4 w-4 text-gray-500" />
                <Input
                    id="search"
                    type="search"
                    placeholder="Search products..."
                    class="pl-10"
                    bind:value={searchQuery}
                />
            </div>
        </div>

        <div class="">
            <Button variant="outline" onclick={() => goto(`/products/create`)}>
                Create
            </Button>
        </div>
    </div>

    {#if $isLoading && $products.length === 0}
        <div class="flex justify-center items-center h-48">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        </div>
    {:else if $error}
        <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
            Error: {$error}
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
                    {#each $products as product (product.id)}
                    <tr class="hover:bg-gray-50 cursor-pointer" data-id={product.id} onclick={() => goto(`/products/${product.id}`, { state: { product } })}>
                        <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                            {product.id}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                            {product.name}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                            {product.description}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap">
                            <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full {product.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
                                {product.is_active ? 'Active' : 'Inactive'}
                            </span>
                        </td>
                    </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>