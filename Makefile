#!make

migrate_ext = sql

ifneq ($(strip $(dbext)),)
	migrate_ext = $(dbext)
endif

migrate:
	migrate create -ext $(migrate_ext) -dir $(CURDIR)/schemes -seq $(name)