CONTAINER_NAME := postgres
CONTAINER_IMAGE := postgres:17.0-alpine3.20
MIGRATIONS_DIR := ./internal/pkg/db/migrations
CONTAINER := docker

# Load environment variables from .env
include .env
export $(shell sed 's/=.*//' .env)

install:
	which migrate || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	which air || go install github.com/air-verse/air@v1.52.2

run:
	go run ./...

dev:
	air

db:
	$(CONTAINER) run -d --rm --network host --name $(CONTAINER_NAME) -e POSTGRES_PASSWORD="$(DB_PASS)" \
		-v ./configs/postgresql.conf:/etc/postgresql/postgresql.conf:Z \
		-v ./configs/psqlrc:/root/.psqlrc:Z \
		$(CONTAINER_IMAGE) -c 'config_file=/etc/postgresql/postgresql.conf'

psql:
	$(CONTAINER) exec -ti $(CONTAINER_NAME) psql -U $(DB_USER) $(DB_NAME)

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
	go test -race ./...

css-watch:
	esbuild ./web/app/css/styles.css --bundle --outdir=./web/static/css --watch

js-watch:
	esbuild ./web/app/js/**/*.js --bundle --outdir=./web/static/js --sourcemap --target=es6 --splitting --format=esm --watch

.PHONY: db psql migrate rollback drop test
