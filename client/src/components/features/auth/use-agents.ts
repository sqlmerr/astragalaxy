import {
  activeShipQueryOptions,
  myShipsQueryOptions,
} from "@/api/queries/ships"
import type { AgentExtended, SchemaCooldown } from "@/api/types"
import { useQueries, useQueryClient } from "@tanstack/react-query"
import { act, useMemo } from "react"
import { useAuth } from "./auth-provider"
import { agentCooldownQueryOptions } from "@/api/queries/agents"
import { queryKeys } from "@/api/query-keys"

export function useAgents() {
  const queryClient = useQueryClient()
  const { agents } = useAuth()

  if (!agents) {
    return {
      data: [],
      isPending: true,
      setCooldown: () => {},
    }
  }

  const shipQueries = useQueries({
    queries: agents.map((agent) => ({
      ...myShipsQueryOptions(agent.id),
    })),
  })

  const cooldownQueries = useQueries({
    queries: agents.map((agent) => ({
      ...agentCooldownQueryOptions(agent.id),
    })),
  })

  const isShipsPending = shipQueries.some((q) => q.isPending)
  const isShipsError = shipQueries.some((q) => q.isError)

  const isCooldownsPending = cooldownQueries.some((q) => q.isPending)
  const isCooldownsError = cooldownQueries.some((q) => q.isError)

  const data = useMemo<AgentExtended[]>(() => {
    return agents.flatMap((agent, index) => {
      const ships = shipQueries[index].data

      if (!ships) return []

      const cooldown = cooldownQueries[index].data
      if (!cooldown) return []

      const activeShip = ships.data.find((s) => s.agent_id === agent.id)

      if (!activeShip) return []

      return [
        {
          agent,
          ship: activeShip,
          cooldown,
          ships: ships.data,
        },
      ]
    })
  }, [agents, shipQueries, cooldownQueries])

  const setCooldown = (agentId: string, cooldown: SchemaCooldown) => {
    queryClient.setQueryData(queryKeys.agents.cooldown(agentId), cooldown)
  }

  return {
    data,
    isPending: isShipsPending || isCooldownsPending,
    isError: isShipsError || isCooldownsError,
    setCooldown,
  }
}
