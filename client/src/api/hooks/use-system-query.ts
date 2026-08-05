import { useQuery } from "@tanstack/react-query"

import { currentSystemQueryOptions } from "@/api/queries/systems"

export function useCurrentSystemQuery(agentID?: string) {
  return useQuery({
    ...currentSystemQueryOptions(agentID ?? ""),
    enabled: !!agentID,
  })
}
