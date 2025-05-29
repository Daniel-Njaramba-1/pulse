<script lang="ts">
    import { onMount } from "svelte";
    import * as Card from "$lib/components/ui/card/index";
    import { dashboardData, loading, error, initializeDashboard } from "$lib/stores/dashboard";
    import { ArrowUp, ArrowDown, TrendingUp, Package, DollarSign, Users, Activity } from "lucide-svelte";

    onMount(async () => {
        await initializeDashboard();
    });
</script>

<div class="container mx-auto py-6 px-4">
    <h1 class="text-2xl font-bold mb-6">Dashboard Overview</h1>

    {#if $loading}
        <div class="flex justify-center items-center h-48">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        </div>
    {:else if $error}
        <div class="border border-red-300 bg-red-50 p-4 rounded-lg text-red-800">
            Error: {$error}
        </div>
    {:else if $dashboardData}
        <!-- Operational Health Metrics -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
            {#each $dashboardData.operational_health as metric}
                <Card.Root class="p-4">
                    <Card.Header>
                        <Card.Title class="text-sm font-medium text-gray-500">{metric.metric}</Card.Title>
                    </Card.Header>
                    <Card.Content>
                        <div class="flex items-center justify-between">
                            <span class="text-2xl font-bold">{metric.value}</span>
                            <span class={`px-2 py-1 rounded-full text-xs ${
                                metric.status === 'good' ? 'bg-green-100 text-green-800' :
                                metric.status === 'warning' ? 'bg-yellow-100 text-yellow-800' :
                                'bg-red-100 text-red-800'
                            }`}>
                                {metric.status}
                            </span>
                        </div>
                    </Card.Content>
                </Card.Root>
            {/each}
        </div>

        <!-- Top Products -->
        <div class="mb-6">
            <h2 class="text-xl font-semibold mb-4">Top Performing Products</h2>
            <div class="bg-white shadow-md rounded-lg overflow-hidden">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Product</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Category</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Total Sales</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Revenue</th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                        {#each $dashboardData.top_products as product}
                            <tr>
                                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                    {product.product_name}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    {product.category_name}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    {product.total_sales}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    ${product.total_revenue.toFixed(2)}
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        </div>

        <!-- Category Revenue -->
        <div class="mb-6">
            <h2 class="text-xl font-semibold mb-4">Category Revenue</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {#each $dashboardData.category_revenue as category}
                    <Card.Root class="p-4">
                        <Card.Header>
                            <Card.Title class="text-sm font-medium text-gray-500">{category.category_name}</Card.Title>
                        </Card.Header>
                        <Card.Content>
                            <div class="flex flex-col">
                                <span class="text-2xl font-bold">${category.total_revenue.toFixed(2)}</span>
                                <span class="text-sm text-gray-500">{category.sales_count} sales</span>
                            </div>
                        </Card.Content>
                    </Card.Root>
                {/each}
            </div>
        </div>

        <!-- Inventory Status -->
        <div class="mb-6">
            <h2 class="text-xl font-semibold mb-4">Low Stock Alert</h2>
            <div class="bg-white shadow-md rounded-lg overflow-hidden">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Product</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Current Stock</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Threshold</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Last Restock</th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                        {#each $dashboardData.inventory_status.filter(item => item.is_low_stock) as item}
                            <tr>
                                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                    {item.product_name}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    {item.current_stock}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    {item.stock_threshold}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    {new Date(item.last_restock).toLocaleDateString()}
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        </div>

        <!-- Model Performance -->
        <div class="mb-6">
            <h2 class="text-xl font-semibold mb-4">Model Performance</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {#each $dashboardData.model_performance as model}
                    <Card.Root class="p-4">
                        <Card.Header>
                            <Card.Title class="text-sm font-medium text-gray-500">Version {model.model_version}</Card.Title>
                            <Card.Description>Training Date: {new Date(model.training_date).toLocaleDateString()}</Card.Description>
                        </Card.Header>
                        <Card.Content>
                            <div class="grid grid-cols-2 gap-2">
                                <div>
                                    <span class="text-xs text-gray-500">R² Score</span>
                                    <p class="text-sm font-medium">{model.r_squared.toFixed(3)}</p>
                                </div>
                                <div>
                                    <span class="text-xs text-gray-500">RMSE</span>
                                    <p class="text-sm font-medium">{model.rmse.toFixed(3)}</p>
                                </div>
                                <div>
                                    <span class="text-xs text-gray-500">MAE</span>
                                    <p class="text-sm font-medium">{model.mae.toFixed(3)}</p>
                                </div>
                                <div>
                                    <span class="text-xs text-gray-500">Sample Size</span>
                                    <p class="text-sm font-medium">{model.sample_size}</p>
                                </div>
                            </div>
                        </Card.Content>
                    </Card.Root>
                {/each}
            </div>
        </div>
    {/if}
</div>
