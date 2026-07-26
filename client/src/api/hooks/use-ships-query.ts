import { useQuery } from "@tanstack/react-query"

import {
  activeShipQueryOptions,
  myShipsQueryOptions,
  shipRadarQueryOptions,
} from "@/api/queries/ships"

export function useMyShipsQuery() {
  return useQuery(myShipsQueryOptions)
}

export function useActiveShipQuery() {
  return useQuery(activeShipQueryOptions)
}

export function useShipRadarQuery() {
  return useQuery(shipRadarQueryOptions)
}
