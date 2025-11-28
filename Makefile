include dev.env

MIGRATE_DB = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
MIGRATE_PATH = ./pkg/database/migrations

run:
	@go run cmd/api/main.go

# ------------------------ Start Migrate -----------------------
# Example: make migrate-create name=example
migrate-create:
	@migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

migrate-up:
	@migrate -database $(MIGRATE_DB) -path $(MIGRATE_PATH) up

migrate-down:
	@migrate -database $(MIGRATE_DB) -path $(MIGRATE_PATH) down 1

# Example: make migrate-force version=1
migrate-force:
	@migrate -database $(MIGRATE_DB) -path $(MIGRATE_PATH) force $(version)

# ------------------------ End Migrate -----------------------
