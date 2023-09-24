#!/bin/sh

set -e

echo "Run DB Migration $dbSource"
source /app/.prod.env
/app/migrate -path /app/migration -database "$dbSource" -verbose up

echo "Start server"
exec "$@"