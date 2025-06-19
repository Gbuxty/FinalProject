-- +goose Up
CREATE TABLE IF NOT EXISTS tokens (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(512) NOT NULL,
    access_token VARCHAR(512) NOT NULL,
    refresh_token_expires_at TIMESTAMPTZ NOT NULL, 
    access_token_expires_at TIMESTAMPTZ NOT NULL,  
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);

-- +goose Down
DROP TABLE IF EXISTS tokens;
DROP INDEX IF EXISTS idx_tokens_user_id;