import type { ShortSystemExtended } from "@/api/types"
import { CELL_SIZE } from "../constants"

export function Labels({ systems }: { systems: ShortSystemExtended[] }) {
  return (
    <pixiContainer>
      {systems.map((s) => {
        if (s.agents.length == 0) {
          return
        }
        const parts: string[] = []
        for (const agent of s.agents) {
          if (parts.length < 2) {
            parts.push(agent.agent.username)
          } else {
            parts.push("...")
            break
          }
        }

        const label = parts.join(", ")
        return (
          <pixiBitmapText
            key={`${s.system.x} ${s.system.y}`}
            x={s.system.x * CELL_SIZE}
            y={s.system.y * CELL_SIZE - 35}
            text={label}
            style={{
              stroke: { color: "black", width: 2 },
              fill: "white",
              fontFamily: ["JetBrains Mono Variable", "sans-serif"],
              fontSize: 1000,
            }}
            anchor={0.5}
            scale={0.03}
          />
        )
      })}
    </pixiContainer>
  )
}
