import { useQuery } from "@tanstack/react-query"

import {
  agentCooldownQueryOptions,
  currentAgentQueryOptions,
  myAgentsQueryOptions,
} from "@/api/queries/agents"

export function useMyAgentsQuery() {
  return useQuery(myAgentsQueryOptions)
}

export function useCurrentAgentQuery() {
  return useQuery(currentAgentQueryOptions)
}

export function useAgentCooldownQuery() {
  return useQuery(agentCooldownQueryOptions)
}
