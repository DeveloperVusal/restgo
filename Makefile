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

GOPATH := $(shell go env GOPATH)

swag:
	@$(GOPATH)/bin/swag init -g cmd/serve/main.go --output docs/swagger --parseInternal true --outputTypes json,yaml,go

serve:
	make swag
	npm run serve