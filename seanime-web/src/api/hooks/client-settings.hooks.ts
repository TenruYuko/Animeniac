import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { buildSeaQuery } from "../client/requests"
import { useState, useEffect } from "react"

// Define the client settings interface
export interface ClientSettings {
  id?: number
  browserId: string
  darkMode?: boolean
  accentColor?: string
  customAccentColor?: string
  uiAnimations?: boolean
  blurEffects?: boolean
  compactView?: boolean
  showAdult?: boolean
  defaultListTab?: string
  defaultSortOrder?: string
  listViewMode?: string
  defaultAudioTrack?: string
  defaultSubTrack?: string
  autoplayNext?: boolean
  preferredResolution?: string
  notifyUpdates?: boolean
  notifyNewEpisodes?: boolean
  extraData?: string
  createdAt?: string
  updatedAt?: string
}

// Query keys
export const clientSettingsKeys = {
  all: ["clientSettings"] as const,
  details: () => [...clientSettingsKeys.all, "details"] as const,
}

// Helper functions for API requests
const getClientSettings = async (): Promise<ClientSettings> => {
  const response = await buildSeaQuery<{ data: ClientSettings }>({
    endpoint: "/api/v1/client-settings",
    method: "GET",
  })
  if (!response || !response.data) {
    throw new Error("Failed to fetch client settings")
  }
  return response.data
}

const updateClientSettings = async (settings: Partial<ClientSettings>): Promise<ClientSettings> => {
  const response = await buildSeaQuery<{ data: ClientSettings }>({
    endpoint: "/api/v1/client-settings",
    method: "PUT",
    data: settings,
  })
  if (!response || !response.data) {
    throw new Error("Failed to update client settings")
  }
  return response.data
}

// Hook for using client settings
export const useClientSettings = () => {
  const queryClient = useQueryClient()

  // Query for fetching client settings
  const {
    data: clientSettings,
    isLoading,
    error,
    isError,
  } = useQuery({
    queryKey: clientSettingsKeys.details(),
    queryFn: getClientSettings,
    staleTime: 1000 * 60 * 60, // 1 hour
    gcTime: 1000 * 60 * 60 * 24, // 24 hours
  })

  // Mutation for updating client settings
  const { mutate: updateSettings, isPending } = useMutation({
    mutationFn: updateClientSettings,
    onSuccess: (data) => {
      queryClient.setQueryData(clientSettingsKeys.details(), data)
      // Success notification can be handled by the component using this hook
    },
    onError: (err) => {
      // Error notification can be handled by the component using this hook
      console.error("Failed to update client settings:", err)
    },
  })

  // Function to update specific settings
  const updateClientSetting = <K extends keyof ClientSettings>(key: K, value: ClientSettings[K]) => {
    if (!clientSettings) return
    
    updateSettings({
      ...clientSettings,
      [key]: value,
    })
  }

  return {
    clientSettings,
    isLoading,
    isError,
    error,
    updateClientSetting,
    updateSettings,
    isUpdating: isPending,
  }
}
