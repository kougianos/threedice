-- V1: Initial schema for ThreeDice game

CREATE TABLE clients
(
    id         BIGSERIAL PRIMARY KEY,
    username   VARCHAR(100)   NOT NULL,
    balance    NUMERIC(12, 2) NOT NULL DEFAULT 1000.00,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE bets
(
    id              BIGSERIAL PRIMARY KEY,
    client_id       BIGINT         NOT NULL REFERENCES clients (id),
    predicted_value INT            NOT NULL,
    stake           NUMERIC(12, 2) NOT NULL,
    status          VARCHAR(4)     NOT NULL CHECK (status IN ('WON', 'LOST')),
    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE draws
(
    id            BIGSERIAL PRIMARY KEY,
    bet_id        BIGINT NOT NULL UNIQUE REFERENCES bets (id),
    die_one       INT    NOT NULL CHECK (die_one BETWEEN 1 AND 6),
    die_two       INT    NOT NULL CHECK (die_two BETWEEN 1 AND 6),
    die_three     INT    NOT NULL CHECK (die_three BETWEEN 1 AND 6),
    product_value INT    NOT NULL
);

CREATE TABLE transactions
(
    id            BIGSERIAL PRIMARY KEY,
    bet_id        BIGINT         NOT NULL UNIQUE REFERENCES bets (id),
    client_id     BIGINT         NOT NULL REFERENCES clients (id),
    type          VARCHAR(6)     NOT NULL CHECK (type IN ('DEBIT', 'CREDIT')),
    amount        NUMERIC(12, 2) NOT NULL,
    balance_after NUMERIC(12, 2) NOT NULL,
    created_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- Indexes for common query patterns
CREATE INDEX idx_bets_client_id ON bets (client_id, created_at DESC);
CREATE INDEX idx_transactions_client_id ON transactions (client_id, created_at DESC);
