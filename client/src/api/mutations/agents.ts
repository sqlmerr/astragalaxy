import { api, getAuthHeaders } from "@/api/client"
import type {
  SchemaCraftRequest,
  SchemaRegisterAgentRequest,
} from "@/api/types"

export async function registerAgent(body: SchemaRegisterAgentRequest) {
  const { data, error } = await api.POST("/api/v1/agents", {
    body,
    headers: getAuthHeaders(),
  })
  if (error) throw error
  return data
}

export async function craftAction(agentId: string, body: SchemaCraftRequest) {
  const { data, error } = await api.POST(
    "/api/v1/agents/current/actions/craft",
    {
      body,
      headers: getAuthHeaders(agentId),
    }
  )
  if (error) throw error
  return data
}
