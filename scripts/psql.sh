#!/usr/bin/env sh

if [ -z $($CONTAINER ps | grep $DB_CONTAINER) ]; then
	. ./.env.testing
fi

$CONTAINER exec -ti $DB_CONTAINER psql -U $DB_USER $DB_NAME
