import { api, getAuthHeaders } from "@/api/client"
import type {
  SchemaNavigatePlanetRequest,
  SchemaNavigateWarpRequest,
  SchemaNavigateWaypointRequest,
} from "@/api/types"

export async function navigateWarp(
  agentID: string,
  body: SchemaNavigateWarpRequest
) {
  const { data, error } = await api.POST("/api/v1/navigation/warp", {
    body,
    headers: getAuthHeaders(agentID),
  })
  if (error) throw error
  return data
}

export async function navigatePlanet(
  agentID: string,
  body: SchemaNavigatePlanetRequest
) {
  const { data, error } = await api.POST("/api/v1/navigation/planet", {
    body,
    headers: getAuthHeaders(agentID),
  })
  if (error) throw error
  return data
}

export async function navigateWaypoint(
  agentID: string,
  body: SchemaNavigateWaypointRequest
) {
  const { data, error } = await api.POST("/api/v1/navigation/waypoint", {
    body,
    headers: getAuthHeaders(agentID),
  })
  if (error) throw error
  return data
}
