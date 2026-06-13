#!/bin/bash

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

set -a
source "$BASE_DIR/.env"
set +a

if [ -z "$1" ]; then
  echo "Uso: ./backup.sh restore backups/archivo.sql.gz"
  exit 1
fi

FILE="$1"

gunzip < "$FILE" | docker exec -i crescendo-music-service-postgresql-1 \
psql -U "$POSTGRES_USER" "$POSTGRES_DB"

echo "✔ Restore PostgreSQL completado"