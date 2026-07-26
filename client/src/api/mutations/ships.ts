import { api, getAuthHeaders } from "@/api/client"
import type { SchemaRenameMyShipRequest } from "@/api/types"

export async function renameShip(id: string, body: SchemaRenameMyShipRequest) {
  const { data, error } = await api.PATCH("/api/v1/ships/my/{id}/rename", {
    params: { path: { id } },
    body,
    headers: getAuthHeaders(),
  })
  if (error) throw error
  return data
}

export async function changeActiveShip(id: string) {
  const { error } = await api.POST("/api/v1/ships/my/{id}/active", {
    params: { path: { id } },
    headers: getAuthHeaders(),
  })
  if (error) throw error
}

export async function dockShip() {
  const { data, error } = await api.POST("/api/v1/ships/my/active/dock", {
    headers: getAuthHeaders(),
  })
  if (error) throw error
  return data
}

export async function orbitShip() {
  const { data, error } = await api.POST("/api/v1/ships/my/active/orbit", {
    headers: getAuthHeaders(),
  })
  if (error) throw error
  return data
}
