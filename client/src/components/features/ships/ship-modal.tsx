import { useRenameShipMutation, useShipInventoryQuery } from "@/api/hooks"
import { queryKeys } from "@/api/query-keys"
import type { SchemaShip } from "@/api/types"
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Json } from "@/components/ui/json"
import { toast } from "@/components/ui/toast"
import { useErrorHandler } from "@/errors/utils"
import { useQueryClient } from "@tanstack/react-query"
import { Check, Pencil, X } from "lucide-react"
import { useEffect, useState } from "react"
import { InventoryModal } from "../inventory/inventory-modal"
import { useShipModulesQuery } from "@/api/hooks/use-ships-query"

interface ShipModalProps {
  ship: SchemaShip | null
  onClose: () => void
}

export function ShipModal({ ship, onClose }: ShipModalProps) {
  const queryClient = useQueryClient()
  const [isRenaming, setIsRenaming] = useState(false)
  const [name, setName] = useState("")
  const renameMutation = useRenameShipMutation()
  const errorHandler = useErrorHandler()
  const [invModalOpen, setInvModalOpen] = useState(false)
  const {
    data: inventory,
    isPending: isInventoryPending,
    isError: isInventoryError,
  } = useShipInventoryQuery(
    ship ? ship.agent_id : "",
    ship ? ship.id : "",
    ship !== null
  )

  const {
    data: modules,
    isPending: isModulesPending,
    isError: isModulesError,
  } = useShipModulesQuery(
    ship ? ship.agent_id : "",
    ship ? ship.id : "",
    ship !== null
  )

  useEffect(() => {
    if (ship) {
      setName(ship.name)
      setIsRenaming(false)
    }
  }, [ship])

  async function renameShip() {
    if (!ship) return

    if (name === ship.name) {
      setIsRenaming(false)
      return
    }

    await renameMutation.mutateAsync(
      {
        agentId: ship.agent_id,
        id: ship.id,
        body: {
          name,
        },
      },
      {
        onSuccess: (renamedShip) => {
          queryClient.invalidateQueries({
            queryKey: queryKeys.ships.my(ship.agent_id),
          })

          toast.add({
            type: "success",
            title: "Success",
            description: "Successfully renamed ship",
          })
        },
        onError: (err) => {
          errorHandler(err, "Failed to rename ship")
        },
      }
    )

    setIsRenaming(false)
  }

  if (
    !inventory ||
    isInventoryPending ||
    isInventoryError ||
    !modules ||
    isModulesError ||
    isModulesPending
  ) {
    return null
  }

  return (
    <>
      <Dialog
        open={ship !== null}
        onOpenChange={(open) => {
          if (!open) onClose()
        }}
      >
        <DialogContent className="max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 pr-10">
              {isRenaming ? (
                <>
                  <Input
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    className="h-8"
                  />

                  <Button
                    size="icon"
                    variant="ghost"
                    disabled={renameMutation.isPending}
                    onClick={renameShip}
                  >
                    <Check className="size-4" />
                  </Button>

                  <Button
                    size="icon"
                    variant="ghost"
                    onClick={() => {
                      setName(ship?.name ?? "")
                      setIsRenaming(false)
                    }}
                  >
                    <X className="size-4" />
                  </Button>
                </>
              ) : (
                <>
                  <span>{ship?.name}</span>

                  {ship?.active && (
                    <>
                      <Badge>ACTIVE</Badge>

                      <Button
                        size="icon"
                        variant="ghost"
                        className="size-7"
                        onClick={() => setIsRenaming(true)}
                      >
                        <Pencil className="size-4" />
                      </Button>
                    </>
                  )}
                </>
              )}
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
                    {ship.location === "none"
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

              <Card className="p-4">
                <h3 className="mb-3 font-semibold">Inventory</h3>

                <div className="flex items-center justify-between gap-4 text-sm">
                  <span className="text-muted-foreground">
                    {inventory.items.length} items, {inventory.resources.length}{" "}
                    resources
                  </span>
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => setInvModalOpen(true)}
                  >
                    Open inventory
                  </Button>
                </div>
              </Card>

              <Card className="p-4">
                <h3 className="mb-3 font-semibold">Modules</h3>
                {modules.data.length > 0 ? (
                  <ul className="space-y-2 text-sm">
                    {modules.data.map((module) => (
                      <li key={module}>
                        <Badge variant="secondary">{module}</Badge>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p className="text-sm text-muted-foreground">
                    This ship has no modules.
                  </p>
                )}
              </Card>

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
      <InventoryModal
        inventory={invModalOpen ? inventory : null}
        onClose={() => setInvModalOpen(false)}
        agentId={ship?.agent_id || ""}
      />
    </>
  )
}
