import {
  useTransferItemsMutation,
  useTransferResourcesMutation,
} from "@/api/hooks"
import { queryKeys } from "@/api/query-keys"
import type { SchemaFullInventory, SchemaResource } from "@/api/types"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { toast } from "@/components/ui/toast"
import { useErrorHandler } from "@/errors/utils"
import { cn } from "@/lib/utils"
import { useQueryClient } from "@tanstack/react-query"
import { MinusIcon, PlusIcon } from "lucide-react"
import { useState } from "react"
import type { SyntheticEvent } from "react"

interface TransferActionDialogProps {
  agentId: string
  inventory: SchemaFullInventory | null
  open: boolean
  onClose: () => void
}

interface ResourceTransferRowProps {
  resource: SchemaResource
  amount: number | null
  onToggle: () => void
  onAmountChange: (value: string) => void
  onStep: (delta: number) => void
  onMax: () => void
}

function ResourceTransferRow({
  resource,
  amount,
  onToggle,
  onAmountChange,
  onStep,
  onMax,
}: ResourceTransferRowProps) {
  const checked = amount !== null
  const inputId = `resource-${resource.resource_type}`

  return (
    <div
      data-checked={checked ? true : undefined}
      className={cn(
        "flex items-center gap-3 rounded-lg border p-3 transition-colors",
        checked
          ? "border-primary/30 bg-primary/5"
          : "border-transparent hover:bg-muted/50"
      )}
    >
      <Checkbox
        id={inputId}
        checked={checked}
        onCheckedChange={onToggle}
        aria-label={resource.resource_type}
      />
      <div className="min-w-0 flex-1">
        <Label
          htmlFor={inputId}
          className="block cursor-pointer truncate text-sm font-semibold"
        >
          {resource.resource_type}
        </Label>
        <p className="mt-1 truncate text-[11px] text-muted-foreground">
          Available: {resource.amount}
        </p>
      </div>
      <div className="flex items-center gap-1">
        <Button
          type="button"
          variant="outline"
          size="icon-sm"
          disabled={!checked}
          onClick={() => onStep(-1)}
          aria-label={`Decrease ${resource.resource_type}`}
        >
          <MinusIcon />
        </Button>
        <Input
          className="w-20 text-center"
          type="number"
          inputMode="numeric"
          min={1}
          max={resource.amount}
          step={1}
          value={amount ?? ""}
          disabled={!checked}
          onChange={(e) => onAmountChange(e.target.value)}
          aria-label={`${resource.resource_type} amount`}
        />
        <Button
          type="button"
          variant="outline"
          size="icon-sm"
          disabled={!checked}
          onClick={() => onStep(1)}
          aria-label={`Increase ${resource.resource_type}`}
        >
          <PlusIcon />
        </Button>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          disabled={!checked}
          onClick={onMax}
        >
          Max
        </Button>
      </div>
    </div>
  )
}

export function TransferActionDialog({
  open,
  agentId,
  inventory,
  onClose,
}: TransferActionDialogProps) {
  const errorHandler = useErrorHandler()
  const queryClient = useQueryClient()
  const transferItemsMutation = useTransferItemsMutation()
  const transferResourcesMutation = useTransferResourcesMutation()

  const [receiverId, setReceiverId] = useState("")
  const [selectedItems, setSelectedItems] = useState<string[]>([])
  const [resources, setResources] = useState<Record<string, number>>({})

  const availableResources =
    inventory?.resources.filter((r) => r.amount > 0) ?? []
  const selectedResourcesCount = Object.values(resources).filter(
    (amount) => amount > 0
  ).length
  const sameInventory =
    receiverId !== "" && receiverId === inventory?.inventory.id

  function reset() {
    setReceiverId("")
    setSelectedItems([])
    setResources({})
  }

  function toggleResource(resource: SchemaResource) {
    setResources((prev) => {
      if (resource.resource_type in prev) {
        const next = { ...prev }
        delete next[resource.resource_type]
        return next
      }
      return { ...prev, [resource.resource_type]: resource.amount }
    })
  }

  function setResourceAmount(raw: string, max: number, resourceType: string) {
    const parsed = Number.parseInt(raw, 10)
    const value = Number.isFinite(parsed)
      ? Math.max(0, Math.min(parsed, max))
      : 0
    setResources((prev) => ({ ...prev, [resourceType]: value }))
  }

  function stepResource(resourceType: string, delta: number, max: number) {
    setResources((prev) => {
      const current = prev[resourceType] ?? 0
      const value = Math.max(1, Math.min(current + delta, max))
      return { ...prev, [resourceType]: value }
    })
  }

  function maxResource(resourceType: string, max: number) {
    setResources((prev) => ({ ...prev, [resourceType]: max }))
  }

  function selectAllResources() {
    setResources(
      Object.fromEntries(
        availableResources.map((r) => [r.resource_type, r.amount])
      )
    )
  }

  async function transferItems(event: SyntheticEvent) {
    event.preventDefault()
    if (!inventory || !receiverId || selectedItems.length === 0) return

    await transferItemsMutation.mutateAsync(
      {
        agentId,
        body: {
          from_inventory_id: inventory.inventory.id,
          to_inventory_id: receiverId,
          items: selectedItems,
        },
      },
      {
        onError(err) {
          errorHandler(err, "Failed to transfer items")
        },
        onSuccess() {
          queryClient.invalidateQueries({ queryKey: queryKeys.inventories.all })
          toast.add({
            type: "success",
            title: "Success",
            description: "Items transferred",
          })
          reset()
          onClose()
        },
      }
    )
  }

  async function transferResources(event: SyntheticEvent) {
    event.preventDefault()
    if (!inventory || !receiverId || selectedResourcesCount === 0) return

    const payload = Object.fromEntries(
      Object.entries(resources).filter(([, amount]) => amount > 0)
    )

    await transferResourcesMutation.mutateAsync(
      {
        agentId,
        body: {
          from_inventory_id: inventory.inventory.id,
          to_inventory_id: receiverId,
          resources: payload,
        },
      },
      {
        onError(err) {
          errorHandler(err, "Failed to transfer resources")
        },
        onSuccess() {
          queryClient.invalidateQueries({ queryKey: queryKeys.inventories.all })
          toast.add({
            type: "success",
            title: "Success",
            description: "Resources transferred",
          })
          reset()
          onClose()
        },
      }
    )
  }

  if (!inventory) return

  const itemsSelectItems = inventory.items.map((i) => ({
    label: i.item_type,
    value: i.id,
  }))

  return (
    <Dialog
      open={open}
      onOpenChange={(o) => {
        if (!o) {
          reset()
          onClose()
        }
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Transfer Request</DialogTitle>
        </DialogHeader>
        <Tabs defaultValue="items">
          <TabsList>
            <TabsTrigger value="items">Items</TabsTrigger>
            <TabsTrigger value="resources">Resources</TabsTrigger>
          </TabsList>
          <TabsContent value="items">
            <form className="space-y-4" onSubmit={transferItems}>
              <div className="grid gap-2 text-sm font-medium">
                <Label htmlFor="items-select">Items</Label>
                <Select
                  items={itemsSelectItems}
                  value={selectedItems}
                  multiple
                  onValueChange={setSelectedItems}
                >
                  <SelectTrigger id="items-select" className="w-full">
                    <SelectValue placeholder="Item" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      {itemsSelectItems.map((item) => (
                        <SelectItem key={item.value} value={item.value}>
                          {item.label}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>

              <div className="grid gap-2 text-sm font-medium">
                <Label htmlFor="items-receiver-id">Receiver Inventory ID</Label>
                <Input
                  id="items-receiver-id"
                  value={receiverId}
                  onChange={(e) => setReceiverId(e.target.value)}
                  placeholder="Enter inventory id"
                  required
                />
                {sameInventory && (
                  <p className="text-xs text-destructive">
                    Source and destination inventories must differ.
                  </p>
                )}
              </div>

              <Button
                className="w-full"
                type="submit"
                disabled={
                  transferItemsMutation.isPending ||
                  selectedItems.length === 0 ||
                  receiverId === "" ||
                  sameInventory
                }
              >
                Transfer
              </Button>
            </form>
          </TabsContent>
          <TabsContent value="resources">
            <form className="space-y-4" onSubmit={transferResources}>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Resources</span>
                <div className="flex items-center gap-1">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={selectAllResources}
                    disabled={availableResources.length === 0}
                  >
                    Select all
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => setResources({})}
                    disabled={selectedResourcesCount === 0}
                  >
                    Clear
                  </Button>
                </div>
              </div>

              {availableResources.length > 0 ? (
                <div className="flex max-h-64 flex-col gap-1 overflow-y-auto pr-1">
                  {availableResources.map((resource) => (
                    <ResourceTransferRow
                      key={resource.resource_type}
                      resource={resource}
                      amount={resources[resource.resource_type] ?? null}
                      onToggle={() => toggleResource(resource)}
                      onAmountChange={(value) =>
                        setResourceAmount(
                          value,
                          resource.amount,
                          resource.resource_type
                        )
                      }
                      onStep={(delta) =>
                        stepResource(
                          resource.resource_type,
                          delta,
                          resource.amount
                        )
                      }
                      onMax={() =>
                        maxResource(resource.resource_type, resource.amount)
                      }
                    />
                  ))}
                </div>
              ) : (
                <div className="flex min-h-20 place-items-center justify-center border border-dashed text-sm text-muted-foreground">
                  Inventory does not contain any resources
                </div>
              )}

              <div className="grid gap-2 text-sm font-medium">
                <Label htmlFor="resources-receiver-id">
                  Receiver Inventory ID
                </Label>
                <Input
                  id="resources-receiver-id"
                  value={receiverId}
                  onChange={(e) => setReceiverId(e.target.value)}
                  placeholder="Enter inventory id"
                  required
                />
                {sameInventory && (
                  <p className="text-xs text-destructive">
                    Source and destination inventories must differ.
                  </p>
                )}
              </div>

              <Button
                className="w-full"
                type="submit"
                disabled={
                  transferResourcesMutation.isPending ||
                  selectedResourcesCount === 0 ||
                  receiverId === "" ||
                  sameInventory
                }
              >
                Transfer
              </Button>
            </form>
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  )
}
