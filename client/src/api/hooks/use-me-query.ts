import { useQuery } from "@tanstack/react-query"

import { meQueryOptions } from "@/api/queries/me"

export function useMeQuery() {
  return useQuery(meQueryOptions)
}
