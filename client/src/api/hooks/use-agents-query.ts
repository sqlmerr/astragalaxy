import { useQuery } from "@tanstack/react-query"

import {
  agentCooldownQueryOptions,
  currentAgentQueryOptions,
  myAgentsQueryOptions,
} from "@/api/queries/agents"

export function useMyAgentsQuery(enabled: boolean = true) {
  return useQuery({ ...myAgentsQueryOptions(), enabled: enabled })
}

export function useCurrentAgentQuery(agentID: string) {
  return useQuery(currentAgentQueryOptions(agentID))
}

export function useAgentCooldownQuery(agentID: string) {
  return useQuery(agentCooldownQueryOptions(agentID))
}
