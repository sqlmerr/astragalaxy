import { api, getAuthHeaders } from "../client"
import type { SchemaPayRequest } from "../types"

export async function pay(
  agentId: string,
  address: string,
  body: SchemaPayRequest
) {
  const { data, error } = await api.POST("/api/v1/wallets/my/{addr}/pay", {
    params: {
      path: {
        addr: address,
      },
    },
    body,
    headers: getAuthHeaders(agentId),
  })
  if (error) throw error
  return data
}
