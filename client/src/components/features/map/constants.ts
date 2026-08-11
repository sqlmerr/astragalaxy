import type { SchemaPlanet, SchemaSystem, SchemaWaypoint } from "@/api/types"

export const CELL_SIZE = 100

export const PLANET_PARAMS: Record<
  SchemaPlanet["type"],
  {
    color: number
    radius: number
    rotationSpeed: number
  }
> = {
  terra: {
    color: 0x4caf50,
    radius: 18,
    rotationSpeed: 0.0035,
  },
  ocean: {
    color: 0x1976d2,
    radius: 20,
    rotationSpeed: 0.0028,
  },
  scorched: {
    color: 0xd84315,
    radius: 14,
    rotationSpeed: 0.006,
  },
  glacial: {
    color: 0xe3f2fd,
    radius: 17,
    rotationSpeed: 0.0018,
  },
  toxic: {
    color: 0x7cb342,
    radius: 16,
    rotationSpeed: 0.0045,
  },
}

export const SYSTEM_PARAMS: Record<
  SchemaSystem["archetype"],
  { color: number }
> = {
  habitable: {
    color: 0x87ae73,
  },
  dead: {
    color: 0xbb0a1e,
  },
  frozen: {
    color: 0x7ec0ee,
  },
}

interface ShapePolygon {
  type: "polygon"
  sides: number // minimum 3 sides for triangle
  radius: number
}

interface ShapeCustom {
  type: "custom"
}

export const WAYPOINT_PARAMS: Record<
  SchemaWaypoint["type"],
  { name: string; shape: ShapePolygon | ShapeCustom; color?: number }
> = {
  asteroid: {
    name: "Asteroid",
    shape: { type: "polygon", sides: 4, radius: 13 }, // TODO: custom waypoint shapes
    color: 0x4a4f6b,
  },
  station: {
    name: "Station",
    shape: { type: "polygon", sides: 5, radius: 15 },
    color: 0xffff00,
  },
}
