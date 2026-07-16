-- name: CreateTransaction :one
INSERT INTO transactions (bet_id, client_id, type, amount, balance_after)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListTransactionHistory :many
SELECT id, bet_id, type, amount, balance_after, created_at
FROM transactions
WHERE client_id = $1
ORDER BY created_at DESC
LIMIT $2;
