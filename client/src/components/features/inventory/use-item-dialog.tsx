import { useUseItemMutation } from "@/api/hooks/use-inventory-mutations"
import { queryKeys } from "@/api/query-keys"
import type { SchemaItem } from "@/api/types"
import { JsonEditor } from "@/components/json-editor"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { toast } from "@/components/ui/toast"
import { useErrorHandler } from "@/errors/utils"
import { useQueryClient } from "@tanstack/react-query"
import { useState, type SyntheticEvent } from "react"
import { useAgents } from "../auth/use-agents"

const DEFAULT_JSONS: Record<SchemaItem["item_type"], string> = {
  contract: '{"address": "your_wallet_address_here"}',
}

const ITEM_ACTION_TIPS: Record<SchemaItem["item_type"], string> = {
  contract:
    "Replace 'your_wallet_address_here' with your wallet address to receive the funds",
  portable_smelter: "This action does not require any input data!",
  portable_printer: "This action does not require any input data!",
  laser_drill: "This action does not require any input data!",
}

interface UseItemDialogParams {
  agentId: string
  item: SchemaItem
  open: boolean
  onClose: () => void
}

export function UseItemDialog({
  agentId,
  item,
  open,
  onClose,
}: UseItemDialogParams) {
  const [input, setInput] = useState(
    item.item_type in DEFAULT_JSONS ? DEFAULT_JSONS[item.item_type] : "{}"
  )

  const { setCooldown } = useAgents()
  const useItemMutation = useUseItemMutation()
  const errorHandler = useErrorHandler()
  const queryClient = useQueryClient()

  const useItem = async (event: SyntheticEvent) => {
    event.preventDefault()

    if (!agentId) return

    await useItemMutation.mutateAsync(
      {
        agentId,
        itemId: item.id,
        input: JSON.parse(input),
      },
      {
        onError(err) {
          errorHandler(err, "Failed to use item")
        },
        onSuccess(data) {
          queryClient.invalidateQueries({ queryKey: queryKeys.inventories.all })
          setCooldown(agentId, data.cooldown)
          toast.add({
            type: "success",
            title: "Success",
            description: "Successfully used item",
          })
          onClose()
        },
      }
    )
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(o) => {
        if (!o) onClose()
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Use Item Request</DialogTitle>
        </DialogHeader>
        {item.item_type in ITEM_ACTION_TIPS ? (
          <span>Tips: {ITEM_ACTION_TIPS[item.item_type]}</span>
        ) : null}
        <form className="space-y-4" onSubmit={useItem}>
          <Label className="grid gap-2 text-sm font-medium" htmlFor="input">
            Input
            <JsonEditor value={input} onChange={(v) => setInput(v)} />
          </Label>

          <Button
            className="w-full"
            type="submit"
            disabled={useItemMutation.isPending}
          >
            Use
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  )
}
