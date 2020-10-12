#!/bin/sh

set -e

host="$1"
shift

until PGUSER=$DB_USER PGPASSWORD=$DB_PASSWORD PGDATABASE=$DB_NAME psql -h "$host" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 5
done

>&2 echo "Postgres is up - executing command"
exec /migrate -path /migrations -database $DATABASE_URL up