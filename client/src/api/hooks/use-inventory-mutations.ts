import { useMutation } from "@tanstack/react-query"

import { transferItems, transferResources } from "@/api/mutations/inventories"

export function useTransferResourcesMutation() {
  return useMutation({
    mutationFn: transferResources,
  })
}

export function useTransferItemsMutation() {
  return useMutation({
    mutationFn: transferItems,
  })
}
