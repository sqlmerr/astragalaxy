# Wallets & Currency

Economic mechanics defining currency types, transaction fees, and cryptographic wallet addressing in **AstraGalaxy**.

## Currencies

The economic system operates on two distinct currency types:

| Currency | Rarity | Primary Use Cases | Source / Generation |
| --- | --- | --- | --- |
| **Units** | Common | General trading, market orders, station services, transaction gas fees. | Planetary discovery payouts, trading, missions. |
| **Quarks** | Rare | High-tier ship upgrades, specialized technology blueprints, rare mechanics. | First-discovery achievements, deep-space exploration milestones. |

## Wallet Addresses

Each agent wallet is identified by a globally unique cryptographic address.

* **Format**: Bech32-style string starting with prefix `ag1q`.
* **Length**: Exactly 43 characters.
* **Format Representation**: `ag1q<39_encoded_alphanumeric_chars>`
* **Validation**: Generated with checksum validation, allowing immediate format verification before database lookup.

## Transfer Fees

Transfers between wallets incur network processing fees.

* **Fee Payment Currency**: All transaction fees are calculated and deducted exclusively in **Units**.
* **Same-Agent Transfers**: Transfers between accounts/wallets owned by the same agent incur **zero fee** (`fee = 0`).
* **Inter-Agent Transfers**: Governed by default rates:
* **Units Fee Rate (`FeeUnits`)**: `0.01` ($1\%$)
* **Quarks Fee Rate (`FeeQuarks`)**: `5.00` ($500\%$, meaning each transferred Quark requires $5$ Units in fees)

### Fee Calculation Formula

$$\text{fee} = (\text{units} \times \text{FeeUnits}) + (\text{quarks} \times \text{FeeQuarks})$$

> **Balance Requirement**: The sending wallet must have sufficient Units to cover both the outgoing Units amount (if transferring Units) and the total computed fee in Units.

## Error Codes

- `INSUFFICIENT_BALANCE`
- `INVALID_ADDRESS`
- `WALLET_NOT_FOUND`
- `INVALID_WALLET_ADDRESS`