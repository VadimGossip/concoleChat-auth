#!/bin/bash
source .env

export MIGRATION_DSN="host=$DB_HOST port=$DB_PORT dbname=$DB_NAME user=$DB_USERNAME password=$DB_PASSWORD sslmode=$DB_SSLMODE"

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v