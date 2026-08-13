import type { SchemaWaypoint, SystemExtended, AgentExtended } from "@/api/types"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { ArrowRight, Pickaxe, X } from "lucide-react"
import { Separator } from "@/components/ui/separator"
import { WAYPOINT_PARAMS } from "../constants"
import { AgentCard } from "../../agents/agent-card"
import { shipLocationIs } from "@/api/utils"
import { ResourceDepositCard } from "../../inventory/resource-deposit-card"

interface WaypointPanelProps {
  currentAgent: AgentExtended
  system: SystemExtended
  waypoint: SchemaWaypoint
  onClose: () => void
  onNavigate: (w: SchemaWaypoint) => void
  onAsteroidMine: () => void
}

export function WaypointPanel({
  currentAgent,
  system,
  waypoint,
  onClose,
  onNavigate,
  onAsteroidMine,
}: WaypointPanelProps) {
  const thisWaypointAgents = system.agents.filter(
    (a) => a.ship.location === "waypoint" && a.ship.location_id === waypoint.id
  )
  const params = WAYPOINT_PARAMS[waypoint.type]

  return (
    <>
      <div className="flex items-center justify-between border-b border-border p-5">
        <div>
          <h2 className="text-xl font-bold">{`${params.name}`}</h2>

          <p className="text-sm text-muted-foreground">Waypoint</p>
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

              <Badge>{waypoint.type}</Badge>
            </div>

            <div className="flex justify-between">
              <span className="text-muted-foreground">ID</span>

              <span>{waypoint.id}</span>
            </div>
          </div>
        </Card>

        <Card className="p-4">
          <h3 className="mb-3 font-semibold">Agents on this waypoint</h3>

          <div className="flex flex-col gap-2">
            {thisWaypointAgents.length > 0 ? (
              thisWaypointAgents.map((a) => (
                <AgentCard key={a.agent.id} agent={a} />
              ))
            ) : (
              <div className="flex min-h-20 place-items-center justify-center border border-dashed text-sm text-muted-foreground">
                No agents on this waypoint
              </div>
            )}
          </div>
        </Card>

        <Card className="p-4">
          <h3 className="mb-3 font-semibold">Available Actions</h3>

          <div className="flex flex-col gap-2">
            {!shipLocationIs(currentAgent.ship, {
              locationType: "waypoint",
              locationId: waypoint.id,
              systemX: currentAgent.ship.system_x,
              systemY: currentAgent.ship.system_y,
            }) ? (
              <>
                <Button
                  className="justify-start"
                  variant={
                    currentAgent.ship.system_x === system.system.x &&
                    currentAgent.ship.system_y === system.system.y
                      ? "default"
                      : "destructive"
                  }
                  onClick={() => onNavigate(waypoint)}
                >
                  <ArrowRight className="mr-2 size-4" />
                  Navigate
                </Button>
              </>
            ) : (
              <>
                {waypoint.type === "asteroid" ? (
                  <Button
                    className="justify-start"
                    variant="secondary"
                    onClick={onAsteroidMine}
                  >
                    <Pickaxe className="mr-2 size-4" />
                    Mine
                  </Button>
                ) : null}
              </>
            )}
          </div>
        </Card>

        <Separator />

        {waypoint.type === "asteroid" && !!waypoint.asteroid ? (
          <Card className="p-4">
            <h3 className="mb-3 font-semibold">Resource deposit</h3>
            <ResourceDepositCard deposit={waypoint.asteroid.deposit} />
          </Card>
        ) : null}

        {waypoint.type === "station" && !!waypoint.station ? (
          <Card className="p-4">
            <h3 className="mb-3 font-semibold">Facilities</h3>
            <ul className="space-y-2 text-sm">
              {waypoint.station.facilities.map((facility) => (
                <li key={facility}>
                  <Badge variant="secondary">{facility}</Badge>
                </li>
              ))}
            </ul>
          </Card>
        ) : null}

        {/* <div>
          <h3 className="mb-3 font-semibold">Description</h3>

          <p className="text-sm text-muted-foreground">
            Additional information about the selected system will appear here
          </p>
        </div> */}
      </div>
    </>
  )
}
