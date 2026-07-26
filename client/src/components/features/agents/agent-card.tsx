import { UserRound } from "lucide-react"

import type { SchemaAgent } from "@/api/types"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"

interface AgentCardProps {
  agent: SchemaAgent
  isActive?: boolean
}

function getInitials(username: string): string {
  return username
    .split(/[\s._-]+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0].toUpperCase())
    .join("")
}

export function AgentCard({ agent, isActive }: AgentCardProps) {
  return (
    <article className="group border border-transparent p-3 transition-colors hover:border-primary/25 hover:bg-primary/[0.04]">
      <div className="flex items-start gap-3">
        <Avatar className="transition-transform duration-300 group-hover:scale-105">
          <AvatarFallback>{getInitials(agent.username)}</AvatarFallback>
        </Avatar>
        <div className="min-w-0 flex-1">
          <div className="flex items-center justify-between gap-2">
            <h3 className="truncate text-sm font-semibold tracking-wide">
              {agent.username}
            </h3>
            {isActive ? <Badge>Active</Badge> : null}
          </div>
          <p className="mt-1 flex items-center gap-1.5 truncate text-[11px] text-muted-foreground">
            <UserRound className="size-3 text-primary/80" aria-hidden="true" />
            ID: {agent.id.slice(0, 8)}...
          </p>
        </div>
      </div>
    </article>
  )
}
