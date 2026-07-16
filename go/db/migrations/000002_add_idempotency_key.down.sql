DROP INDEX IF EXISTS idx_bets_idempotency_key;

ALTER TABLE bets
    DROP COLUMN IF EXISTS idempotency_key;
