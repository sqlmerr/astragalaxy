import { api, getAuthHeaders } from "@/api/client"
import type {
  SchemaMineAsteroidRequest,
  SchemaMinePlanetRequest,
} from "@/api/types"

export async function mineAsteroid(
  agentID: string,
  body: SchemaMineAsteroidRequest
) {
  const { data, error } = await api.POST(
    "/api/v1/agents/current/actions/mine/asteroid",
    {
      body,
      headers: getAuthHeaders(agentID),
    }
  )
  if (error) throw error
  return data
}

export async function minePlanet(
  agentID: string,
  body: SchemaMinePlanetRequest
) {
  const { data, error } = await api.POST(
    "/api/v1/agents/current/actions/mine/planet",
    {
      body,
      headers: getAuthHeaders(agentID),
    }
  )
  if (error) throw error
  return data
}
