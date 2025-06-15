#!/bin/bash

# Cloud SQL Auth Proxy 自動セットアップスクリプト
# 使用方法: ./setup-cloud-sql-complete.sh

set -e

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 TRu-S3 Cloud SQL Auth Proxy 完全セットアップ${NC}"
echo "このスクリプトは以下の手順を自動実行します："
echo "1. Cloud SQLインスタンス情報確認"
echo "2. PostgreSQLクライアントインストール"
echo "3. Cloud SQL Auth Proxyセットアップ"
echo "4. パスワード設定"
echo "5. SSL設定調整"
echo "6. 接続テスト"
echo ""

read -p "続行しますか? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "セットアップをキャンセルしました。"
    exit 1
fi

echo -e "${YELLOW}📋 Step 1: Cloud SQLインスタンス情報確認${NC}"
if ! ./check-cloudsql.sh; then
    echo -e "${RED}❌ Cloud SQLインスタンスの確認に失敗しました${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}📦 Step 2: PostgreSQLクライアント確認・インストール${NC}"
if ! command -v psql &> /dev/null; then
    echo "PostgreSQLクライアントをインストールします..."
    if ! curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg; then
        echo -e "${RED}❌ PostgreSQL GPGキーの追加に失敗しました${NC}"
        exit 1
    fi
    echo "deb https://apt.postgresql.org/pub/repos/apt jammy-pgdg main" | sudo tee /etc/apt/sources.list.d/pgdg.list > /dev/null
    sudo apt update -qq
    sudo apt install -y postgresql-client
    echo -e "${GREEN}✅ PostgreSQLクライアントをインストールしました${NC}"
else
    echo -e "${GREEN}✅ PostgreSQLクライアントは既にインストール済みです${NC}"
fi

echo ""
echo -e "${YELLOW}🔧 Step 3: Cloud SQL Auth Proxyセットアップ${NC}"
if ! ./setup-cloudsql-proxy.sh; then
    echo -e "${RED}❌ Cloud SQL Auth Proxyのセットアップに失敗しました${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}🔐 Step 4: パスワード設定${NC}"
read -s -p "Cloud SQLのpostgresユーザーのパスワードを入力してください: " DB_PASSWORD
echo ""

if [ ${#DB_PASSWORD} -lt 8 ]; then
    echo -e "${RED}❌ パスワードは8文字以上である必要があります${NC}"
    exit 1
fi

echo "Cloud SQLにパスワードを設定中..."
if ! gcloud sql users set-password postgres \
    --instance=prd-db \
    --password="$DB_PASSWORD" \
    --project=zenn-ai-agent-hackathon-460205; then
    echo -e "${RED}❌ Cloud SQLパスワードの設定に失敗しました${NC}"
    exit 1
fi

echo ".env.localファイルにパスワードを設定中..."
sed -i "s/^DB_PASSWORD=.*/DB_PASSWORD=\"$DB_PASSWORD\"/" .env.local
echo -e "${GREEN}✅ パスワードを設定しました${NC}"

echo ""
echo -e "${YELLOW}🔧 Step 5: SSL設定調整${NC}"
sed -i "s/^DB_SSL_MODE=.*/DB_SSL_MODE=disable/" .env.local
echo -e "${GREEN}✅ SSL設定をdisableに変更しました${NC}"

echo ""
echo -e "${YELLOW}🧪 Step 6: 接続テスト${NC}"
echo "環境変数を読み込み中..."
source .env.local

echo "Cloud SQL Auth Proxyを起動中..."
./cloud-sql-proxy $CLOUD_SQL_CONNECTION_NAME --port=5433 &
PROXY_PID=$!

# プロキシの起動を待機
sleep 3

echo "psqlで接続テスト中..."
if PGPASSWORD="$DB_PASSWORD" psql "host=localhost port=5433 sslmode=disable user=postgres dbname=postgres" -c "SELECT version();" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ psql接続テスト成功${NC}"
else
    echo -e "${RED}❌ psql接続テスト失敗${NC}"
    kill $PROXY_PID 2>/dev/null
    exit 1
fi

echo "アプリケーションを起動中..."
go run main.go &
APP_PID=$!

# アプリケーションの起動を待機
sleep 3

echo "API接続テスト中..."
if curl -s http://localhost:8080/health | grep -q "ok"; then
    echo -e "${GREEN}✅ API接続テスト成功${NC}"
else
    echo -e "${RED}❌ API接続テスト失敗${NC}"
    kill $PROXY_PID $APP_PID 2>/dev/null
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 セットアップ完了！${NC}"
echo ""
echo -e "${BLUE}📊 実行中のサービス:${NC}"
echo "• Cloud SQL Auth Proxy: PID $PROXY_PID (ポート5433)"
echo "• Go Application: PID $APP_PID (ポート8080)"
echo ""
echo -e "${BLUE}🔧 確認コマンド:${NC}"
echo "• ヘルスチェック: curl http://localhost:8080/health"
echo "• ファイル一覧: curl http://localhost:8080/api/v1/files"
echo "• psql接続: PGPASSWORD=\"$DB_PASSWORD\" psql \"host=localhost port=5433 sslmode=disable user=postgres dbname=tru_s3\""
echo ""
echo -e "${BLUE}🛑 停止方法:${NC}"
echo "• kill $PROXY_PID $APP_PID"
echo "• または Ctrl+C で現在のターミナルを終了"
echo ""
echo -e "${YELLOW}📝 次のステップ:${NC}"
echo "• 詳細な使用方法: README.md を参照"
echo "• 本番デプロイ: make cloudrun-deploy"
echo "• 詳細ドキュメント: CLOUD_SQL_SETUP.md を参照"
