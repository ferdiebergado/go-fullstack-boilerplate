include .env
export $(shell sed 's/=.*//' .env)

# DB
DB_CONTAINER := postgres
DB_IMAGE := postgres:17.0-alpine3.20

# PROXY
PROXY_CONTAINER := nginx_reverse_proxy
PROXY_IMAGE := nginx:1.27.2-alpine3.20

# MIGRATIONS
MIGRATE_CONTAINER := migrate
MIGRATE_IMAGE := migrate/migrate:v4.17.1
MIGRATIONS_DIR := ./internal/pkg/db/migrations
MIGRATIONS_DIR_REMOTE := /migrations
MIGRATE_CMD := $(CONTAINER) run -it --rm --network host --name $(MIGRATE_CONTAINER) -v $(MIGRATIONS_DIR):/migrations:Z $(MIGRATE_IMAGE)
DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

all: db proxy dev

install:
	which air || go install github.com/air-verse/air@v1.52.2

run:
	go run ./... || true

dev:
	air

db:
	$(CONTAINER) run -d --rm --network host --name $(DB_CONTAINER) -e POSTGRES_PASSWORD="$(DB_PASSWORD)" \
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
	$(MIGRATE_CMD) create -ext sql -dir $(MIGRATIONS_DIR_REMOTE) -seq $(name)

migrate:
	$(MIGRATE_CMD) -database $(DATABASE_URL) -path $(MIGRATIONS_DIR_REMOTE) up $(version)

rollback:
	$(MIGRATE_CMD) -database $(DATABASE_URL) -path $(MIGRATIONS_DIR_REMOTE) down $(version)

drop:
	$(MIGRATE_CMD) -database $(DATABASE_URL) -path $(MIGRATIONS_DIR_REMOTE) drop

force:
	$(MIGRATE_CMD) -database $(DATABASE_URL) -path $(MIGRATIONS_DIR_REMOTE) force $(version)

test:
	go test -v -race ./...

bundle:
	@cd tools && go run tools/bundle.go

watch-css:
	@cd tools && go run tools/bundle.go -watch css || true

watch-ts:
	@cd tools && go run bundle.go -watch ts || true

bundle-prod:
	@cd tools && go run tools/bundle.go -prod

.PHONY: run install dev db psql proxy migrate rollback drop test bundle watch-css watch-ts bundle-prod
