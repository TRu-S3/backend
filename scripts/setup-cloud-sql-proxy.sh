#!/bin/bash

# Cloud SQL Proxy Setup Script
set -e

echo "🚀 Setting up Cloud SQL Proxy..."

# Variables
PROJECT_ID="zenn-ai-agent-hackathon-460205"
INSTANCE_CONNECTION_NAME="zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db"
PROXY_PORT="5433"

# Download Cloud SQL Proxy if not exists
if [ ! -f "./cloud-sql-proxy" ]; then
    echo "📥 Downloading Cloud SQL Proxy..."
    curl -o cloud-sql-proxy https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.8.0/cloud-sql-proxy.linux.amd64
    chmod +x cloud-sql-proxy
fi

# Create proxy script
cat > start-proxy.sh << EOF
#!/bin/bash
echo "🔌 Starting Cloud SQL Proxy..."
./cloud-sql-proxy ${INSTANCE_CONNECTION_NAME} --port ${PROXY_PORT} --credentials-file=./service-account.json &
PROXY_PID=\$!
echo "Cloud SQL Proxy started with PID: \$PROXY_PID"
echo "Connection available at: localhost:${PROXY_PORT}"
echo "To stop: kill \$PROXY_PID"
EOF

chmod +x start-proxy.sh

# Create environment file for proxy connection
cat > .env.proxy << EOF
# TRu-S3 Backend Configuration via Cloud SQL Proxy

# Server Configuration
PORT=8081
GIN_MODE=debug

# GCP Configuration
GCS_BUCKET_NAME=202506-zenn-ai-agent-hackathon
GCS_FOLDER=test
GOOGLE_CLOUD_PROJECT=696136807010

# Database Configuration (via Cloud SQL Proxy)
DB_HOST=localhost
DB_PORT=${PROXY_PORT}
DB_NAME=tru_s3
DB_USER=postgres
DB_PASSWORD=u6"Ml6%XD7cSg9q]
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5

# Cloud SQL Configuration
CLOUD_SQL_CONNECTION_NAME=${INSTANCE_CONNECTION_NAME}
USE_CLOUD_SQL_PROXY=true

# Google API Configuration
GOOGLE_API_KEY=AIzaSyAbPdDElpGSum6FlDKwOAXR1i5U8KASQ08
EOF

echo "✅ Cloud SQL Proxy setup complete!"
echo "📝 Usage:"
echo "  1. Set up service account credentials"
echo "  2. Run: ./start-proxy.sh"
echo "  3. Use .env.proxy for application connection"