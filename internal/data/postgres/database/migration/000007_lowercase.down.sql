UPDATE inventory_resources SET resource_type = UPPER(resource_type);
UPDATE inventory_items SET item_type = UPPER(item_type);
UPDATE resource_deposit_states SET resource_type = UPPER(resource_type);

ALTER TYPE ship_type RENAME VALUE 'trader' TO 'TRADER';
ALTER TYPE ship_type RENAME VALUE 'scout' TO 'SCOUT';
ALTER TYPE ship_type RENAME VALUE 'miner' TO 'MINER';

ALTER TYPE ship_status RENAME VALUE 'orbit' TO 'ORBIT';
ALTER TYPE ship_status RENAME VALUE 'docked' TO 'DOCKED';

ALTER TYPE ship_location RENAME VALUE 'none' TO 'NONE';
ALTER TYPE ship_location RENAME VALUE 'planet' TO 'PLANET';
ALTER TYPE ship_location RENAME VALUE 'waypoint' TO 'WAYPOINT';

ALTER TYPE location_type RENAME VALUE 'planet' TO 'PLANET';
ALTER TYPE location_type RENAME VALUE 'waypoint' TO 'WAYPOINT';
