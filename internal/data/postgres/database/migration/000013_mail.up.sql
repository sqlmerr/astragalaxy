CREATE TABLE IF NOT EXISTS mail_packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_agent UUID REFERENCES agents(id) ON DELETE CASCADE,
    to_agent UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,

    inventory_id UUID NOT NULL REFERENCES inventories(id) ON DELETE CASCADE,

    arrive_time INTERVAL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CHECK (from_agent != to_agent)
);