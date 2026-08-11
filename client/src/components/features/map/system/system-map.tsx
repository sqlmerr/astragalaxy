import { useCallback, useRef } from "react"
import { MapCanvas } from "../map-canvas"
import type { Viewport } from "pixi-viewport"
import type { FederatedPointerEvent, Graphics, Container } from "pixi.js"
import type {
  SchemaPlanet,
  SchemaWaypoint,
  SystemExtended,
  AgentExtended,
} from "@/api/types"
import { OrbitPlanet } from "./orbit"
import { Waypoint } from "./waypoint"

interface SystemMapProps {
  system: SystemExtended
  agents: AgentExtended[]
  selectedPlanet: SchemaPlanet | null
  setSelectedPlanet: (p: SchemaPlanet | null) => void
  selectedWaypoint: SchemaWaypoint | null
  setSelectedWaypoint: (w: SchemaWaypoint | null) => void
}
export function SystemMap({
  system,
  agents,
  selectedPlanet,
  setSelectedPlanet,
  selectedWaypoint,
  setSelectedWaypoint,
}: SystemMapProps) {
  const viewportRef = useRef<Viewport>(null)
  const worldRef = useRef<Container>(null)

  const drawStar = useCallback((g: Graphics) => {
    g.clear()

    g.circle(0, 0, 40).fill({ color: 0xffffff })
  }, [])

  const handlePlanetClick = (orbit: number) => () => {
    const planet = system.system.planets.find((p) => p.orbit === orbit)

    if (planet) {
      setSelectedPlanet(planet)
      setSelectedWaypoint(null)
    }
  }

  const handleWaypointClick = (id: number) => (e: FederatedPointerEvent) => {
    const waypoint = system.system.waypoints.find((w) => w.id === id)
    if (waypoint) {
      setSelectedWaypoint(waypoint)
      setSelectedPlanet(null)
    }
  }

  const agentsInThisSystem = agents.filter(
    (a) =>
      a.ship.system_x === system.system.x && a.ship.system_y === system.system.y
  )

  return (
    <MapCanvas x={0} y={0} worldRef={worldRef} viewportRef={viewportRef}>
      <pixiGraphics draw={drawStar} />
      {system.system.planets.map((p) => {
        return (
          <OrbitPlanet
            key={p.orbit}
            planet={p}
            onClick={handlePlanetClick(p.orbit)}
            isSelected={selectedPlanet?.orbit === p.orbit}
            agents={agentsInThisSystem}
          />
        )
      })}
      {system.system.waypoints.map((w) => (
        <Waypoint
          key={w.id}
          waypoint={w}
          system={system.system}
          onClick={handleWaypointClick(w.id)}
          isSelected={selectedWaypoint?.id === w.id}
          agents={agentsInThisSystem}
        />
      ))}
    </MapCanvas>
  )
}
