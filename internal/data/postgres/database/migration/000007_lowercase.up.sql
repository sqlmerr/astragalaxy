ALTER TYPE ship_type RENAME VALUE 'TRADER' TO 'trader';
ALTER TYPE ship_type RENAME VALUE 'SCOUT' TO 'scout';
ALTER TYPE ship_type RENAME VALUE 'MINER' TO 'miner';

ALTER TYPE ship_status RENAME VALUE 'ORBIT' TO 'orbit';
ALTER TYPE ship_status RENAME VALUE 'DOCKED' TO 'docked';

ALTER TYPE ship_location RENAME VALUE 'NONE' TO 'none';
ALTER TYPE ship_location RENAME VALUE 'PLANET' TO 'planet';
ALTER TYPE ship_location RENAME VALUE 'WAYPOINT' TO 'waypoint';

ALTER TYPE location_type RENAME VALUE 'PLANET' TO 'planet';
ALTER TYPE location_type RENAME VALUE 'WAYPOINT' TO 'waypoint';

UPDATE inventory_resources SET resource_type = LOWER(resource_type);
UPDATE inventory_items SET item_type = LOWER(item_type);
UPDATE resource_deposit_states SET resource_type = LOWER(resource_type);
