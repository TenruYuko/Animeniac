import React, { createContext, ReactNode, useContext, useEffect } from "react"
import { ClientSettings, useClientSettings } from "../api/hooks/client-settings.hooks"
import { useTheme } from "next-themes"

// Define the context type
interface ClientSettingsContextType {
  settings: ClientSettings | undefined
  isLoading: boolean
  updateSetting: <K extends keyof ClientSettings>(key: K, value: ClientSettings[K]) => void
}

// Create the context
const ClientSettingsContext = createContext<ClientSettingsContextType | undefined>(undefined)

// Provider component for the client settings
export const ClientSettingsProvider = ({ children }: { children: ReactNode }) => {
  const { clientSettings, isLoading, updateClientSetting } = useClientSettings()
  const { setTheme } = useTheme()

  // Sync dark mode with Next.js theme
  useEffect(() => {
    if (clientSettings && 'darkMode' in clientSettings && clientSettings.darkMode !== undefined) {
      setTheme(clientSettings.darkMode ? "dark" : "light")
    }
  }, [clientSettings, setTheme])

  // Context value
  const value: ClientSettingsContextType = {
    settings: clientSettings,
    isLoading,
    updateSetting: updateClientSetting,
  }

  return <ClientSettingsContext.Provider value={value}>{children}</ClientSettingsContext.Provider>
}

// Hook for using the client settings context
export const useClientSettingsContext = () => {
  const context = useContext(ClientSettingsContext)
  if (context === undefined) {
    throw new Error("useClientSettingsContext must be used within a ClientSettingsProvider")
  }
  return context
}
