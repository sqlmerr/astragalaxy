import { queryOptions } from "@tanstack/react-query"
import { queryKeys } from "../query-keys"
import { api, getAuthHeaders } from "../client"

export const myWalletsQueryOptions = (agentId: string) =>
  queryOptions({
    queryKey: queryKeys.wallets.my(agentId),
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/wallets/my", {
        headers: getAuthHeaders(agentId),
      })
      if (error) throw error
      return data
    },
  })
