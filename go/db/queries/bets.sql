-- name: BetExistsByIdempotencyKey :one
SELECT EXISTS(SELECT 1 FROM bets WHERE idempotency_key = $1);

-- name: CreateBet :one
INSERT INTO bets (client_id, predicted_value, stake, idempotency_key, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateDraw :one
INSERT INTO draws (bet_id, die_one, die_two, die_three, product_value)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListBetHistory :many
SELECT b.id, b.predicted_value, b.stake, b.status, b.created_at,
       d.die_one, d.die_two, d.die_three, d.product_value
FROM bets b
         JOIN draws d ON d.bet_id = b.id
WHERE b.client_id = $1
ORDER BY b.created_at DESC
LIMIT $2;
