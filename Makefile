include .env
export

ENV_FILE := .env
COMPOSE_DEV_FILE := docker-compose.dev.yml
COMPOSE_PROD_FILE := docker-compose.prod.yml
COMPOSE_INFRA_FILE := docker-compose.infra.yml

import_db:
	docker exec -i postgres-db psql -U postgres -d $(POSTGRES_DB) < ./$(POSTGRES_DB_BACKUP)
export_db:
	docker exec -i postgres-db pg_dump -U postgres -d $(POSTGRES_DB) > ./$(POSTGRES_DB_BACKUP)

create_migration:
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)
migrate_up:
	migrate -path $(MIGRATION_DIR) -database $(DATABASE_URL) up
migrate_down:
	migrate -path $(MIGRATION_DIR) -database $(DATABASE_URL) down 1
migrate_force:
	migrate -path $(MIGRATION_DIR) -database $(DATABASE_URL) force $(VERSION)
migrate_drop:
	migrate -path $(MIGRATION_DIR) -database $(DATABASE_URL) drop
migrate_goto:
	migrate -path $(MIGRATION_DIR) -database $(DATABASE_URL) goto $(VERSION)

sqlc:
	sqlc generate

dev:
	go run ./cmd/api
build:
	go build -o bin/api ./cmd/api
run: build
	./bin/api

dev_up:
	docker compose -f $(COMPOSE_DEV_FILE) --env-file $(ENV_FILE) up --remove-orphans
dev_down:
	docker compose -f $(COMPOSE_DEV_FILE) down

prod_up: prod_down
	docker compose -f $(COMPOSE_PROD_FILE) --env-file $(ENV_FILE) up -d --build
prod_down:
	docker compose -f $(COMPOSE_PROD_FILE) down
prod_logs:
	docker compose -f $(COMPOSE_PROD_FILE) logs -f

# Without the API service
infra_up: infra_down
	docker compose -f $(COMPOSE_INFRA_FILE) --env-file $(ENV_FILE) up -d --build
infra_down:
	docker compose -f $(COMPOSE_INFRA_FILE) down
infra_logs:
	docker compose -f $(COMPOSE_INFRA_FILE) logs -f

.PHONY: import_db export_db \
	create_migration migrate_up migrate_down migrate_force migrate_drop migrate_goto \
	sqlc \
	dev build run dev_up dev_down prod_up prod_down prod_logs infra_up infra_down infra_logs
