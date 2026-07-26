import { api, getAuthHeaders } from "@/api/client"
import type {
  SchemaNavigatePlanetRequest,
  SchemaNavigateWarpRequest,
  SchemaNavigateWaypointRequest,
} from "@/api/types"

export async function navigateWarp(body: SchemaNavigateWarpRequest) {
  const { data, error } = await api.POST("/api/v1/navigation/warp", {
    body,
    headers: getAuthHeaders(),
  })
  if (error) throw error
  return data
}

export async function navigatePlanet(body: SchemaNavigatePlanetRequest) {
  const { data, error } = await api.POST("/api/v1/navigation/planet", {
    body,
    headers: getAuthHeaders(),
  })
  if (error) throw error
  return data
}

export async function navigateWaypoint(body: SchemaNavigateWaypointRequest) {
  const { data, error } = await api.POST("/api/v1/navigation/waypoint", {
    body,
    headers: getAuthHeaders(),
  })
  if (error) throw error
  return data
}
