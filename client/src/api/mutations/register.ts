import { api, getAuthHeaders } from "@/api/client"
import type { SchemaRegisterUserRequest } from "@/api/types"

export async function registerUser(credentials: SchemaRegisterUserRequest) {
  const { data, error } = await api.POST("/api/v1/auth/register", {
    body: credentials,
    headers: getAuthHeaders(),
  })
  if (error) throw error
  return data
}
