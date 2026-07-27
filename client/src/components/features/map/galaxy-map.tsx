import { Application, extend } from "@pixi/react"

import { Container, Graphics, Rectangle } from "pixi.js"
import { useCallback, useRef } from "react"
import { ViewportScene } from "./viewport"
import { Viewport } from "pixi-viewport"
import { type SchemaSystem } from "@/api/types"
import { useShipRadarQuery } from "@/api/hooks"

extend({
  Container,
  Graphics,
  Viewport,
})

const CELL_SIZE: number = 100

export function GalaxyMap() {
  // const { data, isPending, isError } = useShipRadarQuery()

  const systems: SchemaSystem[] = [
    {
      name: "System1",
      x: 1,
      y: 1,
      planets: [],
      waypoints: [],
    },
    {
      name: "System2",
      x: 1,
      y: 2,
      planets: [],
      waypoints: [],
    },
    {
      name: "System3",
      x: 2,
      y: 2,
      planets: [],
      waypoints: [],
    },
    {
      name: "System4",
      x: 2,
      y: 3,
      planets: [],
      waypoints: [],
    },
  ]

  const worldRef = useRef<Container>(null)

  const drawCallback = useCallback((g: Graphics) => {
    g.clear()

    // g.setFillStyle({ color: "red" })
    // g.circle(0, 0, 20)
    // g.fill()

    // g.setFillStyle({ color: "yellow" })
    // g.circle(500, 250, 15)
    // g.fill()

    for (const system of systems) {
      g.circle(system.x * CELL_SIZE, system.y * CELL_SIZE, 15)
      g.fill()
    }
  }, [])

  const hitArea = new Rectangle(-100000, -100000, 200000, 200000)

  return (
    <Application resizeTo={window} className="absoulute inset-0">
      <ViewportScene>
        {/* <ViewportScene /> */}
        <pixiContainer
          ref={worldRef}
          eventMode="static"
          // {...camera.events}
          hitArea={hitArea}
        >
          <pixiGraphics draw={drawCallback} />
        </pixiContainer>
      </ViewportScene>
    </Application>
  )
}
