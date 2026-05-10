-- +goose Up
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- +goose Down
DROP INDEX idx_users_deleted_at;
ALTER TABLE users DROP COLUMN deleted_at;