import type { AgentExtended, SchemaPlanet } from "@/api/types"
import { useTick } from "@pixi/react"
import type { Container, FederatedPointerEvent, Graphics } from "pixi.js"
import { useCallback, useEffect, useRef, useState, type ReactNode } from "react"
import { PLANET_PARAMS } from "../constants"
import { getOrbitRadius } from "../utils"
import { degreesToRadians } from "@/lib/utils"
import { SystemLabel } from "./label"

interface OrbitProps {
  planet: SchemaPlanet
  onClick?: (e: FederatedPointerEvent) => void
  isSelected?: boolean
  agents: AgentExtended[]
}

export function OrbitPlanet({
  planet,
  onClick,
  isSelected,
  agents,
}: OrbitProps) {
  const params = PLANET_PARAMS[planet.type]
  const orbitRadius = getOrbitRadius(planet.orbit)
  const [coords, setCoords] = useState<{ x: number; y: number }>({
    x: orbitRadius,
    y: 0,
  })
  const ref = useRef<Container>(null)
  const angle = useRef(degreesToRadians(Math.random() * 360))

  useEffect(() => {
    const x = Math.cos(angle.current) * orbitRadius
    const y = Math.sin(angle.current) * orbitRadius

    setCoords({ x, y })
  }, [planet])

  useTick((ticker) => {
    if (ref.current) {
      angle.current +=
        params.rotationSpeed * ticker.deltaTime * ((planet.orbit + 1) / 10)
      setCoords({
        x: Math.cos(angle.current) * orbitRadius,
        y: Math.sin(angle.current) * orbitRadius,
      })
    }
  })

  const drawOrbitOutline = useCallback(
    (g: Graphics) => {
      g.clear()
      g.circle(0, 0, orbitRadius).stroke({ alpha: 0.25, color: "white" })
    },
    [planet]
  )

  const drawPlanet = useCallback(
    (g: Graphics) => {
      g.clear()
      g.circle(0, 0, params.radius).fill({ color: params.color })

      if (isSelected) {
        g.circle(0, 0, params.radius * 1.5).stroke({ color: "white" })
      }
    },
    [planet, isSelected]
  )

  return (
    <pixiContainer ref={ref}>
      <pixiGraphics draw={drawOrbitOutline} />
      <pixiContainer x={coords.x} y={coords.y}>
        <pixiContainer cursor="pointer" eventMode="static" onClick={onClick}>
          <pixiGraphics draw={drawPlanet} />
        </pixiContainer>

        <SystemLabel
          x={0}
          y={-35}
          agents={agents}
          location="PLANET"
          locationId={planet.orbit}
        />
      </pixiContainer>
    </pixiContainer>
  )
}
