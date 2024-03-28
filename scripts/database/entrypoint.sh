#!/bin/sh

# stop on errors
set -e

if [ ! "$SKIP_MIGRATIONS" = "true" ]; then
    migrate -path /app/scripts/database/migration -database "$DATABASE_URL?sslmode=disable" up
else
    echo "No DB Migrations"
fi