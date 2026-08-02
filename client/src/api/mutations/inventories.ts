import { api, getAuthHeaders } from "@/api/client"
import type {
  SchemaTransferItemsRequest,
  SchemaTransferResourcesRequest,
} from "@/api/types"

export async function transferResources(
  agentID: string,
  body: SchemaTransferResourcesRequest
) {
  const { error } = await api.POST("/api/v1/inventories/transfer-resources", {
    body,
    headers: getAuthHeaders(agentID),
  })
  if (error) throw error
}

export async function transferItems(
  agentID: string,
  body: SchemaTransferItemsRequest
) {
  const { error } = await api.POST("/api/v1/inventories/transfer-items", {
    body,
    headers: getAuthHeaders(agentID),
  })
  if (error) throw error
}
