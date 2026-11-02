CREATE TABLE IF NOT EXISTS dynamic_waypoints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system_x INT NOT NULL,
    system_y INT NOT NULL,
    owner_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,

    type VARCHAR NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}',

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);