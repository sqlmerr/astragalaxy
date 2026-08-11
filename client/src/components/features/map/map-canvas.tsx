import { extend, Application } from "@pixi/react"
import { Viewport } from "pixi-viewport"
import type {
  FederatedPointerEvent,
  FederatedWheelEvent,FederatedEventHandler} from "pixi.js";
import {
  BitmapText,
  Container,
  Graphics,
  Rectangle
  
} from "pixi.js"
import { CELL_SIZE } from "./constants"
import { ViewportScene } from "./viewport"
import type { ReactNode, RefObject } from "react"

interface MapCanvasProps {
  x: number
  y: number
  worldRef: RefObject<Container | null>
  viewportRef: RefObject<Viewport | null>
  onClick?:
    | FederatedEventHandler<FederatedPointerEvent>
    | FederatedEventHandler<FederatedWheelEvent>
  children: ReactNode | ReactNode[]
}
extend({
  Container,
  Graphics,
  Viewport,
  BitmapText,
  Text,
})

export function MapCanvas({
  x,
  y,
  worldRef,
  viewportRef,
  onClick,
  children,
}: MapCanvasProps) {
  const hitArea = new Rectangle(-100000, -100000, 200000, 200000)

  return (
    <Application resizeTo={window} className="absoulute inset-0">
      <ViewportScene
        viewportRef={viewportRef}
        x={x * CELL_SIZE}
        y={y * CELL_SIZE}
      >
        <pixiContainer
          ref={worldRef}
          eventMode="static"
          hitArea={hitArea}
          onClick={onClick}
        >
          {children}
        </pixiContainer>
      </ViewportScene>
    </Application>
  )
}
