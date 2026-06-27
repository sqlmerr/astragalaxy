-- name: CreateInventory :one
INSERT INTO inventories (max_item_slots, max_resource_volume)
VALUES ($1, $2)
RETURNING *;

-- name: GetInventoryByID :one
SELECT * FROM inventories
WHERE id = $1;

-- name: UpdateInventory :one
UPDATE inventories
SET
    max_item_slots = $2,
    max_resource_volume = $3
WHERE id = $1
RETURNING *;

-- name: GetInventoryOwner :one
SELECT
    id AS owner_id,
    'AGENT' AS owner_type
FROM agents
WHERE agents.inventory_id = $1

UNION ALL

SELECT
    id AS owner_id,
    'SHIP' AS owner_type
FROM ships
WHERE ships.inventory_id = $1;

-- name: CreateInventoryResource :one
INSERT INTO inventory_resources (inventory_id, resource_type, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- -- name: AddInventoryResource :one
-- INSERT INTO inventory_resources (inventory_id, resource_type, amount)
-- VALUES ($1, $2, $3)
-- ON CONFLICT (inventory_id, resource_type)
-- DO UPDATE SET
--     amount = amount + $3
-- RETURNING *;

-- name: GetInventoryResources :many
SELECT * FROM inventory_resources
WHERE inventory_id = $1
ORDER BY amount;

-- name: GetInventoryResource :one
SELECT * FROM inventory_resources
WHERE inventory_id = $1 AND resource_type = $2;

-- name: UpdateInventoryResource :one
UPDATE inventory_resources
SET
    inventory_id = $1,
    amount = $3
WHERE inventory_id = $1 AND resource_type = $2
RETURNING *;

-- name: DeleteInventoryResource :execrows
DELETE FROM inventory_resources WHERE inventory_id = $1 AND resource_type = $2;

-- name: CreateInventoryItem :one
INSERT INTO inventory_items (inventory_id, item_type, metadata)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetInventoryItems :many
SELECT * FROM inventory_items
WHERE inventory_id = $1
ORDER BY created_at;

-- name: GetInventoryItem :one
SELECT * FROM inventory_items
WHERE id = $1;

-- name: UpdateInventoryItem :one
UPDATE inventory_items
SET
    inventory_id = $2,
    metadata = $3
WHERE id = $1
RETURNING *;

-- name: DeleteInventoryItem :execrows
DELETE FROM inventory_items WHERE id = $1;