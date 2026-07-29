import { useMutation } from "@tanstack/react-query"

import {
  navigatePlanet,
  navigateWarp,
  navigateWaypoint,
} from "@/api/mutations/navigation"
import type {
  SchemaNavigatePlanetRequest,
  SchemaNavigateWarpRequest,
  SchemaNavigateWaypointRequest,
} from "@/api/types"

export function useNavigateWarpMutation() {
  return useMutation({
    mutationFn: (data: { agentId: string; body: SchemaNavigateWarpRequest }) =>
      navigateWarp(data.agentId, data.body),
  })
}

export function useNavigatePlanetMutation(agentID: string) {
  return useMutation({
    mutationFn: (body: SchemaNavigatePlanetRequest) =>
      navigatePlanet(agentID, body),
  })
}

export function useNavigateWaypointMutation(agentID: string) {
  return useMutation({
    mutationFn: (body: SchemaNavigateWaypointRequest) =>
      navigateWaypoint(agentID, body),
  })
}
