export const queryKeys = {
  me: ["me"] as const,
  agents: {
    all: ["agents"] as const,
    my: ["agents", "my"] as const,
    current: (agentID: string) => ["agents", "current", agentID] as const,
    cooldown: (agentID: string) =>
      ["agents", "current", "cooldown", agentID] as const,
  },
  ships: {
    my: (agentID: string) => ["ships", "my", agentID] as const,
    active: (agentID: string) => ["ships", "my", "active", agentID] as const,
    radar: (agentID: string) =>
      ["ships", "my", "active", "radar", agentID] as const,
  },
  inventories: {
    my: (agentID: string) => ["inventories", "my", agentID] as const,
    ship: (agentID: string, shipId: string) =>
      ["inventories", "my", "ships", shipId, agentID] as const,
  },
  systems: {
    current: (agentID: string) => ["systems", "current", agentID] as const,
  },
} as const
