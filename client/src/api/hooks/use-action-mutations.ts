import { useMutation } from "@tanstack/react-query"
import { craftAction } from "../mutations/agents"
import type { SchemaCraftRequest } from "../types"

export function useCraftMutation() {
  return useMutation({
    mutationFn: ({
      agentId,
      body,
    }: {
      agentId: string
      body: SchemaCraftRequest
    }) => craftAction(agentId, body),
  })
}
