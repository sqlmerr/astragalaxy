import { useQuery } from "@tanstack/react-query"

import {
  activeShipQueryOptions,
  myShipsQueryOptions,
  shipModulesQueryOptions,
  shipRadarQueryOptions,
} from "@/api/queries/ships"

export function useMyShipsQuery(agentId: string) {
  return useQuery(myShipsQueryOptions(agentId))
}

export function useActiveShipQuery(agentId?: string) {
  return useQuery({
    ...activeShipQueryOptions(agentId ?? ""),
    enabled: !!agentId,
  })
}

export function useShipRadarQuery(agentId?: string) {
  return useQuery({
    ...shipRadarQueryOptions(agentId ?? ""),
    enabled: !!agentId,
  })
}

export function useShipModulesQuery(
  agentId: string,
  shipId: string,
  enabled: boolean
) {
  return useQuery({
    ...shipModulesQueryOptions(agentId, shipId),
    enabled,
  })
}
