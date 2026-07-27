import { ClientOnly, useNavigate } from "@tanstack/react-router"
import { LogOut, Satellite, UserRound } from "lucide-react"

import { useMeQuery } from "@/api/hooks"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { AgentRoster } from "@/components/features/agents/agent-roster"
import { useAuth } from "@/components/features/auth/auth-provider"
import { GalaxyMap } from "@/components/features/map/galaxy-map"

export function StarMapLayout() {
  const { data: user } = useMeQuery()
  const { signOut } = useAuth()
  const navigate = useNavigate()

  function handleSignOut() {
    signOut()
    void navigate({ to: "/login", replace: true })
  }

  return (
    <div className="relative h-screen overflow-hidden bg-background">
      <ClientOnly>
        <GalaxyMap />
      </ClientOnly>

      <div
        className="star-field pointer-events-none absolute inset-0 z-10"
        aria-hidden="true"
      />

      <div className="relative z-20">
        <AgentRoster />

        <div className="fixed top-4 right-4 flex items-center gap-2 lg:top-6 lg:right-6">
          <div className="hidden items-center gap-2 border border-border bg-card/80 px-3 py-2 backdrop-blur-sm sm:flex">
            <UserRound className="size-3.5 text-primary" />
            <span className="text-xs font-semibold tracking-wide">
              {user?.username ?? "Commander"}
            </span>
          </div>

          <Button variant="outline" size="sm" onClick={handleSignOut}>
            <LogOut />
            <span className="hidden sm:inline">Log out</span>
          </Button>
        </div>
      </div>
    </div>
  )
}
