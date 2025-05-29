import { writable } from 'svelte/store';
import type { Writable } from 'svelte/store';
import { authHelpers } from './auth';

// Types
export interface ModelPerformanceData {
    model_version: string;
    training_date: string;
    r_squared: number;
    mse: number;
    rmse: number;
    mae: number;
    sample_size: number;
}

export interface SalesAnalytics {
    date: string;
    product_id: number;
    product_name: string;
    total_sales: number;
    total_revenue: number;
    category_name: string;
    brand_name: string;
}

export interface InventoryStatus {
    product_id: number;
    product_name: string;
    current_stock: number;
    stock_threshold: number;
    is_low_stock: boolean;
    last_restock: string;
    days_since_restock: number;
}

export interface PricingAnalytics {
    product_id: number;
    product_name: string;
    base_price: number;
    adjusted_price: number;
    price_change: number;
    last_adjusted: string;
    model_version: string;
}

export interface CustomerBehavior {
    product_id: number;
    product_name: string;
    average_rating: number;
    review_count: number;
    wishlist_count: number;
    wishlist_to_sales_ratio: number;
    sales_velocity: number;
}

export interface OperationalHealth {
    metric: string;
    value: number;
    status: string;
}

export interface CategoryRevenue {
    category_name: string;
    total_revenue: number;
    sales_count: number;
}

export interface DashboardData {
    model_performance: ModelPerformanceData[];
    sales_analytics: SalesAnalytics[];
    inventory_status: InventoryStatus[];
    pricing_analytics: PricingAnalytics[];
    customer_behavior: CustomerBehavior[];
    operational_health: OperationalHealth[];
    top_products: SalesAnalytics[];
    category_revenue: CategoryRevenue[];
}

// Base API URL
const API_URL = "http://localhost:8080/api/admin";

// Store
export const dashboardData: Writable<DashboardData | null> = writable(null);
export const loading: Writable<boolean> = writable(false);
export const error: Writable<string | null> = writable(null);

// Dashboard data fetching functions
export async function fetchDashboardAnalytics(days: number = 30) {
    loading.set(true);
    error.set(null);
    
    try {
        const response = await fetch(`${API_URL}/dashboard/analytics?days=${days}`, {
            method: 'GET',
            headers: {
                "Content-Type": "application/json",
                ...authHelpers.getAuthHeader()
            },
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to fetch dashboard analytics');
        }

        const data = await response.json();
        dashboardData.set(data.data);
        return data.data;
    } catch (err) {
        console.error('Error fetching dashboard analytics:', err);
        error.set(err instanceof Error ? err.message : 'Failed to fetch dashboard analytics');
        return null;
    } finally {
        loading.set(false);
    }
}

export async function fetchModelPerformance() {
    loading.set(true);
    error.set(null);
    
    try {
        const response = await fetch(`${API_URL}/dashboard/model-performance`, {
            method: 'GET',
            headers: {
                "Content-Type": "application/json",
                ...authHelpers.getAuthHeader()
            },
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to fetch model performance');
        }

        const data = await response.json();
        return data.data;
    } catch (err) {
        console.error('Error fetching model performance:', err);
        error.set(err instanceof Error ? err.message : 'Failed to fetch model performance');
        return null;
    } finally {
        loading.set(false);
    }
}

export async function fetchSalesAnalytics(days: number = 30) {
    loading.set(true);
    error.set(null);
    
    try {
        const response = await fetch(`${API_URL}/dashboard/sales?days=${days}`, {
            method: 'GET',
            headers: {
                "Content-Type": "application/json",
                ...authHelpers.getAuthHeader()
            },
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to fetch sales analytics');
        }

        const data = await response.json();
        return data.data;
    } catch (err) {
        console.error('Error fetching sales analytics:', err);
        error.set(err instanceof Error ? err.message : 'Failed to fetch sales analytics');
        return null;
    } finally {
        loading.set(false);
    }
}

export async function fetchInventoryStatus() {
    loading.set(true);
    error.set(null);
    
    try {
        const response = await fetch(`${API_URL}/dashboard/inventory`, {
            method: 'GET',
            headers: {
                "Content-Type": "application/json",
                ...authHelpers.getAuthHeader()
            },
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to fetch inventory status');
        }

        const data = await response.json();
        return data.data;
    } catch (err) {
        console.error('Error fetching inventory status:', err);
        error.set(err instanceof Error ? err.message : 'Failed to fetch inventory status');
        return null;
    } finally {
        loading.set(false);
    }
}

export async function fetchPricingAnalytics(days: number = 7) {
    loading.set(true);
    error.set(null);
    
    try {
        const response = await fetch(`${API_URL}/dashboard/pricing?days=${days}`, {
            method: 'GET',
            headers: {
                "Content-Type": "application/json",
                ...authHelpers.getAuthHeader()
            },
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to fetch pricing analytics');
        }

        const data = await response.json();
        return data.data;
    } catch (err) {
        console.error('Error fetching pricing analytics:', err);
        error.set(err instanceof Error ? err.message : 'Failed to fetch pricing analytics');
        return null;
    } finally {
        loading.set(false);
    }
}

export async function fetchCustomerBehavior() {
    loading.set(true);
    error.set(null);
    
    try {
        const response = await fetch(`${API_URL}/dashboard/customers`, {
            method: 'GET',
            headers: {
                "Content-Type": "application/json",
                ...authHelpers.getAuthHeader()
            },
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to fetch customer behavior');
        }

        const data = await response.json();
        return data.data;
    } catch (err) {
        console.error('Error fetching customer behavior:', err);
        error.set(err instanceof Error ? err.message : 'Failed to fetch customer behavior');
        return null;
    } finally {
        loading.set(false);
    }
}

export async function fetchOperationalHealth() {
    loading.set(true);
    error.set(null);
    
    try {
        const response = await fetch(`${API_URL}/dashboard/health`, {
            method: 'GET',
            headers: {
                "Content-Type": "application/json",
                ...authHelpers.getAuthHeader()
            },
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to fetch operational health');
        }

        const data = await response.json();
        return data.data;
    } catch (err) {
        console.error('Error fetching operational health:', err);
        error.set(err instanceof Error ? err.message : 'Failed to fetch operational health');
        return null;
    } finally {
        loading.set(false);
    }
}

export async function fetchTopProducts(limit: number = 10) {
    loading.set(true);
    error.set(null);
    
    try {
        const response = await fetch(`${API_URL}/dashboard/top-products?limit=${limit}`, {
            method: 'GET',
            headers: {
                "Content-Type": "application/json",
                ...authHelpers.getAuthHeader()
            },
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to fetch top products');
        }

        const data = await response.json();
        return data.data;
    } catch (err) {
        console.error('Error fetching top products:', err);
        error.set(err instanceof Error ? err.message : 'Failed to fetch top products');
        return null;
    } finally {
        loading.set(false);
    }
}

export async function fetchCategoryRevenue(days: number = 30) {
    loading.set(true);
    error.set(null);
    
    try {
        const response = await fetch(`${API_URL}/dashboard/category-revenue?days=${days}`, {
            method: 'GET',
            headers: {
                "Content-Type": "application/json",
                ...authHelpers.getAuthHeader()
            },
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to fetch category revenue');
        }

        const data = await response.json();
        return data.data;
    } catch (err) {
        console.error('Error fetching category revenue:', err);
        error.set(err instanceof Error ? err.message : 'Failed to fetch category revenue');
        return null;
    } finally {
        loading.set(false);
    }
}

// Initialize dashboard data
export async function initializeDashboard() {
    await fetchDashboardAnalytics();
}
