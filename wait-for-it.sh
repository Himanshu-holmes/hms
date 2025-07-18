#!/bin/sh

set -e

host="$1"
port="$2"

until nc -z "$host" "$port"; do
  >&2 echo "Waiting for $host:$port..."
  sleep 1
done

>&2 echo "$host:$port is up"
