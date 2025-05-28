import { writable } from "svelte/store";
import { authHelpers } from "./auth";

export interface WishlistDetail {
    id: number;
    customer_id: number;
    is_active: boolean;
    items: WishlistItem[];
}

export interface WishlistItem {
    id: number;
    wishlist_id: number;
    product_id: number;
    product_name: string;
    product_image_path: string;
    product_adjusted_price: number;
    product_stock_quantity: number;
}

export interface WishlistResponse {
    success: boolean;
    data?: WishlistDetail;
    error?: string;
}

const API_BASE_URL = "http://localhost:8080/api/customer";

export const wishlistHelpers = {
    fetchWishlist: async (): Promise<WishlistResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/wishlist`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                }
            });
            
            if (!response.ok) {
                const errorData = await response.json()
                const errorMessage = errorData.error || "Failed to fetch wishlist";
                return { success: false, error: errorMessage };
            }

            const data: WishlistDetail = await response.json();
            return { success: true, data };
        } catch (err) {
            return { success: false, error: "Failed to fetch wishlist." };
        }
    },

    removeFromWishlist: async (productId: number): Promise<{ success: boolean; message?: string; error?: string }> => {
        try {
            const response = await fetch(`${API_BASE_URL}/remove-from-wishlist/${productId}`, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                }
            });
            
            if (!response.ok) {
                const errorData = await response.json()
                const errorMessage = errorData.error || "Failed to remove item from wishlist";
                return { success: false, error: errorMessage };
            }
            
            const data = await response.json();
            return { success: true, message: data.message };
        } catch (err) {
            return { success: false, error: "Failed to remove from wishlist." };
        }
    },
}