import { queryOptions } from "@tanstack/react-query"

import { api, getAuthHeaders } from "@/api/client"
import { queryKeys } from "@/api/query-keys"

export const meQueryOptions = queryOptions({
  queryKey: queryKeys.me,
  queryFn: async () => {
    const { data, error } = await api.GET("/api/v1/auth/me", {
      headers: getAuthHeaders(),
    })
    if (error) throw error
    return data
  },
})
