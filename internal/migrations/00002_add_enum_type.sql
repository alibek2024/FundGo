-- +goose Up
CREATE TYPE campaign_status AS ENUM ('active', 'successful', 'failed', 'archived');
CREATE TYPE donation_status AS ENUM ('active', 'refund', 'archived');

-- +goose Down
DROP TYPE IF EXISTS campaign_status;
DROP TYPE IF EXISTS donation_status

