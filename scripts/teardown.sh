#!/usr/bin/env sh

. ./.env.testing

$CONTAINER stop $DB_CONTAINER
