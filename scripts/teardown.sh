#!/usr/bin/env sh

. ./.env.testing

docker stop $DB_CONTAINER
