-- +goose Up
CREATE TABLE campaigns (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    creator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    target_amount DECIMAL(18,2) NOT NULL CHECK (target_amount > 0),
    current_amount DECIMAL(18,2) DEFAULT 0.00,
    status campaign_status NOT NULL DEFAULT 'active',
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS campaigns;
