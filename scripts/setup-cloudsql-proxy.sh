#!/bin/bash

# Cloud SQL Auth Proxy ローカル開発セットアップスクリプト
# Cloud SQLインスタンスにローカルから接続するためのセットアップ

set -e

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 設定
PROJECT_ID=${GOOGLE_CLOUD_PROJECT:-"zenn-ai-agent-hackathon-460205"}
CLOUD_SQL_CONNECTION_NAME=${CLOUD_SQL_CONNECTION_NAME:-"zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db"}
LOCAL_PORT=${LOCAL_PORT:-"5433"}

echo -e "${GREEN}🔗 Cloud SQL Auth Proxy ローカル開発セットアップ${NC}"
echo "プロジェクト: $PROJECT_ID"
echo "Cloud SQL接続名: $CLOUD_SQL_CONNECTION_NAME"
echo "ローカルポート: $LOCAL_PORT"
echo ""

# 必要な設定の確認
if [ "$PROJECT_ID" = "zenn-ai-agent-hackathon-460205" ]; then
    echo -e "${GREEN}✅ プロジェクトID: $PROJECT_ID${NC}"
else
    echo -e "${YELLOW}⚠️  カスタムプロジェクトID: $PROJECT_ID${NC}"
fi

if [ "$CLOUD_SQL_CONNECTION_NAME" = "zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db" ]; then
    echo -e "${GREEN}✅ Cloud SQL接続名: $CLOUD_SQL_CONNECTION_NAME${NC}"
else
    echo -e "${YELLOW}⚠️  カスタムCloud SQL接続名: $CLOUD_SQL_CONNECTION_NAME${NC}"
fi

# gcloud認証確認
echo -e "${YELLOW}🔐 gcloud認証を確認中...${NC}"
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
    echo -e "${RED}❌ gcloud にログインしてください: gcloud auth login${NC}"
    exit 1
fi

# プロジェクト設定
gcloud config set project $PROJECT_ID

# Cloud SQL Auth Proxyのダウンロード（必要に応じて）
PROXY_PATH="./cloud-sql-proxy"
if [ ! -f "$PROXY_PATH" ]; then
    echo -e "${YELLOW}📥 Cloud SQL Auth Proxy をダウンロード中...${NC}"
    
    # OS検出
    OS=$(uname -s)
    ARCH=$(uname -m)
    
    if [ "$OS" = "Linux" ]; then
        if [ "$ARCH" = "x86_64" ]; then
            PROXY_URL="https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.8.0/cloud-sql-proxy.linux.amd64"
        else
            PROXY_URL="https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.8.0/cloud-sql-proxy.linux.386"
        fi
    elif [ "$OS" = "Darwin" ]; then
        if [ "$ARCH" = "arm64" ]; then
            PROXY_URL="https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.8.0/cloud-sql-proxy.darwin.arm64"
        else
            PROXY_URL="https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.8.0/cloud-sql-proxy.darwin.amd64"
        fi
    else
        echo -e "${RED}❌ サポートされていないOS: $OS${NC}"
        exit 1
    fi
    
    curl -o "$PROXY_PATH" "$PROXY_URL"
    chmod +x "$PROXY_PATH"
    echo -e "${GREEN}✅ Cloud SQL Auth Proxy をダウンロードしました${NC}"
fi

# .env.local ファイルの生成
ENV_LOCAL_FILE=".env.local"
echo -e "${YELLOW}📝 $ENV_LOCAL_FILE を生成中...${NC}"

cat > "$ENV_LOCAL_FILE" << EOF
# ローカル開発用Cloud SQL Auth Proxy設定
# 使用方法: source .env.local && ./setup-cloudsql-proxy.sh

# Cloud SQL設定
GOOGLE_CLOUD_PROJECT=$PROJECT_ID
CLOUD_SQL_CONNECTION_NAME=$CLOUD_SQL_CONNECTION_NAME
USE_CLOUD_SQL_PROXY=true

# データベース設定（Cloud SQL Auth Proxy経由）
DB_HOST=localhost
DB_PORT=$LOCAL_PORT
DB_NAME=tru_s3
DB_USER=postgres
DB_PASSWORD=your-cloud-sql-password
DB_SSL_MODE=require

# その他の設定
PORT=8080
GIN_MODE=debug
GCS_BUCKET_NAME=202506-zenn-ai-agent-hackathon
GCS_FOLDER=test
EOF

echo -e "${GREEN}✅ セットアップ完了！${NC}"
echo ""
echo -e "${YELLOW}🚀 使用方法:${NC}"
echo "1. Cloud SQLインスタンスのパスワードを設定:"
echo "   export DB_PASSWORD='your-actual-password'"
echo ""
echo "2. 設定を読み込み:"
echo "   source $ENV_LOCAL_FILE"
echo ""
echo "3. Cloud SQL Auth Proxyを起動:"
echo "   $PROXY_PATH $CLOUD_SQL_CONNECTION_NAME --port=$LOCAL_PORT &"
echo ""
echo "4. 別のターミナルでアプリケーションを起動:"
echo "   go run main.go"
echo ""
echo "5. 接続テスト:"
echo "   curl http://localhost:8080/health"
echo ""
echo -e "${YELLOW}📋 注意事項:${NC}"
echo "- Cloud SQLインスタンスが起動していることを確認してください"
echo "- Cloud SQL Admin APIが有効になっていることを確認してください"
echo "- 適切なIAM権限（Cloud SQL Client）が設定されていることを確認してください"
