import { browser } from "$app/environment";
import { goto } from "$app/navigation";
import { get, writable } from "svelte/store";

interface AuthState {
    token: string | null;
    isAuthenticated: boolean;
    user: UserData | null;
    isLoading: boolean;
}

interface UserData {
    id: number;
    email: string;
    username: string;
}

interface LoginResponse {
    success: boolean;
    token?: string;
    user?: UserData;
    error?: string;
}

interface RegisterResponse {
    success: boolean;
    token?: string;
    user?: UserData;
    error?: string;
}

function setCookie(name:string, value:string, days:number = 7) {
    if (!browser) return;
    const expires = new Date(Date.now() + days * 86400000).toUTCString();
    document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/; SameSite=Strict`;
}

function getCookie(name: string): string | null {
    if (!browser) return null;
    const cookies = document.cookie.split(';');
    for (let cookie of cookies) {
        const [cookieName, cookieValue] = cookie.trim().split('=');
        if (cookieName === name) {
            return decodeURIComponent(cookieValue);
        }
    }
    return null;
}

function deleteCookie(name:string) {
    if (!browser) return;
    document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; SameSite=Strict`;
}

// get intial values from cookie
const storedToken = browser ? getCookie("authToken") : null;
let initialUser: UserData | null = null;

export const PUBLIC_ROUTES = [
    "/login", "/register"
];

export const PROTECTED_ROUTES = [
    "/", "/categories", "/brands", "/products", "/stocks", "/sales", "/customers", "/profile"
];

if (browser && getCookie("userData")) {
    try {
        initialUser = JSON.parse(getCookie("userData") || "null");
    } catch (e) {
        console.error("Failed to parse user data, error:", e);
    }
}   

export const auth = writable<AuthState>({
    token: storedToken,
    isAuthenticated: !!storedToken,
    user: initialUser,
    isLoading: false,
});

export const isAuthInitialized = writable<boolean>(false);

export const authHelpers = {
    login: async (username:string, password:string): Promise<LoginResponse> => {
        try {
            const endpoint = "http://localhost:8080/api/admin/login"
            const response = await fetch(endpoint,{
                method: "POST",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({ username, password }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                return { 
                    success: false, 
                    error: errorData.message || "Failed to login." 
                };
            }

            const data = await response.json();
            setCookie("authToken", data.token);
            if (data.user) {
                setCookie("userData", JSON.stringify(data.user));
            }
            auth.update((state) => ({
                ...state,
                token: data.token,
                isAuthenticated: true,
                user: data.user || null
            }));
            console.log("Login successful:", data);
            return { 
                success: true, 
                token: data.token,
                user: data.user 
            };
        } catch (error) {
            console.error("Login error:", error);
            return { 
                success: false, 
                error: error instanceof Error ? error.message : "Unknown error occurred"
            };
        }
    },

    logout: () => {
        deleteCookie("authToken");
        deleteCookie("userData");
        
        auth.update((state) => ({
            ...state,
            token: null,
            isAuthenticated: false,
            user: null
        }));
        
        goto("/login");
    },

    register: async (email:string, username:string, password:string): Promise<RegisterResponse> => {
        try {
            const endpoint = "http://localhost:8080/api/admin/register";
            const response = await fetch(endpoint, {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({ email, username, password})
            });

            if (!response.ok) {
                const errorData = await response.json();
                return { 
                    success: false, 
                    error: errorData.message || "Failed to register." 
                };
            }
            const data = await response.json();
            setCookie("authToken", data.token);
            if (data.user) {
                setCookie("userData", JSON.stringify(data.user));
            }
            auth.update((state) => ({
                ...state,
                token: data.token,
                isAuthenticated: true,
                user: data.user || null
            }));
            return { 
                success: true,
                token: data.token,
                user: data.user
            };
        } catch (error) {
            console.error("Registration error:", error);
            return { 
                success: false, 
                error: error instanceof Error ? error.message : "Unknown error occurred" 
            };
        }
    },

    getAuthHeader: (): Record<string, string> => {
        const { token } = get(auth);
        return token ? { Authorization: `Bearer ${token}` } : {};
    },


    // Helper to check auth status from cookies
    refreshAuthStateFromCookies: () => {
        if (!browser) return;
        
        const token = getCookie("authToken");
        let user = null;
        
        try {
            const userData = getCookie("userData");
            if (userData) {
                user = JSON.parse(userData);
            }
        } catch (e) {
            console.error("Failed to parse user data from cookie:", e);
        }
        
        auth.update(state => ({
            ...state,
            token,
            isAuthenticated: !!token,
            user
        }));
        
        return !!token;
    }
}