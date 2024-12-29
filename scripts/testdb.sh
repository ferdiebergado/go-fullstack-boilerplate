#!/usr/bin/env sh
set -x

$CONTAINER run -d --rm --network host --name $DB_CONTAINER \
	-e POSTGRES_PASSWORD="$DB_PASSWORD" \
	-e POSTGRES_USER="$DB_USER" -e POSTGRES_DB="$DB_NAME" \
	-v ./configs/postgresql.conf:/etc/postgresql/postgresql.conf:Z \
	-v ./configs/psqlrc:/root/.psqlrc:Z \
	$DB_IMAGE -c 'config_file=/etc/postgresql/postgresql.conf'
