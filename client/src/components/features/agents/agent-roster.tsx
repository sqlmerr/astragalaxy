import { Crosshair, UsersRound } from "lucide-react"

import { useMyAgentsQuery } from "@/api/hooks"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Spinner } from "@/components/ui/spinner"
import { AgentCard } from "@/components/features/agents/agent-card"
import { useState } from "react"
import type { SchemaAgent } from "@/api/types"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"
import { Separator } from "@/components/ui/separator"

export function AgentRoster() {
  const { data, isPending, isError } = useMyAgentsQuery()
  const [selectedAgent, setSelectedAgent] = useState<SchemaAgent | null>(null)
  const agents = data?.data ?? []

  return (
    <>
      <aside
        className="fixed top-4 left-4 z-20 w-[calc(100vw-2rem)] max-w-sm lg:top-6 lg:left-6"
        aria-label="Player agents"
      >
        <Card className="overflow-hidden">
          <CardHeader className="flex-row items-center justify-between border-b border-border">
            <div className="flex items-center gap-2">
              <UsersRound className="size-4 text-primary" aria-hidden="true" />
              <CardTitle className="text-xs tracking-[0.16em] uppercase">
                Your agents
              </CardTitle>
            </div>
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
                  <AgentCard
                    key={agent.id}
                    agent={agent}
                    setSelectedAgent={setSelectedAgent}
                  />
                ))}
              </ScrollArea>
            )}
          </CardContent>
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
              <div className="space-y-3 rounded-lg p-4">
                <div className="grid grid-cols-[120px_1fr] gap-2 text-sm">
                  <span className="text-muted-foreground">ID</span>
                  <code className="font-mono break-all">
                    {selectedAgent.id}
                  </code>

                  <span className="text-muted-foreground">User ID</span>
                  <code className="font-mono break-all">
                    {selectedAgent.user_id}
                  </code>
                </div>
              </div>
              <Separator />
              <Accordion>
                <AccordionItem value="json">
                  <AccordionTrigger>JSON</AccordionTrigger>

                  <AccordionContent>
                    <pre className="overflow-x-auto rounded-lg bg-muted p-4 text-xs">
                      <code>{JSON.stringify(selectedAgent, null, 2)}</code>
                    </pre>
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
