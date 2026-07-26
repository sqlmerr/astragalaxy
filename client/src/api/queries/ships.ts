import { queryOptions } from "@tanstack/react-query"

import { api, getAuthHeaders } from "@/api/client"
import { queryKeys } from "@/api/query-keys"

export const myShipsQueryOptions = queryOptions({
  queryKey: queryKeys.ships.my,
  queryFn: async () => {
    const { data, error } = await api.GET("/api/v1/ships/my", {
      headers: getAuthHeaders(),
    })
    if (error) throw error
    return data
  },
})

export const activeShipQueryOptions = queryOptions({
  queryKey: queryKeys.ships.active,
  queryFn: async () => {
    const { data, error } = await api.GET("/api/v1/ships/my/active", {
      headers: getAuthHeaders(),
    })
    if (error) throw error
    return data
  },
})

export const shipRadarQueryOptions = queryOptions({
  queryKey: queryKeys.ships.radar,
  queryFn: async () => {
    const { data, error } = await api.GET("/api/v1/ships/my/active/radar", {
      headers: getAuthHeaders(),
    })
    if (error) throw error
    return data
  },
})
