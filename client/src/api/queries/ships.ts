import { queryOptions } from "@tanstack/react-query"

import { api, getAuthHeaders } from "@/api/client"
import { queryKeys } from "@/api/query-keys"

export const myShipsQueryOptions = (agentID: string) =>
  queryOptions({
    queryKey: queryKeys.ships.my(agentID),
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/ships/my", {
        headers: getAuthHeaders(agentID),
      })
      if (error) throw error
      return data
    },
  })

export const activeShipQueryOptions = (agentID: string) =>
  queryOptions({
    queryKey: queryKeys.ships.active(agentID),
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/ships/my/active", {
        headers: getAuthHeaders(agentID),
      })
      if (error) throw error
      return data
    },
  })

export const shipRadarQueryOptions = (agentID: string) =>
  queryOptions({
    queryKey: queryKeys.ships.radar(agentID),
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/ships/my/active/radar", {
        headers: getAuthHeaders(agentID),
      })
      if (error) {
        throw error
      }
      return data
    },
  })
