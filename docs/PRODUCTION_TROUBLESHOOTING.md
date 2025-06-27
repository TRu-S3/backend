# 本番環境トラブルシューティングガイド

## 概要

本ドキュメントは、TRu-S3バックエンドの本番環境（Cloud Run + Cloud SQL）でのトラブルシューティング手順を記載します。

## 問題の症状

### 発生していた問題
- ✅ ヘルスチェック（`/health`）は正常
- ✅ ファイル関連API（`/api/v1/files`）は正常
- ❌ データベース関連API（`/api/v1/users`, `/api/v1/tags`等）が404エラー

### エラーの原因
1. **環境変数未設定** - Cloud Runサービスにデータベース接続環境変数が設定されていない
2. **Cloud SQL Proxy未設定** - Cloud SQLインスタンスへの接続設定が不完全
3. **データベーススキーマ不整合** - 既存テーブルのデータ型がGORMと互換性がない

## トラブルシューティング手順

### 1. 現状調査

#### 1.1 ログの確認
```bash
# Cloud Runサービスのログを確認
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=backend" \
  --limit=20 --format="table(timestamp,severity,textPayload)" \
  --project=YOUR_PROJECT_ID
```

#### 1.2 環境変数の確認
```bash
# 現在の環境変数を確認
gcloud run services describe backend \
  --region=asia-northeast1 \
  --project=YOUR_PROJECT_ID \
  --format="value(spec.template.spec.containers[0].env[].name,spec.template.spec.containers[0].env[].value)"
```

#### 1.3 Cloud SQLインスタンスの確認
```bash
# Cloud SQLインスタンス一覧
gcloud sql instances list --project=YOUR_PROJECT_ID

# データベース一覧
gcloud sql databases list --instance=INSTANCE_NAME --project=YOUR_PROJECT_ID

# ユーザー一覧
gcloud sql users list --instance=INSTANCE_NAME --project=YOUR_PROJECT_ID
```

### 2. データベース接続の確認

#### 2.1 データベースへの直接接続テスト
```bash
# Cloud SQLへ直接接続
PGPASSWORD='YOUR_PASSWORD' psql -h INSTANCE_IP -U postgres -d DATABASE_NAME

# テーブル一覧確認
SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';

# 特定テーブルのスキーマ確認
SELECT column_name, data_type FROM information_schema.columns 
WHERE table_name = 'TABLE_NAME' ORDER BY ordinal_position;
```

### 3. スキーマ問題の修正

#### 3.1 問題のあるテーブルの特定
```sql
-- file_metadataテーブルのスキーマ確認例
SELECT column_name, data_type FROM information_schema.columns 
WHERE table_name = 'file_metadata' ORDER BY ordinal_position;
```

#### 3.2 スキーマクリーンアップ（必要な場合）
```sql
-- 全テーブル削除（外部キー制約を考慮した順序で）
DROP TABLE IF EXISTS hackathon_participants CASCADE;
DROP TABLE IF EXISTS hackathons CASCADE;
DROP TABLE IF EXISTS contests CASCADE;
DROP TABLE IF EXISTS bookmarks CASCADE;
DROP TABLE IF EXISTS matchings CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;
DROP TABLE IF EXISTS tags CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS file_metadata CASCADE;
DROP TABLE IF EXISTS verification_results CASCADE;
```

### 4. cloudbuild.yamlの修正

#### 4.1 Secret Managerの設定
```bash
# パスワードをSecret Managerに保存
echo 'YOUR_SECURE_PASSWORD' | gcloud secrets create db-password --data-file=- --project=PROJECT_ID

# Cloud BuildとCloud RunにSecret Manager権限を付与
gcloud projects add-iam-policy-binding PROJECT_ID \
  --member="serviceAccount:PROJECT_NUMBER-compute@developer.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"

gcloud projects add-iam-policy-binding PROJECT_ID \
  --member="serviceAccount:PROJECT_NUMBER@cloudbuild.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

#### 4.2 環境変数とSecretの設定
```yaml
# Deploy to Cloud Run
- name: 'gcr.io/cloud-builders/gcloud'
  args: [
    'run', 'deploy', 'backend',
    '--image', 'gcr.io/$PROJECT_ID/tru-s3-backend:$BUILD_ID',
    '--region', 'asia-northeast1',
    '--platform', 'managed',
    '--allow-unauthenticated',
    '--cpu', '1',
    '--memory', '512Mi',
    '--timeout', '300',
    '--concurrency', '80',
    '--max-instances', '10',
    '--set-env-vars', 'GIN_MODE=release,USE_CLOUD_SQL_PROXY=true,CLOUD_SQL_CONNECTION_NAME=PROJECT_ID:REGION:INSTANCE_NAME,DB_NAME=DATABASE_NAME,DB_USER=postgres,DB_SSL_MODE=require,GCS_BUCKET_NAME=BUCKET_NAME,GCS_FOLDER=FOLDER_NAME,GOOGLE_CLOUD_PROJECT=$PROJECT_ID',
    '--set-secrets', 'DB_PASSWORD=db-password:latest',
    '--add-cloudsql-instances', 'PROJECT_ID:REGION:INSTANCE_NAME'
  ]
```

### 5. 再デプロイ

#### 5.1 Cloud Buildでのビルド&デプロイ
```bash
gcloud builds submit --config=cloudbuild.yaml \
  --project=YOUR_PROJECT_ID \
  --substitutions=_DB_PASSWORD='YOUR_PASSWORD'
```

### 6. 動作確認

#### 6.1 基本ヘルスチェック
```bash
curl -s "https://YOUR_BACKEND_URL/health"
curl -s "https://YOUR_BACKEND_URL/"
```

#### 6.2 データベース関連API確認
```bash
# 一覧取得
curl -s "https://YOUR_BACKEND_URL/api/v1/users"
curl -s "https://YOUR_BACKEND_URL/api/v1/tags"

# データ作成
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"Test User","gmail":"test@example.com"}' \
  "https://YOUR_BACKEND_URL/api/v1/users"

curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"Programming","description":"Programming related tag"}' \
  "https://YOUR_BACKEND_URL/api/v1/tags"
```

## 主要な設定値

### 環境変数一覧
```bash
GIN_MODE=release
USE_CLOUD_SQL_PROXY=true
CLOUD_SQL_CONNECTION_NAME=PROJECT_ID:REGION:INSTANCE_NAME
DB_NAME=tru_s3
DB_USER=postgres
DB_PASSWORD=YOUR_SECURE_PASSWORD
DB_SSL_MODE=require
GCS_BUCKET_NAME=YOUR_BUCKET_NAME
GCS_FOLDER=test
GOOGLE_CLOUD_PROJECT=PROJECT_ID
```

### Cloud SQLの設定
- **インスタンス名**: `prd-db`
- **データベース**: `tru_s3`
- **ユーザー**: `postgres`
- **リージョン**: `asia-northeast1`

## よくある問題と解決法

### 問題1: 環境変数が設定されていない
**症状**: `USE_CLOUD_SQL_PROXY=false`のログが出力される
**解決法**: cloudbuild.yamlに環境変数設定を追加して再ビルド

### 問題2: データベース移行エラー
**症状**: `ERROR: cannot cast type bigint to timestamp with time zone`
**解決法**: 問題のあるテーブルを削除して新規作成

### 問題3: Cloud SQL接続エラー
**症状**: `dial tcp 127.0.0.1:5432: connect: connection refused`
**解決法**: `--add-cloudsql-instances`オプションを追加

### 問題4: 制約エラー
**症状**: `constraint "uni_users_gmail" of relation "users" does not exist`
**解決法**: データベースクリーンアップ後に再マイグレーション

## 成功時の確認事項

### ✅ 正常なログの例
```
2025/06/26 17:37:22 Using Cloud SQL Proxy connection: PROJECT_ID:REGION:INSTANCE_NAME
2025/06/26 17:37:22 Database connection established successfully
2025/06/26 17:37:22 Running database migrations...
2025/06/26 17:37:22 Database migrations completed successfully
```

### ✅ 正常なAPI レスポンス
```json
// GET /api/v1/users
{
  "pagination": {"limit": 10, "page": 1, "total": 0},
  "users": []
}

// POST /api/v1/users
{
  "id": 1,
  "name": "Test User",
  "gmail": "test@example.com",
  "created_at": "2025-06-26T17:43:58.772984Z",
  "updated_at": "2025-06-26T17:43:58.772984Z"
}
```

## 関連ファイル

- `cloudbuild.yaml` - ビルド&デプロイ設定
- `internal/config/config.go` - 設定管理
- `internal/database/connection.go` - データベース接続
- `docker-compose.cloudsql.yml` - ローカルCloud SQL開発環境

## セキュリティベストプラクティス

### Secret Managerの使用
1. **パスワード管理**: データベースパスワードはGoogle Cloud Secret Managerで管理
2. **アクセス制御**: 最小権限の原則に従い、必要なサービスアカウントにのみアクセス許可
3. **バージョン管理**: Secretのバージョン管理でローテーション対応
4. **監査ログ**: Secret Managerのアクセスログを監視

### 環境変数の分離
- `.env` - ローカル開発用（gitignoreで除外）
- `.env.production` - 本番設定テンプレート（機密情報なし、git管理可能）
- Secret Manager - 本番の機密情報

## 注意事項

1. **セキュリティ**: パスワードはSecret Managerで管理し、平文でコード管理しない
2. **タイムアウト**: Cloud Runのタイムアウトは300秒に設定
3. **接続制限**: データベース接続プールの設定を適切に行う
4. **ログ監視**: デプロイ後は必ずログを確認して正常性を検証
5. **権限管理**: IAMロールは最小権限の原則に従う

## 参考コマンド集

```bash
# 現在のCloud Runサービス確認
gcloud run services list --region=asia-northeast1

# ビルド履歴確認
gcloud builds list --limit=10

# ログのリアルタイム監視
gcloud logging tail "resource.type=cloud_run_revision" --format="table(timestamp,severity,textPayload)"

# Cloud SQLへの接続確認
gcloud sql connect INSTANCE_NAME --user=postgres --database=DATABASE_NAME

# Secret Manager関連
# Secretの一覧確認
gcloud secrets list

# Secretの値確認（権限が必要）
gcloud secrets versions access latest --secret="SECRET_NAME"

# Secretの作成
echo 'SECRET_VALUE' | gcloud secrets create SECRET_NAME --data-file=-

# Secretのバージョン履歴
gcloud secrets versions list SECRET_NAME

# IAM権限の確認
gcloud projects get-iam-policy PROJECT_ID --flatten="bindings[].members" --format="table(bindings.role,bindings.members)"
```