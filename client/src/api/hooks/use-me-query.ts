import { useQuery } from "@tanstack/react-query"

import { meQueryOptions } from "@/api/queries/me"

export function useMeQuery(enabled: boolean = true) {
  return useQuery({ ...meQueryOptions(), enabled })
}
