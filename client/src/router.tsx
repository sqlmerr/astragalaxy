import { createRouter as createTanStackRouter } from "@tanstack/react-router"
import { routeTree } from "./routeTree.gen"
import { Spinner } from "@/components/ui/spinner.tsx"

export function getRouter() {
  const router = createTanStackRouter({
    routeTree,
    defaultPendingComponent: () => (
      <div className="grid min-h-svh place-items-center bg-background">
        <Spinner />
      </div>
    ),
    scrollRestoration: true,
    defaultPreload: "intent",
    defaultPreloadStaleTime: 0,
  })

  return router
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof getRouter>
  }
}
