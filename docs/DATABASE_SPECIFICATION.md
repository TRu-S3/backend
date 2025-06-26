# データベース仕様書

## 概要

TRu-S3アプリケーションのデータベース設計仕様書です。PostgreSQLを使用し、GORMによるORM機能を提供しています。

## 最新の実装状況

**データベースインスタンス**: Google Cloud SQL PostgreSQL 17  
**インスタンス名**: `prd-db` (asia-northeast1)  
**接続方式**: Cloud SQL Auth Proxy使用  

**現在のテーブル数**: **10個**  
**実装済みCRUD**: **全テーブルで完全対応**

## データベース基本情報

- **データベース管理システム**: PostgreSQL
- **ORM**: GORM v1.30.0
- **ドライバー**: gorm.io/driver/postgres v1.6.0
- **文字エンコーディング**: UTF-8
- **タイムゾーン**: UTC
- **自動マイグレーション**: 有効

## 現在のテーブル一覧とデータ状況

| テーブル名 | 行数 | 説明 | CRUD実装 |
|------------|------|------|----------|
| **users** | 0行 | ユーザー基本情報 | 完全実装 |
| **tags** | 0行 | プロフィールタグ | 完全実装 |
| **profiles** | 0行 | ユーザープロフィール詳細 | 完全実装 |
| **matchings** | 0行 | ユーザー間マッチング | 完全実装 |
| **bookmarks** | 0行 | ユーザーブックマーク | 完全実装 |
| **contests** | 0行 | プログラミングコンテスト | 完全実装 |
| **hackathons** | **2行** | ハッカソンイベント | 完全実装 |
| **hackathon_participants** | 0行 | ハッカソン参加者 | 完全実装 |
| **file_metadata** | 0行 | ファイルメタデータ | 完全実装 |
| **verification_results** | **3行** | 検証結果（システム用） | - |

## テーブル構造詳細

### 1. ユーザー関連テーブル

#### 1.1 users テーブル
ユーザーの基本情報を格納

| カラム名 | データ型 | 制約 | 説明 |
|---------|---------|------|------|
| id | SERIAL | PRIMARY KEY | ユーザーID（自動採番） |
| name | VARCHAR(100) | NOT NULL | ユーザー名 |
| gmail | VARCHAR(100) | UNIQUE, NOT NULL | Gmail アドレス |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

**インデックス**:
- `idx_users_gmail` (gmail)
- `idx_users_created_at` (created_at)

**関連テーブル**:
- profiles (1:1)
- bookmarks (1:N)
- contests (1:N, author関係)
- hackathon_participants (1:N)

#### 1.2 tags テーブル
プロフィールのタグ情報を格納

| カラム名 | データ型 | 制約 | 説明 |
|---------|---------|------|------|
| id | SERIAL | PRIMARY KEY | タグID（自動採番） |
| name | VARCHAR(50) | UNIQUE, NOT NULL | タグ名 |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

**インデックス**:
- `idx_tags_name` (name)

**更新トリガー**:
- `update_tags_updated_at`: 更新時にupdated_atを自動更新

#### 1.3 profiles テーブル
ユーザーのプロフィール詳細情報を格納

| カラム名 | データ型 | 制約 | 説明 |
|---------|---------|------|------|
| id | SERIAL | PRIMARY KEY | プロフィールID（自動採番） |
| user_id | INTEGER | NOT NULL, UNIQUE, FK | ユーザーID |
| tag_id | INTEGER | FK | タグID |
| bio | TEXT | | 自己紹介 |
| age | INTEGER | | 年齢 |
| location | VARCHAR(100) | | 居住地 |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

**インデックス**:
- `idx_profiles_user_id` (user_id)
- `idx_profiles_tag_id` (tag_id)

**外部キー制約**:
- `profiles_user_id_fkey`: user_id → users(id) ON DELETE CASCADE
- `profiles_tag_id_fkey`: tag_id → tags(id) ON DELETE SET NULL

**更新トリガー**:
- `update_profiles_updated_at`: 更新時にupdated_atを自動更新

#### 1.4 matchings テーブル
ユーザー間のマッチング情報を格納

| カラム名 | データ型 | 制約 | 説明 |
|---------|---------|------|------|
| id | SERIAL | PRIMARY KEY | マッチングID（自動採番） |
| user1_id | INTEGER | NOT NULL, FK | ユーザー1のID |
| user2_id | INTEGER | NOT NULL, FK | ユーザー2のID |
| status | VARCHAR(50) | DEFAULT 'pending' | マッチング状態 |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

**インデックス**:
- `idx_matchings_user1_id` (user1_id)
- `idx_matchings_user2_id` (user2_id)
- `idx_matchings_created_at` (created_at)

**外部キー制約**:
- `matchings_user1_id_fkey`: user1_id → users(id) ON DELETE CASCADE
- `matchings_user2_id_fkey`: user2_id → users(id) ON DELETE CASCADE

**更新トリガー**:
- `update_matchings_updated_at`: 更新時にupdated_atを自動更新

#### 1.5 bookmarks テーブル
ユーザーのブックマーク情報を格納

| カラム名 | データ型 | 制約 | 説明 |
|---------|---------|------|------|
| id | SERIAL | PRIMARY KEY | ブックマークID（自動採番） |
| user_id | INTEGER | NOT NULL, FK | ブックマークするユーザーID |
| bookmarked_user_id | INTEGER | NOT NULL, FK | ブックマークされるユーザーID |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

**インデックス**:
- `idx_bookmarks_user_id` (user_id)
- `idx_bookmarks_bookmarked_user_id` (bookmarked_user_id)
- `idx_bookmarks_created_at` (created_at)

**外部キー制約**:
- `bookmarks_user_id_fkey`: user_id → users(id) ON DELETE CASCADE
- `bookmarks_bookmarked_user_id_fkey`: bookmarked_user_id → users(id) ON DELETE CASCADE

**UNIQUE制約**:
- `bookmarks_user_bookmarked_unique`: (user_id, bookmarked_user_id) - 重複ブックマーク防止

**CHECK制約**:
- `bookmarks_no_self_bookmark`: user_id != bookmarked_user_id - 自己ブックマーク防止

### 2. コンテスト関連テーブル

#### 2.1 contests テーブル
プログラミングコンテストの情報を格納

| カラム名 | データ型 | 制約 | 説明 |
|---------|---------|------|------|
| id | SERIAL | PRIMARY KEY | コンテストID（自動採番） |
| backend_quota | INTEGER | NOT NULL, DEFAULT 0 | バックエンド枠 |
| frontend_quota | INTEGER | NOT NULL, DEFAULT 0 | フロントエンド枠 |
| ai_quota | INTEGER | NOT NULL, DEFAULT 0 | AI枠 |
| application_deadline | TIMESTAMP WITH TIME ZONE | NOT NULL | 応募締切日時 |
| purpose | TEXT | NOT NULL | コンテストの目的 |
| message | TEXT | NOT NULL | メッセージ |
| author_id | INTEGER | NOT NULL, FK | 作成者ID |
| start_time | TIMESTAMP WITH TIME ZONE | | 開始日時 |
| end_time | TIMESTAMP WITH TIME ZONE | | 終了日時 |
| title | VARCHAR(255) | | タイトル |
| description | TEXT | | 詳細説明 |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

**インデックス**:
- `idx_contests_author_id` (author_id)
- `idx_contests_application_deadline` (application_deadline)
- `idx_contests_created_at` (created_at)

**外部キー制約**:
- `contests_author_id_fkey`: author_id → users(id) ON DELETE CASCADE

**更新トリガー**:
- `update_contests_updated_at`: 更新時にupdated_atを自動更新

### 3. ハッカソン関連テーブル

#### 3.1 hackathons テーブル
ハッカソンイベントの情報を格納

| カラム名 | データ型 | 制約 | 説明 |
|---------|---------|------|------|
| id | SERIAL | PRIMARY KEY | ハッカソンID（自動採番） |
| name | VARCHAR(255) | NOT NULL | ハッカソン名 |
| description | TEXT | | 詳細説明 |
| start_date | TIMESTAMP WITH TIME ZONE | NOT NULL | 開始日時 |
| end_date | TIMESTAMP WITH TIME ZONE | NOT NULL | 終了日時 |
| registration_start | TIMESTAMP WITH TIME ZONE | NOT NULL | 参加登録開始日時 |
| registration_deadline | TIMESTAMP WITH TIME ZONE | NOT NULL | 参加登録締切日時 |
| max_participants | INTEGER | DEFAULT 0 | 最大参加者数 |
| location | VARCHAR(255) | | 開催場所 |
| organizer | VARCHAR(255) | NOT NULL | 主催者 |
| contact_email | VARCHAR(255) | | 連絡先メールアドレス |
| prize_info | TEXT | | 賞品情報 |
| rules | TEXT | | 参加規則 |
| tech_stack | TEXT | | 技術スタック（JSON文字列） |
| status | VARCHAR(50) | DEFAULT 'upcoming' | ステータス |
| is_public | BOOLEAN | DEFAULT true | 公開状態 |
| banner_url | TEXT | | バナー画像URL |
| website_url | TEXT | | ウェブサイトURL |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

**インデックス**:
- `idx_hackathons_start_date` (start_date)
- `idx_hackathons_status` (status)

**CHECK制約**:
- `hackathons_date_logic`: end_date > start_date AND registration_deadline <= start_date AND registration_start <= registration_deadline
- `hackathons_max_participants_check`: max_participants >= 0
- `hackathons_status_check`: status IN ('upcoming', 'ongoing', 'completed', 'cancelled')

**更新トリガー**:
- `update_hackathons_updated_at`: 更新時にupdated_atを自動更新

#### 3.2 hackathon_participants テーブル
ハッカソン参加者の情報を格納

| カラム名 | データ型 | 制約 | 説明 |
|---------|---------|------|------|
| id | SERIAL | PRIMARY KEY | 参加者ID（自動採番） |
| hackathon_id | INTEGER | NOT NULL, FK | ハッカソンID |
| user_id | INTEGER | NOT NULL, FK | ユーザーID |
| team_name | VARCHAR(255) | | チーム名 |
| role | VARCHAR(100) | | 役割 |
| registration_date | TIMESTAMP WITH TIME ZONE | AUTO CREATE TIME | 参加登録日時 |
| status | VARCHAR(50) | DEFAULT 'registered' | 参加状態 |
| notes | TEXT | | 備考 |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

**インデックス**:
- `idx_hackathon_participants_hackathon_id` (hackathon_id)
- `idx_hackathon_participants_user_id` (user_id)

**外部キー制約**:
- `hackathon_participants_hackathon_id_fkey`: hackathon_id → hackathons(id) ON DELETE CASCADE
- `hackathon_participants_user_id_fkey`: user_id → users(id) ON DELETE CASCADE

**UNIQUE制約**:
- `hackathon_participants_unique`: (hackathon_id, user_id) - 重複参加防止

**CHECK制約**:
- `hackathon_participants_status_check`: status IN ('registered', 'confirmed', 'cancelled', 'disqualified')

### 4. ファイル管理関連テーブル

#### 4.1 file_metadata テーブル
ファイルのメタデータを格納

| カラム名 | データ型 | 制約 | 説明 |
|---------|---------|------|------|
| id | VARCHAR(255) | PRIMARY KEY | ファイルID（UUID） |
| name | VARCHAR(255) | UNIQUE, NOT NULL | ファイル名 |
| path | VARCHAR(500) | NOT NULL | ファイルパス |
| size | BIGINT | NOT NULL | ファイルサイズ（バイト） |
| content_type | VARCHAR(100) | NOT NULL | MIMEタイプ |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

**インデックス**:
- `idx_file_metadata_name` (name)
- `idx_file_metadata_path` (path)
- `idx_file_metadata_created_at` (created_at)

## データベース関連設定

### 1. 本番環境接続設定 (GCP Cloud SQL)
- **インスタンス**: `prd-db`
- **リージョン**: `asia-northeast1`
- **接続名**: `zenn-ai-agent-hackathon-460205:asia-northeast1:prd-db`
- **エンジン**: PostgreSQL 17
- **ティア**: `db-f1-micro`
- **接続方式**: Cloud SQL Auth Proxy
- **プロキシポート**: 5434 (ローカル開発時)
- **SSL**: 必須

### 2. ローカル開発環境設定
- **ホスト**: 環境変数 `DB_HOST` で設定（デフォルト: localhost）
- **ポート**: 環境変数 `DB_PORT` で設定（デフォルト: 5432）
- **データベース名**: 環境変数 `DB_NAME` で設定（デフォルト: tru_s3）
- **ユーザー名**: 環境変数 `DB_USER` で設定（デフォルト: postgres）
- **パスワード**: 環境変数 `DB_PASSWORD` で設定
- **SSL モード**: 環境変数 `DB_SSLMODE` で設定

### 3. 接続プール設定
- **最大アイドル接続数**: 10
- **最大オープン接続数**: 100
- **接続最大生存時間**: 1時間

### 4. GORM設定
- **ログレベル**: 開発環境では詳細ログ、本番環境では警告レベル
- **自動マイグレーション**: アプリケーション起動時に実行
- **プリペアードステートメント**: 有効
- **外部キー制約**: 有効

## データベース拡張機能

### 1. PostgreSQL拡張
- **uuid-ossp**: UUID生成機能

### 2. トリガー関数
```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';
```

## マイグレーション管理

### 1. 自動マイグレーション
GORMの`AutoMigrate`機能を使用して、アプリケーション起動時に自動的にテーブル構造を更新します。

### 2. マイグレーションファイル
- `migrations/init.sql`: 初期セットアップ
- `migrations/001_create_matching_app_tables.sql`: マッチングアプリ用テーブル
- `migrations/002_create_contests_table.sql`: コンテストテーブル
- `migrations/003_create_hackathons_table.sql`: ハッカソンテーブル
- `migrations/004_complete_schema_gcp.sql`: 完全スキーマ（GCP用）

### 3. インデックス作成
パフォーマンス向上のため、以下のインデックスが自動作成されます：

```sql
-- ファイル関連
CREATE INDEX IF NOT EXISTS idx_file_metadata_name ON file_metadata(name);
CREATE INDEX IF NOT EXISTS idx_file_metadata_path ON file_metadata(path);
CREATE INDEX IF NOT EXISTS idx_file_metadata_created_at ON file_metadata(created_at);

-- コンテスト関連
CREATE INDEX IF NOT EXISTS idx_contests_application_deadline ON contests(application_deadline);
CREATE INDEX IF NOT EXISTS idx_contests_created_at ON contests(created_at);

-- ハッカソン関連
CREATE INDEX IF NOT EXISTS idx_hackathons_start_date ON hackathons(start_date);
CREATE INDEX IF NOT EXISTS idx_hackathons_status ON hackathons(status);
CREATE INDEX IF NOT EXISTS idx_hackathon_participants_hackathon_id ON hackathon_participants(hackathon_id);

-- ユーザー関連
CREATE INDEX IF NOT EXISTS idx_users_gmail ON users(gmail);
CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);
```

## セキュリティ考慮事項

### 1. データ保護
- **パスワード**: 現在未実装（今後ハッシュ化を予定）
- **個人情報**: 最小限の情報のみ格納
- **削除制約**: CASCADE設定により関連データも自動削除

### 2. アクセス制御
- **データベースユーザー**: 最小権限の原則
- **SSL/TLS**: 本番環境では必須
- **バックアップ**: 定期的なバックアップを推奨

### 3. データ整合性
- **外部キー制約**: 全ての関連テーブルで設定
- **UNIQUE制約**: 重複データ防止
- **CHECK制約**: データの論理的整合性保証

## パフォーマンス最適化

### 1. インデックス戦略
- **検索頻度の高いカラム**: インデックス作成済み
- **外部キー**: 自動的にインデックス作成
- **複合インデックス**: 必要に応じて追加可能

### 2. クエリ最適化
- **N+1問題**: GORMのPreload機能で解決
- **ページネーション**: OFFSET/LIMIT使用
- **フィルタリング**: WHERE句の最適化

### 3. 接続プール
- **適切なプールサイズ**: アプリケーション負荷に応じて調整
- **接続タイムアウト**: デッドロック防止

## 今後の拡張予定

### 1. 機能追加
- **ユーザー認証テーブル**: パスワード、トークン管理
- **ログテーブル**: 監査ログ、アクセスログ
- **通知テーブル**: プッシュ通知、メール通知

### 2. パフォーマンス向上
- **読み取り専用レプリカ**: 読み取り負荷分散
- **パーティショニング**: 大量データ対応
- **キャッシュ**: Redis連携

### 3. 運用改善
- **データ移行ツール**: 自動マイグレーション拡張
- **モニタリング**: パフォーマンス監視
- **バックアップ自動化**: 定期バックアップシステム