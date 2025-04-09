import { redirect, type Handle } from "@sveltejs/kit";
import { PUBLIC_ROUTES, PROTECTED_ROUTES } from "$lib/stores/auth";

export const handle: Handle = async ({ event, resolve }) => {
    const path = event.url.pathname;
    const isProtectedRoute = PROTECTED_ROUTES.some(route => path === route || path.startsWith(route + '/'));
    const isPublicRoute = PUBLIC_ROUTES.some(route => path === route || path.startsWith(route + '/'));
    
    // Get the auth token from cookies on the server
    const authToken = event.cookies.get("authToken");
    const isAuthenticated = !!authToken;
    
    // Redirect unauthenticated users from protected routes to login
    if (isProtectedRoute && !isAuthenticated) {
        console.log("User is not authenticated, redirecting to login page");
        redirect(303, "/login");
    }
    
    // Redirect authenticated users from public routes (like login page) to homepage
    if (isPublicRoute && isAuthenticated) {
        console.log("User is authenticated, redirecting to homepage");
        redirect(303, "/");
    }

    // Continue with the request
    return resolve(event);
};