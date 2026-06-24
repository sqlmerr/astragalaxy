-- name: AgentExistsByUsername :one
SELECT EXISTS(SELECT 1 FROM agents WHERE LOWER(username) = $1);

-- name: ChangeAgentToken :execrows
UPDATE agents
SET
    token_hash = $1
WHERE id = $2;

-- name: CountAgentsByUser :one
SELECT COUNT(*) FROM agents WHERE user_id = $1;

-- name: CreateAgent :one
INSERT INTO agents (user_id, username, token_hash) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAgentByID :one
SELECT *
FROM agents
WHERE id = $1;

-- name: GetAgentByToken :one
SELECT *
FROM agents
WHERE token_hash = $1;

-- name: GetAgentsByUser :many
SELECT *
FROM agents
WHERE user_id = $1
ORDER BY created_at;