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

.PHONY: default psql migration migrate rollback drop force test bundle watch-css watch-ts bundle-prod stop restart vulncheck

default:
	$(COMPOSE_CMD) -f $(COMPOSE_DIR)/compose.yml -f $(COMPOSE_DIR)/compose.development.yml up --build

stop:
	$(COMPOSE_CMD) -f $(COMPOSE_DIR)/compose.yml down

restart:
	$(COMPOSE_CMD) -f $(COMPOSE_DIR)/compose.yml up --no-deps -d $(service)

psql:
	@CONTAINER=$(CONTAINER_CMD) ./scripts/psql.sh

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
	@CONTAINER=$(CONTAINER_CMD) ./scripts/test.sh

teardown:
	@CONTAINER=$(CONTAINER_CMD) ./scripts/teardown.sh

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

deploy:
	$(COMPOSE_CMD) -f $(COMPOSE_DIR)/compose.yml -f $(COMPOSE_DIR)/compose.production.yml --env-file /dev/null up --build
