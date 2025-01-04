#!/usr/bin/env sh

. ./.env.testing

if [ -z "$(docker ps | grep $DB_CONTAINER)" ]; then
	echo "test container not running, starting it up..."
	./scripts/testdb.sh
	sleep 3
	DATABASE_URL=postgres://$DB_USER@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE
	docker run -it --rm --network host -v $MIGRATIONS_DIR:/migrations:Z $MIGRATE_IMAGE -database $DATABASE_URL -path $MIGRATIONS_DIR_REMOTE up
fi

go test -v -race -tags integration ./...
