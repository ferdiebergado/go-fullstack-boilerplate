include .env
export $(shell sed 's/=.*//' .env)

DB_CONTAINER := postgres
DB_IMAGE := postgres:17.0-alpine3.20
PROXY_CONTAINER := nginx_reverse_proxy
PROXY_IMAGE := nginx:1.27.2-alpine3.20
JS_RUNTIME_IMAGE := denoland/deno:alpine-2.1.4

MIGRATIONS_DIR := ./internal/pkg/db/migrations
DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)?sslmode=$(DB_SSLMODE)

all: db proxy dev

install:
	which migrate || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1
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
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_DIR) up $(version)

rollback:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_DIR) down $(version)

drop:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_DIR) drop

force:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_DIR) force $(version)

test:
	go test -v -race ./...

bundle:
	go run tools/bundle.go

watch-css:
	go run tools/bundle.go -watch css || true

watch-ts:
	go run tools/bundle.go -watch ts || true

bundle-prod:
	go run tools/bundle.go -prod

.PHONY: run install dev db psql proxy migrate rollback drop test bundle watch-css watch-ts bundle-prod
