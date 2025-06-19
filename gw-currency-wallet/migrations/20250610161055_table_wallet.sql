-- +goose Up
CREATE TABLE IF NOT EXISTS wallet(
    user_id UUID NOT NULL REFERENCES users(id),
    currency VARCHAR(3) NOT NULL, 
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_wallet_user_currency UNIQUE (user_id, currency)
);

CREATE INDEX IF NOT EXISTS idx_wallet_user_id ON wallet(user_id);
CREATE INDEX IF NOT EXISTS idx_wallet_user_currency ON wallet(user_id, currency);
-- +goose Down
DROP TABLE if EXISTS wallet;
DROP INDEX IF EXISTS idx_wallet_user_id;
DROP INDEX IF EXISTS idx_wallet_user_currency;