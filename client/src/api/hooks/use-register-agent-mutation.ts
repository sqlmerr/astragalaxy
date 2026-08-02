import { useMutation } from "@tanstack/react-query"

import { registerAgent } from "@/api/mutations/agents"

export function useRegisterAgentMutation() {
  return useMutation({
    mutationFn: registerAgent,
  })
}
