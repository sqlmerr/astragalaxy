import { toast } from "@/components/ui/toast"
import { ERROR_CODES } from "@/api/errors"
import type { AppError } from "./types"
import type { SchemaErrorResponse } from "@/api/types"
import { useError } from "./store"

export function isApiError(value: unknown): value is SchemaErrorResponse {
  return (
    typeof value === "object" &&
    value !== null &&
    "error" in value &&
    typeof value.error === "string" &&
    "code" in value &&
    typeof value.code === "string" &&
    "message" in value &&
    typeof value.message === "string"
  )
}

export function handleError(
  openError: (error: AppError) => void,
  error: SchemaErrorResponse | object,
  title?: string
): AppError | undefined {
  if (!isApiError(error)) {
    return
  }

  const message: string = ERROR_CODES[error.code as keyof typeof ERROR_CODES]
  const appError: AppError = {
    title: title || "Error",
    description: message,
    timestamp: new Date(Date.now()),
    error: error,
  }

  toast.add({
    type: "error",
    title: appError.title,
    description: appError.description,
    actionProps: {
      children: "Details",
      onClick() {
        openError(appError)
      },
    },
  })

  return appError
}

export function useErrorHandler() {
  const { open } = useError()

  return (error: SchemaErrorResponse | object, title?: string) =>
    handleError(open, error, title)
}
