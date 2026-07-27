import type { Container, FederatedPointerEvent } from "pixi.js"
import { useRef, type RefObject } from "react"

export function useCamera(worldRef: RefObject<Container | null>) {
  const dragging = useRef(false)

  const lastPointer = useRef({
    x: 0,
    y: 0,
  })

  const camera = useRef({
    x: 0,
    y: 0,
    zoom: 1,
  })

  const updateCamera = () => {
    if (!worldRef.current) return

    worldRef.current.position.set(camera.current.x, camera.current.y)

    worldRef.current.scale.set(camera.current.zoom)
  }

  const onPointerDown = (e: FederatedPointerEvent) => {
    dragging.current = true

    lastPointer.current = {
      x: e.global.x,
      y: e.global.y,
    }
  }

  const onPointerMove = (e: FederatedPointerEvent) => {
    if (!dragging.current) return

    const dx = e.global.x - lastPointer.current.x
    const dy = e.global.y - lastPointer.current.y

    camera.current.x += dx
    camera.current.y += dy

    updateCamera()

    lastPointer.current = {
      x: e.global.x,
      y: e.global.y,
    }
  }

  const onPointerUp = () => {
    dragging.current = false
  }

  return {
    events: {
      onPointerDown,
      onPointerMove,
      onPointerUp,
      onPointerUpOutsite: onPointerUp,
    },
  }
}
