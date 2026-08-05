-- name: UpsertResourceDeposit :exec
INSERT INTO resource_deposit_states (system_x, system_y, loc_type, loc_id, resource_type, remaining, last_mined_at) 
VALUES ($1, $2, $3, $4, $5, $6, NOW())
ON CONFLICT (system_x, system_y, loc_type, loc_id, resource_type)
DO UPDATE SET
    remaining = EXCLUDED.remaining,
    last_mined_at = EXCLUDED.last_mined_at;

-- name: GetResourceDeposit :one
SELECT * FROM resource_deposit_states
WHERE system_x = $1 AND system_y = $2 AND loc_type = $3 AND loc_id = $4 AND resource_type = $5;

-- name: GetResourceDepositForUpdate :one
SELECT * FROM resource_deposit_states
WHERE system_x = $1 AND system_y = $2 AND loc_type = $3 AND loc_id = $4 AND resource_type = $5
FOR UPDATE;

-- name: GetSystemResourceDeposits :many
SELECT * FROM resource_deposit_states
WHERE system_x = $1 AND system_y = $2;