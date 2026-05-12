-- +goose Up
CREATE TYPE transaction_type AS ENUM (
    'deposit',
    'donation',
    'refund',
    'withdrawal'
);

CREATE TABLE transactions (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,

    donation_id INTEGER REFERENCES donations(id) ON DELETE SET NULL,

    transaction_type transaction_type NOT NULL,

    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS transactions;
DROP TYPE IF EXISTS transaction_type;

