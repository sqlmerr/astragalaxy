import { extend } from "@pixi/react"

import {
  Container,
  FederatedPointerEvent,
  Graphics,
  Point,
  BitmapText,
  Text,
} from "pixi.js"
import {
  useCallback,
  useImperativeHandle,
  useRef,
  useState,
  type RefObject,
} from "react"
import { Viewport } from "pixi-viewport"
import {
  type SchemaShip,
  type SchemaShortSystem,
  type ShortSystemExtended,
} from "@/api/types"
import { CELL_SIZE, SYSTEM_PARAMS } from "../constants"
import { Labels } from "./labels"
import { MapCanvas } from "../map-canvas"

extend({
  Container,
  Graphics,
  Viewport,
  BitmapText,
  Text,
})

interface GalaxyMapProps {
  onSystemClick: (system: ShortSystemExtended) => void
  ref: RefObject<GalaxyMapRef | null>
  ship: SchemaShip
  systems: ShortSystemExtended[]
}

export interface GalaxyMapRef {
  centerOnSystem(system: SchemaShortSystem): void
  closeSystem(): void
}

export function GalaxyMap({
  onSystemClick,
  ref,
  ship,
  systems,
}: GalaxyMapProps) {
  const [selectedSystem, setSelectedSystem] =
    useState<SchemaShortSystem | null>(null)
  const viewportRef = useRef<Viewport>(null)

  useImperativeHandle(ref, () => ({
    centerOnSystem(system) {
      viewportRef.current?.animate({
        position: new Point(system.x * CELL_SIZE, system.y * CELL_SIZE),
        time: 500,
        ease: "easeInOutSine",
      })
    },
    closeSystem() {
      setSelectedSystem(null)
    },
  }))

  const worldRef = useRef<Container>(null)

  const drawCallback = useCallback(
    (g: Graphics) => {
      g.clear()

      g.circle(0, 0, 100)
      for (const system of systems) {
        const params = SYSTEM_PARAMS[system.system.archetype]
        g.circle(
          system.system.x * CELL_SIZE,
          system.system.y * CELL_SIZE,
          15
        ).fill({ color: params.color })
        if (
          selectedSystem &&
          system.system.x == selectedSystem.x &&
          system.system.y == selectedSystem.y
        ) {
          g.circle(
            system.system.x * CELL_SIZE,
            system.system.y * CELL_SIZE,
            20
          ).stroke({
            width: 1,
          })
        }
      }
    },
    [systems, selectedSystem]
  )

  const handleClick = (e: FederatedPointerEvent) => {
    const pos = e.getLocalPosition(worldRef.current!)

    const clicked = systems.find((system) => {
      const dx = pos.x - system.system.x * CELL_SIZE
      const dy = pos.y - system.system.y * CELL_SIZE

      return dx * dx + dy * dy <= 15 * 15
    })

    if (clicked) {
      onSystemClick(clicked)
      setSelectedSystem(clicked.system)
    }
  }

  return (
    <MapCanvas
      x={ship.system_x}
      y={ship.system_y}
      worldRef={worldRef}
      viewportRef={viewportRef}
      onClick={handleClick}
    >
      <pixiGraphics draw={drawCallback} />
      <Labels systems={systems} />
    </MapCanvas>
  )
}
