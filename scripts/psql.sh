#!/usr/bin/env sh

if [ -z $(docker ps | grep $DB_CONTAINER) ]; then
	. ./.env.testing
fi

docker exec -ti $DB_CONTAINER psql -U $DB_USER $DB_NAME
