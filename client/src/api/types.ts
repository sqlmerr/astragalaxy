import type { components } from "@/api/schema"

export type SchemaAgent = components["schemas"]["Agent"]
export type SchemaUser = components["schemas"]["User"]
export type SchemaShip = components["schemas"]["Ship"]
export type SchemaCooldown = components["schemas"]["Cooldown"]
export type SchemaShortSystem = components["schemas"]["ShortSystem"]
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
export type SchemaMineAsteroidRequest =
  components["schemas"]["MineAsteroidRequest"]
export type SchemaMinePlanetRequest = components["schemas"]["MinePlanetRequest"]
export type SchemaNavigationResponse =
  components["schemas"]["NavigationResponse"]
export type SchemaResetAgentTokenResponse =
  components["schemas"]["ResetAgentTokenResponse"]
export type SchemaErrorResponse = components["schemas"]["ErrorResponse"]
export type SchemaResourceDeposit = components["schemas"]["ResourceDeposit"]
export type SchemaUseItemResponse = components["schemas"]["UseItemResponse"]
export type SchemaInventoryOwner = components["schemas"]["Inventory"]
export type SchemaCraftRequest = components["schemas"]["CraftRequest"]
export type SchemaCraftResponse = components["schemas"]["CraftResponse"]
export type SchemaRecipe = components["schemas"]["Recipe"]

export interface AgentExtended {
  agent: SchemaAgent
  ship: SchemaShip // active ship
  ships: SchemaShip[]
  cooldown: SchemaCooldown
}

export interface SystemExtended {
  system: SchemaSystem
  agents: AgentExtended[]
}

export interface ShortSystemExtended {
  system: SchemaShortSystem
  agents: AgentExtended[]
}

export type AnySystemExtended = SystemExtended | ShortSystemExtended

export function isSystemExtended(
  system: AnySystemExtended | null
): system is SystemExtended {
  return system !== null && "planets" in system.system
}
