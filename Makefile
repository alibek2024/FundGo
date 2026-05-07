DB_DSN := "host=localhost port=5432 user=postgres password=secret dbname=fund_go sslmode=disable"
MIG_DIR := ./migrations

# Команда для проверки статуса
status:
	goose -dir $(MIG_DIR) postgres $(DB_DSN) status

# Команда для наката миграций
up:
	goose -dir $(MIG_DIR) postgres $(DB_DSN) up

# Команда для отката последней миграции
down:
	goose -dir $(MIG_DIR) postgres $(DB_DSN) down