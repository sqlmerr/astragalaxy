import { queryOptions } from "@tanstack/react-query"

import { api, getAuthHeaders } from "@/api/client"
import { queryKeys } from "@/api/query-keys"

export const myAgentsQueryOptions = queryOptions({
  queryKey: queryKeys.agents.my,
  queryFn: async () => {
    const { data, error } = await api.GET("/api/v1/agents/my", {
      headers: getAuthHeaders(),
    })
    if (error) throw error
    return data
  },
})

export const currentAgentQueryOptions = queryOptions({
  queryKey: queryKeys.agents.current,
  queryFn: async () => {
    const { data, error } = await api.GET("/api/v1/agents/current", {
      headers: getAuthHeaders(),
    })
    if (error) throw error
    return data
  },
})

export const agentCooldownQueryOptions = queryOptions({
  queryKey: queryKeys.agents.cooldown,
  queryFn: async () => {
    const { data, error } = await api.GET("/api/v1/agents/current/cooldown", {
      headers: getAuthHeaders(),
    })
    if (error) throw error
    return data
  },
})
