import createClient from "openapi-fetch"

import type { paths } from "@/api/schema"

const AUTH_TOKEN_KEY = "astragalaxy.auth.token"
const API_URL = import.meta.env.VITE_API_URL
console.log(API_URL)

function getAuthToken(): string | null {
  if (typeof window === "undefined") return null
  return window.localStorage.getItem(AUTH_TOKEN_KEY)
}

export const api = createClient<paths>({
  baseUrl: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
})

export function getAuthHeaders(agentID?: string): Record<string, string> {
  const token = getAuthToken()
  if (!token) return {}
  const headers = { Authorization: `Bearer ${token}` }
  if (agentID) {
    return {
      ...headers,
      "X-Agent-ID": agentID,
    }
  }
  return headers
}
