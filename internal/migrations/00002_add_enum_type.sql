-- +goose Up
CREATE TYPE campaign_status AS ENUM ('active', 'successful', 'failed', 'archived');

-- +goose Down
DROP TYPE IF EXISTS campaign_status;
