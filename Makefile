ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

DB_DSN="postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable"
MIG_DIR := internal/migrations

# Команда для проверки статуса
status:
	goose -dir $(MIG_DIR) postgres $(DB_DSN) status

# Команда для наката миграций
up:
	goose -dir $(MIG_DIR) postgres $(DB_DSN) up

# Команда для отката последней миграции
down:
	goose -dir $(MIG_DIR) postgres $(DB_DSN) down

exec:
	docker exec -it my_project_db psql -U postgres -d fundgo