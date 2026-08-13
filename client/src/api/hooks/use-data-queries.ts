import { useQuery } from "@tanstack/react-query"
import { recipesQueryOptions } from "../queries/data"

export function useRecipesQuery() {
  return useQuery(recipesQueryOptions())
}
