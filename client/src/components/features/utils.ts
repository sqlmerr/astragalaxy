import type { SchemaErrorResponse } from "@/api/types"
import { toast } from "../ui/toast"
import { ERROR_CODES } from "@/api/errors"

export function handleError(
  error: SchemaErrorResponse | object,
  title?: string
): SchemaErrorResponse | undefined {
  if (!("code" in error && "error" in error && "message" in error)) {
    return
  }

  const message: string = ERROR_CODES[error.code as keyof typeof ERROR_CODES]
  toast.add({
    type: "error",
    title: title || "Error",
    description: message,
  }) // TODO: click on toast to show full error json

  return error
}
