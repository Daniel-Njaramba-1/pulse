import { writable } from 'svelte/store';
import type { Writable } from 'svelte/store';

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

// Helper functions
async function fetchWithAuth(endpoint: string, params?: Record<string, string>) {
    const token = localStorage.getItem('admin_token');
    if (!token) throw new Error('No authentication token found');

    const queryString = params ? '?' + new URLSearchParams(params).toString() : '';
    const response = await fetch(`${API_URL}${endpoint}${queryString}`, {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        }
    });

    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    if (!data.success) {
        throw new Error(data.error || 'Failed to fetch data');
    }

    return data.data;
}

// Dashboard data fetching functions
export async function fetchDashboardAnalytics(days: number = 30) {
    try {
        loading.set(true);
        error.set(null);
        const data = await fetchWithAuth('/api/admin/dashboard/analytics', { days: days.toString() });
        dashboardData.set(data);
    } catch (err) {
        error.set(err instanceof Error ? err.message : 'An error occurred');
        console.error('Error fetching dashboard analytics:', err);
    } finally {
        loading.set(false);
    }
}

export async function fetchModelPerformance() {
    try {
        loading.set(true);
        error.set(null);
        const data = await fetchWithAuth('/api/admin/dashboard/model-performance');
        dashboardData.update(current => current ? { ...current, model_performance: data } : null);
    } catch (err) {
        error.set(err instanceof Error ? err.message : 'An error occurred');
        console.error('Error fetching model performance:', err);
    } finally {
        loading.set(false);
    }
}

export async function fetchSalesAnalytics(days: number = 30) {
    try {
        loading.set(true);
        error.set(null);
        const data = await fetchWithAuth('/api/admin/dashboard/sales', { days: days.toString() });
        dashboardData.update(current => current ? { ...current, sales_analytics: data } : null);
    } catch (err) {
        error.set(err instanceof Error ? err.message : 'An error occurred');
        console.error('Error fetching sales analytics:', err);
    } finally {
        loading.set(false);
    }
}

export async function fetchInventoryStatus() {
    try {
        loading.set(true);
        error.set(null);
        const data = await fetchWithAuth('/api/admin/dashboard/inventory');
        dashboardData.update(current => current ? { ...current, inventory_status: data } : null);
    } catch (err) {
        error.set(err instanceof Error ? err.message : 'An error occurred');
        console.error('Error fetching inventory status:', err);
    } finally {
        loading.set(false);
    }
}

export async function fetchPricingAnalytics(days: number = 7) {
    try {
        loading.set(true);
        error.set(null);
        const data = await fetchWithAuth('/api/admin/dashboard/pricing', { days: days.toString() });
        dashboardData.update(current => current ? { ...current, pricing_analytics: data } : null);
    } catch (err) {
        error.set(err instanceof Error ? err.message : 'An error occurred');
        console.error('Error fetching pricing analytics:', err);
    } finally {
        loading.set(false);
    }
}

export async function fetchCustomerBehavior() {
    try {
        loading.set(true);
        error.set(null);
        const data = await fetchWithAuth('/api/admin/dashboard/customers');
        dashboardData.update(current => current ? { ...current, customer_behavior: data } : null);
    } catch (err) {
        error.set(err instanceof Error ? err.message : 'An error occurred');
        console.error('Error fetching customer behavior:', err);
    } finally {
        loading.set(false);
    }
}

export async function fetchOperationalHealth() {
    try {
        loading.set(true);
        error.set(null);
        const data = await fetchWithAuth('/api/admin/dashboard/health');
        dashboardData.update(current => current ? { ...current, operational_health: data } : null);
    } catch (err) {
        error.set(err instanceof Error ? err.message : 'An error occurred');
        console.error('Error fetching operational health:', err);
    } finally {
        loading.set(false);
    }
}

export async function fetchTopProducts(limit: number = 10) {
    try {
        loading.set(true);
        error.set(null);
        const data = await fetchWithAuth('/api/admin/dashboard/top-products', { limit: limit.toString() });
        dashboardData.update(current => current ? { ...current, top_products: data } : null);
    } catch (err) {
        error.set(err instanceof Error ? err.message : 'An error occurred');
        console.error('Error fetching top products:', err);
    } finally {
        loading.set(false);
    }
}

export async function fetchCategoryRevenue(days: number = 30) {
    try {
        loading.set(true);
        error.set(null);
        const data = await fetchWithAuth('/api/admin/dashboard/category-revenue', { days: days.toString() });
        dashboardData.update(current => current ? { ...current, category_revenue: data } : null);
    } catch (err) {
        error.set(err instanceof Error ? err.message : 'An error occurred');
        console.error('Error fetching category revenue:', err);
    } finally {
        loading.set(false);
    }
}

// Initialize dashboard data
export async function initializeDashboard() {
    await fetchDashboardAnalytics();
}
