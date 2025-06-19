-- +goose Up
CREATE TABLE IF NOT EXISTS exchange_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    base_currency VARCHAR(3) NOT NULL,    
    target_currency VARCHAR(3) NOT NULL, 
    rate NUMERIC(15, 6) NOT NULL,         
    created_at TIMESTAMPTZ DEFAULT NOW(), 
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT uniq_base_target UNIQUE (base_currency, target_currency)
);

INSERT INTO exchange_rates (base_currency, target_currency, rate) VALUES
  ('USD', 'RUB', 89.00),
  ('RUB', 'USD', 0.0112),
  ('USD', 'EUR', 0.93),
  ('EUR', 'USD', 1.075),
  ('RUB', 'EUR', 0.0104),
  ('EUR', 'RUB', 96.15);


-- +goose Down

DROP TABLE IF EXISTS exchange_rates;