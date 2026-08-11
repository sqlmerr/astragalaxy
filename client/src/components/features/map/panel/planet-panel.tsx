import type { SchemaPlanet, SystemExtended, AgentExtended  } from "@/api/types"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { ArrowRight, Pickaxe, X } from "lucide-react"
import { Separator } from "@/components/ui/separator"
import { AgentCard } from "../../agents/agent-card"
import { ResourceDepositCard } from "../../inventory/resource-deposit-card"
import { shipLocationIs } from "@/api/utils"

interface PlanetPanelProps {
  currentAgent: AgentExtended
  system: SystemExtended
  planet: SchemaPlanet
  onClose: () => void
  onNavigate: (p: SchemaPlanet) => void
  onPlanetMine: () => void
}

export function PlanetPanel({
  currentAgent,
  system,
  planet,
  onClose,
  onNavigate,
  onPlanetMine,
}: PlanetPanelProps) {
  const thisPlanetAgents = system.agents.filter(
    (a) => a.ship.location === "planet" && a.ship.location_id === planet.orbit
  )

  return (
    <>
      <div className="flex items-center justify-between border-b border-border p-5">
        <div>
          <h2 className="text-xl font-bold">{planet.name}</h2>

          <p className="text-sm text-muted-foreground">Planet</p>
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

              <Badge>{planet.type}</Badge>
            </div>

            <div className="flex justify-between">
              <span className="text-muted-foreground">Orbit</span>

              <span>{planet.orbit}</span>
            </div>
          </div>
        </Card>
        <Card className="p-4">
          <h3 className="mb-3 font-semibold">Agents on this planet</h3>

          <div className="flex flex-col gap-2">
            {thisPlanetAgents.length > 0 ? (
              thisPlanetAgents.map((a) => (
                <AgentCard key={a.agent.id} agent={a} />
              ))
            ) : (
              <div className="flex min-h-20 place-items-center justify-center border border-dashed text-sm text-muted-foreground">
                No agents on this planet
              </div>
            )}
          </div>
        </Card>
        <Card className="p-4">
          <h3 className="mb-3 font-semibold">Available Actions</h3>

          <div className="flex flex-col gap-2">
            {!shipLocationIs(currentAgent.ship, {
              locationType: "planet",
              locationId: planet.orbit,
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
                  onClick={() => onNavigate(planet)}
                >
                  <ArrowRight className="mr-2 size-4" />
                  Navigate
                </Button>
              </>
            ) : (
              <>
                <Button
                  className="justify-start"
                  variant="secondary"
                  onClick={onPlanetMine}
                >
                  <Pickaxe className="mr-2 size-4" />
                  Mine
                </Button>
              </>
            )}
          </div>
        </Card>
        <Separator />

        {planet.deposits.length > 0 ? (
          <Card className="p-4">
            <h3 className="mb-3 font-semibold">Resource deposits</h3>
            {planet.deposits.map((d) => (
              <ResourceDepositCard key={d.resource} deposit={d} />
            ))}
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
