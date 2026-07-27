import { useQuery } from "@tanstack/react-query"

import {
  agentCooldownQueryOptions,
  currentAgentQueryOptions,
  myAgentsQueryOptions,
} from "@/api/queries/agents"

export function useMyAgentsQuery() {
  return useQuery(myAgentsQueryOptions())
}

export function useCurrentAgentQuery(agentID: string) {
  return useQuery(currentAgentQueryOptions(agentID))
}

export function useAgentCooldownQuery(agentID: string) {
  return useQuery(agentCooldownQueryOptions(agentID))
}
