#!/bin/sh
set -e

echo "Running database migrations..."
goose -dir /app/migrations postgres "$DB_URL" up

echo "Starting application..."
exec "$@"
