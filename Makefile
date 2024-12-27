include .env
export $(shell sed 's/=.*//' .env)

# DB
ifeq ($(APP_ENV), production)
DB_PASSWORD_HOST := :$(DB_PASSWORD)
endif
DATABASE_URL := postgres://$(DB_USER)$(DB_PASSWORD_HOST)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# MIGRATIONS
MIGRATE_IMAGE := migrate/migrate:v4.17.1
MIGRATIONS_DIR := ./db/migrations
MIGRATIONS_DIR_REMOTE := /migrations
MIGRATE_BASE_CMD := $(CONTAINER) run -it --rm --network host -v $(MIGRATIONS_DIR):/migrations:Z $(MIGRATE_IMAGE)
MIGRATE_CMD := $(MIGRATE_BASE_CMD) -database postgres://$(DB_USER)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE) -path $(MIGRATIONS_DIR_REMOTE)

# ASSETS
BUNDLE_CMD := @cd tools && go run bundle.go

# DEPLOYMENT
COMPOSE_DIR := deployments/docker-compose
COMPOSE_BASE_CMD := $(COMPOSE) -f $(COMPOSE_DIR)/compose.yml -f $(COMPOSE_DIR)/compose.$(APP_ENV).yml

.PHONY: default psql migration migrate rollback drop force test bundle watch-css watch-ts bundle-prod stop restart dump-url vulncheck

default:
	$(COMPOSE_BASE_CMD) up --build

stop:
	$(COMPOSE_BASE_CMD) down

restart:
	$(COMPOSE_BASE_CMD) restart $(service)

psql:
	$(CONTAINER) exec -ti postgres:17.0-alpine3.20 psql -U $(DB_USER) $(DB_NAME)

dump-url:
	@echo DATABASE_URL=$(DATABASE_URL)

migration:
	$(MIGRATE_BASE_CMD) create -ext sql -dir $(MIGRATIONS_DIR_REMOTE) -seq $(name)

migrate:
	$(MIGRATE_CMD) up $(version)

rollback:
	$(MIGRATE_CMD) down $(version)

drop:
	$(MIGRATE_CMD) drop

force:
	$(MIGRATE_CMD) force $(version)

test:
	DB_HOST=localhost go test -v -race ./...

bundle:
	$(BUNDLE_CMD)

watch-css:
	$(BUNDLE_CMD) -watch css || true

watch-ts:
	$(BUNDLE_CMD) -watch ts || true

bundle-prod:
	$(BUNDLE_CMD) -prod

vulncheck:
	@which govulncheck || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck -show verbose ./...
