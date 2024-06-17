#!make

migrate_ext = sql

ifneq ($(strip $(fext)),)
	migrate_ext = $(fext)
endif

migrate:
	migrate create -ext $(migrate_ext) -dir $(CURDIR)/schemes -seq $(name)

migrate_database = postgres://localhost:5432/database?sslmode=enable

ifneq ($(strip $(database)),)
	migrate_database = $(database)
endif

migrate-action:
	migrate -source file://schemes -database $(migrate_database) $(cmd)

swag:
	swag init -g ./cmd/serve/main.go --output ./docs/swagger --parseInternal true

serve:
	make swag
	npm run serve