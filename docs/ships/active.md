# Active Ships

An active ship is the ship to which the agent is currently assigned. The active ship’s coordinates are the same as the player’s coordinates. Coordinates are stored only on ships to prevent desynchronization.

## Changing the Active Ship

To change the active ship, send a `POST` request to `/api/v1/ships/my/{id}/active`, passing the ID of the ship you want to make active.

The new ship must be in the same location as the current ship and have the same status. For example, if the current ship is in `orbit`, the new ship must also be in `orbit`. The same rule applies to the `docked` status.