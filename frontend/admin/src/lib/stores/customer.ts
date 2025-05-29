import { writable } from "svelte/store";
import { authHelpers } from "./auth";

// Define interfaces 
export interface Customer {
    id: number;
    username: string;
    email: string;
    is_active: boolean;
}

export interface CustomerResponse {
    success: boolean;
    data?: Customer | Customer[];
    error?: string;
}

// Create a store for customers
export const customers = writable<Customer[]>([]); 
export const isLoading = writable<boolean>(false);
export const error = writable<string | null>(null);

// Base API URL
const API_BASE_URL = "http://localhost:8080/api/admin";

// Helper functions for customer operations 
export const customerHelpers = {
    // Fetch all customers
    fetchCustomers: async (): Promise<CustomerResponse> => {
        isLoading.set(true);
        error.set(null);

        try {
            const response = await fetch(`${API_BASE_URL}/customers`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.message || "Failed to fetch customers.");
                return { success: false, error: errorData.message };
            }

            const data: Customer[] = await response.json();
            console.log(data)
            customers.set(data);
            return { success: true, data };
        } catch (err) {
            console.error("Error fetching customers:", err);
            error.set("An error occurred while fetching customers.");
            return { success: false, error: "An error occurred while fetching customers." };
        } finally {
            isLoading.set(false);
        }
    }
};
