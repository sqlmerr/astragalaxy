import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { useError } from "./store"
import { Json } from "@/components/ui/json"

export function ErrorDialog() {
  const { error, close } = useError()

  return (
    <Dialog open={!!error} onOpenChange={close}>
      <DialogContent className="max-w-3xl">
        <DialogHeader>
          <DialogTitle>{error?.title}</DialogTitle>
        </DialogHeader>

        {error && <Json data={error.error} />}
      </DialogContent>
    </Dialog>
  )
}
