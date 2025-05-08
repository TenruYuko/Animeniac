/**
 * Tauri environment shim script
 * This script creates mock Tauri objects in non-Tauri environments to prevent errors
 */

// Create enhanced WebSocket implementation for better connection handling
class EnhancedWebSocket extends WebSocket {
  private reconnectAttempts = 0;
  private originalUrl: string;
  private reconnecting = false;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 2000;

  constructor(url: string | URL, protocols?: string | string[]) {
    super(url, protocols);
    
    this.originalUrl = url.toString();
    
    // Handle connection errors
    this.addEventListener('error', () => {
      console.warn(`[WebSocket] Connection error to ${this.originalUrl}`);
      
      // Only try to reconnect if not already reconnecting
      if (!this.reconnecting && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.attemptReconnect();
      }
    });
    
    // Reset reconnect attempts on successful connection
    this.addEventListener('open', () => {
      this.reconnectAttempts = 0;
      this.reconnecting = false;
      console.log(`[WebSocket] Connected to ${this.originalUrl}`);
    });
  }
  
  private attemptReconnect(): void {
    this.reconnecting = true;
    this.reconnectAttempts++;
    
    console.log(`[WebSocket] Reconnection attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts} in ${this.reconnectDelay}ms`);
    
    // Schedule reconnection
    setTimeout(() => {
      try {
        // Create a new connection with the same URL and protocols
        const newSocket = new WebSocket(this.originalUrl);
        
        // Copy event listeners
        if (this.onopen) newSocket.onopen = this.onopen;
        if (this.onclose) newSocket.onclose = this.onclose;
        if (this.onmessage) newSocket.onmessage = this.onmessage;
        if (this.onerror) newSocket.onerror = this.onerror;
        
        // Replace this socket's properties
        for (const prop in newSocket) {
          if (Object.prototype.hasOwnProperty.call(newSocket, prop)) {
            (this as any)[prop] = (newSocket as any)[prop];
          }
        }
        
        this.reconnecting = false;
      } catch (error) {
        console.error('[WebSocket] Reconnection failed:', error);
        this.reconnecting = false;
        
        // Try again if we haven't reached max attempts
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
          this.attemptReconnect();
        }
      }
    }, this.reconnectDelay);
  }
}

// Create a basic implementation of the Tauri internals for browser environments
const createTauriShim = () => {
  // Only run in browser environments
  if (typeof window === 'undefined') return;

  // Apply WebSocket polyfill for better connection handling
  try {
    const OriginalWebSocket = window.WebSocket;
    window.WebSocket = EnhancedWebSocket as any;
    console.log('[WebSocket] Enhanced WebSocket applied for better error handling');
  } catch (error) {
    console.error('[WebSocket] Failed to apply enhanced WebSocket:', error);
  }
  
  // Only define if it doesn't already exist (to avoid overriding actual Tauri implementation)
  if (!('__TAURI_INTERNALS__' in window)) {
    console.log('Creating Tauri shim for browser environment');
    
    // Create a basic mock of Tauri internals
    Object.defineProperty(window, '__TAURI_INTERNALS__', {
      value: {
        // Mock the invoke function
        invoke: (cmd: string, args?: any) => {
          console.log(`[Tauri Shim] invoke called with command: ${cmd}`);
          return Promise.resolve(null);
        },
        // Mock the transformCallback function
        transformCallback: (callback: Function, once = false) => {
          return (...args: any[]) => {
            console.log('[Tauri Shim] transformCallback called');
            return callback(...args);
          };
        },
        // Add other needed Tauri internals
        isTauri: false,
        uid: Math.random().toString(36).substring(2),
      },
      writable: false,
      configurable: false
    });
  }
};

// Execute immediately
createTauriShim();

// Using centralized type definitions from /src/types/global.d.ts

export {};
