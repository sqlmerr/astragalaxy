-- name: CreateShip :one
WITH new_ship AS (
    INSERT INTO ships (agent_id, type, active, system_x, system_y, status, name, inventory_id, location, location_id, class)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    RETURNING *
)
SELECT 
    ns.*,
    COALESCE(
        array_agg(sm.module_type ORDER BY sm.module_type) FILTER (WHERE sm.module_type IS NOT NULL),
        '{}'::text[]
    )::text[] AS modules
FROM new_ship ns
LEFT JOIN ship_modules sm ON sm.ship_id = ns.id
GROUP BY ns.id, ns.agent_id, ns.type, ns.active, ns.system_x, ns.system_y, ns.status, ns.name, ns.inventory_id, ns.location, ns.location_id, ns.class, ns.created_at;

-- name: GetActiveShipByAgent :one
SELECT 
    s.*,
    COALESCE(
        array_agg(sm.module_type ORDER BY sm.module_type) FILTER (WHERE sm.module_type IS NOT NULL),
        '{}'::text[]
    )::text[] AS modules
FROM ships s
LEFT JOIN ship_modules sm ON sm.ship_id = s.id
WHERE s.agent_id = $1 AND s.active = true
GROUP BY s.id;

-- name: GetShipByID :one
SELECT 
    s.*,
    COALESCE(
        array_agg(sm.module_type ORDER BY sm.module_type) FILTER (WHERE sm.module_type IS NOT NULL),
        '{}'::text[]
    )::text[] AS modules
FROM ships s
LEFT JOIN ship_modules sm ON sm.ship_id = s.id
WHERE s.id = $1
GROUP BY s.id;

-- name: GetShipsByAgent :many
SELECT 
    s.*,
    COALESCE(
        array_agg(sm.module_type ORDER BY sm.module_type) FILTER (WHERE sm.module_type IS NOT NULL),
        '{}'::text[]
    )::text[] AS modules
FROM ships s
LEFT JOIN ship_modules sm ON sm.ship_id = s.id
WHERE s.agent_id = $1
GROUP BY s.id, s.created_at;

-- name: SaveShip :one
WITH updated_ship AS (
    UPDATE ships
    SET
        type = $2,
        active = $3,
        system_x = $4,
        system_y = $5,
        status = $6,
        name = $7,
        location = $8,
        location_id = $9,
        class = $10
    WHERE id = $1
    RETURNING *
)
SELECT 
    us.*,
    COALESCE(
        array_agg(sm.module_type ORDER BY sm.module_type) FILTER (WHERE sm.module_type IS NOT NULL),
        '{}'::text[]
    )::text[] AS modules
FROM updated_ship us
LEFT JOIN ship_modules sm ON sm.ship_id = us.id
GROUP BY us.id, us.agent_id, us.type, us.active, us.system_x, us.system_y, us.status, us.name, us.inventory_id, us.location, us.location_id, us.class, us.created_at;
