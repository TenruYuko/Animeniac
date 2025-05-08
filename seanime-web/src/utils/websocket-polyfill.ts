/**
 * WebSocket connection polyfill and error handling
 * 
 * This utility provides error handling for WebSocket connections
 * and helps manage reconnection attempts gracefully
 */

// Original WebSocket reference
const OriginalWebSocket = window.WebSocket;

// WebSocket error handler configuration
const wsConfig = {
  maxReconnectAttempts: 5,
  reconnectDelay: 2000,
  debug: false
};

// Create a proxy for WebSocket connections to handle errors
class EnhancedWebSocket extends OriginalWebSocket {
  private reconnectAttempts = 0;
  private originalUrl: string;
  private reconnecting = false;

  constructor(url: string | URL, protocols?: string | string[]) {
    super(url, protocols);
    
    this.originalUrl = url.toString();
    
    // Handle connection errors
    this.addEventListener('error', (event) => {
      if (wsConfig.debug) {
        console.warn(`[WebSocket] Connection error to ${this.originalUrl}`);
      }
      
      // Only try to reconnect if not already reconnecting
      if (!this.reconnecting && this.reconnectAttempts < wsConfig.maxReconnectAttempts) {
        this.attemptReconnect();
      }
    });
    
    // Reset reconnect attempts on successful connection
    this.addEventListener('open', () => {
      this.reconnectAttempts = 0;
      this.reconnecting = false;
      
      if (wsConfig.debug) {
        console.log(`[WebSocket] Connected to ${this.originalUrl}`);
      }
    });
  }
  
  private attemptReconnect(): void {
    this.reconnecting = true;
    this.reconnectAttempts++;
    
    if (wsConfig.debug) {
      console.log(`[WebSocket] Reconnection attempt ${this.reconnectAttempts}/${wsConfig.maxReconnectAttempts} in ${wsConfig.reconnectDelay}ms`);
    }
    
    // Schedule reconnection
    setTimeout(() => {
      try {
        // Create a new connection
        const newSocket = new OriginalWebSocket(this.originalUrl);
        
        // Copy event listeners from old socket to new one
        // This is a simplified approach
        newSocket.onopen = this.onopen;
        newSocket.onclose = this.onclose;
        newSocket.onmessage = this.onmessage;
        newSocket.onerror = this.onerror;
        
        // Replace this socket's properties with the new one
        // This is a basic implementation - in production you'd need more sophisticated proxy handling
        Object.assign(this, newSocket);
        
        this.reconnecting = false;
      } catch (error) {
        console.error('[WebSocket] Reconnection failed:', error);
        this.reconnecting = false;
        
        // Try again if we haven't reached max attempts
        if (this.reconnectAttempts < wsConfig.maxReconnectAttempts) {
          this.attemptReconnect();
        }
      }
    }, wsConfig.reconnectDelay);
  }
}

// Override the global WebSocket constructor if we're in a browser environment
if (typeof window !== 'undefined') {
  try {
    // Only apply in non-Tauri environments to avoid interference
    if (!window.__TAURI_INTERNALS__ || (window.__TAURI_INTERNALS__ as any).isTauri === false) {
      window.WebSocket = EnhancedWebSocket as any;
      console.log('[WebSocket] Polyfill applied for better error handling');
    }
  } catch (error) {
    console.error('[WebSocket] Failed to apply polyfill:', error);
  }
}

// Export the enhanced WebSocket for explicit usage
export { EnhancedWebSocket };

// Using centralized type definitions from /src/types/global.d.ts
