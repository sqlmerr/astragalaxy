import type { AgentWithShip, SchemaSystem, SchemaWaypoint } from "@/api/types"
import type { Container, FederatedPointerEvent, Graphics } from "pixi.js"
import { useCallback, useEffect, useRef, useState } from "react"
import { getOrbitRadius, regularPolygon } from "../utils"
import { degreesToRadians, randomNumberBetween } from "@/lib/utils"
import { WAYPOINT_PARAMS } from "../constants"
import { useTick } from "@pixi/react"
import { SystemLabel } from "./label"

interface WaypointProps {
  waypoint: SchemaWaypoint
  system: SchemaSystem
  onClick: (e: FederatedPointerEvent) => void
  isSelected?: boolean
  agents: AgentWithShip[]
}

export function Waypoint({
  waypoint,
  system,
  onClick,
  isSelected,
  agents,
}: WaypointProps) {
  const params = WAYPOINT_PARAMS[waypoint.type]
  const ref = useRef<Container>(null)
  const [coords, setCoords] = useState({ x: 0, y: 0 })
  const [rotationSpeed] = useState(() => Math.random() / 1000)
  const [initialRotationAngle] = useState(() => Math.random() * 360)

  useEffect(() => {
    const angle = degreesToRadians(Math.random() * 360)
    const orbitCount = system.planets.length
    if (!orbitCount) {
      return
    }
    const afterOrbit = Math.floor(Math.random() * orbitCount)
    const orbitRadius = getOrbitRadius(afterOrbit)
    const nextOrbitRadius = getOrbitRadius(afterOrbit + 1)
    const m = (orbitRadius + nextOrbitRadius) / 2
    const spread = randomNumberBetween(-25, 25)
    const waypointRadius = m + spread

    const x = Math.cos(angle) * waypointRadius
    const y = Math.sin(angle) * waypointRadius
    setCoords({ x, y })
  }, [waypoint])

  const draw = useCallback(
    (g: Graphics) => {
      g.clear()
      switch (params.shape.type) {
        case "polygon": {
          const points = regularPolygon(
            0,
            0,
            params.shape.radius,
            params.shape.sides
          )
          g.poly(points).fill({ color: params.color })
          if (isSelected) {
            g.circle(0, 0, params.shape.radius * 1.5).stroke({ color: "white" })
          }
          break
        }
        case "custom": // TODO: custom waypoint locations
          break
      }
    },
    [waypoint, isSelected]
  )

  useTick((ticker) => {
    if (ref.current) {
      ref.current.rotation += rotationSpeed * ticker.deltaTime
    }
  })

  return (
    <pixiContainer x={coords.x} y={coords.y}>
      <pixiContainer
        ref={ref}
        rotation={degreesToRadians(initialRotationAngle)}
        eventMode="static"
        cursor="pointer"
        onClick={onClick}
      >
        <pixiGraphics draw={draw} />
      </pixiContainer>
      <SystemLabel
        x={0}
        y={-35}
        agents={agents}
        location="WAYPOINT"
        locationId={waypoint.id}
      />
    </pixiContainer>
  )
}
