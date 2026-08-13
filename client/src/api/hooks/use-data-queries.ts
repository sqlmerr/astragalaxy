import { useQuery } from "@tanstack/react-query"
import {
  facilitiesQueryOptions,
  itemsQueryOptions,
  recipesQueryOptions,
  resourcesQueryOptions,
} from "../queries/data"

export function useRecipesQuery() {
  return useQuery(recipesQueryOptions())
}

export function useItemsQuery() {
  return useQuery(itemsQueryOptions())
}

export function useResourcesQuery() {
  return useQuery(resourcesQueryOptions())
}

export function useFacilitiesQuery() {
  return useQuery(facilitiesQueryOptions())
}
