import { writable } from "svelte/store";
import { authHelpers } from "./auth";

export interface Product {
    id: number;
    brand: string;
    category: string;
    name: string;
    description: string;
    image: File;
    image_path: string;
    is_active: boolean;
    base_price: number
}

export interface ProductDetail {
    // Basic product information
    id: number;
    name: string;
    description: string | null;
    image_path: string | null;
    is_active: boolean;
    created_at: string;
    updated_at: string;

    // Associated brand information
    brand_id: number;
    brand_name: string;

    // Associated category information
    category_id: number;
    category_name: string;

    // Product metrics
    average_rating: number | null;
    review_count: number | null;
    wishlist_count: number | null;
    base_price: number | null;
    adjusted_price: number | null;

    // Stock information
    stock_quantity: number | null;
    stock_threshold: number | null;
}

export interface ProductResponse {
    success: boolean;
    data?: Product | Product[];
    error?: string;
}

export interface ProductDetailResponse {
    success: boolean;
    data?: ProductDetail | ProductDetail[];
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
            const response = await fetch(`${API_BASE_URL}/products`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
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

    getProduct: async(id: number): Promise<ProductDetailResponse> => {
        try {
            const response = await fetch(`${API_BASE_URL}/products/${id}`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                return { success: false, error: errorData.message || "Failed to fetch product." };
            }

            const data: ProductDetail = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error fetching product:", err);
            return { success: false, error: "An error occurred while fetching the product" };
        }
    },

    createProduct: async(brandId:number, categoryId:number, name:string, description: string, base_price:number, initial_stock:number, image: File | null): Promise<ProductResponse> => {
        isLoading.set(true);
        error.set(null);
        
        try {
            // Create FormData object to handle file uploads
            const formData = new FormData();
            formData.append("brand_id", brandId.toString());
            formData.append("category_id", categoryId.toString());
            formData.append("name", name);
            formData.append("description", description);
            formData.append("base_price", base_price.toString())
            formData.append("initial_stock", initial_stock.toString());
            if (image) {
                formData.append("image", image);
            }

            const formDataObj = Object.fromEntries(formData.entries());
            console.log('FormData as object:', formDataObj);
            
            const response = await fetch(`${API_BASE_URL}/products`, {
                method: "POST",
                headers: {
                    ...authHelpers.getAuthHeader()
                },
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

    updateProductDetails: async(id: number, brandId: number, categoryId: number, name: string, description: string, is_active: boolean): Promise<ProductResponse> => {
        isLoading.set(true);
        error.set(null);
        
        try {
            const payload = {
                id: id,
                brand_id: brandId,
                category_id: categoryId,
                name,
                description,
                is_active
            };
            
            const response = await fetch(`${API_BASE_URL}/products/${id}/details`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
                body: JSON.stringify(payload)
            });
            
            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.error || "Failed to update product details");
                return { success: false, error: errorData.error };
            }
            
            const data: Product = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error updating product details:", err);
            error.set("An error occurred while updating the product details");
            return { success: false, error: "An error occurred while updating the product details" };
        } finally {
            isLoading.set(false);
        }
    },

    updateProductImage: async(id: number, image: File): Promise<ProductResponse> => {
        isLoading.set(true);
        error.set(null);
        
        try {
            // Create FormData object to handle file upload
            const formData = new FormData();
            formData.append("image", image);
            
            const response = await fetch(`${API_BASE_URL}/products/${id}/image`, {
                method: "PUT",
                headers: {
                    ...authHelpers.getAuthHeader()
                },
                body: formData
            });
            
            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.error || "Failed to update product image");
                return { success: false, error: errorData.error };
            }
            
            const data: Product = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error updating product image:", err);
            error.set("An error occurred while updating the product image");
            return { success: false, error: "An error occurred while updating the product image" };
        } finally {
            isLoading.set(false);
        }
    },
    
    restock: async(id: number, quantity: number): Promise<ProductResponse> => {
        isLoading.set(true);
        error.set(null);
        
        try {
            const response = await fetch(`${API_BASE_URL}/products/${id}/restock`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
                body: JSON.stringify({ id, quantity })
            });
            
            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.error || "Failed to add product stock");
                return { success: false, error: errorData.error };
            }
            
            const data = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error adding product stock:", err);
            error.set("An error occurred while adding product stock");
            return { success: false, error: "An error occurred while adding product stock" };
        } finally {
            isLoading.set(false);
        }
    },

    updateBasePrice: async(id: number, base_price: number): Promise<ProductResponse> => {
        isLoading.set(true);
        error.set(null);
        
        try {
            const response = await fetch(`${API_BASE_URL}/products/${id}/reprice`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                    ...authHelpers.getAuthHeader()
                },
                body: JSON.stringify({ id, base_price })
            });
            
            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.error || "Failed to update product price");
                return { success: false, error: errorData.error };
            }
            
            const data = await response.json();
            return { success: true, data };
        } catch (err) {
            console.error("Error updating product price:", err);
            error.set("An error occurred while updating the product price");
            return { success: false, error: "An error occurred while updating the product price" };
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
                    ...authHelpers.getAuthHeader()
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
                    ...authHelpers.getAuthHeader()
                },
            });
            
            if (!response.ok) {
                const errorData = await response.json();
                error.set(errorData.error || "Failed to deactivate product");
                return { success: false, error: errorData.error };
            }
            products.update(currentProducts => 
                currentProducts.map(product => 
                    product.id === id ? { ...product, is_active: false } : product
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
                    ...authHelpers.getAuthHeader()
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
                    product.id === id ? { ...product, is_active: true } : product
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