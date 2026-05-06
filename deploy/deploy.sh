#!/bin/bash
set -e
SERVER=${1:-"app@your-server"}
echo "=== Build Backend ==="
cd backend && CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server && cd ..
echo "=== Build Frontend ==="
cd web && npm ci && npm run build && cd ..
echo "=== Deploy ==="
scp backend/server $SERVER:/app/backend/
scp -r web/.next web/package.json web/next.config.js web/public $SERVER:/app/web/
echo "=== Restart ==="
ssh $SERVER "
  cd /app/web && npm ci --production
  sudo systemctl restart campaign-api campaign-web
  echo 'Done'
"
