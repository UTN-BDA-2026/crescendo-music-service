#!/bin/bash

# Ir al directorio base (donde está .env)
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

set -a
source "$BASE_DIR/.env"
set +a

mkdir -p "$BASE_DIR/backups"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
FILE="$BASE_DIR/backups/postgres_$TIMESTAMP.sql.gz"

docker exec crescendo-music-service-postgresql-1 \
pg_dump -U "$POSTGRES_USER" "$POSTGRES_DB" | gzip > "$FILE"

echo "✔ Backup PostgreSQL creado: $FILE"