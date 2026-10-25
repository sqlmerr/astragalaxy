import { useApplication, extend } from "@pixi/react"
import type { PixiReactElementProps } from "@pixi/react"
import { Viewport as PixiViewport } from "pixi-viewport"
import { useEffect } from "react"
import type { ReactNode, RefObject } from "react"

extend({ PixiViewport })

declare module "@pixi/react" {
  interface PixiElements {
    viewport: PixiReactElementProps<typeof PixiViewport>
  }
}

interface ViewportSceneProps {
  viewportRef: RefObject<PixiViewport | null>
  x: number
  y: number
  children: ReactNode | ReactNode[]
}

export function ViewportScene({
  viewportRef,
  x,
  y,
  children,
}: ViewportSceneProps) {
  const { app, isInitialised } = useApplication()

  useEffect(() => {
    const viewport = viewportRef.current
    if (!viewport || !isInitialised) return

    viewport.plugins.removeAll()

    viewport
      .drag()
      .pinch()
      .decelerate()
      .clampZoom({
        minScale: 0.2,
        maxScale: 4,
      })

    const onResize = () => {
      if (!app.renderer) return

      viewport.resize(
        app.screen.width,
        app.screen.height,
        viewport.worldWidth,
        viewport.worldHeight
      )
    }

    window.addEventListener("resize", onResize)
    onResize()
    viewport.moveCenter(x, y)

    return () => {
      window.removeEventListener("resize", onResize)
    }
  }, [app, isInitialised, viewportRef, x, y])

  useEffect(() => {
    const viewport = viewportRef.current
    if (!viewport || !isInitialised || !app.renderer) return

    const canvas = app.canvas

    const onWheel = (event: WheelEvent) => {
      if (event.ctrlKey || event.metaKey) {
        event.preventDefault()

        viewport.zoomPercent(-event.deltaY * 0.01, true)
        return
      }

      event.preventDefault()
      viewport.x -= event.deltaX
      viewport.y -= event.deltaY
    }

    canvas.addEventListener("wheel", onWheel, { passive: false })

    return () => {
      canvas.removeEventListener("wheel", onWheel)
    }
  }, [app, isInitialised, viewportRef])

  if (!isInitialised) return null

  const screenWidth = app.renderer ? app.screen.width : 0
  const screenHeight = app.renderer ? app.screen.height : 0
  const events = app.renderer?.events

  if (!events) return null

  return (
    <viewport
      ref={viewportRef}
      screenWidth={screenWidth}
      screenHeight={screenHeight}
      events={events}
      worldWidth={100000}
      worldHeight={100000}
    >
      {children}
    </viewport>
  )
}
