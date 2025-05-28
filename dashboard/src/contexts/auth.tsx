import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';


interface UserData {
    id: number;
    email: string;
    username: string;
}

interface AuthState {
    token: string | null;
    isAuthenticated: boolean;
    user: UserData | null;
    isLoading: boolean;
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

interface AuthContextType extends AuthState {
    login: (username: string, password: string) => Promise<LoginResponse>;
    logout: () => void;
    register: (email: string, username: string, password: string) => Promise<RegisterResponse>;
    getAuthHeader: () => Record<string, string>;
    refreshAuthStateFromCookies: () => boolean;
    isAuthInitialized: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const PUBLIC_ROUTES = [
    "/login", "/register"
];

export const PROTECTED_ROUTES = [
    "/", 
];

// Cookie utility functions
function setCookie(name: string, value: string, days: number = 365) {
    if (typeof document === 'undefined') return; // SSR safety
    const expires = new Date(Date.now() + days * 86400000).toUTCString();
    document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/; SameSite=Strict`;
}

function getCookie(name: string): string | null {
    if (typeof document === 'undefined') return null; // SSR safety
    const cookies = document.cookie.split(';');
    for (let cookie of cookies) {
        const [cookieName, cookieValue] = cookie.trim().split('=');
        if (cookieName === name) {
            return decodeURIComponent(cookieValue);
        }
    }
    return null;
}

function deleteCookie(name: string) {
    if (typeof document === 'undefined') return; // SSR safety
    document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; SameSite=Strict`;
}

export function AuthProvider({ children }: { children: ReactNode }) {
    const navigate = useNavigate();
    const [isAuthInitialized, setIsAuthInitialized] = useState(false);
    
    // Initialize auth state from cookies
    const [authState, setAuthState] = useState<AuthState>(() => {
        const storedToken = getCookie("authToken");
        let initialUser: UserData | null = null;
        
        if (getCookie("userData")) {
            try {
                initialUser = JSON.parse(getCookie("userData") || "null");
            } catch (e) {
                console.error("Failed to parse user data, error:", e);
            }
        }

        return {
            token: storedToken,
            isAuthenticated: !!storedToken,
            user: initialUser,
            isLoading: false,
        };
    });

    // Initialize auth state on mount (equivalent to Svelte's reactive behavior)
    useEffect(() => {
        refreshAuthStateFromCookies();
        setIsAuthInitialized(true);
    }, []);

    const login = async (username: string, password: string): Promise<LoginResponse> => {
        try {
            const endpoint = "http://localhost:8080/api/admin/login";
            const response = await fetch(endpoint, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
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
            
            setAuthState(prev => ({
                ...prev,
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
    };

    const logout = () => {
        deleteCookie("authToken");
        deleteCookie("userData");
        
        setAuthState(prev => ({
            ...prev,
            token: null,
            isAuthenticated: false,
            user: null
        }));
        
        navigate("/login");
    };

    const register = async (email: string, username: string, password: string): Promise<RegisterResponse> => {
        try {
            const endpoint = "http://localhost:8080/api/admin/register";
            const response = await fetch(endpoint, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, username, password })
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
            
            setAuthState(prev => ({
                ...prev,
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
    };

    const getAuthHeader = (): Record<string, string> => {
        return authState.token ? { Authorization: `Bearer ${authState.token}` } : {};
    };

    // Helper to check auth status from cookies (equivalent to Svelte's refreshAuthStateFromCookies)
    const refreshAuthStateFromCookies = (): boolean => {
        if (typeof document === 'undefined') return false;
        
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
        
        setAuthState(prev => ({
            ...prev,
            token,
            isAuthenticated: !!token,
            user
        }));
        
        return !!token;
    };

    const value: AuthContextType = {
        ...authState,
        login,
        logout,
        register,
        getAuthHeader,
        refreshAuthStateFromCookies,
        isAuthInitialized
    };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}

// Optional: Custom hook for checking if auth is ready (useful for loading states)
export function useAuthReady() {
    const { isAuthInitialized } = useAuth();
    return isAuthInitialized;
}