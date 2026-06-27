CREATE TABLE inventories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    max_item_slots INT NOT NULL DEFAULT 10,
    max_resource_volume INT NOT NULL DEFAULT 1000
);

CREATE TABLE inventory_resources (
    inventory_id UUID NOT NULL REFERENCES inventories(id) ON DELETE CASCADE,
    resource_type VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL DEFAULT 0 CHECK ( amount >= 0 ),
    PRIMARY KEY (inventory_id, resource_type)
);

CREATE TABLE inventory_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inventory_id UUID NOT NULL REFERENCES inventories(id) ON DELETE CASCADE,
    item_type VARCHAR(255) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Agents inventory_id
ALTER TABLE agents ADD COLUMN inventory_id UUID;

DO $$
    DECLARE
        agent_row RECORD;
        new_inv_id UUID;
    BEGIN
        FOR agent_row IN SELECT id FROM agents WHERE inventory_id IS NULL LOOP
            INSERT INTO inventories DEFAULT VALUES
            RETURNING id INTO new_inv_id;
            UPDATE agents SET inventory_id = new_inv_id WHERE id = agent_row.id;

        END LOOP;
END $$;

ALTER TABLE agents ALTER COLUMN inventory_id SET NOT NULL;
ALTER TABLE agents
    ADD CONSTRAINT fk_agents_inventory
    FOREIGN KEY (inventory_id)
    REFERENCES inventories(id)
    ON DELETE CASCADE;

-- Ships inventory_id
ALTER TABLE ships ADD COLUMN inventory_id UUID;

DO $$
    DECLARE
        ship_row RECORD;
        new_inv_id UUID;
    BEGIN
        FOR ship_row IN SELECT id FROM ships WHERE inventory_id IS NULL LOOP
            INSERT INTO inventories (max_item_slots, max_resource_volume)
            VALUES (15, 3000)
            RETURNING id INTO new_inv_id;
            UPDATE ships SET inventory_id = new_inv_id WHERE id = ship_row.id;

        END LOOP;
    END $$;

ALTER TABLE ships ALTER COLUMN inventory_id SET NOT NULL;
ALTER TABLE ships
    ADD CONSTRAINT fk_ships_inventory
    FOREIGN KEY (inventory_id)
    REFERENCES inventories(id)
    ON DELETE CASCADE;
