
# Notification bounded-context make targets
# Usage examples:
#   make -f cmd/cli/makefile/notification/main.mk upGoose
#   make -f cmd/cli/makefile/notification/main.mk upGoose GOOSE_MIGRATION_DIR=db/migrations
#   make -f cmd/cli/makefile/notification/main.mk upGoose GOOSE_DBSTRING="user=... password=..."

-include .env
export

GOOSE ?= goose
YQ ?= yq

CONFIG_NOTIFICATION ?= cmd/notification/config/config.yml

GOOSE_DRIVER ?= postgres
GOOSE_MIGRATION_DIR ?= db/migrations

ifndef GOOSE_DBSTRING
YQ_FOUND := $(strip $(shell where $(YQ) 2>NUL))
ifeq ($(YQ_FOUND),)
GOOSE_DBSTRING := $(shell go run ./cmd/cli/printdsn -config "$(CONFIG_NOTIFICATION)" -sslmode disable)
else
PG_USER := $(shell $(YQ) e -r '.postgresql.username' "$(CONFIG_NOTIFICATION)")
PG_PASSWORD := $(shell $(YQ) e -r '.postgresql.password' "$(CONFIG_NOTIFICATION)")
PG_DBNAME := $(shell $(YQ) e -r '.postgresql.dbname' "$(CONFIG_NOTIFICATION)")
PG_HOST := $(shell $(YQ) e -r '.postgresql.host' "$(CONFIG_NOTIFICATION)")
PG_PORT := $(shell $(YQ) e -r '.postgresql.port' "$(CONFIG_NOTIFICATION)")

GOOSE_DBSTRING := user=$(PG_USER) password=$(PG_PASSWORD) dbname=$(PG_DBNAME) host=$(PG_HOST) port=$(PG_PORT) sslmode=disable
endif
endif

up_by_one:
	set "GOOSE_DRIVER=${GOOSE_DRIVER}" && set "GOOSE_DBSTRING=${GOOSE_DBSTRING}" && ${GOOSE} -dir ${GOOSE_MIGRATION_DIR} up-by-one

create_migration:
	@${GOOSE} -dir=${GOOSE_MIGRATION_DIR} create -s ${name} sql

upGoose:
	set "GOOSE_DRIVER=${GOOSE_DRIVER}" && set "GOOSE_DBSTRING=${GOOSE_DBSTRING}" && ${GOOSE} -dir ${GOOSE_MIGRATION_DIR} up

downGoose:
	set "GOOSE_DRIVER=${GOOSE_DRIVER}" && set "GOOSE_DBSTRING=${GOOSE_DBSTRING}" && ${GOOSE} -dir ${GOOSE_MIGRATION_DIR} down

resetGoose:
	set "GOOSE_DRIVER=${GOOSE_DRIVER}" && set "GOOSE_DBSTRING=${GOOSE_DBSTRING}" && ${GOOSE} -dir ${GOOSE_MIGRATION_DIR} reset

.PHONY: up_by_one create_migration upGoose downGoose resetGoose
