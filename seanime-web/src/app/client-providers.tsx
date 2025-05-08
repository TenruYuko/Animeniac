"use client"
import { WebsocketProvider } from "@/app/websocket-provider"
import { CustomCSSProvider } from "@/components/shared/custom-css-provider"
import { CustomThemeProvider } from "@/components/shared/custom-theme-provider"
import { Toaster } from "@/components/ui/toaster"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { createStore } from "jotai"
import { Provider as JotaiProvider } from "jotai/react"
import { ThemeProvider } from "next-themes"
import { usePathname } from "next/navigation"
import React from "react"
import { CookiesProvider } from "react-cookie"
import { ClientSettingsProvider } from "@/context/client-settings.context"

interface ClientProvidersProps {
    children?: React.ReactNode
}

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            refetchOnWindowFocus: false,
            retry: 0,
        },
    },
})

export const ClientProviders: React.FC<ClientProvidersProps> = ({ children }) => {
    const [store] = React.useState(createStore())
    const pathname = usePathname()

    return (
        <ThemeProvider attribute="class" defaultTheme="dark" forcedTheme={(pathname === "/docs") ? "light" : "dark"}>
            <CookiesProvider>
                <JotaiProvider store={store}>
                    <QueryClientProvider client={queryClient}>
                        <ClientSettingsProvider>
                            <CustomCSSProvider>
                                <WebsocketProvider>
                                    {children}
                                    <CustomThemeProvider />
                                    <Toaster />
                                </WebsocketProvider>
                            </CustomCSSProvider>
                        </ClientSettingsProvider>
                        {/*{process.env.NODE_ENV === "development" && <React.Suspense fallback={null}>*/}
                        {/*    <ReactQueryDevtools />*/}
                        {/*</React.Suspense>}*/}
                    </QueryClientProvider>
                </JotaiProvider>
            </CookiesProvider>
        </ThemeProvider>
    )

}
