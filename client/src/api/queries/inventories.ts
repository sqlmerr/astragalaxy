import { queryOptions } from "@tanstack/react-query"

import { api, getAuthHeaders } from "@/api/client"
import { queryKeys } from "@/api/query-keys"

export const myInventoryQueryOptions = (agentID: string) =>
  queryOptions({
    queryKey: queryKeys.inventories.my(agentID),
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/inventories/my", {
        headers: getAuthHeaders(agentID),
      })
      if (error) throw error
      return data
    },
  })

export const shipInventoryQueryOptions = (agentID: string, shipId: string) =>
  queryOptions({
    queryKey: queryKeys.inventories.ship(agentID, shipId),
    queryFn: async () => {
      const { data, error } = await api.GET(
        "/api/v1/inventories/my/ships/{id}",
        {
          params: { path: { id: shipId } },
          headers: getAuthHeaders(agentID),
        }
      )
      if (error) throw error
      return data
    },
  })
