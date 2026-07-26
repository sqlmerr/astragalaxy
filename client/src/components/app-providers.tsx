import {
  QueryCache,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query"
import { useState } from "react"
import type { ReactNode } from "react"

import { getErrorMessage } from "@/api/errors"
import { Toaster } from "@/components/ui/toaster"
import { AuthProvider } from "@/components/features/auth/auth-provider"
import { ToastProvider } from "@/lib/toast"

export function AppProviders({ children }: { children: ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            retry: false,
            staleTime: 30_000,
          },
        },
        queryCache: new QueryCache({
          onError: (error) => {
            const message = getErrorMessage(error)

            if (typeof window !== "undefined") {
              const AUTH_TOKEN_KEY = "astragalaxy.auth.token"
              const token = window.localStorage.getItem(AUTH_TOKEN_KEY)

              if (
                token &&
                (message.toLowerCase().includes("invalid") ||
                  message.toLowerCase().includes("unauthorized") ||
                  message.toLowerCase().includes("expired"))
              ) {
                window.localStorage.removeItem(AUTH_TOKEN_KEY)
                window.location.href = "/login"
                return
              }
            }

            console.error("[Query Error]", message)
          },
        }),
      })
  )

  return (
    <QueryClientProvider client={queryClient}>
      <ToastProvider>
        <AuthProvider>{children}</AuthProvider>
        <Toaster />
      </ToastProvider>
    </QueryClientProvider>
  )
}
