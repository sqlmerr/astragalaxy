import { createFileRoute } from "@tanstack/react-router"
import { useEffect } from "react"
import { useNavigate } from "@tanstack/react-router"

import { Spinner } from "@/components/ui/spinner"
import { useAuth } from "@/components/features/auth/auth-provider"
import { LoginForm } from "@/components/features/auth/login-form"

export const Route = createFileRoute("/login")({ component: LoginPage })

export function LoginPage() {
  const { isAuthenticated, isReady } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (isReady && isAuthenticated) {
      void navigate({ to: "/", replace: true })
    }
  }, [isAuthenticated, isReady, navigate])

  if (!isReady || isAuthenticated) {
    return (
      <main className="grid min-h-svh place-items-center bg-background">
        <div className="flex items-center gap-3 text-xs font-semibold tracking-widest text-muted-foreground uppercase">
          <Spinner className="text-primary" />
          Synchronizing
        </div>
      </main>
    )
  }

  return (
    <main className="relative grid min-h-svh grid-cols-1 place-items-center overflow-hidden px-4 py-8 sm:px-6 lg:items-center lg:gap-16 lg:px-12 xl:px-20">
      <div
        className="star-field pointer-events-none absolute inset-0"
        aria-hidden="true"
      />
      <section className="relative mx-auto mb-10 w-full max-w-xl">
        <h1 className="mx-auto mt-3 max-w-lg text-center text-4xl leading-tight font-semibold tracking-wide text-balance uppercase sm:text-5xl">
          astragalaxy
        </h1>
      </section>
      <section
        className="relative mx-auto w-full max-w-2xl"
        aria-label="Authentication"
      >
        <LoginForm />
      </section>
    </main>
  )
}
