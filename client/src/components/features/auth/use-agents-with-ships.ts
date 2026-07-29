import { useMyAgentsQuery } from "@/api/hooks"
import {
  activeShipQueryOptions,
  myShipsQueryOptions,
} from "@/api/queries/ships"
import type { AgentWithShip } from "@/api/types"
import { useQueries } from "@tanstack/react-query"
import { useMemo } from "react"
import { useAuth } from "./auth-provider"

export function useAgentsWithShips() {
  const { agents } = useAuth()

  if (!agents) {
    return {
      data: [],
      isPending: true,
    }
  }

  const shipQueries = useQueries({
    queries: agents.map((agent) => ({
      ...activeShipQueryOptions(agent.id),
    })),
  })

  const isShipsPending = shipQueries.some((q) => q.isPending)
  const isShipsError = shipQueries.some((q) => q.isError)

  const data = useMemo<AgentWithShip[]>(() => {
    return agents.flatMap((agent, index) => {
      const ship = shipQueries[index].data

      if (!ship) return []

      return [
        {
          agent,
          ship,
        },
      ]
    })
  }, [agents, shipQueries])

  return {
    data,
    isPending: isShipsPending,
    isError: isShipsError,
  }
}
