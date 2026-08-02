import { api, getAuthHeaders } from "@/api/client"
import type { SchemaRenameMyShipRequest } from "@/api/types"

export async function renameShip(
  agentID: string,
  id: string,
  body: SchemaRenameMyShipRequest
) {
  const { data, error } = await api.PATCH("/api/v1/ships/my/{id}/rename", {
    params: { path: { id } },
    body,
    headers: getAuthHeaders(agentID),
  })
  if (error) throw error
  return data
}

export async function changeActiveShip(agentID: string, id: string) {
  const { error } = await api.POST("/api/v1/ships/my/{id}/active", {
    params: { path: { id } },
    headers: getAuthHeaders(agentID),
  })
  if (error) throw error
}

export async function dockShip(agentID: string) {
  const { data, error } = await api.POST("/api/v1/ships/my/active/dock", {
    headers: getAuthHeaders(agentID),
  })
  if (error) throw error
  return data
}

export async function orbitShip(agentID: string) {
  const { data, error } = await api.POST("/api/v1/ships/my/active/orbit", {
    headers: getAuthHeaders(agentID),
  })
  if (error) throw error
  return data
}
