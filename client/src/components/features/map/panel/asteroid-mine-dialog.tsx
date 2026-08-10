import { useMineAsteroidMutation } from "@/api/hooks"
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
import { useState, type SyntheticEvent } from "react"
import { useAgents } from "../../auth/use-agents"

interface AsteroidMineDialogProps {
  agentId: string
  open: boolean
  onClose: () => void
}

export function AsteroidMineDialog({
  agentId,
  open,
  onClose,
}: AsteroidMineDialogProps) {
  const { setCooldown } = useAgents()
  const [amount, setAmount] = useState("")
  const mineAsteroidMutation = useMineAsteroidMutation()
  const errorHandler = useErrorHandler()

  async function mineAsteroid(event: SyntheticEvent) {
    event.preventDefault()

    const value = Number(amount)
    if (!Number.isInteger(value) || value <= 0) return

    await mineAsteroidMutation.mutateAsync(
      { agentId, body: { amount: value } },
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
        onError: (error) => errorHandler(error, "Failed to mine asteroid"),
      }
    )
  }

  return (
    <Dialog
      open={open}
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
          <label
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
          </label>
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
