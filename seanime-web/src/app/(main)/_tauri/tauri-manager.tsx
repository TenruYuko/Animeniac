"use client"

import React, { useEffect, useState } from "react"
import { safeTauriImport, createSafeTauriProxy } from "@/utils/tauri-utils"
import mousetrap from "mousetrap"

// Safe imports to prevent errors in browser environments
let tauriEvent: any = null;
let tauriWindow: any = null;
let tauriWebview: any = null;

type TauriManagerProps = {
    children?: React.ReactNode
}

// This is only rendered on the Desktop client
export function TauriManager(props: TauriManagerProps) {
    const {
        children,
        ...rest
    } = props

    const [isInitialized, setIsInitialized] = useState(false)

    // Initialize Tauri APIs safely
    useEffect(() => {
        const initTauriApis = async () => {
            try {
                // Safely import Tauri modules
                const eventModule = await import("@tauri-apps/api/event");
                const windowModule = await import("@tauri-apps/api/window");
                const webviewModule = await import("@tauri-apps/api/webviewWindow");
                
                tauriEvent = eventModule;
                tauriWindow = windowModule;
                tauriWebview = webviewModule;
                
                setIsInitialized(true);
            } catch (error) {
                console.error("Failed to initialize Tauri APIs:", error);
                // Create safe proxies that won't throw errors
                tauriEvent = createSafeTauriProxy();
                tauriWindow = createSafeTauriProxy();
                tauriWebview = createSafeTauriProxy();
            }
        };
        
        initTauriApis();
    }, []);

    useEffect(() => {
        if (!isInitialized) return;
        
        let unlisten: any = null;
        
        try {
            // Setup event listeners
            const listenPromise = tauriEvent.listen("message", (event: any) => {
                const message = event.payload;
                console.log("Received message from Rust:", message);
            });
            
            listenPromise.then((unlistenFn: () => void) => {
                unlisten = unlistenFn;
            });

            // Setup keyboard shortcuts
            mousetrap.bind("f11", () => {
                toggleFullscreen();
            });

            mousetrap.bind("esc", () => {
                try {
                    const appWindow = new tauriWindow.Window("main");
                    appWindow.isFullscreen().then((isFullscreen: boolean) => {
                        if (isFullscreen) {
                            toggleFullscreen();
                        }
                    });
                } catch (error) {
                    console.error("Error checking fullscreen state:", error);
                }
            });

            document.addEventListener("fullscreenchange", toggleFullscreen);
        } catch (error) {
            console.error("Error setting up Tauri event listeners:", error);
        }

        return () => {
            try {
                if (unlisten) unlisten();
                mousetrap.unbind("f11");
                document.removeEventListener("fullscreenchange", toggleFullscreen);
            } catch (error) {
                console.error("Error cleaning up Tauri event listeners:", error);
            }
        };
    }, [isInitialized]);

    function toggleFullscreen() {
        try {
            const appWindow = new tauriWindow.Window("main");

            // Only toggle fullscreen on the main window
            const currentWindow = tauriWebview.getCurrentWebviewWindow();
            if (!currentWindow || currentWindow.label !== "main") return;

            appWindow.isFullscreen().then((fullscreen: boolean) => {
                try {
                    // DEVNOTE: When decorations are not shown in fullscreen move there will be a gap at the bottom of the window (Windows)
                    // Hide the decorations when exiting fullscreen
                    // Show the decorations when entering fullscreen
                    appWindow.setDecorations(!fullscreen);
                    appWindow.setFullscreen(!fullscreen);
                } catch (error) {
                    console.error("Error setting fullscreen:", error);
                }
            }).catch((error: any) => {
                console.error("Error checking fullscreen state:", error);
            });
        } catch (error) {
            console.error("Error in toggleFullscreen:", error);
        }
    }

    return (
        <>

        </>
    )
}
