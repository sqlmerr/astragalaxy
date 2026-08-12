import { api, getAuthHeaders } from "@/api/client"
import type {
  SchemaTransferItemsRequest,
  SchemaTransferResourcesRequest,
} from "@/api/types"

export async function transferResources(
  agentId: string,
  body: SchemaTransferResourcesRequest
) {
  const { error } = await api.POST("/api/v1/inventories/transfer-resources", {
    body,
    headers: getAuthHeaders(agentId),
  })
  if (error) throw error
}

export async function transferItems(
  agentId: string,
  body: SchemaTransferItemsRequest
) {
  const { error } = await api.POST("/api/v1/inventories/transfer-items", {
    body,
    headers: getAuthHeaders(agentId),
  })
  if (error) throw error
}

export async function useItem(agentId: string, itemId: string) {
  const { data, error } = await api.POST(
    "/api/v1/inventories/my/items/{id}/use",
    {
      params: { path: { id: itemId } },
      headers: getAuthHeaders(agentId),
    }
  )
  if (error) throw error
  return data
}
