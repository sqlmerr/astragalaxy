# Inventories

Inventories store items and resources.

Each inventory has the following capacity attributes:
- Maximum number of items
- Maximum volume of resources

---

### Items
Items are unstackable and can be used.

### Resources
Resources are stackable.

## Inventory Owner
An inventory owner can be either a ship or an agent. Each agent and each ship has its own inventory.

## Transfers
There are only 3 valid transfer directions:

- Your agent $\to$ your ship (must be active; subject to change in the future)
- Your active ship $\to$ your agent
- Your agent $\to$ another agent (must be in the same location as you)