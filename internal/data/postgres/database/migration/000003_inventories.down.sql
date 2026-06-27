DROP TABLE IF EXISTS inventory_items;
DROP TABLE IF EXISTS inventory_resources;

ALTER TABLE agents DROP CONSTRAINT fk_agents_inventory;
ALTER TABLE agents DROP COLUMN inventory_id;

ALTER TABLE ships DROP CONSTRAINT fk_ships_inventory;
ALTER TABLE ships DROP COLUMN inventory_id;

DROP TABLE IF EXISTS inventories;
