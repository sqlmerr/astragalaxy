import { queryOptions } from "@tanstack/react-query"
import { queryKeys } from "../query-keys"
import { api } from "../client"

export const recipesQueryOptions = () =>
  queryOptions({
    queryKey: queryKeys.data.recipes,
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/data/recipes")
      if (error) throw error
      return data
    },
  })
