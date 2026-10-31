import type { SchemaWallet } from "@/api/types"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { PayActionDialog } from "./pay-action-dialog"
import { useEffect, useState } from "react"

interface WalletModalProps {
  wallet: SchemaWallet | null
  onClose: () => void
}

export function WalletModal({ wallet, onClose }: WalletModalProps) {
  const [payDialogOpened, setPayDialogOpened] = useState(false)

  useEffect(() => {
    if (wallet === null) setPayDialogOpened(false)
  }, [wallet])

  return (
    <>
      <Dialog
        open={wallet !== null}
        onOpenChange={(o) => {
          if (!o) {
            setPayDialogOpened(false)
            onClose()
          }
        }}
      >
        <DialogContent className="max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 pr-10">
              Wallet
            </DialogTitle>
          </DialogHeader>
          {wallet && (
            <div className="space-y-6">
              <div className="grid grid-cols-[80px_1fr] gap-x-3 gap-y-2 text-sm">
                <span className="text-muted-foreground">Address</span>
                <code className="font-mono break-all">{wallet.address}</code>

                <span className="text-muted-foreground">Units</span>
                <code className="font-mono break-all">{wallet.units}</code>

                <span className="text-muted-foreground">Quarks</span>
                <code className="font-mono break-all">{wallet.quarks}</code>
              </div>

              <div className="flex items-center justify-center">
                <Button onClick={() => setPayDialogOpened(true)}>Pay</Button>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
      <PayActionDialog
        wallet={wallet}
        onClose={() => setPayDialogOpened(false)}
        open={payDialogOpened}
      />
    </>
  )
}
