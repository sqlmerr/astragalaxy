import { queryOptions } from "@tanstack/react-query"

import { api, getAuthHeaders } from "@/api/client"
import { queryKeys } from "@/api/query-keys"

export const myShipsQueryOptions = (agentId: string) =>
  queryOptions({
    queryKey: queryKeys.ships.my(agentId),
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/ships/my", {
        headers: getAuthHeaders(agentId),
      })
      if (error) throw error
      return data
    },
  })

export const activeShipQueryOptions = (agentId: string) =>
  queryOptions({
    queryKey: queryKeys.ships.active(agentId),
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/ships/my/active", {
        headers: getAuthHeaders(agentId),
      })
      if (error) throw error
      return data
    },
  })

export const shipRadarQueryOptions = (agentId: string) =>
  queryOptions({
    queryKey: queryKeys.ships.radar(agentId),
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/ships/my/active/radar", {
        headers: getAuthHeaders(agentId),
      })
      if (error) {
        throw error
      }
      return data
    },
  })

export const shipModulesQueryOptions = (agentId: string, shipId: string) =>
  queryOptions({
    queryKey: queryKeys.ships.modules(agentId, shipId),
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/ships/my/{id}/modules", {
        params: { path: { id: shipId } },
        headers: getAuthHeaders(agentId),
      })
      if (error) {
        throw error
      }
      return data
    },
  })
