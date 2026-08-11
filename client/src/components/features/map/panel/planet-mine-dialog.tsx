import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { toast } from "@/components/ui/toast"
import { useErrorHandler } from "@/errors/utils"
import { useState  } from "react"
import type {SyntheticEvent} from "react";
import { useAgents } from "../../auth/use-agents"
import { useMinePlanetMutation } from "@/api/hooks/use-mining-mutations"
import type { SchemaResourceDeposit } from "@/api/types"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

interface PlanetMineDialogProps {
  agentId: string
  deposits: SchemaResourceDeposit[] | null
  onClose: () => void
}

export function PlanetMineDialog({
  agentId,
  deposits,
  onClose,
}: PlanetMineDialogProps) {
  const { setCooldown } = useAgents()
  const [amount, setAmount] = useState("")
  const [selectedDeposit, setSelectedDeposit] = useState(-1)
  const mineAsteroidMutation = useMinePlanetMutation()
  const errorHandler = useErrorHandler()

  async function mineAsteroid(event: SyntheticEvent) {
    event.preventDefault()

    if (!deposits) return

    const value = Number(amount)
    if (!Number.isInteger(value) || value <= 0) return

    if (selectedDeposit === -1) {
      toast.add({
        type: "error",
        title: "Error",
        description: "Select a valid resource deposit",
      })
      return
    }

    const deposit = deposits[selectedDeposit]

    await mineAsteroidMutation.mutateAsync(
      { agentId, body: { amount: value, resource: deposit.resource } },
      {
        onSuccess: (data) => {
          toast.add({
            type: "success",
            title: "Mining",
            description: `Requested ${value} resource${value === 1 ? "" : "s"}.`,
          })
          setAmount("")
          onClose()
          setCooldown(agentId, data)
        },
        onError: (error) => errorHandler(error, "Failed to mine planet"),
      }
    )
  }

  const depositSelectItems =
    deposits?.map((d, index) => ({
      label: `${d.resource} - ${d.amount}`,
      value: index,
    })) || []

  return (
    <Dialog
      open={deposits !== null}
      onOpenChange={(open) => {
        if (!open) {
          onClose()
        }
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Asteroid Mine Request</DialogTitle>
        </DialogHeader>
        <form className="space-y-4" onSubmit={mineAsteroid}>
          <Label
            className="grid gap-2 text-sm font-medium"
            htmlFor="mine-amount"
          >
            Amount to mine
            <Input
              id="mine-amount"
              type="number"
              min="1"
              step="1"
              inputMode="numeric"
              value={amount}
              onChange={(event) => setAmount(event.target.value)}
              placeholder="Enter amount"
              required
              autoFocus
            />
          </Label>
          <Label
            className="grid gap-2 text-sm font-medium"
            htmlFor="mine-deposit"
          >
            Resource deposit
            <Select
              items={depositSelectItems}
              onValueChange={(v) => setSelectedDeposit(v as number)}
            >
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Deposit" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  {depositSelectItems.map((item) => (
                    <SelectItem key={item.value} value={item.value}>
                      {item.label}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </Label>
          <Button
            className="w-full"
            type="submit"
            disabled={mineAsteroidMutation.isPending}
          >
            {mineAsteroidMutation.isPending ? "Submitting..." : "Start mining"}
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  )
}
