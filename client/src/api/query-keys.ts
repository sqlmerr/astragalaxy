export const queryKeys = {
  me: ["me"] as const,
  agents: {
    all: ["agents"] as const,
    my: ["agents", "my"] as const,
    current: (agentId: string) => ["agents", "current", agentId] as const,
    cooldown: (agentId: string) =>
      ["agents", "current", "cooldown", agentId] as const,
  },
  ships: {
    my: (agentId: string) => ["ships", "my", agentId] as const,
    active: (agentId: string) => ["ships", "my", "active", agentId] as const,
    radar: (agentId: string) =>
      ["ships", "my", "active", "radar", agentId] as const,
    modules: (agentId: string, shipId: string) =>
      ["ships", "my", shipId, "modules", agentId] as const,
  },
  inventories: {
    all: ["inventories"] as const,
    my: (agentId: string) => ["inventories", "my", agentId] as const,
    ship: (agentId: string, shipId: string) =>
      ["inventories", "my", "ships", shipId, agentId] as const,
  },
  systems: {
    current: (agentId: string) => ["systems", "current", agentId] as const,
  },
  data: {
    recipes: ["data", "recipes"] as const,
    items: ["data", "items"] as const,
    resources: ["data", "resources"] as const,
    facilities: ["data", "facilities"] as const,
  },
} as const
