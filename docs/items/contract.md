# Contract

This item holds units and quarks that can be claimed into an agent's wallet.

## Usage
```json
{
    "address": "ag1q..." // Target wallet address to receive the funds
}
```

## Response
```json
{
    "ok": true,
    "units": 123, // Amount of units credited to the wallet
    "quarks": 123 // Amount of quarks credited to the wallet
}
```