import type { AgentExtended, SchemaShip } from "@/api/types"
import { useNow } from "@/components/time-provider"
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"
import { Badge } from "@/components/ui/badge"
import { Card } from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Json } from "@/components/ui/json"
import { Progress } from "@/components/ui/progress"
import { Separator } from "@/components/ui/separator"
import { ShipCard } from "../ships/ship-card"
import { useState } from "react"
import { ShipModal } from "../ships/ship-modal"

interface AgentModalProps {
  agent: AgentExtended | null
  setAgent: (a: AgentExtended | null) => void
}

export function AgentModal({ agent, setAgent }: AgentModalProps) {
  const now = useNow()
  let cooldown // in seconds
  let progress
  if (agent) {
    const cooldownExpiresAt =
      new Date(agent.cooldown.set_at).getTime() +
      agent.cooldown.duration_seconds * 1000
    cooldown = Math.max(0, cooldownExpiresAt - now) / 1000
    progress =
      ((agent.cooldown.duration_seconds - cooldown) /
        agent.cooldown.duration_seconds) *
      100
  } else {
    cooldown = 0
    progress = 100
  }

  const [selectedShip, setSelectedShip] = useState<SchemaShip | null>(null)

  return (
    <>
      <Dialog
        open={agent !== null}
        onOpenChange={(open) => {
          if (!open) setAgent(null)
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{agent?.agent.username}</DialogTitle>
          </DialogHeader>
          {agent && (
            <div className="space-y-6">
              <Card className="p-4">
                <h3 className="mb-3 font-semibold">Information</h3>

                <div className="grid grid-cols-[80px_1fr] gap-x-3 gap-y-2 text-sm">
                  <span className="text-muted-foreground">ID</span>
                  <code className="font-mono break-all">{agent.agent.id}</code>

                  <span className="text-muted-foreground">User ID</span>

                  <code className="font-mono break-all">
                    {agent.agent.user_id}
                  </code>
                </div>
              </Card>
              <Separator />

              <Card className="p-4">
                <h3 className="mb-3 font-semibold">Cooldown</h3>

                {cooldown > 0 ? (
                  <div className="grid grid-cols-[100px_1fr] gap-x-3 gap-y-2 text-sm">
                    <span className="text-muted-foreground">Status</span>
                    <Badge variant="secondary">On cooldown</Badge>

                    <span className="text-muted-foreground">Action</span>
                    <span className="font-medium capitalize">
                      {agent.cooldown.action.replaceAll("_", " ")}
                    </span>

                    <Progress value={progress} />
                    <div className="mt-1 flex justify-between text-xs text-muted-foreground">
                      <span>{cooldown}s remaining</span>
                      <span>{Math.round(progress)}%</span>
                    </div>

                    <span className="text-muted-foreground">Started</span>
                    <code className="font-mono">
                      {new Date(agent.cooldown.set_at).toLocaleString()}
                    </code>
                  </div>
                ) : (
                  <div className="flex items-center justify-between rounded-lg border border-green-500/20 bg-green-500/10 px-4 py-3">
                    <span className="text-sm font-medium">
                      Ready for actions
                    </span>
                    <Badge className="bg-green-600 hover:bg-green-600">
                      Ready
                    </Badge>
                  </div>
                )}
              </Card>

              <Separator />

              <Card className="p-4">
                <h3 className="mb-3 font-semibold">Ships</h3>
                {agent.ships.map((s) => (
                  <ShipCard
                    key={s.id}
                    ship={s}
                    onOpen={() => setSelectedShip(s)}
                  />
                ))}
              </Card>

              <Separator />

              <Accordion>
                <AccordionItem value="json">
                  <AccordionTrigger>JSON</AccordionTrigger>
                  <AccordionContent>
                    <div className="rounded-lg bg-muted p-2 text-xs">
                      <Json data={agent.agent} />
                    </div>
                  </AccordionContent>
                </AccordionItem>
              </Accordion>
            </div>
          )}
        </DialogContent>
      </Dialog>
      <ShipModal ship={selectedShip} onClose={() => setSelectedShip(null)} />
    </>
  )
}
