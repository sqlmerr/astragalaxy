import { useQuery } from "@tanstack/react-query"

import {
  myInventoryQueryOptions,
  shipInventoryQueryOptions,
} from "@/api/queries/inventories"

export function useMyInventoryQuery() {
  return useQuery(myInventoryQueryOptions)
}

export function useShipInventoryQuery(shipId: string) {
  return useQuery(shipInventoryQueryOptions(shipId))
}
