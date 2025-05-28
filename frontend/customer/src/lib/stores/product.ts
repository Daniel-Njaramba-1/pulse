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
    is_in_wishlist?: boolean;
}

export interface Review {
    id: number;
    customer_name: string;
    rating: number;
    review_text: string;
    created_at: string;
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

    checkProductInWishlist: async(productId: number): Promise<{success: boolean; inWishlist?: boolean; error?: string }> => {
        error.set(null);

        try {
            const response = await fetch(`${API_BASE_URL}/check-product-in-wishlist/${productId}`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });
            const data = await response.json();
            if (!response.ok) {
                error.set(data.error || "Failed to check wishlist status.");
                return { success: false, error: data.error || "Failed to check wishlist status." };
            }
            return { success: true, inWishlist: data.in_wishlist };
        } catch (err) {
            error.set("Failed to check wishlist status.");
            return { success: false, error: "Failed to check wishlist status." };
        }
    },

    addToWishlist: async (productId: number): Promise<{success: boolean; message?: string; error?: string}> => {
        error.set(null);
        try {
            const response = await fetch(`${API_BASE_URL}/add-to-wishlist/${productId}`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });
            const data = await response.json();
            if (!response.ok) {
                error.set(data.error || "Failed to add to wishlist.");
                return { success: false, error: data.error || "Failed to add to wishlist." };
            }
            return { success: true, message: data.message };
        } catch (err) {
            error.set("Failed to add to wishlist.");
            return { success: false, error: "Failed to add to wishlist." };
        } 
    }, 

    removeFromWishlist: async (productId: number): Promise<{ success: boolean; message?: string; error?: string }> => {
        error.set(null);
        try {
            const response = await fetch(`${API_BASE_URL}/remove-from-wishlist/${productId}`, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                }
            });
            const data = await response.json();
            if (!response.ok) {
                error.set(data.error || "Failed to remove from wishlist.");
                return { success: false, error: data.error || "Failed to remove from wishlist." };
            }
            return { success: true, message: data.message };
        } catch (err) {
            error.set("Failed to remove from wishlist.");
            return { success: false, error: "Failed to remove from wishlist." };
        }
    },

    fetchWishlist: async (): Promise<{ success: boolean; data?: Product[]; error?: string }> => {
        error.set(null);
        try {
            const response = await fetch(`${API_BASE_URL}/wishlist`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                }
            });
            const data = await response.json();
            if (!response.ok) {
                error.set(data.error || "Failed to fetch wishlist.");
                return { success: false, error: data.error || "Failed to fetch wishlist." };
            }
            return { success: true, data };
        } catch (err) {
            error.set("Failed to fetch wishlist.");
            return { success: false, error: "Failed to fetch wishlist." };
        }
    },

    fetchProductReviews: async (productId: number): Promise<{ success: boolean; data?: Review[]; error?: string }> => {
        error.set(null);
        try {
            const response = await fetch(`${API_BASE_URL}/product-reviews/${productId}`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                }
            });
            const data = await response.json();
            if (!response.ok) {
                error.set(data.error || "Failed to fetch product reviews.");
                return { success: false, error: data.error || "Failed to fetch product reviews." };
            }
            return { success: true, data };
        } catch (err) {
            error.set("Failed to fetch product reviews.");
            return { success: false, error: "Failed to fetch product reviews." };
        }
    },

    reviewProduct: async (productId: number, rating: number, reviewText: string): Promise<{ success: boolean; message?: string; error?: string }> => {
        error.set(null);
        try {
            const response = await fetch(`${API_BASE_URL}/review-product`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
                body: JSON.stringify({
                    product_id: productId,
                    rating,
                    review_text: reviewText
                })
            });
            const data = await response.json();
            if (!response.ok) {
                error.set(data.error || "Failed to submit review.");
                return { success: false, error: data.error || "Failed to submit review." };
            }
            return { success: true, message: data.message };
        } catch (err) {
            error.set("Failed to submit review.");
            return { success: false, error: "Failed to submit review." };
        }
    },

    verifyPurchase: async (productId: number): Promise<{ success: boolean; purchased?: boolean; error?: string }> => {
        error.set(null);
        try {
            const response = await fetch(`${API_BASE_URL}/verify-purchase/${productId}`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                }
            });
            const data = await response.json();
            if (!response.ok) {
                error.set(data.error || "Failed to verify purchase.");
                return { success: false, error: data.error || "Failed to verify purchase." };
            }
            return { success: true, purchased: data.purchased };
        } catch (err) {
            error.set("Failed to verify purchase.");
            return { success: false, error: "Failed to verify purchase." };
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
