import { useMutation } from "@tanstack/react-query"

import {
  navigatePlanet,
  navigateWarp,
  navigateWaypoint,
} from "@/api/mutations/navigation"

export function useNavigateWarpMutation() {
  return useMutation({
    mutationFn: navigateWarp,
  })
}

export function useNavigatePlanetMutation() {
  return useMutation({
    mutationFn: navigatePlanet,
  })
}

export function useNavigateWaypointMutation() {
  return useMutation({
    mutationFn: navigateWaypoint,
  })
}
