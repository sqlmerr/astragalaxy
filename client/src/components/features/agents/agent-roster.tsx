import { ChevronLeft, ChevronRight, UsersRound } from "lucide-react"
import { useState } from "react"
import type { SchemaAgent } from "@/api/types"

import { useMyAgentsQuery } from "@/api/hooks"
import { useAuth } from "@/components/features/auth/auth-provider"
import { AgentCard } from "@/components/features/agents/agent-card"
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Separator } from "@/components/ui/separator"
import { Spinner } from "@/components/ui/spinner"

import SyntaxHighlighter from "react-syntax-highlighter"
import { gruvboxDark } from "react-syntax-highlighter/dist/esm/styles/hljs"

export function AgentRoster() {
  const { data, isPending, isError } = useMyAgentsQuery()
  const { currentAgentID, setCurrentAgentID } = useAuth()
  const [selectedAgent, setSelectedAgent] = useState<SchemaAgent | null>(null)
  const [isCollapsed, setIsCollapsed] = useState(false)
  const agents = data?.data ?? []

  return (
    <>
      <aside
        className={`fixed top-4 left-4 z-20 transition-all duration-300 lg:top-6 lg:left-6 ${
          isCollapsed ? "w-10" : "w-[calc(100vw-2rem)] max-w-sm"
        }`}
        aria-label="Player agents"
      >
        <Card className="overflow-hidden">
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
                  <CardTitle className="text-xs tracking-[0.16em] uppercase">
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
                    <AgentCard
                      key={agent.id}
                      agent={agent}
                      isActive={agent.id === currentAgentID}
                      onClick={() => setCurrentAgentID(agent.id)}
                      onInfo={() => setSelectedAgent(agent)}
                    />
                  ))}
                </ScrollArea>
              )}
            </CardContent>
          )}
        </Card>
      </aside>
      <Dialog
        open={selectedAgent !== null}
        onOpenChange={(open) => {
          if (!open) setSelectedAgent(null)
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{selectedAgent?.username}</DialogTitle>
          </DialogHeader>
          {selectedAgent && (
            <div className="space-y-6">
              <Card className="p-4">
                <h3 className="mb-3 font-semibold">Information</h3>

                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">ID</span>

                    <code className="font-mono break-all">
                      {selectedAgent.id}
                    </code>
                  </div>

                  <div className="flex justify-between">
                    <span className="text-muted-foreground">User ID</span>

                    <code className="font-mono break-all">
                      {selectedAgent.user_id}
                    </code>
                  </div>
                </div>
              </Card>
              <Separator />
              <Accordion>
                <AccordionItem value="json">
                  <AccordionTrigger>JSON</AccordionTrigger>
                  <AccordionContent>
                    <div className="overflow-x-auto rounded-lg bg-muted p-2 text-xs">
                      <SyntaxHighlighter language="json" style={gruvboxDark}>
                        {JSON.stringify(selectedAgent, null, 2)}
                      </SyntaxHighlighter>
                    </div>
                  </AccordionContent>
                </AccordionItem>
              </Accordion>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </>
  )
}
