"use client"

import React, { useEffect, useState } from "react"
import dynamic from "next/dynamic"
import { isTauriEnvironment } from "@/utils/tauri-utils"

// Dynamic import with no SSR to ensure Tauri code only runs on client side
const TauriManagerComponent = dynamic(
  () => import("./tauri-manager").then(mod => ({ default: mod.TauriManager })),
  { ssr: false }
)

export function SafeTauriManager() {
  const [isTauriAvailable, setIsTauriAvailable] = useState(false)
  
  useEffect(() => {
    // Only run client-side
    if (typeof window === 'undefined') return;
    
    // Use our utility function to safely check for Tauri
    const available = isTauriEnvironment();
    setIsTauriAvailable(available);
    
    if (!available) {
      console.log("Tauri is not available in this environment - some desktop features will be disabled");
    }
  }, [])
  
  // Don't render anything if Tauri is not available
  if (!isTauriAvailable) {
    return null;
  }
  
  return <TauriManagerComponent />;
}
