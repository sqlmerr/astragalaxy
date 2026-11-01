import { useRegisterWalletMutation } from "@/api/hooks/use-wallets-mutations"
import { queryKeys } from "@/api/query-keys"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { toast } from "@/components/ui/toast"
import { useErrorHandler } from "@/errors/utils"
import { useQueryClient } from "@tanstack/react-query"

interface RegisterWalletDialogProps {
  agentId: string

  open: boolean
  onClose: () => void
}

export function RegisterWalletDialog({
  agentId,
  onClose,
  open,
}: RegisterWalletDialogProps) {
  const queryClient = useQueryClient()
  const errorHandler = useErrorHandler()
  const registerWalletMutation = useRegisterWalletMutation()

  async function registerWallet() {
    await registerWalletMutation.mutateAsync(
      { agentId },
      {
        onError(err) {
          errorHandler(err, "Failed to register wallet")
        },
        onSuccess() {
          queryClient.invalidateQueries({
            queryKey: queryKeys.wallets.my(agentId),
          })

          toast.add({
            type: "success",
            title: "Success",
            description: "Successfully registered new wallet",
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
          <DialogTitle>Register Wallet Request</DialogTitle>
        </DialogHeader>
        <p className="text-sm text-muted-foreground">
          Are you sure you want to register new wallet?
        </p>
        <Button
          className="w-full"
          onClick={registerWallet}
          disabled={registerWalletMutation.isPending}
        >
          {registerWalletMutation.isPending ? "Submitting..." : "Yes"}
        </Button>
      </DialogContent>
    </Dialog>
  )
}
