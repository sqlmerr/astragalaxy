import { api, getAuthHeaders } from "@/api/client"
import type {
  SchemaTransferItemsRequest,
  SchemaTransferResourcesRequest,
} from "@/api/types"

export async function transferResources(body: SchemaTransferResourcesRequest) {
  const { error } = await api.POST("/api/v1/inventories/transfer-resources", {
    body,
    headers: getAuthHeaders(),
  })
  if (error) throw error
}

export async function transferItems(body: SchemaTransferItemsRequest) {
  const { error } = await api.POST("/api/v1/inventories/transfer-items", {
    body,
    headers: getAuthHeaders(),
  })
  if (error) throw error
}
