import { api, getAuthHeaders } from "@/api/client"
import type { SchemaRegisterAgentRequest } from "@/api/types"

export async function registerAgent(body: SchemaRegisterAgentRequest) {
  const { data, error } = await api.POST("/api/v1/agents", {
    body,
    headers: getAuthHeaders(),
  })
  if (error) throw error
  return data
}
