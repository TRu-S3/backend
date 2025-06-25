# Cloud SQL Auth Proxy セットアップガイド

このドキュメントでは、TRu-S3 BackendでCloud SQL Auth Proxyを使用してCloud SQLに接続する詳細手順を説明します。

> 💡 **クイックスタート**: `make cloudsql-setup-complete` で自動セットアップが可能です。

## 📋 前提条件

- Google Cloud Platform アカウント
- 既存のCloud SQLインスタンス（`prd-db`）
- gcloud CLI のインストールと認証
- Go 1.24.2+
- PostgreSQLクライアント（自動インストール対応）

## 🚀 推奨セットアップ方法

### 完全自動セットアップ（推奨）

```bash
# ワンコマンドで全ての手順を実行
make cloudsql-setup-complete

# または直接実行
./setup-cloud-sql-complete.sh
```

**このスクリプトの実行内容:**
1. Cloud SQLインスタンス情報確認
2. PostgreSQLクライアント自動インストール
3. Cloud SQL Auth Proxyダウンロード・設定
4. パスワード設定（対話的）
5. SSL設定調整
6. 接続テスト
7. アプリケーション起動・動作確認

### トラブルシューティング

```bash
# 問題が発生した場合の自動診断
make cloudsql-troubleshoot

# または直接実行
./troubleshoot-cloudsql.sh
```

## 🔧 手動セットアップ（詳細制御が必要な場合）

### 1. 環境の確認

```bash
# プロジェクトとgcloud認証の確認
gcloud auth list --filter=status:ACTIVE --format="value(account)"
gcloud config get-value project

# Cloud SQLインスタンス情報の確認
./check-cloudsql.sh
```

**期待される出力:**
```
✅ Cloud SQL インスタンス一覧:
NAME    REGION           DATABASE_VERSION  TIER         STATUS
prd-db  asia-northeast1  POSTGRES_17       db-f1-micro  RUNNABLE

🔗 Cloud SQL Auth Proxy 接続名:
zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db
```

### 2. 必要なAPIの有効化

```bash
# Cloud SQL Admin APIを有効化
gcloud services enable sqladmin.googleapis.com --project=zenn-ai-agent-hackathon-460205
```

### 3. PostgreSQLクライアントのインストール

```bash
# PostgreSQL公式リポジトリを追加
curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg
echo "deb https://apt.postgresql.org/pub/repos/apt jammy-pgdg main" | sudo tee /etc/apt/sources.list.d/pgdg.list

# パッケージリストを更新してクライアントをインストール
sudo apt update
sudo apt install -y postgresql-client
```

### 4. Cloud SQL Auth Proxyのセットアップ

```bash
# Cloud SQL Auth Proxyをダウンロードして環境を設定
./setup-cloudsql-proxy.sh
```

**実行内容:**
- Cloud SQL Auth Proxyバイナリのダウンロード
- `.env.local` ファイルの生成
- 必要な環境変数の設定

### 5. Cloud SQLパスワードの設定

```bash
# postgresユーザーのパスワードを設定
gcloud sql users set-password postgres \
  --instance=prd-db \
  --password="your-secure-password" \
  --project=zenn-ai-agent-hackathon-460205

# .env.localファイルにパスワードを設定
sed -i "s/^DB_PASSWORD=.*/DB_PASSWORD=\"your-secure-password\"/" .env.local
```

**注意:** セキュリティ上、強力なパスワードを使用してください。

### 6. SSL設定の調整

Cloud SQL Auth Proxyはローカルで暗号化を処理するため、SSL設定を無効にします：

```bash
# .env.localファイルのSSL設定を更新
sed -i "s/^DB_SSL_MODE=.*/DB_SSL_MODE=disable/" .env.local
```

### 7. Cloud SQL Auth Proxyの起動

```bash
# 環境変数を読み込み
source .env.local

# Cloud SQL Auth Proxyをバックグラウンドで起動
./cloud-sql-proxy $CLOUD_SQL_CONNECTION_NAME --port=5433 &
```

**確認:**
```bash
# プロセス確認
ps aux | grep cloud-sql-proxy | grep -v grep

# ポート確認
netstat -ln | grep :5433
```

### 8. psqlでの接続テスト

```bash
# 環境変数を読み込み
source .env.local

# PostgreSQLバージョン確認
PGPASSWORD="$DB_PASSWORD" psql "host=localhost port=5433 sslmode=disable user=postgres dbname=postgres" -c "SELECT version();"

# データベース一覧確認
PGPASSWORD="$DB_PASSWORD" psql "host=localhost port=5433 sslmode=disable user=postgres dbname=postgres" -c "\l"

# 基本情報確認
PGPASSWORD="$DB_PASSWORD" psql "host=localhost port=5433 sslmode=disable user=postgres dbname=postgres" -c "SELECT current_database(), current_user, now();"
```

**期待される出力:**
```
                                         version                                         
-----------------------------------------------------------------------------------------
 PostgreSQL 17.5 on x86_64-pc-linux-gnu, compiled by Debian clang version 12.0.1, 64-bit
(1 row)
```

### 9. アプリケーションの起動とテスト

```bash
# 環境変数を読み込んでアプリケーション起動
source .env.local && go run main.go &

# ヘルスチェック
curl -v http://localhost:8080/health

# API動作確認
curl -s http://localhost:8080/api/v1/files | jq .
```

**期待される出力:**
```json
{"status":"ok"}
```

## 🔧 便利なコマンド

### Makeタスクの使用（推奨）

```bash
# Cloud SQL関連コマンド
make cloudsql-info              # インスタンス情報確認
make cloudsql-setup-complete    # 完全自動セットアップ
make cloudsql-setup             # 基本セットアップのみ
make cloudsql-troubleshoot      # トラブルシューティング
make cloudsql-local             # Cloud SQL Auth Proxyを起動
make cloudsql-compose           # Docker Compose + Cloud SQL

# アプリケーション起動
make local-cloudsql             # Cloud SQL使用でローカル起動
make cloudrun-deploy            # Cloud Runにデプロイ

# テスト
make test                       # 基本APIテスト
make test-all                   # 包括的テスト
```

### 直接スクリプト実行

```bash
# セットアップスクリプト
./setup-cloud-sql-complete.sh   # 完全自動セットアップ
./setup-cloudsql-proxy.sh       # 基本セットアップ
./check-cloudsql.sh             # インスタンス情報確認
./troubleshoot-cloudsql.sh      # 問題診断

# Cloud Runデプロイ
./deploy-cloudrun.sh            # Cloud Runデプロイ
```

### 手動でのデータベース操作

```bash
# tru_s3データベースに接続
source .env.local
PGPASSWORD="$DB_PASSWORD" psql "host=localhost port=5433 sslmode=disable user=postgres dbname=tru_s3"

# SQLクエリ例
# CREATE TABLE test_table (id SERIAL PRIMARY KEY, name VARCHAR(100));
# INSERT INTO test_table (name) VALUES ('Test Data');
# SELECT * FROM test_table;
```

## 🐛 トラブルシューティング

### よくある問題と解決策

#### 1. SSL接続エラー
```
psql: error: server does not support SSL, but SSL was required
```
**解決策:** SSL設定を `disable` に変更
```bash
sed -i "s/^DB_SSL_MODE=.*/DB_SSL_MODE=disable/" .env.local
```

#### 2. Cloud SQL Auth Proxyが起動しない
**確認項目:**
- gcloud認証の確認: `gcloud auth list`
- Cloud SQL Admin APIの有効化
- 正しい接続名の使用

#### 3. アプリケーションがデータベースに接続できない
**確認項目:**
```bash
# Cloud SQL Auth Proxyの状況確認
ps aux | grep cloud-sql-proxy

# ポート使用状況確認
netstat -ln | grep :5433

# 環境変数確認
source .env.local && echo $CLOUD_SQL_CONNECTION_NAME
```

#### 4. PostgreSQLクライアントがない
```bash
# PostgreSQL公式リポジトリから最新版をインストール
curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg
echo "deb https://apt.postgresql.org/pub/repos/apt jammy-pgdg main" | sudo tee /etc/apt/sources.list.d/pgdg.list
sudo apt update && sudo apt install -y postgresql-client
```

## 📁 生成・管理されるファイル

### 自動生成されるファイル
- `.env.local` - ローカル開発用環境変数（`setup-cloudsql-proxy.sh`で生成）
- `cloud-sql-proxy` - Cloud SQL Auth Proxyバイナリ（自動ダウンロード）

### Git管理対象外（.gitignoreで除外）
- `.env` - 本番環境変数
- `.env.local*` - ローカル環境ファイル
- `cloud-sql-proxy*` - バイナリファイル
- `test-upload.txt` - テスト用ファイル
- `*.tmp`, `*.temp` - 一時ファイル

### 削除されたファイル（統合・不要化）
- `backend` - 古いビルド済みバイナリ
- `reset-cloudsql-password.sh` - パスワードリセット機能（完全セットアップに統合）
- `Dockerfile.cloudrun` - Cloud Run用Dockerfile（メインDockerfileに統合）

## 🔒 セキュリティのベストプラクティス

1. **強力なパスワード**: 最低8文字、英数字記号を組み合わせ
2. **環境変数の管理**: `.env.local` ファイルはGitにコミットしない
3. **定期的なパスワード更新**: セキュリティ向上のため定期更新
4. **本番環境**: Secret Managerの使用を推奨

## 🚀 本番環境への移行

ローカルでの動作確認後、Cloud Runへのデプロイ：

```bash
# Secret Managerにパスワードを保存
echo -n "your-cloud-sql-password" | gcloud secrets create tru-s3-db-password --data-file=-

# Cloud Runへデプロイ
make cloudrun-deploy
```

## 📝 注意事項

- Cloud SQL Auth Proxyはローカル接続では暗号化を処理するため、SSL設定は`disable`
- 本番環境（Cloud Run）では自動でSSL接続が有効になります
- パスワードやシークレットは適切に管理してください
- 定期的にCloud SQLインスタンスの監視とメンテナンスを実施してください
