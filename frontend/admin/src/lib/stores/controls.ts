import { writable } from "svelte/store";
import { authHelpers } from "./auth";

// Define interfaces
export interface CoefficientData {
    model_version: string;
    training_date: string;
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
    created_at: string;
    updated_at: string;
}

export interface CoefficientResponse {
    success: boolean;
    data?: CoefficientData;
    error?: string;
}

// Create stores
export const coefficients = writable<CoefficientData | null>(null);
export const isLoading = writable<boolean>(false);
export const error = writable<string | null>(null);
export const isAdjustingPrices = writable<boolean>(false);
export const isTrainingModel = writable<boolean>(false);

// Base API URL
const API_BASE_URL = "http://localhost:8080/api/admin";

// Helper functions for coefficient operations
export const coefficientHelpers = {
    // Fetch coefficients
    fetchCoefficients: async (): Promise<CoefficientResponse> => {
        isLoading.set(true);
        error.set(null);

        try {
            const response = await fetch(`${API_BASE_URL}/dashboard/coefficients`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.message || "Failed to fetch coefficients.");
                return { success: false, error: errorData.message };
            }

            const data = await response.json();
            coefficients.set(data.data);
            return { success: true, data: data.data };
        } catch (err) {
            console.error("Error fetching coefficients:", err);
            error.set("An error occurred while fetching coefficients.");
            return { success: false, error: "An error occurred while fetching coefficients." };
        } finally {
            isLoading.set(false);
        }
    },

    // Trigger price adjustments
    adjustPrices: async (): Promise<{ success: boolean; error?: string }> => {
        isAdjustingPrices.set(true);
        error.set(null);

        try {
            const response = await fetch("http://localhost:5872/adjust-prices", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.message || "Failed to adjust prices.");
                return { success: false, error: errorData.message };
            }

            return { success: true };
        } catch (err) {
            console.error("Error adjusting prices:", err);
            error.set("An error occurred while adjusting prices.");
            return { success: false, error: "An error occurred while adjusting prices." };
        } finally {
            isAdjustingPrices.set(false);
        }
    },

    // Trigger model training
    trainModel: async (): Promise<{ success: boolean; error?: string }> => {
        isTrainingModel.set(true);
        error.set(null);

        try {
            const response = await fetch("http://localhost:5872/train-model", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.message || "Failed to train model.");
                return { success: false, error: errorData.message };
            }

            // Refresh coefficients after training
            await coefficientHelpers.fetchCoefficients();
            return { success: true };
        } catch (err) {
            console.error("Error training model:", err);
            error.set("An error occurred while training model.");
            return { success: false, error: "An error occurred while training model." };
        } finally {
            isTrainingModel.set(false);
        }
    }
};
