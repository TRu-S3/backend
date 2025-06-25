# 🏆 TRu-S3 Backend

**マッチングアプリ + ファイル管理システム**

## 🎯 概要

TRu-S3は、Google Cloud Platform上で動作するマッチングアプリケーションとファイル管理システムを統合したバックエンドAPIです。オニオンアーキテクチャを採用し、高い保守性とテスタビリティを実現しています。

### 🌟 主な機能

- **👥 マッチングシステム**: ユーザープロフィール、マッチング、ブックマーク機能
- **🏆 ハッカソン管理**: ハッカソンイベントと参加者管理
- **📁 ファイル管理**: GCS連携によるファイルアップロード・管理
- **🏁 コンテスト機能**: プログラミングコンテスト管理

### ✨ 技術的特徴

- **🏗️ オニオンアーキテクチャ**: 保守性・テスタビリティを重視した設計
- **☁️ GCP完全対応**: Cloud Storage + Cloud SQL + Cloud Run
- **🔐 セキュア**: 複数レベルのセキュリティ設定（IAM認証、SSL/TLS、プライベートIP）
- **🐳 Docker対応**: ローカル開発から本番デプロイまで統一環境
- **🚀 自動化**: ワンコマンドセットアップ・デプロイ
- **📚 日本語対応**: 包括的な日本語APIドキュメント

### 🛠️ 技術スタック

- **言語**: Go 1.24.2
- **フレームワーク**: Gin (HTTP)、GORM (ORM)
- **データベース**: PostgreSQL 17 (Cloud SQL)
- **ストレージ**: Google Cloud Storage
- **インフラ**: Docker、Cloud Run、Cloud Build
- **認証**: GCP IAM、Cloud SQL Auth Proxy

## 🏗️ アーキテクチャ

```
internal/
├── domain/          # ドメイン層 - エンティティとリポジトリインターフェース
├── application/     # アプリケーション層 - ビジネスロジック
├── infrastructure/ # インフラストラクチャ層 - GCS実装
├── interfaces/     # インターフェース層 - HTTPハンドラー
└── config/         # 設定管理
```

## � 要件

### 開発環境
- Go 1.24.2+
- GCP アカウント
- GCP認証設定

### Docker環境
- Docker
- Docker Compose

## 🚀 セットアップ

### 🛠️ 利用可能なコマンド

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

## 📚 API使用方法

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

## 💡 使用例

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

## 🐳 Docker操作

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

## 🛠️ ローカル開発

Dockerを使用せずにローカルで実行する場合：

```bash
# 依存関係をダウンロード
go mod tidy

# アプリケーションを起動
go run main.go

# または、特定のポートで起動
PORT=8081 go run main.go
```

## ☁️ Cloud SQL Auth Proxy

本番環境では既存のCloud SQL PostgreSQLインスタンス（プロジェクト: `zenn-ai-agent-hackathon-460205`, インスタンス: `prd-db`）を使用してCloud SQL Auth Proxyで接続します。

> 📖 **詳細なセットアップ手順は [CLOUD_SQL_SETUP.md](./CLOUD_SQL_SETUP.md) を参照してください**

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

## 📡 API エンドポイント一覧

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

## 📁 プロジェクト構造

```
TRu-S3/
├── 📖 ドキュメント
│   ├── README.md                    # メインドキュメント
│   └── docs/                       # 詳細ドキュメント
│       ├── README.md               # ドキュメント索引
│       ├── API_DOCUMENTATION.md   # 日本語API仕様書
│       ├── SECURITY.md            # セキュリティガイド
│       ├── DATABASE.md            # データベース設定
│       └── CLOUD_SQL.md           # Cloud SQL詳細ガイド
├── 🐳 Docker設定
│   ├── Dockerfile                   # メインDockerfile（Cloud Run対応）
│   ├── docker-compose.yml          # ローカル開発用
│   ├── docker-compose.cloudsql.yml # Cloud SQL用
│   ├── docker-compose.prod.yml     # 本番用（セキュリティ強化）
│   └── .dockerignore               # Docker除外設定
├── ☁️ Cloud設定
│   ├── cloudbuild.yaml             # Cloud Build設定
│   └── .github/                    # GitHub Actions設定
├── 🔧 スクリプト
│   └── scripts/                    # セットアップ・管理スクリプト
│       ├── setup-cloud-sql-proxy.sh # Cloud SQL Proxyセットアップ
│       ├── setup-iam-auth.sh       # IAM認証セットアップ
│       ├── setup-network-security.sh # ネットワークセキュリティ
│       ├── check-cloudsql.sh       # インスタンス情報確認
│       ├── troubleshoot-cloudsql.sh # トラブルシューティング
│       └── deploy-cloudrun.sh      # Cloud Runデプロイ
├── 🚀 アプリケーション
│   ├── main.go                     # エントリーポイント
│   ├── go.mod / go.sum             # Go依存関係
│   └── internal/                   # アプリケーションコード
│       ├── config/                 # 設定管理
│       ├── domain/                 # ドメイン層
│       ├── application/            # アプリケーション層
│       ├── infrastructure/         # インフラストラクチャ層
│       ├── interfaces/             # インターフェース層
│       └── database/               # データベース管理
├── 🗃️ データベース
│   └── migrations/                 # マイグレーションファイル
│       ├── 001_create_basic_tables_gcp.sql
│       ├── 002_add_missing_tables_gcp.sql
│       ├── 003_create_hackathons_table.sql
│       └── ...
├── 🧪 テスト・ユーティリティ
│   └── test/                       # テスト・検証ファイル
│       ├── test_secure_connection.go
│       ├── verify_gcp_db.go
│       └── migrate_gcp.go
├── ⚙️ 設定ファイル
│   ├── .env.example                # 環境変数テンプレート
│   ├── .gitignore                  # Git除外設定
│   └── Makefile                    # ビルド・タスク管理
└── 🗂️ 自動生成/除外
    ├── .env*                       # 環境変数（Git除外）
    ├── cloud-sql-proxy            # バイナリ（自動DL）
    ├── ssl-certs/                  # SSL証明書（自動生成）
    └── service-account.json        # サービスアカウント（Git除外）
```

## 🐛 トラブルシューティング

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

## 🔄 開発ワークフロー

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

## 📊 パフォーマンス監視

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

## 🤝 貢献

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

## 📖 詳細ドキュメント

- **[📋 API仕様書 (日本語)](./docs/API_DOCUMENTATION.md)** - 全APIの詳細仕様
- **[🔐 セキュリティガイド](./docs/SECURITY.md)** - Cloud SQL接続とセキュリティ設定
- **[💾 データベース設定](./docs/DATABASE.md)** - マッチングアプリDB設定手順
- **[☁️ Cloud SQL設定](./docs/CLOUD_SQL.md)** - Cloud SQL詳細セットアップ

## 🔧 セットアップスクリプト詳細

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

## 📝 ライセンス

MIT License

## 📞 サポート

- 問題報告: GitHub Issues
- 質問: Discussions
- セキュリティ問題: 直接連絡
