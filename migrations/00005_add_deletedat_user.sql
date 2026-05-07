-- +goose Up
ALTER TABLE users
ADD deleted_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- +goose Down
DROP INDEX idx_users_deleted_at;
ALTER TABLE users DROP COLUMN deleted_at;