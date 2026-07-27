import { useMyAgentsQuery } from "@/api/hooks"
import {
  activeShipQueryOptions,
  myShipsQueryOptions,
} from "@/api/queries/ships"
import type { SchemaAgent, SchemaShip } from "@/api/types"
import { useQueries } from "@tanstack/react-query"
import { useMemo } from "react"

export interface AgentWithShip {
  agent: SchemaAgent
  ship: SchemaShip
}

export function useAgentsWithShips() {
  const {
    data: agents = { data: [] },
    isPending: isAgentsPending,
    isError: isAgentsError,
  } = useMyAgentsQuery()

  const shipQueries = useQueries({
    queries: agents.data.map((agent) => ({
      ...activeShipQueryOptions(agent.id),
    })),
  })

  const isShipsPending = shipQueries.some((q) => q.isPending)
  const isShipsError = shipQueries.some((q) => q.isError)

  const data = useMemo<AgentWithShip[]>(() => {
    return agents.data.flatMap((agent, index) => {
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
    isPending: isAgentsPending || isShipsPending,
    isError: isAgentsError || isShipsError,
  }
}
