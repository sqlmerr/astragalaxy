import { useQuery } from "@tanstack/react-query"

import {
  myInventoryQueryOptions,
  shipInventoryQueryOptions,
} from "@/api/queries/inventories"

export function useMyInventoryQuery(agentId: string, enabled: boolean) {
  return useQuery({ ...myInventoryQueryOptions(agentId), enabled })
}

export function useShipInventoryQuery(
  agentId: string,
  shipId: string,
  enabled: boolean
) {
  return useQuery({ ...shipInventoryQueryOptions(agentId, shipId), enabled })
}
