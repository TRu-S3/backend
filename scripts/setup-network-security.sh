#!/bin/bash

# Network Security Setup for Cloud SQL
set -e

echo "🛡️ Setting up Network Security for Cloud SQL..."

# Variables
PROJECT_ID="zenn-ai-agent-hackathon-460205"
INSTANCE_NAME="prd-db"

# 1. Configure Cloud SQL instance for private IP (recommended)
echo "🔒 Configuring private IP access..."
gcloud sql instances patch ${INSTANCE_NAME} \
    --no-assign-ip \
    --network=default \
    --project=${PROJECT_ID} || echo "Network configuration may already be set"

# 2. Create VPC connector for serverless access
echo "🌐 Creating VPC connector for serverless..."
gcloud compute networks vpc-access connectors create cloud-sql-connector \
    --region=asia-northeast1 \
    --subnet-project=${PROJECT_ID} \
    --subnet=default \
    --min-instances=2 \
    --max-instances=3 \
    --project=${PROJECT_ID} || echo "VPC connector may already exist"

# 3. Configure SSL/TLS for secure connections
echo "🔐 Configuring SSL requirements..."
gcloud sql instances patch ${INSTANCE_NAME} \
    --require-ssl \
    --project=${PROJECT_ID}

# 4. Download SSL certificates
echo "📜 Downloading SSL certificates..."
mkdir -p ./ssl-certs
gcloud sql ssl-certs create client-cert ./ssl-certs/client-key.pem \
    --instance=${INSTANCE_NAME} \
    --project=${PROJECT_ID} || echo "Client cert may already exist"

gcloud sql ssl-certs describe client-cert \
    --instance=${INSTANCE_NAME} \
    --format="value(cert)" \
    --project=${PROJECT_ID} > ./ssl-certs/client-cert.pem || echo "Could not download client cert"

gcloud sql instances describe ${INSTANCE_NAME} \
    --format="value(serverCaCert.cert)" \
    --project=${PROJECT_ID} > ./ssl-certs/server-ca.pem

# 5. Create secure environment configuration
cat > .env.secure << EOF
# TRu-S3 Backend Configuration with Enhanced Security

# Server Configuration
PORT=8081
GIN_MODE=release

# GCP Configuration
GCS_BUCKET_NAME=202506-zenn-ai-agent-hackathon
GCS_FOLDER=test
GOOGLE_CLOUD_PROJECT=${PROJECT_ID}

# Database Configuration (Secure SSL Connection)
DB_HOST=localhost
DB_PORT=5433
DB_NAME=tru_s3
DB_USER=postgres
DB_PASSWORD=u6"Ml6%XD7cSg9q]
DB_SSL_MODE=require
DB_SSL_CERT=./ssl-certs/client-cert.pem
DB_SSL_KEY=./ssl-certs/client-key.pem
DB_SSL_ROOT_CERT=./ssl-certs/server-ca.pem
DB_MAX_OPEN_CONNS=10
DB_MAX_IDLE_CONNS=2

# Cloud SQL Configuration
CLOUD_SQL_CONNECTION_NAME=${PROJECT_ID}:asia-northeast1:${INSTANCE_NAME}
USE_CLOUD_SQL_PROXY=true
REQUIRE_SSL=true

# Google Application Credentials
GOOGLE_APPLICATION_CREDENTIALS=./service-account.json

# Google API Configuration
GOOGLE_API_KEY=AIzaSyAbPdDElpGSum6FlDKwOAXR1i5U8KASQ08

# Security Settings
TRUSTED_PROXIES=127.0.0.1
CORS_ALLOWED_ORIGINS=https://yourdomain.com
RATE_LIMIT_ENABLED=true
EOF

# 6. Create production deployment configuration
cat > docker-compose.prod.yml << EOF
version: '3.8'

services:
  tru-s3-backend:
    build: .
    ports:
      - "8081:8081"
    environment:
      - ENV_FILE=.env.secure
    volumes:
      - ./ssl-certs:/app/ssl-certs:ro
      - ./service-account.json:/app/service-account.json:ro
    depends_on:
      - cloud-sql-proxy
    networks:
      - secure-network

  cloud-sql-proxy:
    image: gcr.io/cloud-sql-connectors/cloud-sql-proxy:2.8.0
    command: >
      /cloud-sql-proxy
      --private-ip
      --port=5432
      --credentials-file=/credentials/service-account.json
      ${PROJECT_ID}:asia-northeast1:${INSTANCE_NAME}
    volumes:
      - ./service-account.json:/credentials/service-account.json:ro
    networks:
      - secure-network

networks:
  secure-network:
    driver: bridge
EOF

echo "✅ Network Security setup complete!"
echo "🔐 Security features enabled:"
echo "  - Private IP access"
echo "  - SSL/TLS encryption"
echo "  - VPC connector for serverless"
echo "  - Client certificates"
echo "  - Production-ready Docker configuration"