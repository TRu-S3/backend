# TRu-S3 Backend

**包括的なマッチングアプリケーション + ファイル管理システム**

## 📋 目次

- [概要](#概要)
- [🏗️ アーキテクチャ詳細](#️-アーキテクチャ詳細)
- [🚀 クイックスタート](#-クイックスタート)
- [🛠️ 環境構築詳細](#️-環境構築詳細)
- [☁️ クラウド環境接続](#️-クラウド環境接続)
- [🔧 開発・運用ガイド](#-開発運用ガイド)
- [📊 API仕様](#-api仕様)
- [🔍 トラブルシューティング](#-トラブルシューティング)

## 概要

TRu-S3は、Google Cloud Platform上で動作する包括的なマッチングアプリケーションとファイル管理システムを統合したバックエンドAPIです。Clean Architecture（オニオンアーキテクチャ）を採用し、高い保守性とテスタビリティを実現しています。

## **最新の実装状況**

**完全実装済み機能**:
- **ユーザー管理**: 完全CRUD + 認証準備
- **タグ管理**: プロフィールカテゴリ化
- **プロフィール管理**: 詳細ユーザー情報 + 年齢・地域フィルタ  
- **マッチング機能**: ユーザー間マッチング + ステータス管理
- **ブックマーク機能**: お気に入りユーザー管理
- **コンテスト管理**: プログラミングコンテスト運営
- **ハッカソン管理**: イベント + 参加者管理
- **ファイル管理**: GCS連携ファイルシステム

**データベース**: PostgreSQL 17 on Google Cloud SQL (10テーブル、完全正規化)  
**API エンドポイント**: 45+ の REST API  
**本番環境**: Google Cloud Run + Cloud SQL 運用中

### 主な機能

- **ユーザー管理システム**: ユーザー登録・プロフィール・タグ分類
- **マッチング機能**: ユーザー間マッチング・ブックマーク・承認フロー
- **イベント管理**: ハッカソン・コンテスト運営システム
- **ファイル管理**: GCS連携ファイルアップロード・管理
- **高度検索**: 年齢・地域・タグ・ステータス別フィルタリング

### 技術的特徴

- **オニオンアーキテクチャ**: 保守性・テスタビリティを重視した設計
- **GCP完全対応**: Cloud Storage + Cloud SQL + Cloud Run
- **セキュア**: 複数レベルのセキュリティ設定（IAM認証、SSL/TLS、プライベートIP）
- **Docker対応**: ローカル開発から本番デプロイまで統一環境
- **自動化**: ワンコマンドセットアップ・デプロイ
- **日本語対応**: 包括的な日本語APIドキュメント

### 技術スタック

- **言語**: Go 1.24.2
- **フレームワーク**: Gin (HTTP)、GORM (ORM)
- **データベース**: PostgreSQL 17 (Cloud SQL)
- **ストレージ**: Google Cloud Storage
- **インフラ**: Docker、Cloud Run、Cloud Build
- **認証**: GCP IAM、Cloud SQL Auth Proxy

## 🏗️ アーキテクチャ詳細

### Clean Architecture（オニオンアーキテクチャ）設計

```
┌─────────────────────────────────────────────────────────────┐
│                    外部インターフェース                        │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              インターフェース層                          │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │            アプリケーション層                        │ │ │
│  │  │  ┌─────────────────────────────────────────────────┐ │ │ │
│  │  │  │              ドメイン層                         │ │ │ │
│  │  │  │    ┌─────────────────────────────────────────┐   │ │ │ │
│  │  │  │    │         エンティティ                    │   │ │ │ │
│  │  │  │    │   - User, Profile, File                │   │ │ │ │
│  │  │  │    │   - Matching, Contest, Hackathon       │   │ │ │ │
│  │  │  │    └─────────────────────────────────────────┘   │ │ │ │
│  │  │  │    ┌─────────────────────────────────────────┐   │ │ │ │
│  │  │  │    │      リポジトリインターフェース           │   │ │ │ │
│  │  │  │    │   - FileRepository                     │   │ │ │ │
│  │  │  │    └─────────────────────────────────────────┘   │ │ │ │
│  │  │  └─────────────────────────────────────────────────┘ │ │ │
│  │  │  ┌─────────────────────────────────────────────────┐ │ │ │
│  │  │  │           ユースケース/サービス                  │ │ │ │
│  │  │  │   - FileService                               │ │ │ │
│  │  │  │   - UserService (将来実装)                     │ │ │ │
│  │  │  └─────────────────────────────────────────────────┘ │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │              HTTPハンドラー                         │ │ │
│  │  │   - FileHandler, UserHandler                      │ │ │
│  │  │   - ContestHandler, HackathonHandler              │ │ │
│  │  │   - MatchingHandler, BookmarkHandler              │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │            インフラストラクチャ層                        │ │
│  │   - GCSClient, GCSFileRepository                      │ │
│  │   - PostgreSQL Database                               │ │
│  │   - Cloud SQL Auth Proxy                              │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### ディレクトリ構造詳細

```
internal/
├── domain/                    # ドメイン層
│   ├── file.go               # ファイルエンティティ
│   └── repository.go         # リポジトリインターフェース
├── application/               # アプリケーション層
│   └── file_service.go       # ファイル管理ビジネスロジック
├── infrastructure/            # インフラストラクチャ層
│   ├── gcs_client.go         # Google Cloud Storage クライアント
│   └── gcs_file_repository.go # GCS実装
├── interfaces/                # インターフェース層
│   ├── routes.go             # ルート定義
│   ├── base_handler.go       # 共通ハンドラー機能
│   ├── file_handler.go       # ファイル管理API
│   ├── user_handler.go       # ユーザー管理API
│   ├── contest_handler.go    # コンテスト管理API
│   ├── hackathon_handler.go  # ハッカソン管理API
│   ├── matching_handler.go   # マッチング機能API
│   └── bookmark_handler.go   # ブックマーク機能API
├── database/                  # データベース管理
│   ├── connection.go         # DB接続管理
│   ├── models.go            # 統合モデル定義
│   ├── user/                # ユーザー関連モデル
│   ├── file/                # ファイル関連モデル
│   ├── contest/             # コンテスト関連モデル
│   └── hackathon/           # ハッカソン関連モデル
├── config/                   # 設定管理
│   └── config.go            # 環境変数・設定管理
└── utils/                    # ユーティリティ
    ├── http.go              # HTTP共通処理
    └── validation.go        # バリデーション
```

### 依存関係の流れ

```
外部リクエスト → インターフェース層 → アプリケーション層 → ドメイン層
                     ↓                    ↓
              インフラストラクチャ層 ← ← ← ←
```

**重要な設計原則:**
- **依存関係逆転**: 内側の層は外側の層を知らない
- **単一責任**: 各層は明確な責任を持つ
- **テスタビリティ**: インターフェースによる抽象化でテストが容易
- **拡張性**: 新機能追加時の影響範囲を最小化

## 要件

### 開発環境
- Go 1.24.2+
- GCP アカウント
- GCP認証設定

### Docker環境
- Docker
- Docker Compose

## 🚀 クイックスタート

### 最速セットアップ（推奨）

```bash
# 1. リポジトリクローン
git clone <repository-url>
cd TRu-S3

# 2. Cloud SQL完全自動セットアップ（推奨）
make cloudsql-setup-complete

# 3. 動作確認
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/files
```

### Docker環境での起動

```bash
# 1. 基本セットアップ
make build && make run

# 2. 動作確認
make test
```

## 🛠️ 環境構築詳細

### 利用可能なコマンド

```bash
# ヘルプ表示（全コマンド確認）
make help

# 基本的なセットアップと起動
make build                     # Dockerイメージビルド
make run                       # Docker環境で起動
make dev                       # 開発モード（ログ表示）

# ローカル開発
make local-run                 # ローカル環境で起動

# Cloud SQL使用
make cloudsql-setup-complete   # Cloud SQL完全自動セットアップ
make cloudsql-troubleshoot     # トラブルシューティング
make local-cloudsql           # Cloud SQL使用でローカル起動

# Cloud Run本番環境
make cloudrun-deploy          # Cloud Runにデプロイ

# テスト
make test                     # 基本APIテスト
make test-all                 # 包括的APIテスト
```

### 1. 基本セットアップ（Docker使用）

```bash
# リポジトリのクローン
git clone <repository-url>
cd TRu-S3

# 環境変数の設定
cp .env.example .env
# .env ファイルを編集してGCP設定を入力

# アプリケーションの起動
make build
make run

# 動作確認
curl http://localhost:8080/health
```

### 2. Cloud SQL使用（推奨）

```bash
# 完全自動セットアップ（推奨）
make cloudsql-setup-complete

# または段階的セットアップ
make cloudsql-info              # Cloud SQL情報確認
make cloudsql-setup             # Cloud SQL Auth Proxyセットアップ
# パスワード設定とSSL調整
make local-cloudsql             # アプリケーション起動
```

### 3. 本番環境デプロイ

```bash
# Cloud Runへのデプロイ
make cloudrun-deploy

# 手動ビルド・デプロイ
make cloudrun-build
```

`.env` ファイルの内容例：
```bash
# Server Configuration
PORT=8080
GIN_MODE=debug

# GCP Configuration
GCS_BUCKET_NAME=202506-zenn-ai-agent-hackathon
GCS_FOLDER=test
GOOGLE_CLOUD_PROJECT=your-project-id

# Optional: Service Account Key
# GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
```

### 3. GCP認証

以下のいずれかの方法で認証を設定：

**方法1: gcloud CLI認証**
```bash
gcloud auth application-default login
```

**方法2: サービスアカウントキー**
```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
```

### 4. アプリケーション起動

**Go直接実行:**
```bash
go mod tidy
go run main.go
```

**Docker使用:**
```bash
docker-compose up --build -d
```

## API使用方法

### 基本情報
- **ベースURL**: `http://localhost:8080`
- **Content-Type**: `application/json` (レスポンス), `multipart/form-data` (ファイルアップロード)

### 1. ヘルスチェック

```bash
# サーバー稼働確認
curl http://localhost:8080/health

# レスポンス例
{"status":"ok"}
```

### 2. ファイル一覧取得

```bash
# 全ファイル一覧
curl -s http://localhost:8080/api/v1/files | jq .

# パラメータ付き取得
curl -s "http://localhost:8080/api/v1/files?limit=10&offset=0&prefix=test" | jq .

# レスポンス例
{
  "count": 2,
  "files": [
    {
      "id": "example.txt",
      "name": "example.txt",
      "path": "test/example.txt",
      "size": 1024,
      "content_type": "text/plain",
      "created_at": "2025-06-14T10:00:00Z",
      "updated_at": "2025-06-14T10:00:00Z"
    }
  ]
}
```

### 3. ファイルアップロード

```bash
# テキストファイルのアップロード
echo "Hello, TRu-S3!" > test.txt
curl -X POST -F "file=@test.txt" http://localhost:8080/api/v1/files

# 画像ファイルのアップロード
curl -X POST -F "file=@image.jpg" http://localhost:8080/api/v1/files

# レスポンス例
{
  "id": "test.txt",
  "name": "test.txt",
  "path": "test/test.txt",
  "size": 15,
  "content_type": "text/plain",
  "created_at": "2025-06-14T10:00:00Z",
  "updated_at": "2025-06-14T10:00:00Z"
}
```

### 4. ファイル情報取得

```bash
# ファイルメタデータの取得
curl -s http://localhost:8080/api/v1/files/test.txt | jq .

# レスポンス例
{
  "id": "test.txt",
  "name": "test.txt",
  "path": "test/test.txt",
  "size": 15,
  "content_type": "text/plain",
  "created_at": "2025-06-14T10:00:00Z",
  "updated_at": "2025-06-14T10:00:00Z"
}
```

### 5. ファイルダウンロード

```bash
# ファイル内容のダウンロード
curl -O http://localhost:8080/api/v1/files/test.txt/download

# または標準出力に表示
curl http://localhost:8080/api/v1/files/test.txt/download
```

### 6. ファイル更新

```bash
# ファイル内容の更新
echo "Updated content!" > updated.txt
curl -X PUT -F "file=@updated.txt" http://localhost:8080/api/v1/files/test.txt

# ファイル名の変更
curl -X PUT -F "name=new_name.txt" http://localhost:8080/api/v1/files/test.txt

# ファイル内容と名前の両方を更新
curl -X PUT -F "file=@updated.txt" -F "name=new_name.txt" http://localhost:8080/api/v1/files/test.txt
```

### 7. ファイル削除

```bash
# ファイルの削除
curl -X DELETE http://localhost:8080/api/v1/files/test.txt

# レスポンス例
{"message":"File deleted successfully"}
```

## 使用例

### 完全なワークフロー例

```bash
# 1. サーバー起動確認
curl http://localhost:8080/health

# 2. テストファイル作成
echo "Hello, World!" > hello.txt

# 3. ファイルアップロード
curl -X POST -F "file=@hello.txt" http://localhost:8080/api/v1/files

# 4. ファイル一覧確認
curl -s http://localhost:8080/api/v1/files | jq .

# 5. ファイル内容確認
curl http://localhost:8080/api/v1/files/hello.txt/download

# 6. ファイル更新
echo "Updated Hello, World!" > hello_updated.txt
curl -X PUT -F "file=@hello_updated.txt" http://localhost:8080/api/v1/files/hello.txt

# 7. 更新後の内容確認
curl http://localhost:8080/api/v1/files/hello.txt/download

# 8. ファイル削除
curl -X DELETE http://localhost:8080/api/v1/files/hello.txt

# 9. 削除確認
curl -s http://localhost:8080/api/v1/files | jq .
```

### バッチ処理例

```bash
# 複数ファイルの一括アップロード
for file in *.txt; do
  echo "Uploading $file..."
  curl -X POST -F "file=@$file" http://localhost:8080/api/v1/files
done

# 全ファイルの一括ダウンロード
curl -s http://localhost:8080/api/v1/files | jq -r '.files[].id' | while read filename; do
  echo "Downloading $filename..."
  curl -o "downloaded_$filename" http://localhost:8080/api/v1/files/$filename/download
done
```

## Docker操作

### 基本コマンド

```bash
# イメージをビルド
docker-compose build

# バックグラウンドで起動
docker-compose up -d

# 開発モード（ログを表示しながら起動）
docker-compose up

# ログを確認
docker-compose logs -f

# アプリケーションを停止
docker-compose down

# 完全クリーンアップ
docker-compose down --rmi all --volumes --remove-orphans

# 実行中のコンテナにシェルでアクセス
docker-compose exec app sh
```

## ローカル開発

Dockerを使用せずにローカルで実行する場合：

```bash
# 依存関係をダウンロード
go mod tidy

# アプリケーションを起動
go run main.go

# または、特定のポートで起動
PORT=8081 go run main.go
```

## ☁️ クラウド環境接続

### GCP環境構成

```
┌─────────────────────────────────────────────────────────────┐
│                        GCP プロジェクト                       │
│  ┌─────────────────┐    ┌─────────────────┐    ┌───────────┐ │
│  │   Cloud Run     │───▶│   Cloud SQL     │    │    GCS    │ │
│  │  (Backend API)  │    │  (PostgreSQL)   │    │(Storage)  │ │
│  │                 │    │                 │    │           │ │
│  │ - Auto Scaling  │    │ - Auth Proxy    │    │ - Files   │ │
│  │ - HTTPS/SSL     │    │ - Private IP    │    │ - Backup  │ │
│  │ - IAM Auth      │    │ - SSL/TLS       │    │           │ │
│  └─────────────────┘    └─────────────────┘    └───────────┘ │
│           │                       │                    │     │
│           └───────────────────────┼────────────────────┘     │
│                                   │                          │
│  ┌─────────────────────────────────┼─────────────────────────┐ │
│  │              セキュリティ層      │                         │ │
│  │  - IAM認証                     │                         │ │
│  │  - Cloud SQL Auth Proxy        │                         │ │
│  │  - VPC ネットワーク             │                         │ │
│  │  - Secret Manager               │                         │ │
│  └─────────────────────────────────┼─────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 環境別設定

#### 1. ローカル開発環境
```bash
# 設定ファイル: .env
GIN_MODE=debug
DB_HOST=localhost
DB_PORT=5432
DB_SSL_MODE=disable
USE_CLOUD_SQL_PROXY=false

# 起動方法
make local-run
```

#### 2. Cloud SQL接続環境
```bash
# 設定ファイル: .env.local
GIN_MODE=debug
USE_CLOUD_SQL_PROXY=true
CLOUD_SQL_CONNECTION_NAME=zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db
DB_SSL_MODE=disable

# 起動方法
make cloudsql-setup-complete
```

#### 3. 本番環境（Cloud Run）
```bash
# 自動設定される環境変数
GIN_MODE=release
USE_CLOUD_SQL_PROXY=true
CLOUD_SQL_CONNECTION_NAME=zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db
DB_PASSWORD=<Secret Manager から取得>

# デプロイ方法
make cloudrun-deploy
```

## Cloud SQL Auth Proxy

本番環境では既存のCloud SQL PostgreSQLインスタンス（プロジェクト: `zenn-ai-agent-hackathon-460205`, インスタンス: `prd-db`）を使用してCloud SQL Auth Proxyで接続します。

> **詳細なセットアップ手順は [CLOUD_SQL_SETUP.md](./CLOUD_SQL_SETUP.md) を参照してください**

### クイックセットアップ

```bash
# 1. Cloud SQLインスタンス情報の確認
make cloudsql-info

# 2. 完全自動セットアップ（推奨）
make cloudsql-setup-complete

# 3. テスト
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/files
```

### 手動セットアップ

```bash
# 1. Cloud SQLインスタンス情報の確認
./check-cloudsql.sh

# 2. Cloud SQL Auth Proxyのセットアップ
./setup-cloudsql-proxy.sh

# 3. パスワードの設定
gcloud sql users set-password postgres --instance=prd-db --password="your-password" --project=zenn-ai-agent-hackathon-460205
sed -i "s/^DB_PASSWORD=.*/DB_PASSWORD=\"your-password\"/" .env.local

# 4. SSL設定の調整
sed -i "s/^DB_SSL_MODE=.*/DB_SSL_MODE=disable/" .env.local

# 5. 起動
source .env.local
./cloud-sql-proxy $CLOUD_SQL_CONNECTION_NAME --port=5433 &
go run main.go &

# 6. テスト
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/files
```

### 3. Cloud SQL Auth Proxy用Docker Compose

```bash
# 接続名は既にデフォルト設定済み（prd-db）
export CLOUD_SQL_CONNECTION_NAME="zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db"

# Cloud SQL Auth Proxy用のDocker Compose設定
cp .env.example .env.cloudsql
# .env.cloudsqlを編集してCloud SQL設定を入力

# Cloud SQL Auth Proxy付きで起動
make cloudsql-compose
```

### 4. Cloud Runへのデプロイ

```bash
# 1. 接続名は既にデフォルト設定済み（prd-db）
export CLOUD_SQL_CONNECTION_NAME="zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db"

# 2. Secret Managerにパスワードを保存
echo -n "your-cloud-sql-password" | gcloud secrets create tru-s3-db-password --data-file=-

# 3. Cloud Runにデプロイ
make cloudrun-deploy
```

### 本番環境の設定

Cloud Run環境では以下の環境変数が自動設定されます：

- `USE_CLOUD_SQL_PROXY=true`
- `CLOUD_SQL_CONNECTION_NAME`: Cloud SQLインスタンスの接続名
- `DB_PASSWORD`: Secret Managerから自動取得

### 接続確認

```bash
# ヘルスチェック
curl https://your-service-url/health

# ファイル一覧取得
curl https://your-service-url/api/v1/files
```

## 🔧 開発・運用ガイド

### 新機能開発フロー

#### 1. データベーススキーマ変更

**新しいテーブル追加例:**
```go
// internal/database/[domain]/models.go
type NewEntity struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"not null;size:255" json:"name"`
    Status    string    `gorm:"default:'active'" json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// AutoMigrate関数に追加
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(&NewEntity{})
}
```

**マイグレーション実行:**
```bash
# アプリケーション起動時に自動実行
go run main.go

# または手動実行
go run migrations/migrate.go
```

#### 2. API エンドポイント追加

**Step 1: ハンドラー作成**
```go
// internal/interfaces/new_handler.go
type NewHandler struct {
    db *gorm.DB
}

func NewNewHandler(db *gorm.DB) *NewHandler {
    return &NewHandler{db: db}
}

func (h *NewHandler) CreateNew(c *gin.Context) {
    var req CreateNewRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // ビジネスロジック実装
    entity := &NewEntity{
        Name:   req.Name,
        Status: "active",
    }
    
    if err := h.db.Create(entity).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to create entity"})
        return
    }
    
    c.JSON(201, entity)
}
```

**Step 2: ルート追加**
```go
// internal/interfaces/routes.go
func SetupRoutes(r *gin.Engine, ..., newHandler *NewHandler) {
    v1 := r.Group("/api/v1")
    {
        // 既存ルート...
        
        // 新しいルート
        news := v1.Group("/news")
        {
            news.POST("", newHandler.CreateNew)
            news.GET("", newHandler.ListNews)
            news.GET("/:id", newHandler.GetNew)
            news.PUT("/:id", newHandler.UpdateNew)
            news.DELETE("/:id", newHandler.DeleteNew)
        }
    }
}
```

**Step 3: main.goでの登録**
```go
// main.go
func main() {
    // 既存の初期化...
    
    // 新しいハンドラー作成
    newHandler := interfaces.NewNewHandler(database.GetDB())
    
    // ルート設定
    interfaces.SetupRoutes(r, fileHandler, ..., newHandler)
}
```

#### 3. 設定変更

**環境変数追加:**
```go
// internal/config/config.go
type Config struct {
    // 既存フィールド...
    NewFeatureEnabled bool   `json:"new_feature_enabled"`
    NewServiceURL     string `json:"new_service_url"`
}

func Load() *Config {
    return &Config{
        // 既存設定...
        NewFeatureEnabled: getEnvBoolWithDefault("NEW_FEATURE_ENABLED", false),
        NewServiceURL:     getEnvWithDefault("NEW_SERVICE_URL", ""),
    }
}
```

**バリデーション追加:**
```go
func (c *Config) Validate() error {
    var errors []string
    
    // 既存バリデーション...
    
    if c.NewFeatureEnabled && c.NewServiceURL == "" {
        errors = append(errors, "NEW_SERVICE_URL is required when NEW_FEATURE_ENABLED is true")
    }
    
    // エラーハンドリング...
}
```

### デプロイメント戦略

#### 1. 段階的デプロイ

```bash
# 1. ローカルテスト
make test-all

# 2. Cloud SQL接続テスト
make cloudsql-setup-complete
curl http://localhost:8080/api/v1/new-endpoint

# 3. ステージング環境（オプション）
GIN_MODE=release make local-cloudsql

# 4. 本番デプロイ
make cloudrun-deploy
```

#### 2. ロールバック手順

```bash
# 前のバージョンに戻す
gcloud run services update tru-s3-backend \
    --image gcr.io/zenn-ai-agent-hackathon-460205/tru-s3-backend:previous-tag \
    --region asia-northeast1

# または特定のリビジョンに戻す
gcloud run services update-traffic tru-s3-backend \
    --to-revisions=tru-s3-backend-00001-abc=100 \
    --region asia-northeast1
```

### 監視・ログ管理

#### 1. アプリケーション監視

```bash
# Cloud Run ログ監視
gcloud logs tail tru-s3-backend --format='default'

# エラーログのみ表示
gcloud logs read "resource.type=cloud_run_revision AND severity>=ERROR" --limit=50

# 特定の期間のログ
gcloud logs read "resource.type=cloud_run_revision" \
    --since="2025-01-15T10:00:00Z" \
    --until="2025-01-15T11:00:00Z"
```

#### 2. パフォーマンス監視

```bash
# レスポンス時間測定
time curl -s http://localhost:8080/api/v1/files > /dev/null

# 負荷テスト（簡易版）
for i in {1..100}; do
    curl -s http://localhost:8080/health > /dev/null &
done
wait

# メモリ使用量監視
docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"
```

#### 3. データベース監視

```bash
# Cloud SQL インスタンス状態
gcloud sql instances describe prd-db

# 接続数確認
gcloud sql operations list --instance=prd-db --limit=10

# パフォーマンス確認（psql接続後）
SELECT * FROM pg_stat_activity;
SELECT * FROM pg_stat_database;
```

### セキュリティ管理

#### 1. 認証・認可

```bash
# サービスアカウント権限確認
gcloud projects get-iam-policy zenn-ai-agent-hackathon-460205

# Cloud SQL IAM認証設定
gcloud sql users create service-account@project.iam.gserviceaccount.com \
    --instance=prd-db \
    --type=cloud_iam_service_account
```

#### 2. シークレット管理

```bash
# Secret Manager でパスワード管理
echo -n "new-secure-password" | gcloud secrets create new-secret --data-file=-

# シークレット更新
echo -n "updated-password" | gcloud secrets versions add new-secret --data-file=-

# シークレット確認
gcloud secrets versions access latest --secret="tru-s3-db-password"
```

## 📊 API仕様

### 基本API
| エンドポイント | メソッド | 説明 |
|-------------|--------|------|
| `/` | GET | メイン エンドポイント |
| `/health` | GET | ヘルスチェック |

### ファイル管理API
| エンドポイント | メソッド | 説明 |
|-------------|--------|------|
| `/api/v1/files` | GET | ファイル一覧取得 |
| `/api/v1/files` | POST | ファイルアップロード |
| `/api/v1/files/:id` | GET | ファイル情報取得 |
| `/api/v1/files/:id/download` | GET | ファイルダウンロード |
| `/api/v1/files/:id` | PUT | ファイル更新 |
| `/api/v1/files/:id` | DELETE | ファイル削除 |

### マッチングAPI
| エンドポイント | メソッド | 説明 |
|-------------|--------|------|
| `/api/v1/bookmarks` | GET | ブックマーク一覧 |
| `/api/v1/bookmarks` | POST | ブックマーク作成 |
| `/api/v1/bookmarks/:id` | DELETE | ブックマーク削除 |

### コンテストAPI
| エンドポイント | メソッド | 説明 |
|-------------|--------|------|
| `/api/v1/contests` | GET | コンテスト一覧 |
| `/api/v1/contests` | POST | コンテスト作成 |
| `/api/v1/contests/:id` | GET | コンテスト詳細 |
| `/api/v1/contests/:id` | PUT | コンテスト更新 |
| `/api/v1/contests/:id` | DELETE | コンテスト削除 |

### ハッカソンAPI
| エンドポイント | メソッド | 説明 |
|-------------|--------|------|
| `/api/v1/hackathons` | GET | ハッカソン一覧 |
| `/api/v1/hackathons` | POST | ハッカソン作成 |
| `/api/v1/hackathons/:id` | GET | ハッカソン詳細 |
| `/api/v1/hackathons/:id` | PUT | ハッカソン更新 |
| `/api/v1/hackathons/:id` | DELETE | ハッカソン削除 |
| `/api/v1/hackathons/:id/participants` | GET | 参加者一覧 |
| `/api/v1/hackathons/:id/participants` | POST | 参加者登録 |
| `/api/v1/hackathons/:id/participants/:participant_id` | DELETE | 参加者削除 |

## プロジェクト構造

```
TRu-S3/
├── ドキュメント
│   ├── README.md                    # メインドキュメント
│   └── docs/                       # 詳細ドキュメント
│       ├── README.md               # ドキュメント索引
│       ├── API_DOCUMENTATION.md   # 日本語API仕様書
│       ├── SECURITY.md            # セキュリティガイド
│       ├── DATABASE.md            # データベース設定
│       └── CLOUD_SQL.md           # Cloud SQL詳細ガイド
├── Docker設定
│   ├── Dockerfile                   # メインDockerfile（Cloud Run対応）
│   ├── docker-compose.yml          # ローカル開発用
│   ├── docker-compose.cloudsql.yml # Cloud SQL用
│   ├── docker-compose.prod.yml     # 本番用（セキュリティ強化）
│   └── .dockerignore               # Docker除外設定
├── Cloud設定
│   ├── cloudbuild.yaml             # Cloud Build設定
│   └── .github/                    # GitHub Actions設定
├── スクリプト
│   └── scripts/                    # セットアップ・管理スクリプト
│       ├── setup-cloud-sql-proxy.sh # Cloud SQL Proxyセットアップ
│       ├── setup-iam-auth.sh       # IAM認証セットアップ
│       ├── setup-network-security.sh # ネットワークセキュリティ
│       ├── check-cloudsql.sh       # インスタンス情報確認
│       ├── troubleshoot-cloudsql.sh # トラブルシューティング
│       └── deploy-cloudrun.sh      # Cloud Runデプロイ
├── アプリケーション
│   ├── main.go                     # エントリーポイント
│   ├── go.mod / go.sum             # Go依存関係
│   └── internal/                   # アプリケーションコード
│       ├── config/                 # 設定管理
│       ├── domain/                 # ドメイン層
│       ├── application/            # アプリケーション層
│       ├── infrastructure/         # インフラストラクチャ層
│       ├── interfaces/             # インターフェース層
│       └── database/               # データベース管理
├── データベース
│   └── migrations/                 # マイグレーションファイル
│       ├── 001_create_basic_tables_gcp.sql
│       ├── 002_add_missing_tables_gcp.sql
│       ├── 003_create_hackathons_table.sql
│       └── ...
├── テスト・ユーティリティ
│   └── test/                       # テスト・検証ファイル
│       ├── test_secure_connection.go
│       ├── verify_gcp_db.go
│       └── migrate_gcp.go
├── 設定ファイル
│   ├── .env.example                # 環境変数テンプレート
│   ├── .gitignore                  # Git除外設定
│   └── Makefile                    # ビルド・タスク管理
└── 自動生成/除外
    ├── .env*                       # 環境変数（Git除外）
    ├── cloud-sql-proxy            # バイナリ（自動DL）
    ├── ssl-certs/                  # SSL証明書（自動生成）
    └── service-account.json        # サービスアカウント（Git除外）
```

## 🔍 トラブルシューティング

### 自動診断ツール

```bash
# Cloud SQL関連の問題を自動診断
make cloudsql-troubleshoot

# または直接実行
./troubleshoot-cloudsql.sh
```

### よくある問題と解決策

1. **ポートが既に使用されている**
   ```bash
   # 使用中のプロセスを確認
   lsof -i :8080
   
   # プロセスを停止
   pkill -f "main.go"
   
   # または、.envファイルでポートを変更
   PORT=8081
   ```

2. **GCP認証エラー**
   ```bash
   # 認証状態を確認
   gcloud auth list
   
   # 認証を設定
   gcloud auth application-default login
   
   # プロジェクトを設定
   gcloud config set project your-project-id
   ```

3. **GCSバケットにアクセスできない**
   ```bash
   # バケットの存在確認
   gsutil ls gs://202506-zenn-ai-agent-hackathon/
   
   # 権限確認
   gsutil iam get gs://202506-zenn-ai-agent-hackathon/
   ```

4. **依存関係のエラー**
   ```bash
   # Go modulesをクリーンアップ
   go clean -modcache
   go mod tidy
   ```

5. **Dockerビルドエラー**
   ```bash
   # キャッシュをクリアしてリビルド
   docker-compose build --no-cache
   
   # Dockerシステムをクリーンアップ
   docker system prune -a
   ```

### ログ確認

```bash
# アプリケーションログ
docker-compose logs -f app

# 特定のエラーを検索
docker-compose logs app | grep ERROR

# リアルタイムログ監視
tail -f /var/log/app.log
```

### デバッグ情報

```bash
# 設定確認
curl -s http://localhost:8080/ | jq .

# 詳細なヘルス情報
curl -s http://localhost:8080/health | jq .

# GCS接続確認
gsutil ls gs://202506-zenn-ai-agent-hackathon/test/
```

## 開発ワークフロー

### 機能開発

1. **ブランチ作成**
   ```bash
   git checkout -b feature/new-feature
   ```

2. **ローカル開発**
   ```bash
   go run main.go
   # 別ターミナルでテスト
   curl http://localhost:8080/health
   ```

3. **テスト実行**
   ```bash
   go test ./...
   ```

4. **Docker確認**
   ```bash
   docker-compose up --build
   ```

### デプロイ準備

1. **本番設定確認**
   ```bash
   GIN_MODE=release go run main.go
   ```

2. **Docker本番ビルド**
   ```bash
   docker build -t tru-s3:latest .
   ```

## パフォーマンス監視

### 基本的な監視

```bash
# アプリケーション応答時間
time curl http://localhost:8080/health

# メモリ使用量
docker stats

# CPU使用量
top -p $(pgrep main)
```

### ファイル操作パフォーマンス

```bash
# 大量ファイルアップロードテスト
for i in {1..10}; do
  echo "Test file $i" > test$i.txt
  time curl -X POST -F "file=@test$i.txt" http://localhost:8080/api/v1/files
done

# 一覧取得パフォーマンス
time curl -s http://localhost:8080/api/v1/files
```

## 貢献

プロジェクトへの貢献を歓迎します！

### 貢献方法

1. Issueを確認または作成
2. フォークしてブランチを作成
3. 変更を実装
4. テストを追加・実行
5. プルリクエストを作成

### コーディング規約

- Go標準のフォーマットを使用 (`go fmt`)
- リンターを通す (`golangci-lint run`)
- テストを書く
- コミットメッセージは明確に

## 詳細ドキュメント

- **[API仕様書 (日本語)](./docs/API_DOCUMENTATION.md)** - 全APIの詳細仕様
- **[セキュリティガイド](./docs/SECURITY.md)** - Cloud SQL接続とセキュリティ設定
- **[データベース設定](./docs/DATABASE.md)** - マッチングアプリDB設定手順
- **[Cloud SQL設定](./docs/CLOUD_SQL.md)** - Cloud SQL詳細セットアップ

## セットアップスクリプト詳細

### 利用可能なスクリプト

- **`setup-cloud-sql-complete.sh`** - 完全自動セットアップ（推奨）
  - PostgreSQLクライアント自動インストール
  - Cloud SQL Auth Proxy自動ダウンロード・設定
  - パスワード設定（対話的）
  - 接続テスト・動作確認

- **`setup-cloudsql-proxy.sh`** - 基本セットアップ
  - Cloud SQL Auth Proxyダウンロード
  - `.env.local` ファイル生成

- **`check-cloudsql.sh`** - Cloud SQL情報確認
  - インスタンス詳細情報表示
  - 接続名生成
  - データベース・ユーザー一覧

- **`troubleshoot-cloudsql.sh`** - トラブルシューティング
  - システム状況自動診断
  - 問題の自動検出
  - 解決方法の提示

- **`deploy-cloudrun.sh`** - Cloud Runデプロイ
  - Cloud Build実行
  - 必要なAPI有効化
  - サービス情報表示

## 📋 技術仕様詳細

### データベーススキーマ

#### 主要テーブル構成

```sql
-- ユーザー管理
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    gmail VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- プロフィール
CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    age INTEGER,
    location VARCHAR(255),
    bio TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ファイル管理
CREATE TABLE file_metadata (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(500) NOT NULL,
    size BIGINT NOT NULL,
    content_type VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- マッチング
CREATE TABLE matchings (
    id SERIAL PRIMARY KEY,
    user1_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    user2_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### API レスポンス形式

#### 成功レスポンス
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "name": "example",
    "created_at": "2025-01-15T10:00:00Z"
  },
  "message": "Operation completed successfully"
}
```

#### エラーレスポンス
```json
{
  "status": "error",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "field": "email",
      "reason": "Invalid email format"
    }
  }
}
```

### 環境変数一覧

#### 必須環境変数
```bash
# サーバー設定
PORT=8080                           # サーバーポート
GIN_MODE=debug|release             # Ginモード

# データベース設定
DB_HOST=localhost                   # DBホスト
DB_PORT=5432                       # DBポート
DB_NAME=tru_s3                     # DB名
DB_USER=postgres                   # DBユーザー
DB_PASSWORD=password               # DBパスワード
DB_SSL_MODE=disable|require        # SSL設定

# GCP設定
GCS_BUCKET_NAME=bucket-name        # GCSバケット名
GCS_FOLDER=folder-name             # GCSフォルダ
GOOGLE_CLOUD_PROJECT=project-id    # GCPプロジェクトID
```

#### オプション環境変数
```bash
# Cloud SQL設定
USE_CLOUD_SQL_PROXY=true|false     # Cloud SQL Proxy使用
CLOUD_SQL_CONNECTION_NAME=conn     # Cloud SQL接続名

# データベース接続プール
DB_MAX_OPEN_CONNS=25              # 最大接続数
DB_MAX_IDLE_CONNS=5               # アイドル接続数

# 認証設定
GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json
```

### パフォーマンス仕様

#### レスポンス時間目標
- **ヘルスチェック**: < 50ms
- **ファイル一覧取得**: < 200ms
- **ファイルアップロード**: < 2s (10MB以下)
- **データベースクエリ**: < 100ms

#### スケーラビリティ
- **Cloud Run**: 最大10インスタンス
- **同時接続数**: インスタンスあたり80
- **データベース**: 最大25接続
- **ファイルサイズ制限**: 100MB

### セキュリティ仕様

#### 認証・認可
```bash
# GCP IAM ロール
- Cloud SQL Client
- Storage Object Admin
- Secret Manager Secret Accessor
- Cloud Run Invoker

# ネットワークセキュリティ
- HTTPS/TLS 1.2+
- Cloud SQL Auth Proxy
- VPC ネットワーク
- Private IP 接続
```

#### データ保護
```bash
# 暗号化
- 転送時: TLS 1.2+
- 保存時: Google Cloud 標準暗号化
- パスワード: Secret Manager

# アクセス制御
- IAM ベース認証
- リソースレベル権限
- 監査ログ
```

### 運用仕様

#### 監視項目
```bash
# アプリケーション
- レスポンス時間
- エラー率
- スループット
- メモリ使用量

# インフラストラクチャ
- Cloud Run インスタンス数
- Cloud SQL 接続数
- GCS ストレージ使用量
- ネットワーク帯域
```

#### ログ設定
```bash
# ログレベル
- ERROR: エラー情報
- WARN: 警告情報
- INFO: 一般情報（本番）
- DEBUG: 詳細情報（開発）

# ログ出力先
- 開発環境: 標準出力
- 本番環境: Cloud Logging
```

### 開発環境仕様

#### 必要なツール
```bash
# 開発ツール
Go 1.24.2+                        # Go言語
Docker 20.10+                     # コンテナ
Docker Compose 2.0+               # オーケストレーション
Git 2.30+                         # バージョン管理

# GCP ツール
gcloud CLI 400.0+                 # GCP CLI
Cloud SQL Auth Proxy 2.8+         # DB接続

# 推奨ツール
VS Code + Go Extension            # IDE
Postman                          # API テスト
pgAdmin                          # DB管理
```

#### 開発フロー
```bash
# 1. 機能開発
git checkout -b feature/new-feature
go run main.go                    # ローカル開発
go test ./...                     # テスト実行

# 2. 統合テスト
make cloudsql-setup-complete      # Cloud SQL接続
make test-all                     # 包括テスト

# 3. デプロイ
make cloudrun-deploy              # 本番デプロイ
```

## ライセンス

MIT License

## サポート

- **問題報告**: [GitHub Issues](https://github.com/TRu-S3/backend/issues)
- **質問・議論**: [GitHub Discussions](https://github.com/TRu-S3/backend/discussions)
- **セキュリティ問題**: security@tru-s3.com
- **ドキュメント**: [docs/](./docs/) フォルダ内の詳細ドキュメント

## 貢献者

このプロジェクトは以下の方針で開発されています：

- **Clean Architecture**: 保守性重視の設計
- **Test-Driven Development**: テスト駆動開発
- **Cloud-Native**: GCP完全対応
- **Security-First**: セキュリティ優先
- **Documentation-First**: ドキュメント重視

貢献をお待ちしています！
