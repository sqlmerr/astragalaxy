import { useQuery } from "@tanstack/react-query"

import {
  myInventoryQueryOptions,
  shipInventoryQueryOptions,
} from "@/api/queries/inventories"

export function useMyInventoryQuery(agentID: string) {
  return useQuery(myInventoryQueryOptions(agentID))
}

export function useShipInventoryQuery(agentID: string, shipId: string) {
  return useQuery(shipInventoryQueryOptions(agentID, shipId))
}
