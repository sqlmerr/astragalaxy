import type { SchemaFullInventory } from "@/api/types"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { ItemCard } from "./item-card"
import { ResourceCard } from "./resource-card"
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"
import { Json } from "@/components/ui/json"
import { CraftActionDialog } from "./craft-action-dialog"
import { useState } from "react"

interface InventoryModalProps {
  inventory: SchemaFullInventory | null
  onClose: () => void
  agentId: string
}

export function InventoryModal({
  inventory,
  onClose,
  agentId,
}: InventoryModalProps) {
  const [craftOpened, setCraftOpened] = useState(false)

  return (
    <>
      <Dialog
        open={inventory !== null}
        onOpenChange={(open) => {
          if (!open) onClose()
        }}
      >
        <DialogContent className="max-h-[90vh] overflow-y-auto">
          <DialogHeader className="flex-row items-center justify-between gap-2 pr-10">
            <DialogTitle className="flex items-center gap-2 pr-10">
              Inventory
            </DialogTitle>
            <Button
              type="button"
              size="sm"
              onClick={() => setCraftOpened(true)}
            >
              Craft
            </Button>
          </DialogHeader>
          {inventory && (
            <div className="space-y-6">
              <Card className="p-4">
                <h3 className="mb-3 font-semibold">Information</h3>

                <div className="grid grid-cols-[80px_1fr] gap-x-3 gap-y-2 text-sm">
                  <span className="text-muted-foreground">ID</span>
                  <code className="font-mono break-all">
                    {inventory.inventory.id}
                  </code>

                  <span className="text-muted-foreground">Items</span>
                  <code className="font-mono break-all">
                    {inventory.items.length}
                  </code>

                  <span className="text-muted-foreground">Resources</span>
                  <code className="font-mono break-all">
                    {inventory.resources.length}
                  </code>

                  <span className="text-muted-foreground">Max Item Slots</span>
                  <code className="font-mono break-all">
                    {inventory.inventory.max_item_slots}
                  </code>

                  <span className="text-muted-foreground">
                    Max Resource Volume
                  </span>
                  <code className="font-mono break-all">
                    {inventory.inventory.max_resource_volume}
                  </code>
                </div>
              </Card>
              <Card className="p-4">
                <h3 className="mb-3 font-semibold">Items</h3>

                <div className="flex flex-col gap-2">
                  {inventory.items.length > 0 ? (
                    inventory.items.map((i) => (
                      <ItemCard
                        key={i.id}
                        item={i}
                        inventoryAgentId={agentId}
                      />
                    ))
                  ) : (
                    <div className="flex min-h-20 place-items-center justify-center border border-dashed text-sm text-muted-foreground">
                      Inventory does not contain any items
                    </div>
                  )}
                </div>
              </Card>
              <Card className="p-4">
                <h3 className="mb-3 font-semibold">Resources</h3>

                <div className="flex flex-col gap-2">
                  {inventory.resources.filter((v) => v.amount > 0).length >
                  0 ? (
                    inventory.resources
                      .filter((v) => v.amount > 0)
                      .map((r) => (
                        <ResourceCard key={r.resource_type} resource={r} />
                      ))
                  ) : (
                    <div className="flex min-h-20 place-items-center justify-center border border-dashed text-sm text-muted-foreground">
                      Inventory does not contain any resources
                    </div>
                  )}
                </div>
              </Card>

              <Accordion>
                <AccordionItem value="json">
                  <AccordionTrigger>JSON</AccordionTrigger>
                  <AccordionContent>
                    <div className="rounded-lg bg-muted p-2 text-xs">
                      <Json data={inventory} />
                    </div>
                  </AccordionContent>
                </AccordionItem>
              </Accordion>
            </div>
          )}
        </DialogContent>
      </Dialog>
      <CraftActionDialog
        agentId={agentId}
        inventoryId={inventory?.inventory.id || null}
        open={craftOpened}
        onClose={() => setCraftOpened(false)}
      />
    </>
  )
}
