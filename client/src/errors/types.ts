import type { SchemaErrorResponse } from "@/api/types"

export type AppError = {
  title: string
  description?: string
  error: SchemaErrorResponse
  timestamp: Date
}
