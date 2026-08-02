import { createFileRoute } from "@tanstack/react-router"
import { AuthGate } from "@/components/features/auth/auth-gate"
import { StarMapLayout } from "@/layouts/star-map-layout"
import { toast } from "@/components/ui/toast"

export const Route = createFileRoute("/")({ component: MainPage })

function MainPage() {
  return (
    <AuthGate>
      <StarMapLayout />
    </AuthGate>
  )
}
