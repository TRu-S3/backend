#!/bin/bash

# Cloud SQLインスタンス情報確認スクリプト
# 既存のCloud SQLインスタンスの詳細情報を取得します

set -e

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

PROJECT_ID="zenn-ai-agent-hackathon-460205"

echo -e "${GREEN}🔍 Cloud SQL インスタンス情報確認${NC}"
echo "プロジェクト: $PROJECT_ID"
echo ""

# gcloud認証確認
echo -e "${YELLOW}🔐 gcloud認証を確認中...${NC}"
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
    echo -e "${RED}❌ gcloud にログインしてください: gcloud auth login${NC}"
    exit 1
fi

# プロジェクト設定
gcloud config set project $PROJECT_ID

echo -e "${YELLOW}📋 Cloud SQL インスタンス一覧を取得中...${NC}"
INSTANCES=$(gcloud sql instances list --format="table(name,region,databaseVersion,settings.tier,state)" 2>/dev/null)

if [ -z "$INSTANCES" ] || [ "$INSTANCES" = "Listed 0 items." ]; then
    echo -e "${RED}❌ Cloud SQL インスタンスが見つかりません${NC}"
    echo -e "${YELLOW}💡 Cloud SQL インスタンスを作成するか、プロジェクトIDを確認してください${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Cloud SQL インスタンス一覧:${NC}"
echo "$INSTANCES"
echo ""

# 最初のインスタンスの詳細を取得
FIRST_INSTANCE=$(gcloud sql instances list --format="value(name)" --limit=1 2>/dev/null)

if [ -n "$FIRST_INSTANCE" ]; then
    echo -e "${BLUE}📊 インスタンス '$FIRST_INSTANCE' の詳細情報:${NC}"
    
    # インスタンス詳細
    REGION=$(gcloud sql instances describe $FIRST_INSTANCE --format="value(region)" 2>/dev/null)
    DATABASE_VERSION=$(gcloud sql instances describe $FIRST_INSTANCE --format="value(databaseVersion)" 2>/dev/null)
    TIER=$(gcloud sql instances describe $FIRST_INSTANCE --format="value(settings.tier)" 2>/dev/null)
    
    echo "  インスタンス名: $FIRST_INSTANCE"
    echo "  リージョン: $REGION"
    echo "  データベースバージョン: $DATABASE_VERSION"
    echo "  ティア: $TIER"
    echo ""
    
    # 接続名を生成
    CONNECTION_NAME="$PROJECT_ID:$REGION:$FIRST_INSTANCE"
    echo -e "${GREEN}🔗 Cloud SQL Auth Proxy 接続名:${NC}"
    echo -e "${YELLOW}$CONNECTION_NAME${NC}"
    echo ""
    
    # 環境変数設定例
    echo -e "${BLUE}💡 環境変数設定例:${NC}"
    echo "export GOOGLE_CLOUD_PROJECT=\"$PROJECT_ID\""
    echo "export CLOUD_SQL_CONNECTION_NAME=\"$CONNECTION_NAME\""
    echo ""
    
    # データベース一覧
    echo -e "${YELLOW}📋 データベース一覧を取得中...${NC}"
    DATABASES=$(gcloud sql databases list --instance=$FIRST_INSTANCE --format="table(name)" 2>/dev/null)
    if [ -n "$DATABASES" ]; then
        echo -e "${GREEN}✅ データベース一覧:${NC}"
        echo "$DATABASES"
        echo ""
    fi
    
    # ユーザー一覧
    echo -e "${YELLOW}👥 ユーザー一覧を取得中...${NC}"
    USERS=$(gcloud sql users list --instance=$FIRST_INSTANCE --format="table(name,type)" 2>/dev/null)
    if [ -n "$USERS" ]; then
        echo -e "${GREEN}✅ ユーザー一覧:${NC}"
        echo "$USERS"
        echo ""
    fi
    
    # 次のステップ
    echo -e "${GREEN}🚀 次のステップ:${NC}"
    echo "1. 環境変数を設定:"
    echo "   export CLOUD_SQL_CONNECTION_NAME=\"$CONNECTION_NAME\""
    echo ""
    echo "2. データベースとユーザーを確認・作成 (必要に応じて):"
    echo "   gcloud sql databases create tru_s3 --instance=$FIRST_INSTANCE"
    echo "   gcloud sql users create appuser --instance=$FIRST_INSTANCE --password=your-password"
    echo ""
    echo "3. Cloud SQL Auth Proxyのセットアップ:"
    echo "   ./setup-cloudsql-proxy.sh"
    echo ""
    echo "4. ローカル開発での接続テスト:"
    echo "   ./cloud-sql-proxy $CONNECTION_NAME --port=5433 &"
    echo "   psql \"host=localhost port=5433 sslmode=require user=postgres dbname=postgres\""
    
else
    echo -e "${RED}❌ インスタンス詳細の取得に失敗しました${NC}"
fi
