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
            <div class="text-gray-500">Loading products...</div>
        </div>
    {:else if $error}
        <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
            Error: {$error}
        </div>
    {:else}
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
                {#each $products as product (product.id)}
                <Table.Row class="hover:bg-gray-50 cursor-pointer" data-id={product.id} onclick={() => goto(`/products/${product.id}`, { state: { product } })}>
                    <Table.Cell>{product.id}</Table.Cell>
                    <Table.Cell>{product.name}</Table.Cell>
                    <Table.Cell>{product.description}</Table.Cell>
                    <Table.Cell>{product.is_active ? "Yes" : "No"}</Table.Cell>
                </Table.Row>
                {/each}
            </Table.Body>
        </Table.Root>
    {/if}
</div>