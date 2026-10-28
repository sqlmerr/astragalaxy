# Mining Mechanics

Mining is one of the foundational resource-gathering activities in AstraGalaxy. Agents extract ores, minerals, and volatile compounds from planetary deposits and asteroid belts to craft items, refine alloys, and produce warp cells.

---

## Resource Deposits

Resources can be gathered from celestial bodies within solar systems:

* **Asteroids:** Contain a **single deposit** per body.
* **Planets:** Host **multiple deposits** across their surface, determined by the planet's biome and orbital tier.

### Deposit Attributes

Each deposit has two primary parameters:
* **Resource Type:** The specific resource available (e.g., `iron`, `crystals`, `carbon`, `ice`, `titanium`, etc.).
* **Richness (`richness`):** A normalized float value between $0.0$ and $1.0$ reflecting the quality or density of the deposit.

### Regeneration

Resource deposits are non-exhaustive:
* If a deposit is depleted, it naturally regenerates over time.
* **Base Regeneration Rate:** `1 resource unit per minute`.
* *Note: Regeneration speeds and replenishment scaling may change in future balance passes.*

---

## Extraction Limits & Mining Rates

Mining operations are subject to rate limiting and per-action extraction caps to ensure server balance and reward specialized fleet loadouts.

### Per-Action Batch Limit

By default, an extraction operation can extract up to a fixed volume per request:

$$\text{Base Extraction Limit} = 500 \text{ units}$$

This limit scales depending on the ship's **class** and **ship type**.

| Factor | Effect |
| :--- | :--- |
| **Ship Class** ($S, A, B, C$) | Directly scales the maximum batch limit per mining cycle. |
| **Miner Type** | Multiplies the extraction cap by $\mathbf{1.5\times}$. |

$$\text{Total Batch Limit} = \text{Base Limit} \times \text{Class Multiplier} \times \text{Type Multiplier}$$

---

## Mining Speed & Cooldowns

### Base Mining Speed

* **Rate:** $4 \text{ resource units per second}$.
* Operations incur an action cooldown based on the volume extracted.

### Ship Type Modifiers

Specialized ship hulls grant inherent operational bonuses:

* **Miner Type Ships:**
  * **Batch Limit:** $\times 1.5$ bonus.
  * **Action Cooldown:** Reduced by a factor of $\mathbf{0.8\times}$ (i.e., $20\%$ faster recovery after mining cycles).

---

## Equipment & Upgrades

Ships can mount specialized mining modules in their equipment slots to modify extraction parameters.

### `laser_drill`

The `laser_drill` module improves extraction efficiency through focused beam technology:

* **Speed Modification:** Multiplies extraction speed by **$1.2$**.
