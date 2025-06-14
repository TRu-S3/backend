# TRu-S3 Backend

## 🎯 概要

- **バケット**: `202506-zenn-ai-agent-hackathon`
- **フォルダ**: `test`
- **アーキテクチャ**: オニオンアーキテクチャ
- **API**: RESTful API
- **認証**: GCP認証

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

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd TRu-S3
```

### 2. 環境変数の設定

```bash
# .env.example から .env をコピー
cp .env.example .env

# .env ファイルを編集
vim .env
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

## 📡 API エンドポイント一覧

| エンドポイント | メソッド | 説明 |
|-------------|--------|------|
| `/` | GET | メイン エンドポイント |
| `/health` | GET | ヘルスチェック |
| `/api/v1/files` | GET | ファイル一覧取得 |
| `/api/v1/files` | POST | ファイルアップロード |
| `/api/v1/files/:id` | GET | ファイル情報取得 |
| `/api/v1/files/:id/download` | GET | ファイルダウンロード |
| `/api/v1/files/:id` | PUT | ファイル更新 |
| `/api/v1/files/:id` | DELETE | ファイル削除 |

## 📁 プロジェクト構造

```
.
├── main.go                    # エントリーポイント
├── .env                       # 環境変数設定 (作成が必要)
├── .env.example               # 環境変数テンプレート
├── internal/
│   ├── config/
│   │   └── config.go         # 設定管理
│   ├── domain/
│   │   ├── file.go           # ファイルエンティティ
│   │   └── repository.go     # リポジトリインターフェース
│   ├── application/
│   │   └── file_service.go   # ファイルサービス
│   ├── infrastructure/
│   │   ├── gcs_client.go     # GCSクライアント
│   │   └── gcs_file_repository.go # GCS実装
│   └── interfaces/
│       ├── file_handler.go   # HTTPハンドラー
│       └── routes.go         # ルート設定
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── API_DOCS.md               # 詳細API仕様
└── README.md
```

## 🐛 トラブルシューティング

### よくある問題

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

## 📖 詳細仕様

詳細なAPI仕様については `API_DOCS.md` を参照してください。

## 📝 ライセンス

MIT License

## 📞 サポート

- 問題報告: GitHub Issues
- 質問: Discussions
- セキュリティ問題: 直接連絡
