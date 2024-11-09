stage := $(or $(STAGE), dev)

include Makefile.$(stage)

goose-install:
	go get -u github.com/pressly/goose/v3/cmd/goose

MIGRATION_NAME=$(or $(MIGRATION), init)

# Example: make migrate-create MIGRATION=init
migrate-create:
	mkdir -p ./repository/db/migrations
	goose -dir ./repository/db/migrations -table schema_migrations postgres $(POSTGRES_URI) create $(MIGRATION_NAME) sql

migrate-up:
	goose -dir ./repository/db/migrations -table schema_migrations postgres $(POSTGRES_URI) up
migrate-redo:
	goose -dir ./repository/db/migrations -table schema_migrations postgres $(POSTGRES_URI) redo
migrate-down:
	goose -dir ./repository/db/migrations -table schema_migrations postgres $(POSTGRES_URI) down
migrate-reset:
	goose -dir ./repository/db/migrations -table schema_migrations postgres $(POSTGRES_URI) reset
migrate-status:
	goose -dir ./repository/db/migrations -table schema_migrations postgres $(POSTGRES_URI) status

sqlc-generate:
	cd ./repository/db/ && sqlc generate
sqlc-vet:
	cd ./repository/db/ && sqlc vet
