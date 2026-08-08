import { Info, UserRound } from "lucide-react"

import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import type { AgentExtended } from "@/api/types"
import { AgentActionsMenu } from "./actions-menu"
import { useNow } from "@/components/time-provider"

interface AgentCardProps {
  agent: AgentExtended
  isActive?: boolean
  onClick?: () => void
  onInfo?: () => void
  expanded?: boolean
}

function getInitials(username: string): string {
  return username
    .split(/[\s._-]+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0].toUpperCase())
    .join("")
}

export function AgentCard({
  agent,
  isActive,
  onClick,
  onInfo,
  expanded,
}: AgentCardProps) {
  const now = useNow()
  const cooldownExpiresAt =
    new Date(agent.cooldown.set_at).getTime() +
    agent.cooldown.duration_seconds * 1000
  const cooldown = Math.max(0, cooldownExpiresAt - now)
  return (
    <article
      className={`group border p-3 transition-colors ${
        isActive
          ? "border-primary/40 bg-primary/8"
          : "border-transparent hover:border-primary/25 hover:bg-primary/4"
      }`}
      onClick={onClick}
    >
      <div className="flex items-start gap-3">
        <Avatar className="transition-transform duration-300 group-hover:scale-105">
          <AvatarFallback>{getInitials(agent.agent.username)}</AvatarFallback>
        </Avatar>
        <div className="min-w-0 flex-1">
          <div className="flex items-center justify-between gap-2">
            <h3 className="truncate text-sm font-semibold tracking-wide">
              {agent.agent.username}
            </h3>
            <div className="flex items-center gap-1.5">
              {isActive ? <Badge>Active</Badge> : null}
              {expanded && (
                <>
                  <button
                    type="button"
                    className="rounded-md p-1 text-muted-foreground transition-colors hover:text-foreground"
                    onClick={(e) => {
                      e.stopPropagation()
                      onInfo?.()
                    }}
                  >
                    <Info className="size-3.5" />
                  </button>
                  <AgentActionsMenu
                    btnClassName="rounded-md p-1 text-muted-foreground transition-colors hover:text-foreground"
                    agent={agent}
                  />
                </>
              )}
            </div>
          </div>
          <div className="flex items-center justify-between gap-2">
            <p className="mt-1 flex items-center gap-1.5 truncate text-[11px] text-muted-foreground">
              <UserRound
                className="size-3 text-primary/80"
                aria-hidden="true"
              />
              ID: {agent.agent.id.slice(0, 8)}...
            </p>
            {cooldown > 0 && <p>{Math.floor(cooldown / 1000)}s</p>}
          </div>
        </div>
      </div>
    </article>
  )
}
