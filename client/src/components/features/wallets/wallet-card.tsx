import type { SchemaWallet } from "@/api/types"
import { Button } from "@/components/ui/button"
import { Wallet } from "lucide-react"

interface WalletCardProps {
  wallet: SchemaWallet

  onOpen?: () => void
}

export function WalletCard({ wallet, onOpen }: WalletCardProps) {
  return (
    <div
      className={`"hover:border-primary/20" rounded-lg border p-3 transition-colors`}
    >
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <Wallet className="size-4 shrink-0 text-primary" />

            <div className="min-w-0 flex-1">
              <span className="block font-medium break-all">
                {wallet.address.slice(0, 18) + "..."}
              </span>
            </div>
          </div>

          <div className="mt-2 space-y-1 text-sm text-muted-foreground">
            <div className="flex items-center gap-2">
              <span>units: {wallet.units}</span>
            </div>

            <div className="flex items-center gap-2">
              <span>quarks: {wallet.quarks}</span>
            </div>
          </div>
        </div>

        <div className="flex flex-col gap-2">
          <Button size="sm" variant="outline" onClick={onOpen}>
            Details
          </Button>
        </div>
      </div>
    </div>
  )
}
