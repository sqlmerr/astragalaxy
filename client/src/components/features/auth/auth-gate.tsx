import { useEffect } from "react"
import type { ReactNode } from "react"
import { useNavigate } from "@tanstack/react-router"

import { Spinner } from "@/components/ui/spinner"
import { useAuth } from "@/components/features/auth/auth-provider"

export function AuthGate({ children }: { children: ReactNode }) {
  const { isAuthenticated, isReady } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (isReady && !isAuthenticated) {
      void navigate({ to: "/login", replace: true })
    }
  }, [isAuthenticated, isReady, navigate])

  if (!isReady || !isAuthenticated) {
    return (
      <main className="grid min-h-svh place-items-center bg-background">
        <div className="flex items-center gap-3 text-xs font-semibold tracking-widest text-muted-foreground uppercase">
          <Spinner className="text-primary" />
          Loading...
        </div>
      </main>
    )
  }

  return <>{children}</>
}
