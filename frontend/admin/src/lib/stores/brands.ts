import { writable } from "svelte/store";

// Define interfaces 
export interface Brand {
    id: number;
    name: string;
    description: string;
}

export interface BrandResponse {
    success: boolean;
    data?: Brand | Brand[];
    error?: string;
}

// Create a store for brands
export const brands = writable<Brand[]>([]); 
export const isLoading = writable<boolean>(false);
export const error = writable<string | null>(null);

// Base API URL
const API_BASE_URL = "http://localhost:8080/api/admin";

// Helper functions for brand operations 
export const brandHelpers = {
    // Fetch all brands
    fetchBrands: async (): Promise<BrandResponse> => {
        isLoading.set(true);
        error.set(null);

        try {
            const response = await fetch(`${API_BASE_URL}/brands`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.message || "Failed to fetch brands.");
                return { success: false, error: errorData.message };
            }

            const data: Brand[] = await response.json();
            brands.set(data);
            return { success: true, data };
        } catch (err) {
            console.error("Error fetching brands:", err);
            error.set("An error occurred while fetching brands.");
            return { success: false, error: "An error occurred while fetching brands." };
        } finally {
            isLoading.set(false);
        }
    },

    // Get a single brand by ID
    getBrand: async (id: number): Promise<BrandResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/brands/${id}`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                return { success: false, error: errorData.message || "Failed to fetch brand." };
            }

            const data: Brand = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error fetching brand:", err);
            return { success: false, error: "An error occurred while fetching the brand." };
        }
    },

    // Create a new brand
    createBrand: async (name: string, description: string): Promise<BrandResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/brands`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ name, description }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                return { success: false, error: errorData.message || "Failed to create brand." };
            }

            const data: Brand = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error creating brand:", err);
            return { success: false, error: "An error occurred while creating the brand." };
        }
    },

    // Update an existing brand
    updateBrand: async (id: number, name: string, description: string): Promise<BrandResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/brands/${id}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ id, name, description }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                return { success: false, error: errorData.message || "Failed to update brand." };
            }

            const data: Brand = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error updating brand:", err);
            return { success: false, error: "An error occurred while updating the brand." };
        }
    },

    // Delete a brand
    deleteBrand: async (id: number): Promise<BrandResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/brands/${id}`, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                return { success: false, error: errorData.message || "Failed to delete brand." };
            }

            return { success: true };
        } catch (err) {
            console.error("Error deleting brand:", err);
            return { success: false, error: "An error occurred while deleting the brand." };
        }
    },
};