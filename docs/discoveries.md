# Discoveries

Mechanic that tracks the first visit of star systems and planets by agents.

## Overview

* **First Visit Registration**: When an agent arrives at an unvisited system or planet, the server records the agent as its first discoverer.
* **Immutability**: Once recorded, the discoverer ID cannot be changed or overwritten.
* **Scope**:
  * **System Discovery**: Triggered upon warp jump arrival at system coordinates `(x, y)`.
  * **Planet Discovery**: Triggered upon entering orbit or interacting with the planet.
