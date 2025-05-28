import { useState, useCallback } from 'react';

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

// Custom hook for dashboard data
export function useDashboard() {
    const [dashboardData, setDashboardData] = useState<DashboardData | null>(null);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    // Helper functions
    const fetchWithAuth = useCallback(async (endpoint: string, params?: Record<string, string>) => {
        const token = getCookie('authToken');
        if (!token) throw new Error('No authentication token found');

        const queryString = params ? '?' + new URLSearchParams(params).toString() : '';
        const response = await fetch(`${API_URL}${endpoint}${queryString}`, {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
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
    }, []);

    // Dashboard data fetching functions
    const fetchDashboardAnalytics = useCallback(async (days: number = 30) => {
        try {
            setLoading(true);
            setError(null);
            const data = await fetchWithAuth('/dashboard/analytics', { days: days.toString() });
            setDashboardData(data);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            console.error('Error fetching dashboard analytics:', err);
        } finally {
            setLoading(false);
        }
    }, [fetchWithAuth]);

    const fetchModelPerformance = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);
            const data = await fetchWithAuth('/api/admin/dashboard/model-performance');
            setDashboardData(current => current ? { ...current, model_performance: data } : null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            console.error('Error fetching model performance:', err);
        } finally {
            setLoading(false);
        }
    }, [fetchWithAuth]);

    const fetchSalesAnalytics = useCallback(async (days: number = 30) => {
        try {
            setLoading(true);
            setError(null);
            const data = await fetchWithAuth('/api/admin/dashboard/sales', { days: days.toString() });
            setDashboardData(current => current ? { ...current, sales_analytics: data } : null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            console.error('Error fetching sales analytics:', err);
        } finally {
            setLoading(false);
        }
    }, [fetchWithAuth]);

    const fetchInventoryStatus = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);
            const data = await fetchWithAuth('/api/admin/dashboard/inventory');
            setDashboardData(current => current ? { ...current, inventory_status: data } : null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            console.error('Error fetching inventory status:', err);
        } finally {
            setLoading(false);
        }
    }, [fetchWithAuth]);

    const fetchPricingAnalytics = useCallback(async (days: number = 7) => {
        try {
            setLoading(true);
            setError(null);
            const data = await fetchWithAuth('/api/admin/dashboard/pricing', { days: days.toString() });
            setDashboardData(current => current ? { ...current, pricing_analytics: data } : null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            console.error('Error fetching pricing analytics:', err);
        } finally {
            setLoading(false);
        }
    }, [fetchWithAuth]);

    const fetchCustomerBehavior = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);
            const data = await fetchWithAuth('/api/admin/dashboard/customers');
            setDashboardData(current => current ? { ...current, customer_behavior: data } : null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            console.error('Error fetching customer behavior:', err);
        } finally {
            setLoading(false);
        }
    }, [fetchWithAuth]);

    const fetchOperationalHealth = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);
            const data = await fetchWithAuth('/api/admin/dashboard/health');
            setDashboardData(current => current ? { ...current, operational_health: data } : null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            console.error('Error fetching operational health:', err);
        } finally {
            setLoading(false);
        }
    }, [fetchWithAuth]);

    const fetchTopProducts = useCallback(async (limit: number = 10) => {
        try {
            setLoading(true);
            setError(null);
            const data = await fetchWithAuth('/api/admin/dashboard/top-products', { limit: limit.toString() });
            setDashboardData(current => current ? { ...current, top_products: data } : null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            console.error('Error fetching top products:', err);
        } finally {
            setLoading(false);
        }
    }, [fetchWithAuth]);

    const fetchCategoryRevenue = useCallback(async (days: number = 30) => {
        try {
            setLoading(true);
            setError(null);
            const data = await fetchWithAuth('/api/admin/dashboard/category-revenue', { days: days.toString() });
            setDashboardData(current => current ? { ...current, category_revenue: data } : null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            console.error('Error fetching category revenue:', err);
        } finally {
            setLoading(false);
        }
    }, [fetchWithAuth]);

    const initializeDashboard = useCallback(async () => {
        await fetchDashboardAnalytics();
    }, [fetchDashboardAnalytics]);

    return {
        dashboardData,
        loading,
        error,
        fetchDashboardAnalytics,
        fetchModelPerformance,
        fetchSalesAnalytics,
        fetchInventoryStatus,
        fetchPricingAnalytics,
        fetchCustomerBehavior,
        fetchOperationalHealth,
        fetchTopProducts,
        fetchCategoryRevenue,
        initializeDashboard
    };
}

function getCookie(name: string): string | null {
    if (typeof document === 'undefined') return null; // SSR safety
    const cookies = document.cookie.split(';');
    for (let cookie of cookies) {
        const [cookieName, cookieValue] = cookie.trim().split('=');
        if (cookieName === name) {
            return decodeURIComponent(cookieValue);
        }
    }
    return null;
}

