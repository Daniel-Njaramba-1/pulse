import { writable } from "svelte/store";
import { authHelpers } from "./auth";

export interface Category {
    id: number;
    name: string;
    description: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface CategoryResponse {
    success: boolean;
    data?: Category | Category[];
    error?: string;
}

// Create a store for categories
export const categories = writable<Category[]>([]);
export const isLoading = writable<boolean>(false);
export const error = writable<string | null>(null);

// Base API URL
const API_BASE_URL = "http://localhost:8080/api/admin";

// Helper functions for category operations
export const categoryHelpers = {
    // Fetch all categories
    fetchCategories: async (): Promise<CategoryResponse> => {
        isLoading.set(true);
        error.set(null);

        try {
            const response = await fetch(`${API_BASE_URL}/categories`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.message || "Failed to fetch categories.");
                return { success: false, error: errorData.message };
            }

            const data: Category[] = await response.json();
            categories.set(data);
            return { success: true, data };
        } catch (err) {
            console.error("Error fetching categories:", err);
            error.set("An error occurred while fetching categories.");
            return { success: false, error: "An error occurred while fetching categories." };
        } finally {
            isLoading.set(false);
        }
    },

    // Get a single category by ID
    getCategory: async (id: number): Promise<CategoryResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/categories/${id}`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                return { success: false, error: errorData.message || "Failed to fetch category." };
            }

            const data: Category = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error fetching category:", err);
            return { success: false, error: "An error occurred while fetching the category." };
        }
    },

    // Create a new category
    createCategory: async (name: string, description: string): Promise<CategoryResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/categories`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
                body: JSON.stringify({ name, description}),
            });

            if (!response.ok) {
                const errorData = await response.json();
                return { success: false, error: errorData.message || "Failed to create category." };
            }

            const data: Category = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error creating category:", err);
            return { success: false, error: "An error occurred while creating the category." };
        }
    },

    // Update an existing category
    updateCategory: async (id: number, name: string, description: string): Promise<CategoryResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/categories/${id}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
                body: JSON.stringify({ id, name, description }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                return { success: false, error: errorData.message || "Failed to update category." };
            }

            const data: Category = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error updating category:", err);
            return { success: false, error: "An error occurred while updating the category." };
        }
    },

    // Delete a category
    deleteCategory: async (id: number): Promise<CategoryResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/categories/${id}`, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                return { success: false, error: errorData.message || "Failed to delete category." };
            }
            
            return { success: true };
        } catch (err) {
            console.error("Error deleting category:", err);
            return { success: false, error: "An error occurred while deleting the category." };
        }
    },
};