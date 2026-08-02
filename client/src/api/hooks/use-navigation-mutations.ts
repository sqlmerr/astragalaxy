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

export function useNavigatePlanetMutation() {
  return useMutation({
    mutationFn: ({
      body,
      agentId,
    }: {
      agentId: string
      body: SchemaNavigatePlanetRequest
    }) => navigatePlanet(agentId, body),
  })
}

export function useNavigateWaypointMutation() {
  return useMutation({
    mutationFn: ({
      body,
      agentId,
    }: {
      agentId: string
      body: SchemaNavigateWaypointRequest
    }) => navigateWaypoint(agentId, body),
  })
}
