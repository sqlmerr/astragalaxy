-- name: CreateDynamicWaypoint :one
INSERT INTO dynamic_waypoints (system_x, system_y, owner_id, type, metadata)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetDynamicWaypoint :one
SELECT * FROM dynamic_waypoints
WHERE id = $1;

-- name: GetDynamicWaypointsBySystem :many
SELECT * FROM dynamic_waypoints
WHERE system_x = $1 AND system_y = $2;

-- name: GetDynamicWaypointsByAgent :many
SELECT * FROM dynamic_waypoints
WHERE owner_id = $1;

-- name: GetDynamicWaypointsByAgentAndSystem :many
SELECT * FROM dynamic_waypoints
WHERE system_x = $1 AND system_y = $2 AND owner_id = $3;

-- name: SaveDynamicWaypoint :one
UPDATE dynamic_waypoints
SET
    system_x = $2,
    system_y = $3,
    owner_id = $4,
    type = $5,
    metadata = $6
WHERE id = $1
RETURNING *;

-- name: DeleteDynamicWaypoint :execrows
DELETE FROM dynamic_waypoints WHERE id = $1;