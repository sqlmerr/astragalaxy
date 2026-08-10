-- name: CreateShipModule :one
INSERT INTO ship_modules (ship_id, module_type)
VALUES ($1, $2)
RETURNING *;

-- name: GetShipModules :many
SELECT * FROM ship_modules
WHERE ship_id = $1;

-- name: DeleteShipModule :exec
DELETE FROM ship_modules
WHERE ship_id = $1 AND module_type = $2;

-- name: GetShipModule :one
SELECT * FROM ship_modules
WHERE ship_id = $1 AND module_type = $2;
