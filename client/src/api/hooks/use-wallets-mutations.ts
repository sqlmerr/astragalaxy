import { useMutation } from "@tanstack/react-query"
import type { SchemaPayRequest } from "../types"
import { pay } from "../mutations/wallets"

export function usePayMutation() {
  return useMutation({
    mutationFn: ({
      agentId,
      address,
      body,
    }: {
      agentId: string
      address: string
      body: SchemaPayRequest
    }) => pay(agentId, address, body),
  })
}
