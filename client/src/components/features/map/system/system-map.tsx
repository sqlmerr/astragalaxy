import {
  useCallback,
  useImperativeHandle,
  useRef,
  useState,
  type RefObject,
} from "react"
import { MapCanvas } from "../map-canvas"
import type { Viewport } from "pixi-viewport"
import { FederatedPointerEvent, type Container, type Graphics } from "pixi.js"
import type { SchemaPlanet, SystemExtended } from "@/api/types"
import type { AgentWithShip } from "@/api/types"
import { OrbitPlanet } from "./orbit"

interface SystemMapProps {
  system: SystemExtended
  agents: AgentWithShip[]
  onPlanetClick: (planet: SchemaPlanet) => void
  ref: RefObject<SystemMapRef | null>
}

export interface SystemMapRef {
  selectPlanet(planet: SchemaPlanet | null): void
}

export function SystemMap({
  system,
  agents,
  onPlanetClick,
  ref,
}: SystemMapProps) {
  const viewportRef = useRef<Viewport>(null)
  const worldRef = useRef<Container>(null)

  const [selectedPlanet, setSelectedPlanet] = useState<SchemaPlanet | null>(
    null
  )

  useImperativeHandle(ref, () => ({
    selectPlanet(planet: SchemaPlanet | null) {
      setSelectedPlanet(planet)
    },
  }))

  const drawCallback = useCallback((g: Graphics) => {
    g.clear()

    g.circle(0, 0, 40).fill({ color: 0xffffff })
  }, [])

  const handleClick = (orbit: number) => (e: FederatedPointerEvent) => {
    // e.preventDefault()
    const planet = system.system.planets.find((p) => p.orbit === orbit)
    if (planet) {
      setSelectedPlanet(planet)
      onPlanetClick(planet)
    }
  }

  return (
    <MapCanvas x={0} y={0} worldRef={worldRef} viewportRef={viewportRef}>
      <pixiGraphics draw={drawCallback} />
      {system.system.planets.map((p) => {
        return (
          <OrbitPlanet
            key={p.orbit}
            planet={p}
            onClick={handleClick(p.orbit)}
            isSelected={selectedPlanet?.orbit === p.orbit}
          />
        )
      })}
    </MapCanvas>
  )
}
