import { useMutation } from "@tanstack/react-query"

import { transferItems, transferResources } from "@/api/mutations/inventories"
import type {
  SchemaTransferItemsRequest,
  SchemaTransferResourcesRequest,
} from "@/api/types"

export function useTransferResourcesMutation(agentID: string) {
  return useMutation({
    mutationFn: (body: SchemaTransferResourcesRequest) =>
      transferResources(agentID, body),
  })
}

export function useTransferItemsMutation(agentID: string) {
  return useMutation({
    mutationFn: (body: SchemaTransferItemsRequest) =>
      transferItems(agentID, body),
  })
}
