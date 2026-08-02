import { ChevronLeft, ChevronRight, UsersRound } from "lucide-react"
import { useState } from "react"

import { useAuth } from "@/components/features/auth/auth-provider"
import { AgentCard } from "@/components/features/agents/agent-card"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Spinner } from "@/components/ui/spinner"

import { useAgents } from "../auth/use-agents"
import { AgentModal } from "./agent-modal"

export function AgentRoster() {
  const { data: agents = [], isPending, isError } = useAgents()
  const { currentAgentID, setCurrentAgentID } = useAuth()
  const [selectedAgentId, setSelectedAgentId] = useState<string | null>(null)
  const [isCollapsed, setIsCollapsed] = useState(false)

  const selectedAgent =
    agents.find((s) => s.agent.id === selectedAgentId) ?? null

  return (
    <>
      <aside
        className={`fixed top-4 left-4 z-20 transition-all duration-300 lg:top-6 lg:left-6 ${
          isCollapsed ? "w-10" : "w-[calc(100vw-2rem)] max-w-sm"
        }`}
        aria-label="Player agents"
      >
        <Card className="overflow-hidde2">
          <CardHeader
            className={`w-full flex-row items-center border-b border-border p-3 ${
              isCollapsed ? "justify-center border-b-0" : ""
            }`}
          >
            <div className="flex w-full items-center justify-between">
              {!isCollapsed ? (
                <div className="flex items-center gap-2">
                  <UsersRound
                    className="size-4 text-primary"
                    aria-hidden="true"
                  />
                  <CardTitle className="tracking-[0.16em] uppercase">
                    Your agents
                  </CardTitle>
                </div>
              ) : (
                <span />
              )}

              <Button
                variant="ghost"
                size="icon"
                className="size-7 shrink-0"
                onClick={() => setIsCollapsed((prev) => !prev)}
              >
                {isCollapsed ? (
                  <ChevronRight className="size-3.5" />
                ) : (
                  <ChevronLeft className="size-3.5" />
                )}
              </Button>
            </div>
          </CardHeader>
          {!isCollapsed && (
            <CardContent className="">
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
                <ScrollArea className="flex max-h-60 flex-col divide-y divide-border">
                  {agents.map((agent) => (
                    <AgentCard
                      key={agent.agent.id}
                      agent={agent}
                      isActive={agent.agent.id === currentAgentID}
                      onClick={() => setCurrentAgentID(agent.agent.id)}
                      onInfo={() => setSelectedAgentId(agent.agent.id)}
                      expanded={true}
                    />
                  ))}
                </ScrollArea>
              )}
            </CardContent>
          )}
        </Card>
      </aside>
      <AgentModal
        agent={selectedAgent}
        setAgent={(a) => setSelectedAgentId(a?.agent.id || null)}
      />
    </>
  )
}
