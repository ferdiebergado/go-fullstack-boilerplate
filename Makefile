include .env
export

# Migrate
MIGRATE_BASE_CMD := docker run -it --rm --network host -v $(MIGRATIONS_DIR):/migrations:Z $(MIGRATE_IMAGE)
MIGRATE_CMD := $(MIGRATE_BASE_CMD) -database postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE) \
-path $(MIGRATIONS_DIR_REMOTE)

# Bundler
BUNDLE_CMD := @cd tools && go run bundle.go

.PHONY: $(wildcard *)

default:
	@sed -n 's/^##//p' Makefile | column -t -s ':' | sed -e 's/^//'

## dev: Deploy for development
dev:
	docker compose -f $(COMPOSE_DIR)/compose.yml -f $(COMPOSE_DIR)/compose.development.yml up --build

## stop: Stop all running services
stop:
	docker compose -f $(COMPOSE_DIR)/compose.yml down

## restart: Restart a service (make restart service=proxy)
restart:
	docker compose build $(service)
	docker compose -f $(COMPOSE_DIR)/compose.yml up --no-deps -d $(service)

## psql: Invoke psql on the running database instance
psql:
	./scripts/psql.sh

## migration: Create a migration (make migration name=create_users_table)
migration:
	$(MIGRATE_BASE_CMD) create -ext sql -dir $(MIGRATIONS_DIR_REMOTE) -seq $(name)

## migrate: Run the migrations
migrate:
	$(MIGRATE_CMD) up $(version)

## rollback: Rollback all migrations
rollback:
	$(MIGRATE_CMD) down $(version)

## drop: Drop all tables in the database
drop:
	$(MIGRATE_CMD) drop

## force: Force a migration (make force version=1)
force:
	$(MIGRATE_CMD) force $(version)

## test: Run the unit tests
test:
	go test -v -race ./...

## integration: Run the integration tests
integration:
	./scripts/test.sh

## teardown: Clean up after the tests
teardown:
	./scripts/teardown.sh

## bundle: Bundle the assets
bundle:
	$(BUNDLE_CMD)

## watch-css: Bundle css in watch mode
watch-css:
	$(BUNDLE_CMD) -watch css || true

## watch-ts: Bundle typescript in watch mode
watch-ts:
	$(BUNDLE_CMD) -watch ts || true

## bundle-prod: Bundle assets for production
bundle-prod:
	$(BUNDLE_CMD) -prod

## vulncheck: Check for vulnerable packages
vulncheck:
	@which govulncheck || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck -show verbose ./...

## deploy: Deploy for production
deploy:
	docker compose -f $(COMPOSE_DIR)/compose.yml -f $(COMPOSE_DIR)/compose.production.yml --env-file /dev/null up --build
