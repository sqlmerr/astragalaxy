import { cn } from "@/lib/utils"
import { Checkbox } from "./ui/checkbox"
import { Label } from "./ui/label"
import { Button } from "./ui/button"
import { MinusIcon, PlusIcon } from "lucide-react"
import { Input } from "./ui/input"

interface AdvancedNumberInputProps {
  label: string
  inputId: string
  maxValue: number
  value: number | null // current value

  onToggle: () => void
  onValueChange: (value: string) => void
  onStep: (delta: number) => void
  onMax: () => void
}

export function AdvancedNumberInput({
  label,
  inputId,
  maxValue,
  value,
  onToggle,
  onValueChange,
  onStep,
  onMax,
}: AdvancedNumberInputProps) {
  const checked = value !== null

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
        aria-label={label}
      />
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
          onChange={(e) => onValueChange(e.target.value)}
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
          onClick={onMax}
        >
          Max
        </Button>
      </div>
    </div>
  )
}
