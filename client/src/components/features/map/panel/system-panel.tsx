import type { SystemExtended } from "@/api/types"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { ArrowRight, Orbit, SquareArrowOutUpRight, X } from "lucide-react"
import { AgentCard } from "../../agents/agent-card"
import { Separator } from "@/components/ui/separator"
import type { AgentExtended } from "@/api/types"

interface SystemPanelProps {
  currentAgent: AgentExtended
  system: SystemExtended
  onClose: () => void
  onCenterCamera: () => void
  onSystemOpen: () => void
  onWarp: () => void
}

export function SystemPanel({
  currentAgent,
  system,
  onClose,
  onCenterCamera,
  onSystemOpen,
  onWarp,
}: SystemPanelProps) {
  return (
    <>
      <div className="flex items-center justify-between border-b border-border p-5">
        <div>
          <h2 className="text-xl font-bold">{system.system.name}</h2>

          <p className="text-sm text-muted-foreground">Star System</p>
        </div>

        <Button variant="ghost" size="icon" onClick={onClose}>
          <X />
        </Button>
      </div>

      <div className="space-y-6 p-5">
        <Card className="p-4">
          <h3 className="mb-3 font-semibold">Information</h3>

          <div className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Type</span>

              <Badge>system.type</Badge>
            </div>

            <div className="flex justify-between">
              <span className="text-muted-foreground">Coordinates</span>

              <span>
                {system.system.x}, {system.system.y}
              </span>
            </div>

            <div className="flex justify-between">
              <span className="text-muted-foreground">Waypoints</span>

              <span>{system.system.waypoints.length}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Planets</span>

              <span>{system.system.planets.length}</span>
            </div>
          </div>
        </Card>

        <Card className="p-4">
          <h3 className="mb-3 font-semibold">Agents in this system</h3>

          <div className="flex flex-col gap-2">
            {system.agents.length > 0 ? (
              system.agents.map((a) => <AgentCard key={a.agent.id} agent={a} />)
            ) : (
              <div className="flex min-h-20 place-items-center justify-center border border-dashed text-sm text-muted-foreground">
                No agents in this system
              </div>
            )}
          </div>
        </Card>

        <Card className="p-4">
          <h3 className="mb-3 font-semibold">Available Actions</h3>

          <div className="flex flex-col gap-2">
            <Button
              className="justify-start"
              hidden={
                currentAgent.ship.system_x === system.system.x &&
                currentAgent.ship.system_y === system.system.y
                  ? true
                  : false
              }
              onClick={onWarp}
            >
              <ArrowRight className="mr-2 size-4" />
              Warp
            </Button>

            <Button
              variant="secondary"
              className="justify-start"
              onClick={onCenterCamera}
            >
              <Orbit className="mr-2 size-4" />
              Center Camera
            </Button>

            <Button
              variant="secondary"
              className="justify-start"
              onClick={onSystemOpen}
            >
              <SquareArrowOutUpRight className="mr-2 size-4" />
              Open
            </Button>
          </div>
        </Card>

        <Separator />

        <div>
          <h3 className="mb-3 font-semibold">Description</h3>

          <p className="text-sm text-muted-foreground">
            Additional information about the selected system will appear here
          </p>
        </div>
      </div>
    </>
  )
}
