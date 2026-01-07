include .env
export

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

server:
	air

.PHONY: import_db export_db create_migration migrate_up migrate_down migrate_force migrate_drop migrate_goto sqlc server
