import { useQuery } from "@tanstack/react-query"

import {
  activeShipQueryOptions,
  myShipsQueryOptions,
  shipRadarQueryOptions,
} from "@/api/queries/ships"

export function useMyShipsQuery(agentID: string) {
  return useQuery(myShipsQueryOptions(agentID))
}

export function useActiveShipQuery(agentID?: string) {
  return useQuery({
    ...activeShipQueryOptions(agentID ?? ""),
    enabled: !!agentID,
  })
}

export function useShipRadarQuery(agentID?: string) {
  return useQuery({
    ...shipRadarQueryOptions(agentID ?? ""),
    enabled: !!agentID,
  })
}
