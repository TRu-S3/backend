# デプロイメント・運用ガイド

## 概要

TRu-S3 Backendの本番環境デプロイメント、CI/CD、監視、運用に関する包括的なガイドです。

## インフラストラクチャ構成

### 現在の本番環境
- **Cloud Run**: アプリケーションサーバー
- **Cloud SQL**: PostgreSQL 17データベース
- **Google Cloud Storage**: ファイルストレージ
- **Cloud SQL Auth Proxy**: セキュアDB接続
- **Cloud Build**: CI/CDパイプライン

### アーキテクチャ図
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client        │───▶│   Cloud Run     │───▶│   Cloud SQL     │
│   (Frontend)    │    │   (Backend)     │    │   (PostgreSQL)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │       GCS       │
                       │  (File Storage) │
                       └─────────────────┘
```

## デプロイメント手順

### 1. 手動デプロイ

#### 前提条件確認
```bash
# gcloud認証確認
gcloud auth list
gcloud config get-value project

# Docker動作確認
docker --version

# 権限確認
gcloud projects get-iam-policy zenn-ai-agent-hackathon-460205
```

#### ビルド・デプロイ
```bash
# 1. イメージビルド
docker build -t gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend:latest .

# 2. イメージテスト（ローカル）
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5434 \
  gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend:latest

# 3. イメージプッシュ
docker push gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend:latest

# 4. Cloud Runデプロイ
gcloud run deploy tru-s3-backend \
  --image gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend:latest \
  --platform managed \
  --region asia-northeast1 \
  --allow-unauthenticated \
  --port 8080 \
  --memory 512Mi \
  --cpu 1 \
  --concurrency 100 \
  --max-instances 10 \
  --add-cloudsql-instances zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db \
  --set-env-vars "GIN_MODE=release,GOOGLE_CLOUD_PROJECT=zenn-ai-agent-hackathon-460205" \
  --set-env-vars "CLOUD_SQL_CONNECTION_NAME=zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db" \
  --set-env-vars "USE_CLOUD_SQL_PROXY=true,DB_HOST=/cloudsql/zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db" \
  --set-env-vars "DB_NAME=tru_s3,DB_USER=postgres,DB_SSL_MODE=disable"
```

### 2. スクリプトを使用したデプロイ

既存のデプロイスクリプトを使用:
```bash
# デプロイスクリプト実行
./scripts/deploy-cloudrun.sh

# スクリプト内容確認
cat scripts/deploy-cloudrun.sh
```

### 3. 環境別デプロイ設定

#### 開発環境
```bash
gcloud run deploy tru-s3-backend-dev \
  --image gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend:dev \
  --region asia-northeast1 \
  --memory 256Mi \
  --cpu 0.5 \
  --max-instances 3
```

#### ステージング環境
```bash
gcloud run deploy tru-s3-backend-staging \
  --image gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend:staging \
  --region asia-northeast1 \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 5
```

#### 本番環境
```bash
gcloud run deploy tru-s3-backend \
  --image gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend:latest \
  --region asia-northeast1 \
  --memory 1Gi \
  --cpu 2 \
  --max-instances 20 \
  --min-instances 1
```

## CI/CD パイプライン

### 1. GitHub Actions設定

`.github/workflows/ci-cd.yml`:
```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  PROJECT_ID: zenn-ai-agent-hackathon-460205
  SERVICE_NAME: tru-s3-backend
  REGION: asia-northeast1

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:17
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: tru_s3_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.24.2
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      env:
        DB_HOST: localhost
        DB_PORT: 5432
        DB_NAME: tru_s3_test
        DB_USER: postgres
        DB_PASSWORD: postgres
        DB_SSL_MODE: disable
      run: |
        go test -v -cover ./...
    
    - name: Run linter
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest

  deploy:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Cloud SDK
      uses: google-github-actions/setup-gcloud@v1
      with:
        project_id: ${{ env.PROJECT_ID }}
        service_account_key: ${{ secrets.GCP_SA_KEY }}
        export_default_credentials: true
    
    - name: Configure Docker
      run: gcloud auth configure-docker
    
    - name: Build and push Docker image
      run: |
        docker build -t gcr.io/$PROJECT_ID/$SERVICE_NAME:$GITHUB_SHA .
        docker push gcr.io/$PROJECT_ID/$SERVICE_NAME:$GITHUB_SHA
        docker tag gcr.io/$PROJECT_ID/$SERVICE_NAME:$GITHUB_SHA gcr.io/$PROJECT_ID/$SERVICE_NAME:latest
        docker push gcr.io/$PROJECT_ID/$SERVICE_NAME:latest
    
    - name: Deploy to Cloud Run
      run: |
        gcloud run deploy $SERVICE_NAME \
          --image gcr.io/$PROJECT_ID/$SERVICE_NAME:$GITHUB_SHA \
          --platform managed \
          --region $REGION \
          --allow-unauthenticated \
          --memory 512Mi \
          --cpu 1 \
          --max-instances 10 \
          --add-cloudsql-instances $PROJECT_ID:$REGION:prd-db
```

### 2. Cloud Build設定

`cloudbuild.yaml`:
```yaml
steps:
# Test
- name: 'golang:1.24.2'
  entrypoint: 'bash'
  args:
  - '-c'
  - |
    go mod download
    go test -v ./...
  env:
  - 'GO111MODULE=on'

# Build
- name: 'gcr.io/cloud-builders/docker'
  args: ['build', '-t', 'gcr.io/$PROJECT_ID/tru-s3-backend:$SHORT_SHA', '.']

# Push
- name: 'gcr.io/cloud-builders/docker'
  args: ['push', 'gcr.io/$PROJECT_ID/tru-s3-backend:$SHORT_SHA']

# Deploy
- name: 'gcr.io/cloud-builders/gcloud'
  args:
  - 'run'
  - 'deploy'
  - 'tru-s3-backend'
  - '--image=gcr.io/$PROJECT_ID/tru-s3-backend:$SHORT_SHA'
  - '--region=asia-northeast1'
  - '--platform=managed'
  - '--allow-unauthenticated'

images:
- 'gcr.io/$PROJECT_ID/tru-s3-backend:$SHORT_SHA'

options:
  logging: CLOUD_LOGGING_ONLY
```

### 3. 自動トリガー設定
```bash
# GitHub連携トリガー作成
gcloud builds triggers create github \
  --repo-name=backend \
  --repo-owner=TRu-S3 \
  --branch-pattern="^main$" \
  --build-config=cloudbuild.yaml
```

## 監視・ログ

### 1. Cloud Monitoring設定

#### メトリクス監視
```bash
# CPU使用率アラート
gcloud alpha monitoring policies create \
  --policy-from-file=monitoring/cpu-alert.yaml

# メモリ使用率アラート
gcloud alpha monitoring policies create \
  --policy-from-file=monitoring/memory-alert.yaml

# レスポンス時間アラート
gcloud alpha monitoring policies create \
  --policy-from-file=monitoring/latency-alert.yaml
```

#### カスタムメトリクス
```go
// metrics.go
package main

import (
    "contrib.go.opencensus.io/exporter/stackdriver"
    "go.opencensus.io/stats"
    "go.opencensus.io/stats/view"
)

var (
    RequestCount = stats.Int64("request_count", "Number of requests", stats.UnitDimensionless)
    RequestLatency = stats.Float64("request_latency", "Request latency", stats.UnitMilliseconds)
)

func initMetrics() {
    exporter, _ := stackdriver.NewExporter(stackdriver.Options{
        ProjectID: "zenn-ai-agent-hackathon-460205",
    })
    view.RegisterExporter(exporter)
}
```

### 2. ログ管理

#### 構造化ログ
```go
// logging.go
package main

import (
    "encoding/json"
    "log"
    "os"
)

type LogEntry struct {
    Severity  string      `json:"severity"`
    Message   string      `json:"message"`
    Component string      `json:"component"`
    UserID    string      `json:"user_id,omitempty"`
    RequestID string      `json:"request_id,omitempty"`
    Data      interface{} `json:"data,omitempty"`
}

func LogInfo(message string, data interface{}) {
    entry := LogEntry{
        Severity:  "INFO",
        Message:   message,
        Component: "tru-s3-backend",
        Data:      data,
    }
    json.NewEncoder(os.Stdout).Encode(entry)
}
```

#### ログ検索クエリ例
```bash
# エラーログ検索
gcloud logging read 'resource.type="cloud_run_revision" AND severity="ERROR"' --limit=50

# API レスポンス時間監視
gcloud logging read 'resource.type="cloud_run_revision" AND jsonPayload.latency>1000' --limit=20

# 特定ユーザーのアクション
gcloud logging read 'resource.type="cloud_run_revision" AND jsonPayload.user_id="123"' --limit=30
```

### 3. アラート設定

#### Slack通知設定
```bash
# Slack チャンネル設定
gcloud alpha monitoring channels create \
  --display-name="TRu-S3 Alerts" \
  --type=slack \
  --channel-labels=channel_name=#alerts,url=https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
```

#### Email通知設定
```bash
# Email通知設定
gcloud alpha monitoring channels create \
  --display-name="TRu-S3 Admin" \
  --type=email \
  --channel-labels=email_address=admin@example.com
```

## 運用タスク

### 1. 定期メンテナンス

#### データベースメンテナンス
```bash
# Cloud SQL メンテナンス確認
gcloud sql instances describe prd-db --format="value(settings.maintenanceWindow)"

# メンテナンスウィンドウ設定
gcloud sql instances patch prd-db \
  --maintenance-window-day=SUN \
  --maintenance-window-hour=3 \
  --maintenance-release-channel=production
```

#### バックアップ確認
```bash
# 自動バックアップ設定確認
gcloud sql instances describe prd-db --format="value(settings.backupConfiguration)"

# 手動バックアップ作成
gcloud sql backups create --instance=prd-db --description="Manual backup $(date +%Y%m%d)"

# バックアップ一覧
gcloud sql backups list --instance=prd-db
```

### 2. スケーリング

#### 自動スケーリング設定
```bash
# Cloud Run 自動スケーリング
gcloud run services update tru-s3-backend \
  --region=asia-northeast1 \
  --min-instances=1 \
  --max-instances=20 \
  --concurrency=100 \
  --cpu-throttling
```

#### データベーススケーリング
```bash
# Cloud SQL インスタンスタイプ変更
gcloud sql instances patch prd-db \
  --tier=db-custom-2-4096 \
  --availability-type=REGIONAL
```

### 3. セキュリティ

#### IAM監査
```bash
# サービスアカウント監査
gcloud iam service-accounts list

# 権限確認
gcloud projects get-iam-policy zenn-ai-agent-hackathon-460205

# 不要な権限削除
gcloud projects remove-iam-policy-binding PROJECT_ID \
  --member=serviceAccount:old-service@PROJECT_ID.iam.gserviceaccount.com \
  --role=roles/editor
```

#### シークレット管理
```bash
# Secret Manager 使用
echo -n "production_password" | gcloud secrets create db-password --data-file=-

# Cloud Run からシークレット参照
gcloud run services update tru-s3-backend \
  --update-secrets=DB_PASSWORD=db-password:latest
```

## 障害対応

### 1. 障害検知

#### 症状別診断
```bash
# Cloud Run サービス状態確認
gcloud run services describe tru-s3-backend --region=asia-northeast1

# コンテナログ確認
gcloud run services logs read tru-s3-backend --region=asia-northeast1 --limit=100

# データベース接続確認
gcloud sql instances describe prd-db

# GCS接続確認
gsutil ls gs://202506-zenn-ai-agent-hackathon/
```

### 2. 復旧手順

#### サービス再起動
```bash
# 新しいリビジョンデプロイ（強制再起動）
gcloud run deploy tru-s3-backend \
  --image gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend:latest \
  --region asia-northeast1
```

#### ロールバック
```bash
# 以前のリビジョンに戻す
gcloud run services update-traffic tru-s3-backend \
  --to-revisions=tru-s3-backend-00001-abc=100 \
  --region=asia-northeast1
```

#### データベース復旧
```bash
# バックアップからの復元
gcloud sql backups restore BACKUP_ID \
  --restore-instance=prd-db-restore \
  --backup-instance=prd-db
```

### 3. 障害報告

#### インシデント管理
1. **障害検知** → Slack通知
2. **初期対応** → 症状確認・一次対応
3. **詳細調査** → ログ分析・根本原因特定
4. **復旧作業** → 修正・テスト・デプロイ
5. **事後分析** → ポストモーテム作成

## パフォーマンス最適化

### 1. 監視項目
- レスポンス時間: 平均 < 200ms
- CPU使用率: < 70%
- メモリ使用率: < 80%
- データベース接続数: < 70%
- エラー率: < 1%

### 2. 最適化施策
```bash
# Cloud Run 設定最適化
gcloud run services update tru-s3-backend \
  --memory=1Gi \
  --cpu=2 \
  --concurrency=80 \
  --execution-environment=gen2

# データベース最適化
gcloud sql instances patch prd-db \
  --insights-config-query-insights-enabled \
  --insights-config-record-application-tags \
  --insights-config-record-client-address
```

## コスト管理

### 1. コスト監視
```bash
# 予算アラート設定
gcloud billing budgets create \
  --billing-account=BILLING_ACCOUNT_ID \
  --display-name="TRu-S3 Monthly Budget" \
  --budget-amount=100USD \
  --threshold-rule=percent=50,basis=CURRENT_SPEND \
  --threshold-rule=percent=90,basis=CURRENT_SPEND
```

### 2. リソース最適化
- 未使用リビジョンの削除
- ログ保存期間の調整
- 開発環境の自動停止
- Cloud SQL 接続プールの最適化

これらの運用手順により、安定したサービス運用が可能になります。