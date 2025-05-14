import { writable } from "svelte/store";
import { authHelpers } from "./auth";

export interface Cart {
    id: number;
    customer_id: number;
    is_active: boolean;
    items: CartItem[];
    total_items: number;
    total_price: number;
}

export interface CartItem {
    id: number;
    product_id: number;
    cart_id: number;
    quantity: number;
    price: number;
    is_processed: boolean;
    created_at: string;
    updated_at: string;
    product_name: string;
    product_image_path: string;
    product_adjusted_price: number;
    product_stock_quantity: number;
}

export interface CartResponse {
    success: boolean;
    data?: Cart;
    error?: string;
}

const API_BASE_URL = "http://localhost:8080/api/customer";

export const cartHelpers = {
    // get cart & cart items
    fetchCartWithItems: async (): Promise<CartResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/cart-with-items`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                const errorMessage = errorData.error || "Failed to fetch cart with items.";                
                return { success: false, error: errorMessage };
            }

            const data: Cart = await response.json();
            return { success: true, data };
        } catch (err) {
            return { success: false, error: "Failed to fetch cart." };
        }
    },

    addToCart: async (productId: number, quantity: number): Promise<CartResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/add-to-cart`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
                body: JSON.stringify({ product_id: productId, quantity: quantity })
            });

            if (!response.ok) {
                const errorData = await response.json();
                const errorMessage = errorData.error || "Failed to add item to cart.";                
                return { success: false, error: errorMessage };
            }

            await cartHelpers.fetchCartWithItems();
            return { success: true}
        } catch (err) {
            return { success: false, error: "Failed to add item to cart"}
        }
    },
    
    removeFromCart: async (itemId: number): Promise<CartResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/remove-from-cart`, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
                body: JSON.stringify({ itemId: itemId})
            });

            if (!response.ok) {
                const errorData = await response.json();
                const errorMessage = errorData.error || "Failed to remove item from cart.";                
                return { success: false, error: errorMessage };
            }

            await cartHelpers.fetchCartWithItems();
            return { success: true}
        } catch (err) {
            return { success: false, error: "Failed to remove item from cart"}
        }
    },

    clearCart: async (): Promise<CartResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/clear-cart`, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                const errorMessage = errorData.error || "Failed to clear cart.";                
                return { success: false, error: errorMessage };
            }

            await cartHelpers.fetchCartWithItems();
            return { success: true}
        } catch (err) {
            return { success: false, error: "Failed to clear cart"}
        }
    },
}
