-- +goose Up
CREATE TYPE transaction_type AS ENUM (
    'deposit',
    'donation',
    'refund',
    'withdrawal'
);

CREATE TABLE transactions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,

    donation_id BIGINT REFERENCES donations(id) ON DELETE SET NULL,

    transaction_type transaction_type NOT NULL,

    balance_before DECIMAL(18,2) NOT NULL CHECK (balance_before >= 0),

    balance_after DECIMAL(18,2) NOT NULL CHECK (balance_after >= 0),

    amount DECIMAL(18,2) NOT NULL CHECK (amount > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS transactions;
DROP TYPE IF EXISTS transaction_type;

