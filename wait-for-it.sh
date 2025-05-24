#!/bin/sh
# wait-for-it.sh

set -e

host="$1"
port="$2"
MIGRATIONS_DIR="$3"
DB_DIALECT="$4"
DB_DSN="$5"
GOOSE_COMMAND="$6"
shift 2
cmd="$@"

until nc -z "$host" "$port"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"

echo "Running Goose migrations from: $MIGRATIONS_DIR"
goose -dir "$MIGRATIONS_DIR" "$DB_DIALECT" "$DB_DSN" "$GOOSE_COMMAND"

exec $cmd