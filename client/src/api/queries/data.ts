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

export const itemsQueryOptions = () =>
  queryOptions({
    queryKey: queryKeys.data.items,
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/data/items")
      if (error) throw error
      return data
    },
  })

export const resourcesQueryOptions = () =>
  queryOptions({
    queryKey: queryKeys.data.resources,
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/data/resources")
      if (error) throw error
      return data
    },
  })

export const facilitiesQueryOptions = () =>
  queryOptions({
    queryKey: queryKeys.data.facilities,
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/data/facilities")
      if (error) throw error
      return data
    },
  })
