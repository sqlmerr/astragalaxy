import { useCraftMutation } from "@/api/hooks/use-action-mutations"
import { Button } from "@/components/ui/button"
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
import { toast } from "@/components/ui/toast"
import { useErrorHandler } from "@/errors/utils"
import { useData } from "@/hooks/use-data"
import { useState, type SyntheticEvent } from "react"
import { useAgents } from "../auth/use-agents"
import { useQueryClient } from "@tanstack/react-query"
import { queryKeys } from "@/api/query-keys"

interface CraftActionDialogProps {
  agentId: string
  inventoryId: string | null
  open: boolean
  onClose: () => void
}

export function CraftActionDialog({
  agentId,
  inventoryId,
  open,
  onClose,
}: CraftActionDialogProps) {
  const [amount, setAmount] = useState("")
  const [recipeId, setRecipeId] = useState<string | null>(null)

  const queryClient = useQueryClient()
  const { setCooldown } = useAgents()
  const craftActionMutation = useCraftMutation()
  const errorHandler = useErrorHandler()
  const { recipes } = useData()

  async function craft(event: SyntheticEvent) {
    event.preventDefault()
    if (!agentId || !recipeId || !inventoryId) {
      return
    }

    const amountNumber = Number(amount)
    if (!Number.isInteger(amountNumber) || amountNumber <= 0) return

    await craftActionMutation.mutateAsync(
      {
        agentId,
        body: {
          amount: amountNumber,
          recipe_id: recipeId,
          target_inventory_id: inventoryId,
        },
      },
      {
        onError(err) {
          errorHandler(err, "Failed to craft recipe")
        },
        onSuccess(data) {
          setCooldown(agentId, data.cooldown)
          queryClient.invalidateQueries({ queryKey: queryKeys.inventories.all })

          toast.add({
            type: "success",
            title: "Success",
            description: `Crafted recipe ${recipeId}`,
          })
          setAmount("")
          setRecipeId(null)
          onClose()
        },
      }
    )
  }

  const recipeSelectItems = recipes.map((r) => ({
    label: r.id,
    value: r.id,
  }))

  return (
    <Dialog
      open={open}
      onOpenChange={(o) => {
        if (!o) onClose()
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Craft Action Request</DialogTitle>
        </DialogHeader>
        <form className="space-y-4" onSubmit={craft}>
          <Label className="grid gap-2 text-sm font-medium" htmlFor="amount">
            Recipe
            <Select
              items={recipeSelectItems}
              onValueChange={setRecipeId}
              required
            >
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Recipe" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  {recipeSelectItems.map((item) => (
                    <SelectItem key={item.value} value={item.value}>
                      {item.label}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </Label>

          <Label className="grid gap-2 text-sm font-medium" htmlFor="amount">
            Amount to craft
            <Input
              id="amount"
              type="number"
              min="1"
              step="1"
              inputMode="numeric"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              placeholder="Enter amount"
              required
              autoFocus
            />
          </Label>

          <Button
            className="w-full"
            type="submit"
            disabled={craftActionMutation.isPending}
          >
            Craft
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  )
}
