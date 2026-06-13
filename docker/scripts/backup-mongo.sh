#!/bin/bash

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

set -a
source "$BASE_DIR/.env"
set +a

mkdir -p "$BASE_DIR/backups"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
FILE="$BASE_DIR/backups/mongo_$TIMESTAMP.archive.gz"

docker exec crescendo-music-service-mongodb-1 \
mongodump \
  --host "$MONGODB_HOST" \
  --username "$MONGODB_USER" \
  --password "$MONGODB_PASSWORD" \
  --authenticationDatabase admin \
  --archive \
  --gzip > "$FILE"

echo "✔ Backup MongoDB creado: $FILE"