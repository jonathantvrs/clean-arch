#!/bin/sh

set -e

echo "Running database migrations..."

sleep 5

migrate -path /migrations -database "$DATABASE_URL" up

echo "Migrations completed successfully!"

exec ./order-service
