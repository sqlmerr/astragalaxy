import { useMutation } from "@tanstack/react-query"

import {
  transferItems,
  transferResources,
  useItem,
} from "@/api/mutations/inventories"
import type {
  SchemaTransferItemsRequest,
  SchemaTransferResourcesRequest,
} from "@/api/types"

export function useTransferResourcesMutation() {
  return useMutation({
    mutationFn: ({
      agentId,
      body,
    }: {
      agentId: string
      body: SchemaTransferResourcesRequest
    }) => transferResources(agentId, body),
  })
}

export function useTransferItemsMutation() {
  return useMutation({
    mutationFn: ({
      agentId,
      body,
    }: {
      agentId: string
      body: SchemaTransferItemsRequest
    }) => transferItems(agentId, body),
  })
}

export function useUseItemMutation() {
  return useMutation({
    mutationFn: ({
      agentId,
      itemId,
      input,
    }: {
      agentId: string
      itemId: string
      input: any
    }) => useItem(agentId, itemId, input),
  })
}
