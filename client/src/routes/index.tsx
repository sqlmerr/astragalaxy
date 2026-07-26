import { createFileRoute } from "@tanstack/react-router"
import { AuthGate } from "@/components/features/auth/auth-gate"
import { StarMapLayout } from "@/layouts/star-map-layout"

export const Route = createFileRoute("/")({ component: MainPage })

export function MainPage() {
  return (
    <AuthGate>
      <StarMapLayout />
    </AuthGate>
  )
}
