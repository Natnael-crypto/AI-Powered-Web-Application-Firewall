#!/bin/bash

# Create Docker network if it doesn't exist
if ! docker network ls --format '{{.Name}}' | grep -wq waf_network; then
    echo "[+] Creating Docker network: waf_network"
    docker network create waf_network
else
    echo "[*] Docker network 'waf_network' already exists"
fi

echo "[+] Starting PostgreSQL..."
docker run -d \
  --name postgres \
  --network waf_network \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=waf_db \
  natnaelcrypto/waf-db:latest

sleep 15

echo "[+] Starting Backend (no external exposure)..."
docker run -d \
  --name backend \
  --network waf_network \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_NAME=waf_db \
  -e WSKEY=lknlsdkclksdlcnlsdkvjnfjnvljlkcmlskmskdnfjnsdkfs \
  natnaelcrypto/waf-backend:latest

echo "[+] Starting Frontend (localhost only)..."
docker run -d \
  --name frontend \
  --network waf_network \
  --publish 127.0.0.1:5173:3000 \
  natnaelcrypto/waf-frontend:latest

echo "[+] Starting Ml API (internal only)..."
docker run -d \
  --name Ml \
  --network waf_network \
  natnaelcrypto/ml:latest

echo "[+] Starting Interceptor (public on port 80)..."
docker run -d \
  --name interceptor \
  --network waf_network \
  -e BACKENDHOST=backend \
  -e BACKENDPORT=8080 \
  -e MLHOST=cnnapi \
  -e MLPORT=5000 \
  -e WSKEY=lknlsdkclksdlcnlsdkvjnfjnvljlkcmlskmskdnfjnsdkfs \
  -p 80:80 \
  natnaelcrypto/interceptor:latest

echo "[✓] All containers started securely."
