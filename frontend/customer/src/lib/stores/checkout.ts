import { authHelpers } from "./auth"

export interface Order {
    id: number;
    customer_id: number;
    total_price: number;
    status: string;
    items: OrderItem[];
    price_valid_until: string;
    created_at: string;
}

export interface OrderItem {
    id: number;
    product_id: number;
    order_id: number;
    quantity: number;
    price: number;
    created_at: string;
    updated_at: string;
    product_name: string;
    product_image_path: string;
}

export interface OrderResponse {
    success: boolean;
    data?: Order;
    error?: string;
}

const API_BASE_URL = "http://localhost:8080/api/customer"

export const checkoutHelpers = {
    generateOrder: async (): Promise<OrderResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/order`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                const errorMessage = errorData.error || "Failed to generate order";
                return { success: false, error: errorMessage}
            }

            return {success: true}
        } catch (err) {
            return { success: false, error: "Failed to generate order" }
        }
    },

    fetchOrderWithItems: async (): Promise<OrderResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/order-with-items`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                const errorMessage = errorData.error || "Failed to fetch order with items";
                return { success: false, error: errorMessage}
            }

            const data: Order = await response.json();
            return { success: true, data };
        } catch {
            return { success: false, error: "Failed to fetch order" };
        }
    }, 

    payment: async (): Promise<OrderResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/payment`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                const errorMessage = errorData.error || "Failed to pay for order";
                return { success: false, error: errorMessage}
            }

            return {success: true}
        } catch (err) {
            return { success: false, error: "Failed to pay for order" }
        }
    },
}