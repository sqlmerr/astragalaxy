import { usePayMutation } from "@/api/hooks/use-wallets-mutations"
import type { SchemaWallet } from "@/api/types"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useErrorHandler } from "@/errors/utils"
import { cn } from "@/lib/utils"
import { MinusIcon, PlusIcon } from "lucide-react"
import { useEffect, useState } from "react"
import type { SyntheticEvent } from "react"
import { useAgents } from "../auth/use-agents"
import { useQueryClient } from "@tanstack/react-query"
import { queryKeys } from "@/api/query-keys"
import { toast } from "@/components/ui/toast"

interface AdvancedNumberInputProps {
  label: string
  inputId: string
  maxValue: number
  value: number | null // current value

  onValueChange: (value: number) => void
}

function AdvancedNumberInput({
  label,
  inputId,
  maxValue,
  value,
  onValueChange,
}: AdvancedNumberInputProps) {
  const checked = value !== null

  const onValueChangeRaw = (raw: string) => {
    const parsed = Number.parseInt(raw, 10)
    const newValue = Number.isFinite(parsed)
      ? Math.max(0, Math.min(parsed, maxValue))
      : 0
    onValueChange(newValue)
  }

  const onStep = (delta: number) => {
    const currentValue = value ?? 0
    const newValue = Math.max(1, Math.min(currentValue + delta, maxValue))
    onValueChange(newValue)
  }

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
      <div className="min-w-0 flex-1">
        <Label
          htmlFor={inputId}
          className="block cursor-pointer truncate text-sm font-semibold"
        >
          {label}
        </Label>
        <p className="mt-1 truncate text-[11px] text-muted-foreground">
          Available: {maxValue}
        </p>
      </div>
      <div className="flex items-center gap-1">
        <Button
          type="button"
          variant="outline"
          size="icon-sm"
          disabled={!checked}
          onClick={() => onStep(-1)}
          aria-label={`Decrease ${label}`}
        >
          <MinusIcon />
        </Button>
        <Input
          className="w-20 text-center"
          type="number"
          inputMode="numeric"
          min={1}
          max={maxValue}
          step={1}
          value={value ?? ""}
          disabled={!checked}
          onChange={(e) => onValueChangeRaw(e.target.value)}
          aria-label={`${label} value`}
        />
        <Button
          type="button"
          variant="outline"
          size="icon-sm"
          disabled={!checked}
          onClick={() => onStep(1)}
          aria-label={`Increase ${label}`}
        >
          <PlusIcon />
        </Button>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          disabled={!checked}
          onClick={() => onValueChange(maxValue)}
        >
          Max
        </Button>
      </div>
    </div>
  )
}

interface PayActionDialogProps {
  wallet: SchemaWallet | null

  open: boolean
  onClose: () => void
}

export function PayActionDialog({
  wallet,
  open,
  onClose,
}: PayActionDialogProps) {
  const [toAddress, setToAddress] = useState("")
  const [units, setUnits] = useState(0)
  const [quarks, setQuarks] = useState(0)

  const queryClient = useQueryClient()
  const { setCooldown } = useAgents()
  const payMutation = usePayMutation()
  const errorHandler = useErrorHandler()

  const isOpen = open && wallet !== null

  useEffect(() => {
    if (!isOpen) {
      setToAddress("")
      setUnits(0)
      setQuarks(0)
    }
  }, [isOpen])

  async function pay(event: SyntheticEvent) {
    event.preventDefault()

    if (!wallet) return

    await payMutation.mutateAsync(
      {
        agentId: wallet.owner_id,
        address: wallet.address,
        body: {
          to: toAddress,
          units: units,
          quarks: quarks,
        },
      },
      {
        onError(err) {
          errorHandler(err, "Failed to pay")
        },
        onSuccess(data) {
          setCooldown(wallet.owner_id, data.cooldown)
          queryClient.invalidateQueries({
            queryKey: queryKeys.wallets.my(wallet.owner_id),
          })

          toast.add({
            type: "success",
            title: "Success",
            description: `Successfully payed to ${toAddress}`,
          })
          setUnits(0)
          setQuarks(0)
          setToAddress("")
          onClose()
        },
      }
    )
  }

  return (
    <Dialog
      open={isOpen}
      onOpenChange={(o) => {
        if (!o) onClose()
      }}
    >
      <DialogContent className="max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Pay Request</DialogTitle>
        </DialogHeader>
        <form className="space-y-4" onSubmit={pay}>
          <Label className="grid gap-2 text-sm font-medium" htmlFor="to">
            To address
            <Input
              id="to"
              value={toAddress}
              onChange={(e) => setToAddress(e.target.value)}
              placeholder="Enter address"
              required
              autoFocus
            />
          </Label>

          <Label className="grid gap-2 text-sm font-medium" htmlFor="units">
            Units
            <AdvancedNumberInput
              label="Units"
              inputId="units"
              maxValue={wallet?.units ?? 0}
              value={units}
              onValueChange={(v) => setUnits(v)}
            />
          </Label>

          <Label className="grid gap-2 text-sm font-medium" htmlFor="quarks">
            Quarks
            <AdvancedNumberInput
              label="Quarks"
              inputId="quarks"
              maxValue={wallet?.quarks ?? 0}
              value={quarks}
              onValueChange={(v) => setQuarks(v)}
            />
          </Label>

          <Button
            className="w-full"
            type="submit"
            disabled={payMutation.isPending}
          >
            Pay
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  )
}
