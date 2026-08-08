import type { SchemaShip } from "./types"

export function shipLocationIs(
  ship: SchemaShip,
  location: {
    systemX: number
    systemY: number
    locationType: SchemaShip["location"]
    locationId?: SchemaShip["location_id"]
  }
): boolean {
  return (
    ship.system_x === location.systemX &&
    ship.system_y === location.systemY &&
    ship.location === location.locationType &&
    (location.locationId !== undefined
      ? ship.location_id === location.locationId
      : true)
  )
}
