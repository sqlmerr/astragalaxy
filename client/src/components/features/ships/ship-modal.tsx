import type { SchemaShip } from "@/api/types"
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
import { Separator } from "@/components/ui/separator"

interface ShipModalProps {
  ship: SchemaShip | null
  onClose: () => void
}

export function ShipModal({ ship, onClose }: ShipModalProps) {
  return (
    <Dialog
      open={ship !== null}
      onOpenChange={(open) => {
        if (!open) onClose()
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            {ship?.name} {ship?.active && <Badge>ACTIVE</Badge>}
          </DialogTitle>
        </DialogHeader>
        {ship && (
          <div className="space-y-6">
            <Card className="p-4">
              <h3 className="mb-3 font-semibold">Information</h3>

              <div className="grid grid-cols-[80px_1fr] gap-x-3 gap-y-2 text-sm">
                <span className="text-muted-foreground">ID</span>
                <code className="font-mono break-all">{ship.id}</code>

                <span className="text-muted-foreground">Agent ID</span>
                <code className="font-mono break-all">{ship.agent_id}</code>

                <span className="text-muted-foreground">Type</span>
                <Badge>{ship.type}</Badge>

                <span className="text-muted-foreground">Location</span>
                <code className="font-mono break-all">
                  {ship.location === "NONE"
                    ? ship.location
                    : `${ship.location} #${ship.location_id}`}
                </code>

                <span className="text-muted-foreground">System</span>
                <code className="font-mono break-all">
                  <span className="text-muted-foreground">x=</span>
                  {ship.system_x}{" "}
                  <span className="text-muted-foreground">y=</span>
                  {ship.system_y}
                </code>

                <span className="text-muted-foreground">Status</span>
                <Badge>{ship.status}</Badge>
              </div>
            </Card>

            <Separator />

            <Accordion>
              <AccordionItem value="json">
                <AccordionTrigger>JSON</AccordionTrigger>
                <AccordionContent>
                  <div className="rounded-lg bg-muted p-2 text-xs">
                    <Json data={ship} />
                  </div>
                </AccordionContent>
              </AccordionItem>
            </Accordion>
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}
