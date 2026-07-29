import type { SchemaPlanet } from "@/api/types"
import { useTick } from "@pixi/react"
import type { Container, FederatedPointerEvent, Graphics } from "pixi.js"
import { useCallback, useEffect, useRef, useState, type ReactNode } from "react"
import { PLANET_PARAMS } from "../constants"
import { getOrbitRadius } from "../utils"

interface OrbitProps {
  planet: SchemaPlanet
  onClick?: (e: FederatedPointerEvent) => void
  isSelected?: boolean
}

export function OrbitPlanet({ planet, onClick, isSelected }: OrbitProps) {
  const params = PLANET_PARAMS[planet.type]
  const orbitRadius = getOrbitRadius(planet.orbit)
  const [coords, setCoords] = useState<{ x: number; y: number }>({
    x: orbitRadius,
    y: 0,
  })
  const ref = useRef<Container>(null)

  useEffect(() => {
    const angle = Math.random() * 360
    const x = Math.cos(angle) * orbitRadius
    const y = Math.sin(angle) * orbitRadius

    setCoords({ x, y })
  }, [planet])

  useTick((ticker) => {
    if (ref.current) {
      ref.current.rotation +=
        params.rotationSpeed * ticker.deltaTime * ((planet.orbit + 1) / 10)
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
      <pixiContainer
        x={coords.x}
        y={coords.y}
        onClick={onClick}
        cursor="pointer"
        eventMode="static"
      >
        <pixiGraphics draw={drawPlanet} />
      </pixiContainer>
    </pixiContainer>
  )
}
