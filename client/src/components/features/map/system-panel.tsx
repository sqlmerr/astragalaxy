import type { SchemaSystem, SystemExtended } from "@/api/types"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Separator } from "@/components/ui/separator"
import { cn } from "@/lib/utils"
import {
  ArrowRight,
  Hammer,
  MapPin,
  Orbit,
  ShoppingCart,
  X,
} from "lucide-react"
import { AgentCard } from "../agents/agent-card"

interface SystemPanelProps {
  system: SystemExtended | null
  onClose: () => void
  onCenterCamera: () => void
}

export function SystemPanel({
  system,
  onClose,
  onCenterCamera,
}: SystemPanelProps) {
  return (
    <aside
      className={`fixed top-0 right-0 z-50 h-screen w-96 transform border-l border-border bg-card/95 backdrop-blur-md transition-transform duration-300 ${
        system ? "translate-x-0" : "translate-x-full"
      }`}
    >
      {system && (
        <ScrollArea className="h-full">
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
                  system.agents.map((a) => <AgentCard key={a.id} agent={a} />)
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
                <Button className="justify-start">
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
              </div>
            </Card>

            <Separator />

            <div>
              <h3 className="mb-3 font-semibold">Description</h3>

              <p className="text-sm text-muted-foreground">
                Additional information about the selected system will appear
                here
              </p>
            </div>
          </div>
        </ScrollArea>
      )}
    </aside>
  )
}
