import type { SchemaShip } from "@/api/types"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { ArrowRightLeft, MapPin, Rocket, ShipWheel } from "lucide-react"

interface ShipCardProps {
  ship: SchemaShip

  onOpen?: () => void
  onSwitch?: () => void
}

export function ShipCard({ ship, onOpen, onSwitch }: ShipCardProps) {
  return (
    <div
      className={`rounded-lg border p-3 transition-colors ${
        ship.active
          ? "border-primary/40 bg-primary/5"
          : "hover:border-primary/20"
      }`}
    >
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <Rocket className="size-4 shrink-0 text-primary" />

            <div className="min-w-0 flex-1">
              <span className="block font-medium break-all">{ship.name}</span>
            </div>

            {ship.active && (
              <Badge variant="default" className="shrink-0 uppercase">
                Active
              </Badge>
            )}
          </div>

          <div className="mt-2 space-y-1 text-sm text-muted-foreground">
            <div className="flex items-center gap-2">
              <ShipWheel className="size-3.5" />
              <span>{ship.type}</span>
            </div>

            <div className="flex items-center gap-2">
              <MapPin className="size-3.5" />
              <span>
                [{ship.system_x}; {ship.system_y}] {ship.location}{" "}
                {ship.location !== "none" && ship.location_id}
              </span>
            </div>
          </div>
        </div>

        <div className="flex flex-col gap-2">
          <Button size="sm" variant="outline" onClick={onOpen}>
            Details
          </Button>

          {!ship.active && (
            <Button size="sm" onClick={onSwitch}>
              <ArrowRightLeft className="mr-2 size-3.5" />
              Switch
            </Button>
          )}
        </div>
      </div>
    </div>
  )
}
