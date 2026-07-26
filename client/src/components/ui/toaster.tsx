import { AlertCircle, CheckCircle2, Info, X } from "lucide-react"

import { useToast } from "@/lib/toast"
import type { ToastVariant } from "@/lib/toast"
import { cn } from "@/lib/utils"

const variantStyles: Record<ToastVariant, string> = {
  success: "border-emerald-500/40 bg-emerald-950/80 text-emerald-200",
  error: "border-destructive/40 bg-destructive/15 text-destructive",
  info: "border-primary/30 bg-primary/10 text-primary",
}

const variantIcons: Record<
  ToastVariant,
  React.ComponentType<{ className?: string }>
> = {
  success: CheckCircle2,
  error: AlertCircle,
  info: Info,
}

export function Toaster() {
  const { toasts, removeToast } = useToast()

  if (toasts.length === 0) return null

  return (
    <div
      className="fixed right-4 bottom-4 z-50 flex w-full max-w-sm flex-col gap-2 sm:right-6 sm:bottom-6"
      aria-live="polite"
      aria-label="Notifications"
    >
      {toasts.map((toast) => {
        const Icon = variantIcons[toast.variant]

        return (
          <div
            key={toast.id}
            role="alert"
            className={cn(
              "group flex animate-in items-start gap-3 border p-4 shadow-lg backdrop-blur-sm transition-all duration-300 slide-in-from-bottom-5 fade-in",
              variantStyles[toast.variant]
            )}
          >
            <Icon className="mt-0.5 size-4 shrink-0" aria-hidden="true" />
            <div className="min-w-0 flex-1">
              <p className="text-xs font-semibold tracking-wide">
                {toast.title}
              </p>
              {toast.description ? (
                <p className="mt-1 text-[11px] leading-relaxed opacity-80">
                  {toast.description}
                </p>
              ) : null}
            </div>
            <button
              type="button"
              onClick={() => removeToast(toast.id)}
              className="mt-0.5 shrink-0 opacity-60 transition-opacity hover:opacity-100"
              aria-label="Dismiss notification"
            >
              <X className="size-3.5" aria-hidden="true" />
            </button>
          </div>
        )
      })}
    </div>
  )
}
