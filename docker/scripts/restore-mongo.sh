#!/bin/bash

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

set -a
source "$BASE_DIR/.env"
set +a

if [ -z "$1" ]; then
  echo "Uso: ./backup.sh restore backups/archivo.archive.gz"
  exit 1
fi

FILE="$1"

cat "$FILE" | docker exec -i crescendo-music-service-mongodb-1 \
mongorestore \
  --host "$MONGODB_HOST" \
  --username "$MONGODB_USER" \
  --password "$MONGODB_PASSWORD" \
  --authenticationDatabase admin \
  --drop \
  --archive --gzip

echo "✔ Restore MongoDB completado"