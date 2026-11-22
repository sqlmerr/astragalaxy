-- name: CreateMailPackage :one
INSERT INTO mail_packages (from_agent, to_agent, inventory_id, arrive_time) 
VALUES ($1, $2, $3, $4)
RETURNING *; 

-- name: GetMailPackage :one
SELECT * FROM mail_packages 
WHERE id = $1;

-- name: GetMailPackagesByFromAgent :many
SELECT * FROM mail_packages
WHERE from_agent = $1
ORDER BY created_at;

-- name: GetMailPackagesByToAgent :many
SELECT * FROM mail_packages
WHERE to_agent = $1
ORDER BY created_at;

-- name: GetMailPackagesCountByFromAgent :one
SELECT COUNT(*)::INTEGER FROM mail_packages
WHERE from_agent = $1;

-- name: DeleteMailPackage :execrows
DELETE FROM mail_packages WHERE id = $1;