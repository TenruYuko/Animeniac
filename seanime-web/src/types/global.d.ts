// Global type declarations for the entire application

// Declare Tauri-specific window properties for TypeScript
interface Window {
  // Use a specific type for Tauri internals
  __TAURI_INTERNALS__?: {
    invoke: (cmd: string, args?: any) => Promise<any>;
    transformCallback: (callback: Function, once?: boolean) => (...args: any[]) => any;
    isTauri?: boolean;
    uid?: string;
    [key: string]: any;
  };
  
  // Add WebSocket type for our polyfill
  WebSocket: typeof WebSocket;
}
