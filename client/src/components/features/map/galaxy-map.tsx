import { Application, extend } from "@pixi/react"

import {
  Container,
  FederatedPointerEvent,
  Graphics,
  Point,
  Rectangle,
  Text,
} from "pixi.js"
import {
  useCallback,
  useEffect,
  useImperativeHandle,
  useRef,
  useState,
  type RefObject,
} from "react"
import { ViewportScene } from "./viewport"
import { Viewport } from "pixi-viewport"
import { type SchemaSystem } from "@/api/types"
import { useActiveShipQuery, useShipRadarQuery } from "@/api/hooks"
import { useAuth } from "../auth/auth-provider"
import { toast } from "@/components/ui/toast"

extend({
  Container,
  Graphics,
  Viewport,
  Text,
})

const CELL_SIZE: number = 100

interface GalaxyMapProps {
  onSystemClick: (system: SchemaSystem) => void
  ref: RefObject<GalaxyMapRef | null>
}

export interface GalaxyMapRef {
  centerOnSystem(system: SchemaSystem): void
  closeSystem(): void
}

export function GalaxyMap({ onSystemClick, ref }: GalaxyMapProps) {
  const [selectedSystem, setSelectedSystem] = useState<SchemaSystem | null>(
    null
  )
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

  const { currentAgentID, isReady, agents, currentAgent } = useAuth()

  const { data, isPending, isError } = useShipRadarQuery(
    currentAgentID || undefined
  )

  const {
    data: ship,
    isPending: isShipPending,
    isError: isShipError,
  } = useActiveShipQuery(currentAgentID || undefined)

  const systems: SchemaSystem[] = data?.data ?? []

  const worldRef = useRef<Container>(null)

  const drawCallback = useCallback(
    (g: Graphics) => {
      g.clear()

      g.circle(0, 0, 100)
      for (const system of systems) {
        g.circle(system.x * CELL_SIZE, system.y * CELL_SIZE, 15)
        g.fill()
        if (
          selectedSystem &&
          system.x == selectedSystem.x &&
          system.y == selectedSystem.y
        ) {
          g.circle(system.x * CELL_SIZE, system.y * CELL_SIZE, 20).stroke({
            width: 1,
          })
        }
      }
    },
    [systems, selectedSystem]
  )

  useEffect(() => {
    if (isError) {
      toast.add({
        type: "error",
        title: "Error loading radar data",
      })
    }
    if (isShipError) {
      toast.add({
        type: "error",
        title: "Error loading agent ship data",
      })
    }
  }, [isError, isShipError, toast])

  if (
    isPending ||
    isShipPending ||
    isError ||
    isShipError ||
    !currentAgentID ||
    agents === null ||
    ship == null
  ) {
    return null
  }

  const hitArea = new Rectangle(-100000, -100000, 200000, 200000)

  const ready =
    isReady &&
    !isPending &&
    !isShipPending &&
    !isError &&
    !isShipError &&
    currentAgentID &&
    agents &&
    ship

  const handleClick = (e: FederatedPointerEvent) => {
    const pos = e.getLocalPosition(worldRef.current!)

    const clicked = systems.find((system) => {
      const dx = pos.x - system.x * CELL_SIZE
      const dy = pos.y - system.y * CELL_SIZE

      return dx * dx + dy * dy <= 15 * 15
    })

    if (clicked) {
      onSystemClick(clicked)
      setSelectedSystem(clicked)
    }
  }

  return (
    <Application resizeTo={window} className="absoulute inset-0">
      {ready ? (
        <ViewportScene
          viewportRef={viewportRef}
          x={ship.system_x * CELL_SIZE}
          y={ship.system_y * CELL_SIZE}
        >
          {/* <ViewportScene /> */}
          <pixiContainer
            ref={worldRef}
            eventMode="static"
            // {...camera.events}
            hitArea={hitArea}
            onClick={handleClick}
          >
            <pixiGraphics draw={drawCallback} />
            {/* <pixiText text={"Hello"} /> */}
          </pixiContainer>
        </ViewportScene>
      ) : null}
    </Application>
  )
}
