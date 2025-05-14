import { writable } from "svelte/store";
import { authHelpers } from "./auth";

export interface Product {
    id: number;
    name: string;
    description: string;
    image_path: string;
    is_active: boolean;
    brand_id: number;
    brand_name: string;
    category_id: number;
    category_name: string;
    average_rating: number;
    review_count: number;
    wishlist_count: number;
    base_price: number;
    adjusted_price: number;
    stock_quantity: number;
    stock_threshold: number;
}
export interface ProductResponse {
    success: boolean;
    data?: Product | Product[];
    error?: string;
}

export const products = writable<Product[]>([]);
export const isLoading = writable<boolean>(false);
export const error = writable<string | null>(null);

const API_BASE_URL = "http://localhost:8080/api/customer";

export const productHelpers = {
    fetchProducts: async (): Promise<ProductResponse> => {
        isLoading.set(true);
        error.set(null);

        try {
            const response = await fetch(`${API_BASE_URL}/products`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });


            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.message || "Failed to fetch products.");
                return { success: false, error: errorData.message || "Failed to fetch products." };
            }

            const data: Product[] = await response.json();
            products.set(data);
            return {success: true, data};
        } catch (err) {
            error.set("Failed to fetch products.");
            return { success: false, error: "Failed to fetch products." };
        } finally {
            isLoading.set(false);
        }
    },

    getProduct: async (productId: number): Promise<ProductResponse> => {
        error.set(null);

        try {
            const response = await fetch(`${API_BASE_URL}/product-by-id/${productId}`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.message || "Failed to fetch product.");
                return { success: false, error: errorData.message || "Failed to fetch product." };
            }

            const data: Product = await response.json();
            return {success: true, data};
        } catch (err) {
            error.set("Failed to fetch product.");
            return { success: false, error: "Failed to fetch product." };
        }  
    },

    // Helper to create URL-friendly slug for a product
    createProductSlug: (product: Product): string => {
        // Convert product name to lowercase, replace spaces and special chars with dashes
        const nameSlug = product.name
            .toLowerCase()
            .replace(/[^\w\s-]/g, '') // Remove special characters
            .replace(/\s+/g, '-')     // Replace spaces with dashes
            .replace(/-+/g, '-');     // Replace multiple dashes with single dash
        
        // Append product ID to ensure uniqueness
        return `${nameSlug}-${product.id}`;
    },

    extractIdFromSlug: (slug: string): number | null => {
        const parts = slug.split('-');
        const idString = parts[parts.length - 1];

        if (parts.length > 1 && idString ) {
            const id = parseInt(idString, 10)
            if (!isNaN(id)) {
                return id;
            }
        }

        console.error(`Could not extract valid ID from slug: "${slug}"`);
        return null;
    }

}
