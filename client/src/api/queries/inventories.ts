import { queryOptions } from "@tanstack/react-query"

import { api, getAuthHeaders } from "@/api/client"
import { queryKeys } from "@/api/query-keys"

export const myInventoryQueryOptions = queryOptions({
  queryKey: queryKeys.inventories.my,
  queryFn: async () => {
    const { data, error } = await api.GET("/api/v1/inventories/my", {
      headers: getAuthHeaders(),
    })
    if (error) throw error
    return data
  },
})

export function shipInventoryQueryOptions(shipId: string) {
  return queryOptions({
    queryKey: queryKeys.inventories.ship(shipId),
    queryFn: async () => {
      const { data, error } = await api.GET(
        "/api/v1/inventories/my/ships/{id}",
        {
          params: { path: { id: shipId } },
          headers: getAuthHeaders(),
        }
      )
      if (error) throw error
      return data
    },
  })
}
