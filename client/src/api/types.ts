import type { components } from "@/api/schema"

export type SchemaAgent = components["schemas"]["Agent"]
export type SchemaUser = components["schemas"]["User"]
export type SchemaShip = components["schemas"]["Ship"]
export type SchemaCooldown = components["schemas"]["Cooldown"]
export type SchemaSystem = components["schemas"]["System"]
export type SchemaPlanet = components["schemas"]["Planet"]
export type SchemaWaypoint = components["schemas"]["Waypoint"]
export type SchemaFullInventory = components["schemas"]["FullInventory"]
export type SchemaResource = components["schemas"]["Resource"]
export type SchemaItem = components["schemas"]["Item"]
export type SchemaLoginUserRequest = components["schemas"]["LoginUserRequest"]
export type SchemaLoginUserResponse = components["schemas"]["LoginUserResponse"]
export type SchemaRegisterUserRequest =
  components["schemas"]["RegisterUserRequest"]
export type SchemaRegisterAgentRequest =
  components["schemas"]["RegisterAgentRequest"]
export type SchemaRegisterAgentResponse =
  components["schemas"]["RegisterAgentResponse"]
export type SchemaRenameMyShipRequest =
  components["schemas"]["RenameMyShipRequest"]
export type SchemaNavigateWarpRequest =
  components["schemas"]["NavigateWarpRequest"]
export type SchemaNavigatePlanetRequest =
  components["schemas"]["NavigatePlanetRequest"]
export type SchemaNavigateWaypointRequest =
  components["schemas"]["NavigateWaypointRequest"]
export type SchemaTransferResourcesRequest =
  components["schemas"]["TransferResourcesRequest"]
export type SchemaTransferItemsRequest =
  components["schemas"]["TransferItemsRequest"]
export type SchemaNavigationResponse =
  components["schemas"]["NavigationResponse"]
export type SchemaResetAgentTokenResponse =
  components["schemas"]["ResetAgentTokenResponse"]
export type SchemaErrorResponse = components["schemas"]["ErrorResponse"]
