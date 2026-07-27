import { queryOptions } from "@tanstack/react-query"

import { api, getAuthHeaders } from "@/api/client"
import { queryKeys } from "@/api/query-keys"

export const myAgentsQueryOptions = () =>
  queryOptions({
    queryKey: queryKeys.agents.my,
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/agents/my", {
        headers: getAuthHeaders(),
      })
      if (error) throw error
      return data
    },
  })

export const currentAgentQueryOptions = (agentID: string) =>
  queryOptions({
    queryKey: queryKeys.agents.current(agentID),
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/agents/current", {
        headers: getAuthHeaders(agentID),
      })
      if (error) throw error
      return data
    },
  })

export const agentCooldownQueryOptions = (agentID: string) =>
  queryOptions({
    queryKey: queryKeys.agents.cooldown(agentID),
    queryFn: async () => {
      const { data, error } = await api.GET(
        "/api/v1/agents/current/cooldown",
        {
          headers: getAuthHeaders(agentID),
        }
      )
      if (error) throw error
      return data
    },
  })
