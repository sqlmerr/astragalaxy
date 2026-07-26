export const queryKeys = {
  me: ["me"] as const,
  agents: {
    all: ["agents"] as const,
    my: ["agents", "my"] as const,
    current: ["agents", "current"] as const,
    cooldown: ["agents", "current", "cooldown"] as const,
  },
  ships: {
    all: ["ships"] as const,
    my: ["ships", "my"] as const,
    active: ["ships", "my", "active"] as const,
    radar: ["ships", "my", "active", "radar"] as const,
  },
  inventories: {
    my: ["inventories", "my"] as const,
    ship: (id: string) => ["inventories", "my", "ships", id] as const,
  },
  systems: {
    current: ["systems", "current"] as const,
  },
} as const
