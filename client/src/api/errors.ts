export interface ApiError {
  status: number
  message: string
  code?: string
}

export async function extractApiError(response: Response): Promise<ApiError> {
  const status = response.status
  let message = response.statusText
  let code: string | undefined

  try {
    const body = (await response.json()) as {
      message?: string
      code?: string
    }
    if (body.message) message = body.message
    if (body.code) code = body.code
  } catch {
    // body may not be JSON
  }

  return { status, message, code }
}

function hasMessage(error: unknown): error is { message: string } {
  return (
    typeof error === "object" &&
    error !== null &&
    "message" in error &&
    typeof (error as Record<string, unknown>).message === "string"
  )
}

export function getErrorMessage(error: unknown): string {
  if (hasMessage(error)) {
    return error.message
  }
  return "An unexpected error occurred"
}
