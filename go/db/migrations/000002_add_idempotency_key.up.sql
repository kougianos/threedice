-- V2: Add idempotency key to bets table for duplicate submission prevention

ALTER TABLE bets
    ADD COLUMN idempotency_key VARCHAR(64);

-- Backfill existing rows with a generated unique value
UPDATE bets SET idempotency_key = 'legacy-' || id WHERE idempotency_key IS NULL;

-- Now make it NOT NULL and add unique constraint
ALTER TABLE bets
    ALTER COLUMN idempotency_key SET NOT NULL;

CREATE UNIQUE INDEX idx_bets_idempotency_key ON bets (idempotency_key);
