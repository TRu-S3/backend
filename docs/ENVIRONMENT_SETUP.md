# 環境構築ガイド

## 概要

TRu-S3 Backend の開発環境構築からGCP本番環境デプロイまでの包括的なガイドです。

## 前提条件

### 必須ツール
- **Go**: 1.24.2以降
- **PostgreSQL**: 17.x (ローカル開発用)
- **Git**: バージョン管理
- **Docker**: コンテナ化（オプション）

### GCP関連ツール
- **gcloud CLI**: Google Cloud SDKの最新版
- **Cloud SQL Auth Proxy**: データベース接続用

### 推奨ツール
- **VS Code**: Go拡張機能付き
- **Postman**: API テスト用
- **pgAdmin**: PostgreSQL GUI管理ツール

## クイックスタート

### 1. リポジトリクローン
```bash
git clone https://github.com/TRu-S3/backend.git
cd TRu-S3
```

### 2. 依存関係インストール
```bash
go mod download
```

### 3. 環境変数設定
`.env`ファイルを作成:
```bash
# サーバー設定
PORT=8080
GIN_MODE=debug

# GCP設定
GCS_BUCKET_NAME=202506-zenn-ai-agent-hackathon
GCS_FOLDER=test
GOOGLE_CLOUD_PROJECT=zenn-ai-agent-hackathon-460205

# ローカルデータベース設定
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tru_s3
DB_USER=postgres
DB_PASSWORD=your_password
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5

# GCP Cloud SQL設定（本番用）
CLOUD_SQL_CONNECTION_NAME=zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db
USE_CLOUD_SQL_PROXY=false
```

### 4. ローカルデータベース起動
```bash
# PostgreSQL起動（macOS）
brew services start postgresql

# PostgreSQL起動（Linux）
sudo systemctl start postgresql

# データベース作成
createdb tru_s3
```

### 5. アプリケーション起動
```bash
go run main.go
```

サーバーが http://localhost:8080 で起動します。

## 詳細な環境構築

### ローカル開発環境

#### 1. Go環境設定
```bash
# Go version確認
go version

# Go 1.24.2以降が必要
# アップデートが必要な場合:
# https://golang.org/dl/ からダウンロード
```

#### 2. PostgreSQL セットアップ

**macOS (Homebrew):**
```bash
# PostgreSQL インストール
brew install postgresql@17

# サービス開始
brew services start postgresql@17

# データベース作成
createdb tru_s3

# 接続テスト
psql tru_s3
```

**Linux (Ubuntu/Debian):**
```bash
# PostgreSQL インストール
sudo apt-get update
sudo apt-get install postgresql-17 postgresql-client-17

# PostgreSQL開始
sudo systemctl start postgresql
sudo systemctl enable postgresql

# データベース作成
sudo -u postgres createdb tru_s3
```

**Windows:**
```bash
# PostgreSQL公式インストーラーを使用
# https://www.postgresql.org/download/windows/

# または Chocolatey
choco install postgresql
```

#### 3. 環境変数設定の詳細

**必須環境変数:**
```bash
# .envファイル例（ローカル開発）
PORT=8080
GIN_MODE=debug

# データベース
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tru_s3
DB_USER=postgres
DB_PASSWORD=your_local_password
DB_SSL_MODE=disable

# GCP設定（開発時は不要）
GCS_BUCKET_NAME=test-bucket
GCS_FOLDER=dev
```

**本番環境変数:**
```bash
# .envファイル例（本番）
PORT=8080
GIN_MODE=release

# Cloud SQL
DB_HOST=localhost  # プロキシ使用時
DB_PORT=5434
DB_NAME=tru_s3
DB_USER=postgres
DB_PASSWORD=production_password
DB_SSL_MODE=disable  # プロキシ使用時
CLOUD_SQL_CONNECTION_NAME=project:region:instance
USE_CLOUD_SQL_PROXY=true

# GCP
GOOGLE_CLOUD_PROJECT=your-project-id
GCS_BUCKET_NAME=your-bucket
GCS_FOLDER=production
```

### GCP開発環境

#### 1. Google Cloud SDK セットアップ
```bash
# gcloud CLI インストール
curl https://sdk.cloud.google.com | bash
exec -l $SHELL

# 認証
gcloud auth login
gcloud auth application-default login

# プロジェクト設定
gcloud config set project zenn-ai-agent-hackathon-460205
```

#### 2. Cloud SQL Auth Proxy セットアップ
```bash
# Cloud SQL Auth Proxy ダウンロード
curl -o cloud-sql-proxy https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64
chmod +x cloud-sql-proxy

# プロキシ起動
./cloud-sql-proxy zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db --port=5434 &

# 接続テスト
PGPASSWORD='your_password' psql -h localhost -p 5434 -U postgres -d tru_s3
```

#### 3. GCS 設定
```bash
# バケット作成（必要に応じて）
gsutil mb gs://your-bucket-name

# 権限設定
gsutil iam ch serviceAccount:your-service-account@project.iam.gserviceaccount.com:objectAdmin gs://your-bucket-name
```

## テスト環境

### 単体テスト実行
```bash
# 全テスト実行
go test ./...

# カバレッジ付きテスト
go test -cover ./...

# 詳細なテストレポート
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 統合テスト
```bash
# テスト用データベース設定
export DB_NAME=tru_s3_test

# テスト実行
go test -tags=integration ./...
```

### API テスト
```bash
# サーバー起動
go run main.go

# ヘルスチェック
curl http://localhost:8080/health

# API テスト例
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "テストユーザー", "gmail": "test@example.com"}'
```

## 本番環境デプロイ

### 1. Cloud Run デプロイ
```bash
# Docker イメージビルド
docker build -t gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend .

# イメージプッシュ
docker push gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend

# Cloud Run デプロイ
gcloud run deploy tru-s3-backend \
  --image gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend \
  --platform managed \
  --region asia-northeast1 \
  --allow-unauthenticated \
  --add-cloudsql-instances zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db
```

### 2. 自動デプロイ（GitHub Actions）
```yaml
# .github/workflows/deploy.yml
name: Deploy to Cloud Run

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Setup Cloud SDK
      uses: google-github-actions/setup-gcloud@v0
      with:
        project_id: zenn-ai-agent-hackathon-460205
        service_account_key: ${{ secrets.GCP_SA_KEY }}
        export_default_credentials: true
    
    - name: Build and Deploy
      run: |
        gcloud builds submit --tag gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend
        gcloud run deploy tru-s3-backend \
          --image gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend \
          --platform managed \
          --region asia-northeast1
```

## トラブルシューティング

### よくある問題と解決策

#### 1. データベース接続エラー
```bash
# 問題: connection refused
# 解決: PostgreSQLサービス確認
brew services restart postgresql  # macOS
sudo systemctl restart postgresql  # Linux

# 問題: authentication failed
# 解決: パスワード確認
psql -U postgres -h localhost -p 5432
\password postgres
```

#### 2. Cloud SQL接続エラー
```bash
# 問題: timeout expired
# 解決: Auth Proxy確認
ps aux | grep cloud-sql-proxy
./cloud-sql-proxy --help

# 問題: permission denied
# 解決: 認証確認
gcloud auth list
gcloud auth application-default login
```

#### 3. GCS権限エラー
```bash
# 問題: access denied
# 解決: IAM権限確認
gcloud projects get-iam-policy zenn-ai-agent-hackathon-460205
gcloud projects add-iam-policy-binding zenn-ai-agent-hackathon-460205 \
  --member="user:your-email@example.com" \
  --role="roles/storage.admin"
```

#### 4. メモリ不足エラー
```bash
# 問題: out of memory
# 解決: 接続プール調整
export DB_MAX_OPEN_CONNS=10
export DB_MAX_IDLE_CONNS=2
```

### デバッグ方法

#### 1. ログレベル設定
```bash
# 詳細ログ有効化
export GIN_MODE=debug
export LOG_LEVEL=debug
```

#### 2. データベース接続デバッグ
```bash
# 接続テストツール実行
go run scripts/check_tables.go
```

#### 3. API レスポンス確認
```bash
# ヘルスチェック
curl -v http://localhost:8080/health

# 詳細レスポンス
curl -v http://localhost:8080/api/v1/users
```

## 開発ツール

### 推奨VS Code拡張機能
- Go (Google)
- PostgreSQL (Chris Kolkman)
- REST Client (Huachao Mao)
- GitLens (GitKraken)

### 有用なコマンド
```bash
# コード整形
go fmt ./...

# 依存関係整理
go mod tidy

# セキュリティチェック
go list -json -m all | nancy sleuth

# パフォーマンス測定
go test -bench=. -benchmem
```

### データベース管理
```bash
# マイグレーション実行
go run main.go  # 自動マイグレーション

# データベースリセット
dropdb tru_s3 && createdb tru_s3

# バックアップ
pg_dump tru_s3 > backup.sql

# リストア
psql tru_s3 < backup.sql
```

## 本番運用注意事項

### セキュリティ
- 環境変数にシークレットを含めない
- HTTPS通信を強制
- CORS設定を適切に行う
- SQL Injection対策を確認

### パフォーマンス
- 接続プールサイズを調整
- インデックスを適切に設定
- レスポンス時間を監視

### 監視
- ログ集約システム導入
- メトリクス収集
- アラート設定

これで開発環境から本番環境まで、包括的な環境構築が可能です。