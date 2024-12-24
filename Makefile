include .env
export $(shell sed 's/=.*//' .env)

# DB
DB_CONTAINER := gfb-db
DB_IMAGE := postgres:17.0-alpine3.20
DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# PROXY
PROXY_CONTAINER := gfb-proxy
PROXY_IMAGE := nginx:1.27.2-alpine3.20

# MIGRATIONS
MIGRATE_CONTAINER := gfb-migrate
MIGRATE_IMAGE := migrate/migrate:v4.17.1
MIGRATIONS_DIR := ./internal/pkg/db/migrations
MIGRATIONS_DIR_REMOTE := /migrations
MIGRATE_BASE_CMD := $(CONTAINER) run -it --rm --network host --name $(MIGRATE_CONTAINER) -v $(MIGRATIONS_DIR):/migrations:Z $(MIGRATE_IMAGE)
MIGRATE_CMD := $(MIGRATE_BASE_CMD) -database $(DATABASE_URL) -path $(MIGRATIONS_DIR_REMOTE)

# ASSETS
BUNDLE_CMD := @cd tools && go run bundle.go

# APP
DEV_CMD := $(CONTAINER) run -it --rm --network host --name air -w "/app" -v ./:/app:Z -p $(SERVER_PORT):$(SERVER_PORT) cosmtrek/air

# DEPLOYMENT
COMPOSE_DIR := deployments/docker-compose
COMPOSE_BASE_CMD := $(COMPOSE) -f $(COMPOSE_DIR)/compose.yml -f $(COMPOSE_DIR)/compose.dev.yml

.PHONY: all run dev db proxy psql migration migrate rollback drop force test bundle watch-css watch-ts bundle-prod stop restart

all: db proxy dev

run:
	go run ./... || true

dev:
	$(COMPOSE_BASE_CMD) up --build

stop:
	$(COMPOSE_BASE_CMD) down

restart:
	$(COMPOSE_BASE_CMD) restart $(service)

db:
	$(CONTAINER) run -d --rm --network host --name $(DB_CONTAINER) \
		-e POSTGRES_PASSWORD="$(DB_PASSWORD)" \
		-e POSTGRES_USER="$(DB_USER)" -e POSTGRES_DB="$(DB_NAME)" \
		-v ./configs/postgresql.conf:/etc/postgresql/postgresql.conf:Z \
		-v ./configs/psqlrc:/root/.psqlrc:Z \
		$(DB_IMAGE) -c 'config_file=/etc/postgresql/postgresql.conf'

proxy:
	$(CONTAINER) run -d --rm --network host --name $(PROXY_CONTAINER) \
		-v ./configs/nginx.conf:/etc/nginx/nginx.conf:Z \
		-v ./web/static:/usr/share/nginx/html:Z \
		$(PROXY_IMAGE)

psql:
	$(CONTAINER) exec -ti $(DB_CONTAINER) psql -U $(DB_USER) $(DB_NAME)

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
	go test -v -race ./...

bundle:
	$(BUNDLE_CMD)

watch-css:
	$(BUNDLE_CMD) -watch css || true

watch-ts:
	$(BUNDLE_CMD) -watch ts || true

bundle-prod:
	$(BUNDLE_CMD) -prod

