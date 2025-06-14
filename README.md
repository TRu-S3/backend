# TRu-S3 Backend

TRu-S3のバックエンドアプリケーションです。Go言語とGinフレームワークを使用して構築されています。

## 📋 要件

- Docker
- Docker Compose
- Make（オプション、便利なコマンドを使用する場合）

## 🚀 クイックスタート

### Dockerを使用して起動

1. **リポジトリをクローン**
   ```bash
   git clone <repository-url>
   cd TRu-S3
   ```

2. **アプリケーションをビルドして起動**
   ```bash
   # Makefileを使用する場合
   make build
   make run
   
   # または、docker-compose を直接使用
   docker-compose up --build -d
   ```

3. **アプリケーションが起動していることを確認**
   ```bash
   # Makefileを使用する場合
   make test
   
   # または、直接curlを使用
   curl http://localhost:8080
   curl http://localhost:8080/health
   ```

## 🐳 Docker操作

### 基本コマンド

```bash
# イメージをビルド
make build

# バックグラウンドで起動
make run

# 開発モード（ログを表示しながら起動）
make dev

# ログを確認
make logs

# アプリケーションを停止
make stop

# コンテナとイメージをクリーンアップ
make clean

# 再起動（リビルド込み）
make restart

# 実行中のコンテナにシェルでアクセス
make shell
```

### Docker Compose を直接使用する場合

```bash
# ビルドして起動
docker-compose up --build

# バックグラウンドで起動
docker-compose up -d

# ログを確認
docker-compose logs -f

# 停止
docker-compose down

# 完全クリーンアップ
docker-compose down --rmi all --volumes --remove-orphans
```

## 🛠️ ローカル開発

Dockerを使用せずにローカルで実行する場合：

```bash
# 依存関係をダウンロード
go mod download

# アプリケーションを起動
go run main.go

# または、Makefileを使用
make local-run
```

## 📡 API エンドポイント

| エンドポイント | メソッド | 説明 |
|-------------|--------|------|
| `/` | GET | メイン エンドポイント - アプリケーションの状態を返す |
| `/health` | GET | ヘルスチェック エンドポイント |

### レスポンス例

**GET /**
```json
{
  "message": "TRu-S3 Backend is running!",
  "status": "healthy"
}
```

**GET /health**
```json
{
  "status": "ok"
}
```

## 🔧 設定

### 環境変数

現在、特別な環境変数は必要ありませんが、将来的に追加される可能性があります。

### ポート

- デフォルトポート: `8080`
- Docker内部ポート: `8080`
- ホストマシンポート: `8080`

## 📁 プロジェクト構造

```
.
├── main.go              # メインアプリケーションファイル
├── go.mod               # Go モジュール定義
├── go.sum               # 依存関係のチェックサム
├── Dockerfile           # Dockerイメージの定義
├── docker-compose.yml   # Docker Compose設定
├── .dockerignore        # Dockerビルド時の除外ファイル
├── Makefile            # 便利なコマンド集
└── README.md           # このファイル
```

## 🐛 トラブルシューティング

### よくある問題

1. **ポート8080が既に使用されている**
   ```bash
   # 使用中のプロセスを確認
   lsof -i :8080
   
   # または、別のポートを使用
   docker-compose up --build -p 8081:8080
   ```

2. **Dockerイメージのビルドに失敗**
   ```bash
   # キャッシュをクリアしてリビルド
   docker-compose build --no-cache
   ```

3. **アプリケーションにアクセスできない**
   ```bash
   # コンテナの状態を確認
   docker-compose ps
   
   # ログを確認
   docker-compose logs
   ```

## 🤝 貢献

1. このリポジトリをフォーク
2. 機能ブランチを作成 (`git checkout -b feature/AmazingFeature`)
3. 変更をコミット (`git commit -m 'Add some AmazingFeature'`)
4. ブランチにプッシュ (`git push origin feature/AmazingFeature`)
5. プルリクエストを作成

## 📝 ライセンス

このプロジェクトのライセンス情報については、LICENSE ファイルを参照してください。

## 📞 サポート

問題や質問がある場合は、Issueを作成してください。
