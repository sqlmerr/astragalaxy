CREATE TABLE IF NOT EXISTS wallet_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_address VARCHAR(255) NOT NULL REFERENCES wallets(address),
    to_address VARCHAR(255) NOT NULL REFERENCES wallets(address),

    units BIGINT NOT NULL,
    quarks BIGINT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CHECK (from_address != to_address)
);