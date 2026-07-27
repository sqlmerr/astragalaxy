import { createFileRoute } from "@tanstack/react-router"
import { AuthGate } from "@/components/features/auth/auth-gate"
import { StarMapLayout } from "@/layouts/star-map-layout"

export const Route = createFileRoute("/")({ component: MainPage })

function MainPage() {
  return (
    <AuthGate>
      <StarMapLayout />
    </AuthGate>
  )
}
