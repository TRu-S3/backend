# マッチングアプリ データベースセットアップ手順書

## 📋 概要

TRu-S3プロジェクトにマッチングアプリのデータベース機能を追加するためのセットアップ手順書です。

### 追加されるテーブル
- **Users** - ユーザーマスタ
- **Tags** - タグマスタ 
- **Profiles** - ユーザープロフィール
- **Matchings** - マッチング情報
- **Bookmarks** - ブックマーク

## 🚀 セットアップ手順

### 1. 前提条件の確認

既存のTRu-S3プロジェクトが正常に動作していることを確認してください。

```bash
# プロジェクトディレクトリに移動
cd /path/to/TRu-S3

# 既存のアプリケーションが動作することを確認
make test
```

### 2. データベース環境の選択

以下のいずれかの環境でセットアップを行います：

#### オプション A: Docker環境（推奨）
```bash
# Docker環境でPostgreSQLを起動
make build
make run

# コンテナが起動していることを確認
docker-compose ps
```

#### オプション B: Cloud SQL環境
```bash
# Cloud SQL環境のセットアップ
make cloudsql-setup-complete

# 接続確認
curl http://localhost:8080/health
```

#### オプション C: ローカルPostgreSQL環境
```bash
# ローカルのPostgreSQLが起動していることを確認
pg_isready -h localhost -p 5432

# アプリケーションを起動
make local-run
```

### 3. マイグレーション実行

#### 方法1: PostgreSQL CLIを使用（推奨）

```bash
# Docker環境の場合
docker-compose exec db psql -U postgres -d tru_s3 -f /docker-entrypoint-initdb.d/001_create_matching_app_tables.sql

# ローカル環境の場合
psql -h localhost -p 5432 -U postgres -d tru_s3 -f migrations/001_create_matching_app_tables.sql

# Cloud SQL環境の場合
psql "host=localhost port=5433 dbname=tru_s3 user=postgres sslmode=disable" -f migrations/001_create_matching_app_tables.sql
```

#### 方法2: マイグレーションファイルをコピーして実行

Docker環境の場合：
```bash
# SQLファイルをコンテナにコピー
docker cp migrations/001_create_matching_app_tables.sql $(docker-compose ps -q db):/tmp/

# コンテナ内で実行
docker-compose exec db psql -U postgres -d tru_s3 -f /tmp/001_create_matching_app_tables.sql
```

### 4. テーブル作成の確認

```bash
# テーブルが正しく作成されたか確認
docker-compose exec db psql -U postgres -d tru_s3 -c "\dt"

# または、SQL直接実行で確認
docker-compose exec db psql -U postgres -d tru_s3 -c "
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name IN ('users', 'tags', 'profiles', 'matchings', 'bookmarks')
ORDER BY table_name;
"
```

期待される出力：
```
 table_name 
-----------
 bookmarks
 matchings
 profiles
 tags
 users
(5 rows)
```

### 5. 初期データ投入

```bash
# 初期データを投入
docker-compose exec db psql -U postgres -d tru_s3 -f /docker-entrypoint-initdb.d/seed_data.sql

# または、ローカルファイルから
psql -h localhost -p 5432 -U postgres -d tru_s3 -f seed_data.sql
```

### 6. データ投入の確認

```bash
# 各テーブルのレコード数を確認
docker-compose exec db psql -U postgres -d tru_s3 -c "
SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'tags' as table_name, COUNT(*) as count FROM tags
UNION ALL
SELECT 'profiles' as table_name, COUNT(*) as count FROM profiles
UNION ALL
SELECT 'bookmarks' as table_name, COUNT(*) as count FROM bookmarks
UNION ALL
SELECT 'matchings' as table_name, COUNT(*) as count FROM matchings;
"
```

期待される出力：
```
 table_name | count 
-----------+-------
 users     |    10
 tags      |    10
 profiles  |    10
 bookmarks |    13
 matchings |     5
(5 rows)
```

## 🔧 Go アプリケーションとの統合

### 1. モデルの読み込み

既存の `internal/database/models.go` ファイルに以下を追加するか、新しいファイルとして作成済みの `internal/database/matching_models.go` を使用してください。

### 2. データベース接続の更新

`internal/database/connection.go` を更新して、新しいモデルのマイグレーションを含めます：

```go
// 既存のコードに追加
func AutoMigrateMatchingApp(db *gorm.DB) error {
    return db.AutoMigrate(MatchingAppModels...)
}
```

### 3. アプリケーション起動時のマイグレーション

`main.go` または適切な初期化ファイルで：

```go
// データベース接続後
err = database.AutoMigrateMatchingApp(db)
if err != nil {
    log.Fatal("Failed to migrate matching app tables:", err)
}
```

## 📊 動作確認

### 1. 基本的なクエリテスト

```bash
# ユーザー一覧の確認
docker-compose exec db psql -U postgres -d tru_s3 -c "
SELECT u.name, u.gmail, p.bio, t.name as tag_name 
FROM users u 
LEFT JOIN profiles p ON u.id = p.user_id 
LEFT JOIN tags t ON p.tag_id = t.id 
LIMIT 5;
"

# ブックマーク関係の確認
docker-compose exec db psql -U postgres -d tru_s3 -c "
SELECT u1.name as user_name, u2.name as bookmarked_user 
FROM bookmarks b 
JOIN users u1 ON b.user_id = u1.id 
JOIN users u2 ON b.bookmarked_user_id = u2.id 
LIMIT 5;
"
```

### 2. アプリケーション動作確認

```bash
# アプリケーションが正常に起動することを確認
curl http://localhost:8080/health

# 既存のファイルAPIが動作することを確認
curl http://localhost:8080/api/v1/files
```

## 🛠️ トラブルシューティング

### よくある問題と解決策

#### 1. テーブル作成時のエラー

**エラー**: `relation "users" already exists`
```bash
# 既存のテーブルを確認
docker-compose exec db psql -U postgres -d tru_s3 -c "\dt"

# 必要に応じてテーブルを削除（注意：既存データも削除されます）
docker-compose exec db psql -U postgres -d tru_s3 -c "
DROP TABLE IF EXISTS bookmarks CASCADE;
DROP TABLE IF EXISTS matchings CASCADE; 
DROP TABLE IF EXISTS profiles CASCADE;
DROP TABLE IF EXISTS tags CASCADE;
DROP TABLE IF EXISTS users CASCADE;
"
```

#### 2. データベース接続エラー

```bash
# PostgreSQLの状態確認
docker-compose logs db

# 接続テスト
docker-compose exec app psql -h db -U postgres -d tru_s3 -c "SELECT version();"
```

#### 3. 初期データ投入エラー

**エラー**: `duplicate key value violates unique constraint`
```bash
# 既存データをクリア（注意：全データ削除）
docker-compose exec db psql -U postgres -d tru_s3 -c "
TRUNCATE TABLE bookmarks CASCADE;
TRUNCATE TABLE matchings CASCADE;
TRUNCATE TABLE profiles CASCADE;
TRUNCATE TABLE users CASCADE;
TRUNCATE TABLE tags CASCADE;
"

# 初期データを再投入
docker-compose exec db psql -U postgres -d tru_s3 -f /docker-entrypoint-initdb.d/seed_data.sql
```

### 4. パフォーマンス確認

```bash
# インデックスの存在確認
docker-compose exec db psql -U postgres -d tru_s3 -c "
SELECT indexname, tablename, indexdef 
FROM pg_indexes 
WHERE tablename IN ('users', 'tags', 'profiles', 'matchings', 'bookmarks')
ORDER BY tablename, indexname;
"
```

## 📚 次のステップ

### API エンドポイントの実装

マッチングアプリ用のAPIエンドポイントを実装する場合：

1. **User Management API** - ユーザー登録・更新・取得
2. **Profile API** - プロフィール管理
3. **Matching API** - マッチング機能
4. **Bookmark API** - ブックマーク機能
5. **Tag API** - タグ管理

### サンプル実装の場所

- `internal/interfaces/` - HTTPハンドラー
- `internal/application/` - ビジネスロジック  
- `internal/infrastructure/` - データアクセス層

## 📞 サポート

問題が発生した場合：

1. **ログ確認**: `docker-compose logs -f`
2. **データベース状態確認**: 上記のトラブルシューティング手順を実行
3. **既存機能への影響確認**: `make test-all`

セットアップが完了したら、既存のTRu-S3機能とマッチングアプリ機能の両方が正常に動作することを確認してください。