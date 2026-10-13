# Ship Classes

A ship class is a single letter (`S`, `A`, `B`, `C`) that acts as a **modifier bundle** for a
ship's base stats. A ship itself is defined by its type (`trader`, `scout`, `miner`) — the class
sits on top and tunes how strong that ship is.

## What a class affects

| Field              | Meaning                                                    |
| ------------------ | ---------------------------------------------------------- |
| `max_module_slots` | How many modules a ship of this class can install          |
| `max_warp_distance`| Furthest distance the ship can travel in one warp jump     |

Classes are **data-driven**: they live in `data/ship-classes.json`. Adding or
tweaking a class should never require a code change or a redeploy.
