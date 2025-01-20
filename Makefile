ifeq ("$(wildcard .env)","")
# .env file does not exist
else
include .env
export
endif

PROJECT_NAME=order_service

MIGRATE_CMD := $(shell which migrate)
MIGRATE_VERSION := v4.15.2

PROJECT_ROOT_PATH=$(CURDIR)
MIGRATIONS_DIR=$(PROJECT_ROOT_PATH)/migrations
DOCKER_COMPOSE_PATH=$(PROJECT_ROOT_PATH)/docker-compose.yaml

SCRIPTS_DIR=$(PROJECT_ROOT_PATH)/scripts
CREATE_DEFAULT_ENV_SCRIPT=$(SCRIPTS_DIR)/create_default_env.sh
CREATE_LOCAL_ENV_SCRIPT=$(SCRIPTS_DIR)/create_local_env.sh

PUBLISHER_APP_PATH=$(PROJECT_ROOT_PATH)/cmd/publisher/main.go
CONSUMER_APP_PATH=$(PROJECT_ROOT_PATH)/cmd/consumer/main.go

DATABASE_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable



.PHONY: env
env:
	$(CREATE_DEFAULT_ENV_SCRIPT)

.PHONY: local-env
local-env:
	$(CREATE_LOCAL_ENV_SCRIPT)



.PHONY: build-service
build-service:
	docker-compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_PATH) build 

.PHONY: up-service
up-service: 
	docker-compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_PATH) up -d

.PHONY: down-service
down-service:
	docker-compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_PATH) down

.PHONY: run-service
run-publisher:
	go run $(PUBLISHER_APP_PATH)

.PHONY: run-service
run-consumer:
	go run $(CONSUMER_APP_PATH)


.PHONY: .check-migrate
.check-migrate:
ifeq ($(MIGRATE_CMD),)
	@echo "migrate not found, installing..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
endif

.PHONY: migration-up
migration-up: .check-migrate
	@echo "Running migrations up..."
	migrate -path $(MIGRATIONS_DIR) -database $(DATABASE_URL) up

.PHONY: migration-down
migration-down: .check-migrate
	@echo "Running migrations down..."
	migrate -path $(MIGRATIONS_DIR) -database $(DATABASE_URL) down

MIGRATION_NAME := $(name)
.PHONY: migration-create
migration-create: .check-migrate
	@if [ -z "$(MIGRATION_NAME)" ]; then \
		echo "Migration name is required. Use: make migration-create name=<migration_name>"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(MIGRATION_NAME)"
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(MIGRATION_NAME)



.PHONY: lines-count
lines-count:
	@echo 	Number of lines in GO files:
	@echo 	""[${shell find $(CURDIR) -name '*.go' -type f -print0 | xargs -0 cat | wc -l}]
