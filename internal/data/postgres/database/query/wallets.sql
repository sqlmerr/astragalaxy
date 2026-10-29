-- name: CreateWallet :one
INSERT INTO wallets (address, units, quarks, owner_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWalletByAddress :one
SELECT * FROM wallets
WHERE address = $1;

-- name: GetWalletsByAgent :many
SELECT * FROM wallets
WHERE owner_id = $1;

-- name: SaveWallet :one
UPDATE wallets
SET
    units = $2,
    quarks = $3
WHERE address = $1
RETURNING *;

-- name: DeleteWallet :execrows
DELETE FROM wallets
WHERE address = $1;