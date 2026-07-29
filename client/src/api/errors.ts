// export async function extractApiError(
//   response: Response
// ): Promise<SchemaErrorResponse> {
//   const status = response.status
//   let message = response.statusText
//   let code: string | undefined

//   try {
//     const body = (await response.json()) as {
//       message?: string
//       code?: string
//     }
//     if (body.message) message = body.message
//     if (body.code) code = body.code
//   } catch {
//     // body may not be JSON
//   }

//   return { status, message, code }
// }

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

export const ERROR_CODES = {
  INTERNAL_SERVER_ERROR: "Internal server error",
  ANOMALY: "Anomaly: something unexpected happened",
  UNKNOWN: "Unknown error",
  DECODE_ERROR: "Decode error",
  VALIDATION_ERROR: "Validation error",
  INVALID_JWT_TOKEN: "Invalid JWT token",
  USER_USERNAME_ALREADY_OCCUPIED: "This username is already occupied",
  USER_NOT_FOUND: "User not found",
  INVALID_CREDENTIALS: "Invalid credentials",
  ACCESS_DENIED: "Access denied",
  AGENT_NOT_FOUND: "Agent not found",
  AGENT_USERNAME_ALREADY_OCCUPIED: "Agent's username is already occupied",
  INVALID_AGENT_TOKEN: "Invalid agent token",
  AGENT_LIMIT_EXCEEDED: "Agent limit exceeded",
  RADAR_AREA_TOO_LARGE: "Radar area is too large",
  SHIP_NOT_FOUND: "Ship not found",
  INVENTORY_NOT_FOUND: "Inventory not found",
  RESOURCE_NOT_FOUND: "Resource not found",
  ITEM_NOT_FOUND: "Item not found",
  SHIP_MUST_BE_ACTIVE: "Ship must be active",
  INVALID_TRANSFER_DIRECTION: "Invalid transfer direction",
  NOT_ENOUGH_RESOURCES: "Not enough resources",
  INVENTORY_IS_FULL: "Inventory is full",
  ITEM_NOT_IN_INVENTORY: "Item is not in inventory",
  AGENT_IN_COOLDOWN: "Agent is in cooldown",
  INVALID_WARP_PATH: "Invalid warp path",
  ALREADY_AT_DESTINATION: "Already at destination",
  INVALID_COORDINATES: "Invalid coordinates",
  INVALID_SHIP_STATE: "Invalid ship state",
  SHIP_ALREADY_IN_THIS_STATE: "Ship is already in this state",
  CANNOT_DOCK_HERE: "Can't dock here",
  INVALID_UUID: "Invalid uuid",
}
