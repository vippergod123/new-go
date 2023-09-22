#!/bin/sh

set -e

echo "Run DB Migration $dbSource"
/app/migrate -path /app/migration -database "$dbSource" -verbose up

echo "Start server"
exec "$@"