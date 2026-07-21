-- name: CreateShip :one
INSERT INTO ships (agent_id, type, active, system_x, system_y, status, name, inventory_id, location, location_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetActiveShipByAgent :one
SELECT *
FROM ships
WHERE agent_id = $1 AND active = true;

-- name: GetShipByID :one
SELECT *
FROM ships
WHERE id = $1;

-- name: GetShipsByAgent :many
SELECT *
FROM ships
WHERE agent_id = $1
ORDER BY created_at DESC;

-- name: SaveShip :one
UPDATE ships
SET
    type = $2,
    active = $3,
    system_x = $4,
    system_y = $5,
    status = $6,
    name = $7,
    location = $8,
    location_id = $9
WHERE id = $1
RETURNING *;
