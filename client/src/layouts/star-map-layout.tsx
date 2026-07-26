import { useNavigate } from "@tanstack/react-router"
import { LogOut, Satellite, UserRound } from "lucide-react"

import { useMeQuery } from "@/api/hooks"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { AgentRoster } from "@/components/features/agents/agent-roster"
import { useAuth } from "@/components/features/auth/auth-provider"

export function StarMapLayout() {
  const { data: user } = useMeQuery()
  const { signOut } = useAuth()
  const navigate = useNavigate()

  function handleSignOut() {
    signOut()
    void navigate({ to: "/login", replace: true })
  }

  return (
    <div className="relative min-h-svh overflow-hidden bg-background">
      <div
        className="star-field pointer-events-none absolute inset-0"
        aria-hidden="true"
      />
      <AgentRoster />

      <div className="fixed top-4 right-4 z-20 flex items-center gap-2 lg:top-6 lg:right-6">
        <div className="hidden items-center gap-2 border border-border bg-card/80 px-3 py-2 backdrop-blur-sm sm:flex">
          <UserRound className="size-3.5 text-primary" aria-hidden="true" />
          <span className="text-xs font-semibold tracking-wide">
            {user?.username ?? "Commander"}
          </span>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={handleSignOut}
          aria-label="Sign out"
        >
          <LogOut data-icon="inline-start" aria-hidden="true" />
          <span className="hidden sm:inline">Log out</span>
        </Button>
      </div>

      <main className="relative z-10 min-h-svh p-4 pt-[21.5rem] sm:p-6 sm:pt-[21rem] lg:pt-6 lg:pl-[22rem]">
        <Card className="space-grid flex h-[calc(100svh-22.5rem)] min-h-80 items-center justify-center overflow-hidden border-primary/20 bg-card/35 p-6 backdrop-blur-[2px] sm:h-[calc(100svh-22rem)] lg:h-[calc(100svh-3rem)]">
          <div className="relative flex max-w-sm flex-col items-center text-center">
            <div
              className="absolute -top-10 size-32 rounded-full border border-primary/15"
              aria-hidden="true"
            />
            <div
              className="absolute -top-4 size-20 rounded-full border border-primary/20"
              aria-hidden="true"
            />
            <Satellite
              className="relative size-8 text-primary"
              aria-hidden="true"
            />
            <p className="relative mt-6 text-[10px] font-bold tracking-[0.24em] text-primary uppercase">
              Navigation array
            </p>
            <h1 className="relative mt-2 text-xl font-semibold tracking-wide">
              Star Map (Coming Soon)
            </h1>
            <p className="relative mt-3 text-sm leading-relaxed text-muted-foreground">
              This sector is reserved for the live stellar navigation interface.
            </p>
          </div>
        </Card>
      </main>
    </div>
  )
}
