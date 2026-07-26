import { api } from "@/api/client"
import type { SchemaLoginUserRequest, SchemaLoginUserResponse } from "@/api/types"

export async function loginWithPassword(
  credentials: SchemaLoginUserRequest,
): Promise<SchemaLoginUserResponse> {
  const { data, error } = await api.POST("/api/v1/auth/login", {
    body: credentials,
  })
  if (error) throw error
  return data as SchemaLoginUserResponse
}
