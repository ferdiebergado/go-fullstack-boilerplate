include .env
export

# Check if podman exists; if not, fallback to docker
CONTAINER_CMD := $(shell command -v podman 2>/dev/null || command -v docker)
COMPOSE_CMD := $(shell command -v podman-compose 2>/dev/null || command -v docker-compose)

# Migrate
MIGRATE_BASE_CMD := $(CONTAINER) run -it --rm --network host -v $(MIGRATIONS_DIR):/migrations:Z $(MIGRATE_IMAGE)
MIGRATE_CMD := $(MIGRATE_BASE_CMD) -database postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE) \
-path $(MIGRATIONS_DIR_REMOTE)

# Bundler
BUNDLE_CMD := @cd tools && go run bundle.go

.PHONY: default psql migration migrate rollback drop force test bundle watch-css watch-ts bundle-prod stop restart vulncheck teardown deploy

# deploy for development
default:
	$(COMPOSE_CMD) -f $(COMPOSE_DIR)/compose.yml -f $(COMPOSE_DIR)/compose.development.yml up --build

# stop all running services
stop:
	$(COMPOSE_CMD) -f $(COMPOSE_DIR)/compose.yml down

# restart a service (make restart service=proxy)
restart:
	$(COMPOSE_CMD) -f $(COMPOSE_DIR)/compose.yml up --no-deps -d $(service)

# interact with the database
psql:
	@CONTAINER=$(CONTAINER_CMD) ./scripts/psql.sh

# create a migration (make migration name=create_users_table)
migration:
	$(MIGRATE_BASE_CMD) create -ext sql -dir $(MIGRATIONS_DIR_REMOTE) -seq $(name)

# run the migrations
migrate:
	$(MIGRATE_CMD) up $(version)

# rollback all migrations
rollback:
	$(MIGRATE_CMD) down $(version)

# drop all tables in the database
drop:
	$(MIGRATE_CMD) drop

# force a migration (make force version=1)
force:
	$(MIGRATE_CMD) force $(version)

# run tests
test:
	@CONTAINER=$(CONTAINER_CMD) ./scripts/test.sh

# clean up after the tests
teardown:
	@CONTAINER=$(CONTAINER_CMD) ./scripts/teardown.sh

# bundle the assets
bundle:
	$(BUNDLE_CMD)

# bundle css in watch mode
watch-css:
	$(BUNDLE_CMD) -watch css || true

# bundle typescript in watch mode
watch-ts:
	$(BUNDLE_CMD) -watch ts || true

# bundle assets for production
bundle-prod:
	$(BUNDLE_CMD) -prod

# check for vulnerable packages
vulncheck:
	@which govulncheck || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck -show verbose ./...

# deploy for production
deploy:
	$(COMPOSE_CMD) -f $(COMPOSE_DIR)/compose.yml -f $(COMPOSE_DIR)/compose.production.yml --env-file /dev/null up --build
