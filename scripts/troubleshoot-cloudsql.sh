#!/bin/bash

# Cloud SQL トラブルシューティングスクリプト
# 使用方法: ./troubleshoot-cloudsql.sh

set -e

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Cloud SQL トラブルシューティング${NC}"
echo "システム状況を確認しています..."
echo ""

# 1. gcloud認証確認
echo -e "${YELLOW}1. gcloud認証状況${NC}"
if gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
    ACCOUNT=$(gcloud auth list --filter=status:ACTIVE --format="value(account)")
    echo -e "${GREEN}✅ 認証済み: $ACCOUNT${NC}"
else
    echo -e "${RED}❌ gcloud認証が必要です${NC}"
    echo "解決方法: gcloud auth login"
    echo ""
fi

# 2. プロジェクト設定確認
echo -e "${YELLOW}2. プロジェクト設定${NC}"
PROJECT_ID=$(gcloud config get-value project 2>/dev/null || echo "")
if [ "$PROJECT_ID" = "zenn-ai-agent-hackathon-460205" ]; then
    echo -e "${GREEN}✅ 正しいプロジェクト: $PROJECT_ID${NC}"
else
    echo -e "${RED}❌ プロジェクトが正しくありません: $PROJECT_ID${NC}"
    echo "解決方法: gcloud config set project zenn-ai-agent-hackathon-460205"
fi
echo ""

# 3. Cloud SQL Auth Proxyプロセス確認
echo -e "${YELLOW}3. Cloud SQL Auth Proxyプロセス${NC}"
if ps aux | grep cloud-sql-proxy | grep -v grep > /dev/null; then
    echo -e "${GREEN}✅ Cloud SQL Auth Proxy実行中${NC}"
    ps aux | grep cloud-sql-proxy | grep -v grep
else
    echo -e "${RED}❌ Cloud SQL Auth Proxyが実行されていません${NC}"
    echo "解決方法: ./cloud-sql-proxy \$CLOUD_SQL_CONNECTION_NAME --port=5433 &"
fi
echo ""

# 4. ポート使用状況確認
echo -e "${YELLOW}4. ポート使用状況${NC}"
if netstat -ln | grep :5433 > /dev/null; then
    echo -e "${GREEN}✅ ポート5433使用中${NC}"
    netstat -ln | grep :5433
else
    echo -e "${RED}❌ ポート5433が使用されていません${NC}"
    echo "解決方法: Cloud SQL Auth Proxyを起動してください"
fi

if netstat -ln | grep :8080 > /dev/null; then
    echo -e "${GREEN}✅ ポート8080使用中${NC}"
    netstat -ln | grep :8080
else
    echo -e "${RED}❌ ポート8080が使用されていません${NC}"
    echo "解決方法: go run main.go でアプリケーションを起動してください"
fi
echo ""

# 5. 環境変数確認
echo -e "${YELLOW}5. 環境変数確認${NC}"
if [ -f ".env.local" ]; then
    echo -e "${GREEN}✅ .env.localファイル存在${NC}"
    source .env.local 2>/dev/null || true
    
    if [ -n "$CLOUD_SQL_CONNECTION_NAME" ]; then
        echo -e "${GREEN}✅ CLOUD_SQL_CONNECTION_NAME: $CLOUD_SQL_CONNECTION_NAME${NC}"
    else
        echo -e "${RED}❌ CLOUD_SQL_CONNECTION_NAME が設定されていません${NC}"
    fi
    
    if [ -n "$DB_PASSWORD" ]; then
        echo -e "${GREEN}✅ DB_PASSWORD: 設定済み${NC}"
    else
        echo -e "${RED}❌ DB_PASSWORD が設定されていません${NC}"
    fi
    
    if [ "$DB_SSL_MODE" = "disable" ]; then
        echo -e "${GREEN}✅ DB_SSL_MODE: disable${NC}"
    else
        echo -e "${YELLOW}⚠️  DB_SSL_MODE: $DB_SSL_MODE (推奨: disable)${NC}"
        echo "解決方法: sed -i 's/^DB_SSL_MODE=.*/DB_SSL_MODE=disable/' .env.local"
    fi
else
    echo -e "${RED}❌ .env.localファイルが存在しません${NC}"
    echo "解決方法: ./setup-cloudsql-proxy.sh を実行してください"
fi
echo ""

# 6. Cloud SQLインスタンス確認
echo -e "${YELLOW}6. Cloud SQLインスタンス状況${NC}"
if gcloud sql instances describe prd-db --project=zenn-ai-agent-hackathon-460205 --format="value(state)" 2>/dev/null | grep -q "RUNNABLE"; then
    echo -e "${GREEN}✅ Cloud SQLインスタンス 'prd-db' 実行中${NC}"
else
    echo -e "${RED}❌ Cloud SQLインスタンス 'prd-db' に問題があります${NC}"
    echo "解決方法: Google Cloud Consoleでインスタンス状況を確認してください"
fi
echo ""

# 7. 接続テスト
echo -e "${YELLOW}7. 接続テスト${NC}"

# psql接続テスト
if command -v psql &> /dev/null && [ -f ".env.local" ]; then
    source .env.local 2>/dev/null || true
    if [ -n "$DB_PASSWORD" ] && netstat -ln | grep :5433 > /dev/null; then
        echo "psql接続テスト実行中..."
        if PGPASSWORD="$DB_PASSWORD" psql "host=localhost port=5433 sslmode=disable user=postgres dbname=postgres" -c "SELECT 1;" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ psql接続成功${NC}"
        else
            echo -e "${RED}❌ psql接続失敗${NC}"
            echo "- パスワードが正しいか確認してください"
            echo "- Cloud SQL Auth Proxyが起動しているか確認してください"
        fi
    else
        echo -e "${YELLOW}⚠️  psql接続テストをスキップ（前提条件不足）${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  psqlまたは.env.localが見つかりません${NC}"
fi

# HTTP接続テスト
if netstat -ln | grep :8080 > /dev/null; then
    echo "HTTP接続テスト実行中..."
    if curl -s http://localhost:8080/health | grep -q "ok"; then
        echo -e "${GREEN}✅ HTTP接続成功${NC}"
    else
        echo -e "${RED}❌ HTTP接続失敗${NC}"
        echo "- アプリケーションが正常に起動しているか確認してください"
    fi
else
    echo -e "${YELLOW}⚠️  HTTPテストをスキップ（ポート8080未使用）${NC}"
fi
echo ""

# 8. 推奨アクション
echo -e "${BLUE}🔧 推奨アクション${NC}"
echo "問題が見つかった場合、以下の順序で解決してください："
echo ""
echo "1. 基本認証:"
echo "   gcloud auth login"
echo "   gcloud config set project zenn-ai-agent-hackathon-460205"
echo ""
echo "2. 完全セットアップ:"
echo "   ./setup-cloud-sql-complete.sh"
echo ""
echo "3. 手動起動:"
echo "   source .env.local"
echo "   ./cloud-sql-proxy \$CLOUD_SQL_CONNECTION_NAME --port=5433 &"
echo "   go run main.go &"
echo ""
echo "4. テスト:"
echo "   curl http://localhost:8080/health"
echo "   curl http://localhost:8080/api/v1/files"
echo ""
echo "詳細な手順: CLOUD_SQL_SETUP.md を参照してください"
