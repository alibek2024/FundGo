-- +goose Up
CREATE UNIQUE INDEX idx_users_email ON users(email);

-- +goose Down
DROP INDEX idx_users_email;
