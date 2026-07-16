-- name: GetClient :one
SELECT * FROM clients WHERE id = $1;

-- Pessimistic lock, mirroring ClientRepository.findByIdForUpdate:
-- prevents concurrent transactions from reading a stale balance.
-- name: GetClientForUpdate :one
SELECT * FROM clients WHERE id = $1 FOR UPDATE;

-- name: ClientExists :one
SELECT EXISTS(SELECT 1 FROM clients WHERE id = $1);

-- name: CountClients :one
SELECT COUNT(*) FROM clients;

-- name: CreateClient :one
INSERT INTO clients (username, balance) VALUES ($1, $2) RETURNING *;

-- name: UpdateClientBalance :exec
UPDATE clients SET balance = $2 WHERE id = $1;
