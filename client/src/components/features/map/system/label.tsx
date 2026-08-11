import type { AgentExtended } from "@/api/types"

interface SystemLabelProps {
  agents: AgentExtended[]
  location: "planet" | "waypoint"
  locationId: number
  x: number
  y: number
}

export function SystemLabel({
  agents,
  location,
  locationId,
  x,
  y,
}: SystemLabelProps) {
  const parts: string[] = []
  for (const agent of agents) {
    if (
      agent.ship.location === location &&
      agent.ship.location_id === locationId
    ) {
      if (parts.length < 2) {
        parts.push(agent.agent.username)
      } else {
        parts.push("...")
        break
      }
    }
  }
  const label = parts.join(", ")
  return (
    <pixiBitmapText
      x={x}
      y={y}
      text={label}
      style={{
        stroke: { color: "black", width: 2 },
        fill: "white",
        fontFamily: ["JetBrains Mono Variable", "sans-serif"],
        fontSize: 1000,
      }}
      scale={0.03}
      anchor={0.5}
      rotation={0}
    />
  )
}
