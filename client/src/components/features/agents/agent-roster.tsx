import { Crosshair, UsersRound } from "lucide-react"

import { useMyAgentsQuery } from "@/api/hooks"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Spinner } from "@/components/ui/spinner"
import { AgentCard } from "@/components/features/agents/agent-card"

export function AgentRoster() {
  const { data, isPending, isError } = useMyAgentsQuery()
  const agents = data?.data ?? []

  return (
    <aside
      className="fixed top-4 left-4 z-20 w-[calc(100vw-2rem)] max-w-sm lg:top-6 lg:left-6"
      aria-label="Player agents"
    >
      <Card className="surface-panel overflow-hidden">
        <CardHeader className="flex-row items-center justify-between border-b border-border py-4">
          <div className="flex items-center gap-2">
            <UsersRound className="size-4 text-primary" aria-hidden="true" />
            <CardTitle className="text-xs tracking-[0.16em] uppercase">
              Field agents
            </CardTitle>
          </div>
          <span className="text-[10px] font-bold tracking-widest text-primary">
            {isPending ? "..." : agents.length}
          </span>
        </CardHeader>
        <CardContent className="p-0">
          {isPending ? (
            <div className="flex h-48 items-center justify-center gap-2 text-xs text-muted-foreground">
              <Spinner className="text-primary" />
              Loading agents
            </div>
          ) : isError ? (
            <div className="flex h-32 items-center justify-center px-5 text-center text-xs text-muted-foreground">
              Unable to load agents. Check your connection and try again.
            </div>
          ) : agents.length === 0 ? (
            <div className="flex h-32 items-center justify-center px-5 text-center text-xs text-muted-foreground">
              No agents registered yet. Deploy your first agent to begin
              operations.
            </div>
          ) : (
            <ScrollArea className="max-h-60 divide-y divide-border">
              {agents.map((agent) => (
                <AgentCard key={agent.id} agent={agent} />
              ))}
            </ScrollArea>
          )}
        </CardContent>
      </Card>
    </aside>
  )
}
