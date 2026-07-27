import { useApplication, extend, type PixiReactElementProps } from "@pixi/react"
import { Viewport as PixiViewport } from "pixi-viewport"
import { useEffect, useRef, type ReactNode, type RefObject } from "react"

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
    if (!viewportRef.current || !isInitialised) return

    const viewport = viewportRef.current

    viewport.plugins.removeAll()

    viewport.drag().pinch().wheel().decelerate().clampZoom({
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
  }, [app, isInitialised, x, y])

  if (!isInitialised) {
    return null
  }

  const screenWidth = app.renderer ? app.screen.width : 0
  const screenHeight = app.renderer ? app.screen.height : 0
  const events = app.renderer?.events

  if (!events) {
    return null
  }

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
