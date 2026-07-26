import { queryOptions } from "@tanstack/react-query"

import { api, getAuthHeaders } from "@/api/client"
import { queryKeys } from "@/api/query-keys"

export const currentSystemQueryOptions = queryOptions({
  queryKey: queryKeys.systems.current,
  queryFn: async () => {
    const { data, error } = await api.GET("/api/v1/systems/current", {
      headers: getAuthHeaders(),
    })
    if (error) throw error
    return data
  },
})
