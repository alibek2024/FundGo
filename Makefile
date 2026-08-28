ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

DB_DSN=$(DATABASE_URL)

# Команда для проверки статуса
status:
	goose -dir internal/migrations postgres $(DB_DSN) status

# Команда для наката миграций
up:
	goose -dir internal/migrations postgres "$(DB_DSN)" up

# Команда для отката последней миграции
down:
	goose -dir internal/migrations postgres $(DB_DSN) down

exec:
	docker exec -it my_project_db psql -U postgres -d fundgo

exec-redis:
	docker exec -it f259a7ac9275 redis-cli