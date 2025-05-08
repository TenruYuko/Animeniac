/**
 * Utility functions for Tauri compatibility and safety checks
 */

/**
 * Checks if the current environment is running in Tauri
 */
export function isTauriEnvironment(): boolean {
  try {
    return typeof window !== 'undefined' && 
           window !== null && 
           'object' === typeof window && 
           '__TAURI_INTERNALS__' in window;
  } catch (error) {
    console.log("Error checking Tauri environment:", error);
    return false;
  }
}

/**
 * Safely imports a Tauri module only if running in Tauri
 * Falls back to a mock implementation for browser environments
 * 
 * @param importFn Function that imports the Tauri module
 * @param mockImplementation Mock implementation for browser environments
 */
export async function safeTauriImport<T>(
  importFn: () => Promise<T>,
  mockImplementation: Partial<T> = {} as Partial<T>
): Promise<Partial<T>> {
  if (isTauriEnvironment()) {
    try {
      return await importFn();
    } catch (error) {
      console.error("Failed to import Tauri module:", error);
      return mockImplementation;
    }
  }
  return mockImplementation;
}

/**
 * Creates a safe proxy for Tauri imports that won't throw errors in browser environments
 */
export function createSafeTauriProxy<T extends object>(
  mockImplementation: Partial<T> = {} as Partial<T>
): T {
  return new Proxy({} as T, {
    get: (target, prop) => {
      // If the property exists in the mock, return it
      if (prop in mockImplementation) {
        return mockImplementation[prop as keyof typeof mockImplementation];
      }
      
      // Otherwise return a safe function that does nothing
      return typeof prop === 'string' ? 
        (...args: any[]) => {
          console.log(`Called Tauri API '${String(prop)}' in browser environment`);
          return Promise.resolve(null);
        } : 
        undefined;
    }
  });
}

// Import types from types.d.ts - don't redeclare Window interface here
