import { writable } from "svelte/store";

export interface Product {
    id: number;
    brand: string;
    category: string;
    name: string;
    description: string;
    image: File;
    imagePath: string;
    isActive: boolean;
}

export interface ProductResponse {
    success: boolean;
    data?: Product | Product[];
    error?: string;
}

export const products = writable<Product[]>([]);
export const isLoading = writable<boolean>(false);
export const error = writable<string | null>(null);

// Base API URL
const API_BASE_URL = "http://localhost:8080/api/admin";

export const productHelpers = {
    fetchProducts: async (): Promise<ProductResponse> => {
        isLoading.set(true);
        error.set(null);
        
        try {
            const response = await fetch (`${API_BASE_URL}/products`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            });
            
            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.message || "Failed to fetch products");
                return { success:false, error:errorData.message};
            }
        
            const data: Product[] = await response.json();
            products.set(data);
            return {success:true, data};
        } catch (err) {
            console.error("Error fetching products", err)
            error.set("An error occured while fetching products");
            return {success:false, error:"An error occured while fetching products"};
        } finally {
            isLoading.set(false);
        }
    },

    getProduct: async(id: number): Promise<ProductResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/products/${id}`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                return { success: false, error: errorData.message || "Failed to fetch product." };
            }

            const data: Product = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error fetching product:", err);
            return { success: false, error: "An error occurred while fetching the product" };
        }
    },

    createProduct: async(brandId:number, categoryId:number, name:string, description: string, image: File | null): Promise<ProductResponse> => {
        isLoading.set(true);
        error.set(null);
        
        try {
            // Create FormData object to handle file uploads
            const formData = new FormData();
            formData.append("brand_id", brandId.toString());
            formData.append("category_id", categoryId.toString());
            formData.append("name", name);
            formData.append("description", description);
            if (image) {
                formData.append("image", image);
            }

            const formDataObj = Object.fromEntries(formData.entries());
            console.log('FormData as object:', formDataObj);
            
            const response = await fetch(`${API_BASE_URL}/products`, {
                method: "POST",
                body: formData
            });
            
            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.error || "Failed to create product");
                return { success: false, error: errorData.error };
            }
            
            const data: Product = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error creating product:", err);
            error.set("An error occurred while creating the product");
            return { success: false, error: "An error occurred while creating the product" };
        } finally {
            isLoading.set(false);
        }
    },

    updateProduct: async(id: number, brand:any, category:any, name:string, description: string, image:any): Promise<ProductResponse> => {
        isLoading.set(true);
        error.set(null);
        
        try {
            // Create FormData object to handle file uploads
            const formData = new FormData();
            formData.append("brand", typeof brand === 'object' ? brand.id.toString() : brand.toString());
            formData.append("category", typeof category === 'object' ? category.id.toString() : category.toString());
            formData.append("name", name);
            formData.append("description", description);
            
            // Only append image if it exists and is a new file
            if (image instanceof File) {
                formData.append("image", image);
            }
            
            const response = await fetch(`${API_BASE_URL}/products/${id}`, {
                method: "PUT",
                body: formData
            });
            
            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.error || "Failed to update product");
                return { success: false, error: errorData.error };
            }
            
            const data: Product = await response.json();

            return { success: true, data };
        } catch (err) {
            console.error("Error updating product:", err);
            error.set("An error occurred while updating the product");
            return { success: false, error: "An error occurred while updating the product" };
        } finally {
            isLoading.set(false);
        }
    },

    deleteProduct: async(id: number): Promise<ProductResponse> => {
        isLoading.set(true);
        error.set(null);
        
        try {
            const response = await fetch(`${API_BASE_URL}/products/${id}`, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                },
            });
            
            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.error || "Failed to delete product");
                return { success: false, error: errorData.error };
            }
            
            return { success: true };
        } catch (err) {
            console.error("Error deleting product:", err);
            error.set("An error occurred while deleting the product");
            return { success: false, error: "An error occurred while deleting the product" };
        } finally {
            isLoading.set(false);
        }
    },

    deactivateProduct: async(id: number): Promise<ProductResponse> => {
        isLoading.set(true);
        error.set(null);
        
        try {
            const response = await fetch(`${API_BASE_URL}/products/${id}/deactivate`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
            });
            
            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.error || "Failed to deactivate product");
                return { success: false, error: errorData.error };
            }
            products.update(currentProducts => 
                currentProducts.map(product => 
                    product.id === id ? { ...product, isActive: false } : product
                )
            );
            return { success: true };
        } catch (err) {
            console.error("Error deactivating product:", err);
            error.set("An error occurred while deactivating the product");
            return { success: false, error: "An error occurred while deactivating the product" };
        } finally {
            isLoading.set(false);
        }
    },

    reactivateProduct: async(id: number): Promise<ProductResponse> => {
        isLoading.set(true);
        error.set(null);
        
        try {
            const response = await fetch(`${API_BASE_URL}/products/${id}/reactivate`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
            });
            
            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.error || "Failed to reactivate product");
                return { success: false, error: errorData.error };
            }
            
            // Update the product's status in the store
            products.update(currentProducts => 
                currentProducts.map(product => 
                    product.id === id ? { ...product, isActive: true } : product
                )
            );
            
            return { success: true };
        } catch (err) {
            console.error("Error reactivating product:", err);
            error.set("An error occurred while reactivating the product");
            return { success: false, error: "An error occurred while reactivating the product" };
        } finally {
            isLoading.set(false);
        }
    },
}
