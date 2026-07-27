import { Application, extend } from "@pixi/react"

import {
  Container,
  FederatedPointerEvent,
  Graphics,
  Point,
  Rectangle,
  BitmapText,
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
import {
  type SchemaAgent,
  type SchemaSystem,
  type SystemExtended,
} from "@/api/types"
import { useActiveShipQuery, useShipRadarQuery } from "@/api/hooks"
import { useAuth } from "../auth/auth-provider"
import { toast } from "@/components/ui/toast"
import { useAgentsWithShips } from "../auth/use-agents-with-ships"

extend({
  Container,
  Graphics,
  Viewport,
  BitmapText,
  Text,
})

const CELL_SIZE: number = 100

interface GalaxyMapProps {
  onSystemClick: (system: SystemExtended) => void
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

  const { currentAgentID, isReady, agents } = useAuth()

  const { data, isPending, isError } = useShipRadarQuery(
    currentAgentID || undefined
  )

  const {
    data: ship,
    isPending: isShipPending,
    isError: isShipError,
  } = useActiveShipQuery(currentAgentID || undefined)

  const {
    data: agentsWithShips,
    isPending: isAgentsPending,
    isError: isAgentsError,
  } = useAgentsWithShips()

  const systems: SchemaSystem[] = data?.data ?? []

  const agentLocations: Map<string, SchemaAgent[]> = new Map()
  for (const a of agentsWithShips) {
    const key = `${a.ship.system_x}:${a.ship.system_y}`
    const l = agentLocations.get(key) ?? []
    l.push(a.agent)
    agentLocations.set(key, l)
    console.log(key, l)
  }
  console.log(agentLocations)

  const extendedSystems: SystemExtended[] = []
  systems.forEach((s) => {
    const agentsInSystem = agentLocations.get(`${s.x}:${s.y}`) || []

    extendedSystems.push({
      system: s,
      agents: agentsInSystem,
    })
  })

  const worldRef = useRef<Container>(null)

  const drawCallback = useCallback(
    (g: Graphics) => {
      g.clear()

      g.circle(0, 0, 100)
      for (const system of systems) {
        g.circle(system.x * CELL_SIZE, system.y * CELL_SIZE, 15).fill()
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
    if (isAgentsError) {
      toast.add({
        type: "error",
        title: "Error loading agents data",
      })
    }
  }, [isError, isShipError, isAgentsError, toast])

  const ready =
    isReady &&
    !isPending &&
    !isShipPending &&
    !isError &&
    !isShipError &&
    currentAgentID &&
    agents &&
    ship &&
    !isAgentsError &&
    !isAgentsPending

  const hitArea = new Rectangle(-100000, -100000, 200000, 200000)

  const handleClick = (e: FederatedPointerEvent) => {
    const pos = e.getLocalPosition(worldRef.current!)

    const clicked = extendedSystems.find((system) => {
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
            <pixiContainer>
              {extendedSystems.map((s) => {
                if (s.agents.length == 0) {
                  return
                }
                const parts: string[] = []
                for (const agent of s.agents) {
                  if (parts.length < 2) {
                    parts.push(agent.username)
                  } else {
                    parts.push("...")
                    break
                  }
                }

                const label = parts.join(", ")
                return (
                  <pixiText
                    key={`${s.system.x} ${s.system.y}`}
                    x={s.system.x * CELL_SIZE}
                    y={s.system.y * CELL_SIZE - 35}
                    text={label}
                    style={{
                      stroke: { color: "black", width: 2 },
                      fill: "white",
                      fontFamily: ["Jetbrains Mono Variable", "sans-serif"],
                      fontSize: 1000,
                    }}
                    anchor={0.5}
                    scale={0.03}
                  />
                )
              })}
            </pixiContainer>
          </pixiContainer>
        </ViewportScene>
      ) : null}
    </Application>
  )
}
