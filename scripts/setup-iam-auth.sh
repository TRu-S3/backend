#!/bin/bash

# IAM Authentication Setup for Cloud SQL
set -e

echo "🔐 Setting up IAM Authentication for Cloud SQL..."

# Variables
PROJECT_ID="zenn-ai-agent-hackathon-460205"
INSTANCE_NAME="prd-db"
SERVICE_ACCOUNT_NAME="cloud-sql-client"
SERVICE_ACCOUNT_EMAIL="${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

# 1. Create service account for Cloud SQL access
echo "👤 Creating service account..."
gcloud iam service-accounts create ${SERVICE_ACCOUNT_NAME} \
    --display-name="Cloud SQL Client Service Account" \
    --description="Service account for secure Cloud SQL access" \
    --project=${PROJECT_ID} || echo "Service account already exists"

# 2. Grant necessary IAM roles
echo "🔑 Granting IAM roles..."
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
    --role="roles/cloudsql.client"

gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
    --role="roles/cloudsql.instanceUser"

# 3. Create and download service account key
echo "🔐 Creating service account key..."
gcloud iam service-accounts keys create ./service-account.json \
    --iam-account=${SERVICE_ACCOUNT_EMAIL} \
    --project=${PROJECT_ID}

# 4. Add database user for IAM authentication
echo "👥 Adding database user for IAM auth..."
gcloud sql users create ${SERVICE_ACCOUNT_EMAIL} \
    --instance=${INSTANCE_NAME} \
    --type=cloud_iam_service_account \
    --project=${PROJECT_ID} || echo "IAM user already exists"

# 5. Create IAM-based environment configuration
cat > .env.iam << EOF
# TRu-S3 Backend Configuration with IAM Authentication

# Server Configuration
PORT=8081
GIN_MODE=debug

# GCP Configuration
GCS_BUCKET_NAME=202506-zenn-ai-agent-hackathon
GCS_FOLDER=test
GOOGLE_CLOUD_PROJECT=${PROJECT_ID}

# Database Configuration (IAM Authentication)
DB_HOST=localhost
DB_PORT=5433
DB_NAME=tru_s3
DB_USER=${SERVICE_ACCOUNT_EMAIL}
DB_PASSWORD=""
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5

# Cloud SQL Configuration
CLOUD_SQL_CONNECTION_NAME=${PROJECT_ID}:asia-northeast1:${INSTANCE_NAME}
USE_CLOUD_SQL_PROXY=true
USE_IAM_AUTH=true

# Google Application Credentials
GOOGLE_APPLICATION_CREDENTIALS=./service-account.json

# Google API Configuration
GOOGLE_API_KEY=AIzaSyAbPdDElpGSum6FlDKwOAXR1i5U8KASQ08
EOF

echo "✅ IAM Authentication setup complete!"
echo "📝 Next steps:"
echo "  1. Use service-account.json for authentication"
echo "  2. Use .env.iam for IAM-based connection"
echo "  3. Start Cloud SQL Proxy with IAM credentials"