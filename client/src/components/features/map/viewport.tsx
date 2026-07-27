import { useApplication, extend, type PixiReactElementProps } from "@pixi/react"
import { Viewport as PixiViewport } from "pixi-viewport"
import { useEffect, useRef, type ReactNode } from "react"

extend({ PixiViewport })

declare module "@pixi/react" {
  interface PixiElements {
    viewport: PixiReactElementProps<typeof PixiViewport>
  }
}

export function ViewportScene({
  children,
}: {
  children: ReactNode | ReactNode[]
}) {
  const { app, isInitialised } = useApplication()

  const viewportRef = useRef<PixiViewport>(null)
  useEffect(() => {
    if (!viewportRef.current || !isInitialised) return

    const viewport = viewportRef.current

    viewport.plugins.removeAll()

    viewport.drag().pinch().wheel().decelerate().clampZoom({
      minScale: 0.2,
      maxScale: 4,
    })

    const onResize = () => {
      viewport.resize(
        app.screen.width,
        app.screen.height,
        viewport.worldWidth,
        viewport.worldHeight
      )
    }

    window.addEventListener("resize", onResize)

    // сразу после создания
    onResize()

    return () => {
      window.removeEventListener("resize", onResize)
    }
  }, [app, isInitialised])

  if (!isInitialised) {
    return null
  }

  return (
    <viewport
      ref={viewportRef}
      screenHeight={app.screen.height}
      screenWidth={app.screen.width}
      worldHeight={100000}
      worldWidth={100000}
      events={app.renderer.events}
    >
      {children}
    </viewport>
  )
}
