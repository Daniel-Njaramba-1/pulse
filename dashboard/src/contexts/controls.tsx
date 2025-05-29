import { createContext, useContext, useState, type ReactNode } from "react";
import { useCallback } from "react";

export type CoefficientData = {
    model_version: string;
    training_date: string; // ISO string
    sample_size: number;
    r_squared: number;
    mse: number;
    rmse: number;
    mae: number;
    days_since_last_sale_coef: number;
    sales_velocity_coef: number;
    total_sales_count_coef: number;
    total_sales_value_coef: number;
    category_percentile_coef: number;
    review_score_coef: number;
    wishlist_to_sales_ratio_coef: number;
    days_since_restock_coef: number;
    created_at: string; // ISO string
    updated_at: string; // ISO string
};

type ModelMetrics = {
    r_squared: number;
    rmse: number;
    mae: number;
    sample_size: number;
};

type ControlsContextType = {
    coefficients: CoefficientData | null;
    metrics: ModelMetrics | null;
    loading: boolean;
    error: string | null;
    fetchCoefficients: () => Promise<void>;
    triggerAdjustAllPrices: () => Promise<void>;
    triggerTrainModel: () => Promise<void>;
};

const ControlsContext = createContext<ControlsContextType | undefined>(undefined);

export const ControlsProvider = ({ children }: { children: ReactNode }) => {
    const [coefficients, setCoefficients] = useState<CoefficientData | null>(null);
    const [metrics, setMetrics] = useState<ModelMetrics | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Helper to get cookie value
    const getCookie = (name: string) => {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${name}=`);
        if (parts.length === 2) return parts.pop()!.split(';').shift();
        return null;
    };

    const API_URL = "http://localhost:8080/api/admin";

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
        if (!data.Success) {
            throw new Error(data.Error || 'Failed to fetch data');
        }

        return data.Data;
    }, []);

    const fetchCoefficients = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const data: CoefficientData = await fetchWithAuth("/dashboard/coefficients");
            console.log("API response:", data);
            setCoefficients(data);

            // Extract metrics from the coefficient data
            setMetrics({
                r_squared: data.r_squared,
                rmse: data.rmse,
                mae: data.mae,
                sample_size: data.sample_size
            });
        } catch (e: any) {
            setError(e.message || "Failed to fetch coefficients");
        } finally {
            setLoading(false);
        }
    }, [fetchWithAuth]);

    const triggerAdjustAllPrices = async () => {
        setLoading(true);
        setError(null);
        try {
            const res = await fetch("http://localhost:5872/adjust-prices", { method: "POST" });
            if (!res.ok) throw new Error("Failed to adjust all prices");
        } catch (e: any) {
            setError(e.message || "Failed to adjust all prices");
        } finally {
            setLoading(false);
        }
    };

    const triggerTrainModel = async () => {
        setLoading(true);
        setError(null);
        try {
            const res = await fetch("http://localhost:5872/train-model", { method: "POST" });
            const data = await res.json();
            if (data.metrics) setMetrics(data.metrics);
            else setError(data.message || "Failed to train model");
        } catch (e: any) {
            setError(e.message || "Failed to train model");
        } finally {
            setLoading(false);
        }
    };

    return (
        <ControlsContext.Provider
            value={{
                coefficients,
                metrics,
                loading,
                error,
                fetchCoefficients,
                triggerAdjustAllPrices,
                triggerTrainModel,
            }}
        >
            {children}
        </ControlsContext.Provider>
    );
};

export const useControls = () => {
    const ctx = useContext(ControlsContext);
    if (!ctx) throw new Error("useControls must be used within ControlsProvider");
    return ctx;
};