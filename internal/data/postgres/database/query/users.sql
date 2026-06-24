-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE LOWER(username) = $1;

-- name: UserExistsByUsername :one
SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = $1);