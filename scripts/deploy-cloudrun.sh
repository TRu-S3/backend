#!/bin/bash

# Cloud Run デプロイスクリプト
# 使用方法: ./deploy-cloudrun.sh

set -e

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 設定（環境変数またはここで直接指定）
PROJECT_ID=${GOOGLE_CLOUD_PROJECT:-"zenn-ai-agent-hackathon-460205"}
REGION=${CLOUD_RUN_REGION:-"asia-northeast1"}
SERVICE_NAME=${CLOUD_RUN_SERVICE:-"tru-s3-backend"}
CLOUD_SQL_CONNECTION_NAME=${CLOUD_SQL_CONNECTION_NAME:-"zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db"}
DB_NAME=${DB_NAME:-"tru_s3"}
DB_USER=${DB_USER:-"postgres"}
DB_PASSWORD_SECRET=${DB_PASSWORD_SECRET:-"tru-s3-db-password"}
GCS_BUCKET_NAME=${GCS_BUCKET_NAME:-"202506-zenn-ai-agent-hackathon"}
GCS_FOLDER=${GCS_FOLDER:-"test"}

echo -e "${GREEN}🚀 TRu-S3 Backend - Cloud Run デプロイ開始${NC}"
echo "プロジェクト: $PROJECT_ID"
echo "リージョン: $REGION"
echo "サービス名: $SERVICE_NAME"
echo "Cloud SQL接続名: $CLOUD_SQL_CONNECTION_NAME"
echo ""

# 必要な設定の確認
if [ "$PROJECT_ID" = "your-project-id" ]; then
    echo -e "${RED}❌ GOOGLE_CLOUD_PROJECT 環境変数を設定してください${NC}"
    exit 1
fi

if [ "$CLOUD_SQL_CONNECTION_NAME" = "your-project:asia-northeast1:your-instance" ]; then
    echo -e "${RED}❌ CLOUD_SQL_CONNECTION_NAME 環境変数を設定してください${NC}"
    exit 1
fi

# gcloudのプロジェクト設定確認
echo -e "${YELLOW}📋 プロジェクト設定を確認中...${NC}"
gcloud config set project $PROJECT_ID

# 必要なAPIの有効化
echo -e "${YELLOW}🔧 必要なAPIを有効化中...${NC}"
gcloud services enable cloudbuild.googleapis.com \
    run.googleapis.com \
    sql-component.googleapis.com \
    sqladmin.googleapis.com \
    storage-api.googleapis.com \
    storage-component.googleapis.com

# Cloud Buildでイメージをビルドしてデプロイ
echo -e "${YELLOW}🏗️  Cloud Buildでイメージをビルド・デプロイ中...${NC}"
gcloud builds submit \
    --config=cloudbuild.yaml \
    --substitutions=_CLOUD_SQL_CONNECTION_NAME=$CLOUD_SQL_CONNECTION_NAME,_DB_NAME=$DB_NAME,_DB_USER=$DB_USER,_DB_PASSWORD_SECRET=$DB_PASSWORD_SECRET,_GCS_BUCKET_NAME=$GCS_BUCKET_NAME,_GCS_FOLDER=$GCS_FOLDER

# デプロイ完了後の情報取得
echo -e "${YELLOW}📊 デプロイ情報を取得中...${NC}"
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region=$REGION --format="value(status.url)")

echo -e "${GREEN}✅ デプロイ完了！${NC}"
echo ""
echo "🌐 サービスURL: $SERVICE_URL"
echo "🏥 ヘルスチェック: $SERVICE_URL/health"
echo ""
echo -e "${YELLOW}📝 使用例:${NC}"
echo "curl $SERVICE_URL/health"
echo "curl $SERVICE_URL/api/files"
echo ""
echo -e "${YELLOW}🔍 ログ確認:${NC}"
echo "gcloud logs tail $SERVICE_NAME --format='default'"
echo ""
echo -e "${YELLOW}🛠️  追加の設定が必要な場合:${NC}"
echo "1. Cloud SQL インスタンスが起動していることを確認"
echo "2. Secret Manager に DB_PASSWORD が設定されていることを確認:"
echo "   gcloud secrets create $DB_PASSWORD_SECRET --data-file=-"
echo "3. Cloud Run サービスアカウントに必要な権限が付与されていることを確認"
echo "   - Cloud SQL Client"
echo "   - Storage Object Admin"
echo "   - Secret Manager Secret Accessor"
