import { useQuery } from "@tanstack/react-query"

import { currentSystemQueryOptions } from "@/api/queries/systems"

export function useCurrentSystemQuery() {
  return useQuery(currentSystemQueryOptions)
}
