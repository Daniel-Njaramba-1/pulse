import { writable } from "svelte/store";

export interface PriceUpdate {
    product_id: number;
    new_price: number;
    changed_at: string;
    price_change: number;
    change_type: 'increase' | 'decrease' | 'unchanged';
    product_name: string;
}

export const recentPriceUpdates = writable<PriceUpdate[]>([]);
export const productPrices = writable<Map<number, number>>(new Map());
export const connectionStatus = writable<'connected' | 'disconnected' | 'connecting' | 'error'>('disconnected');

// Helper function to get current connection status
function getConnectionStatus(): 'connected' | 'disconnected' | 'connecting' | 'error' {
    let status: 'connected' | 'disconnected' | 'connecting' | 'error' = 'disconnected';
    connectionStatus.subscribe(value => {
        status = value;
    })();
    return status;
}

let eventSource: EventSource | null = null;
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
let connectionAttempts = 0;
const MAX_RECONNECT_ATTEMPTS = 5;

export function connectToSSE() {
    if (eventSource) {
        eventSource.close();
    }

    connectionStatus.set('connecting');
    connectionAttempts++;
    
    try {
        eventSource = new EventSource('http://localhost:8080/api/price-adjustments');
        console.log('EventSource created, waiting for connection...');

        const connectionTimeout = setTimeout(() => {
            connectionStatus.set('connected');
            console.log('SSE connection established');
            connectionAttempts = 0; // Reset attempts on successful connection
        }, 1000);

        eventSource.onmessage = (event) => {
            connectionStatus.set('connected');

            if (reconnectTimer) {
                clearTimeout(reconnectTimer);
                reconnectTimer = null;
            }

            try {
                const update: PriceUpdate = JSON.parse(event.data);
                
                // Add to recent updates (keep last 10)
                recentPriceUpdates.update(updates => {
                    const updated = [update, ...updates.slice(0, 9)];
                    return updated;
                });

                // Update products store with new prices
                productPrices.update(map => {
                    const newMap = new Map(map);
                    newMap.set(update.product_id, update.new_price);
                    return newMap;
                });
                
                // Optional: show notification
                showPriceChangeNotification(update);
            } catch (error) {
                console.error('Error processing SSE message:', error);
            }
        };

        eventSource.onerror = (error) => {
            console.error('SSE connection error:', error);
            connectionStatus.set('error');
            
            // Close the current connection
            if (eventSource) {
                eventSource.close();
                eventSource = null;
            }
            
            // Attempt to reconnect after 5 seconds
            if (!reconnectTimer && connectionAttempts < MAX_RECONNECT_ATTEMPTS) {
                const delay = Math.min(30000, 1000 * Math.pow(2, connectionAttempts));
                console.log(`Will attempt to reconnect in ${delay/1000} seconds...`);
                
                reconnectTimer = setTimeout(() => {
                    reconnectTimer = null;
                    console.log('Attempting to reconnect to SSE...');
                    connectToSSE();
                }, delay);
            } else if (connectionAttempts >= MAX_RECONNECT_ATTEMPTS) {
                console.error('Maximum reconnection attempts reached. Please try again later.');
                connectionStatus.set('disconnected');
            }
        };
    } catch (err) {
        console.error('Failed to create EventSource:', err);
        connectionStatus.set('error');
        
        // Attempt to reconnect after delay if under max attempts
        if (!reconnectTimer && connectionAttempts < MAX_RECONNECT_ATTEMPTS) {
            const delay = Math.min(30000, 1000 * Math.pow(2, connectionAttempts));
            
            reconnectTimer = setTimeout(() => {
                reconnectTimer = null;
                connectToSSE();
            }, delay);
        } else if (connectionAttempts >= MAX_RECONNECT_ATTEMPTS) {
            console.error('Maximum reconnection attempts reached. Please try again later.');
            connectionStatus.set('disconnected');
        }
    }
}

function showPriceChangeNotification(update: PriceUpdate) {
    // Format the price change for display
    const priceChangeAbs = Math.abs(update.price_change).toFixed(2);
    const formattedNewPrice = update.new_price.toFixed(2);
    
    // Create custom event that components can listen for
    const event = new CustomEvent('priceChange', { 
        detail: {
            ...update,
            formattedChange: `$${priceChangeAbs}`,
            formattedPrice: `$${formattedNewPrice}`
        } 
    });
    document.dispatchEvent(event);
    
    // Log to console
    console.log(`${update.product_name} price ${update.change_type}d by $${priceChangeAbs} to $${formattedNewPrice}`);
}

export function disconnectFromSSE() {
    if (eventSource) {
        eventSource.close();
        eventSource = null;
    }
    
    if (reconnectTimer) {
        clearTimeout(reconnectTimer);
        reconnectTimer = null;
    }
    
    connectionStatus.set('disconnected');
}

// Function to help components react to price updates for a specific product
export function subscribeToProductPrice(productId: number, callback: (price: number) => void) {
    // Initial value
    const unsubscribe = productPrices.subscribe(prices => {
        if (prices.has(productId)) {
            callback(prices.get(productId)!);
        }
    });
    
    // Listen for future updates
    const handler = (event: Event) => {
        const detail = (event as CustomEvent<PriceUpdate & { formattedChange: string, formattedPrice: string }>).detail;
        if (detail.product_id === productId) {
            callback(detail.new_price);
        }
    };
    
    document.addEventListener('priceChange', handler as EventListener);
    
    // Return unsubscribe function
    return () => {
        unsubscribe();
        document.removeEventListener('priceChange', handler as EventListener);
    };
}