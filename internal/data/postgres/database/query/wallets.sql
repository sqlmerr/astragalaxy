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

-- name: AddFundsToWallet :one
UPDATE wallets
SET
    units = units + $2,
    quarks = quarks + $3
WHERE address = $1
RETURNING *;

-- name: DeleteWallet :execrows
DELETE FROM wallets
WHERE address = $1;


-- name: CreateWalletTransaction :one
INSERT INTO wallet_transactions (from_address, to_address, units, quarks)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWalletTransaction :one
SELECT * FROM wallet_transactions
WHERE id = $1;

-- name: GetWalletTransactionsByFromAddress :many
SELECT * FROM wallet_transactions
WHERE from_address = $1
ORDER BY created_at;

-- name: GetWalletTransactionsByToAddress :many
SELECT * FROM wallet_transactions
WHERE to_address = $1
ORDER BY created_at;

-- name: GetWalletTransactionsByAddress :many
SELECT * FROM wallet_transactions
WHERE from_address = $1 OR to_address = $1
ORDER BY created_at;

-- name: DeleteWalletTransaction :execrows
DELETE FROM wallet_transactions
WHERE id = $1;