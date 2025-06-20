-- マッチングアプリ用テーブル作成マイグレーション
-- 作成日: 2025-06-20

-- Users テーブル
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    gmail VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    icon_url TEXT
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_users_gmail ON users(gmail);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Tag テーブル
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);

-- Profile テーブル
CREATE TABLE IF NOT EXISTS profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bio TEXT,
    tag_id INTEGER REFERENCES tags(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id) -- 1ユーザーにつき1プロフィール
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_profiles_tag_id ON profiles(tag_id);

-- Matching テーブル
CREATE TABLE IF NOT EXISTS matchings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    notify_id INTEGER, -- 通知ID（他システムとの連携用）
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_matchings_user_id ON matchings(user_id);
CREATE INDEX IF NOT EXISTS idx_matchings_created_at ON matchings(created_at);
CREATE INDEX IF NOT EXISTS idx_matchings_notify_id ON matchings(notify_id);

-- Bookmark テーブル
CREATE TABLE IF NOT EXISTS bookmarks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bookmarked_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, bookmarked_user_id), -- 同じユーザーを重複ブックマーク防止
    CHECK (user_id != bookmarked_user_id) -- 自分自身をブックマーク防止
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_bookmarked_user_id ON bookmarks(bookmarked_user_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_created_at ON bookmarks(created_at);

-- トリガー関数: updated_at を自動更新
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- updated_at 自動更新トリガー
CREATE TRIGGER update_tags_updated_at 
    BEFORE UPDATE ON tags 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_profiles_updated_at 
    BEFORE UPDATE ON profiles 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_matchings_updated_at 
    BEFORE UPDATE ON matchings 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- コメント追加
COMMENT ON TABLE users IS 'ユーザーマスタテーブル';
COMMENT ON COLUMN users.id IS 'ユーザーID（主キー）';
COMMENT ON COLUMN users.gmail IS 'Gmailアドレス（ユニーク）';
COMMENT ON COLUMN users.name IS 'ユーザー名';
COMMENT ON COLUMN users.created_at IS 'アカウント作成日時';
COMMENT ON COLUMN users.icon_url IS 'プロフィール画像URL';

COMMENT ON TABLE tags IS 'タグマスタテーブル';
COMMENT ON COLUMN tags.id IS 'タグID（主キー）';
COMMENT ON COLUMN tags.name IS 'タグ名（ユニーク）';

COMMENT ON TABLE profiles IS 'ユーザープロフィールテーブル';
COMMENT ON COLUMN profiles.id IS 'プロフィールID（主キー）';
COMMENT ON COLUMN profiles.user_id IS 'ユーザーID（外部キー）';
COMMENT ON COLUMN profiles.bio IS '自己紹介文';
COMMENT ON COLUMN profiles.tag_id IS 'タグID（外部キー）';

COMMENT ON TABLE matchings IS 'マッチング情報テーブル';
COMMENT ON COLUMN matchings.id IS 'マッチングID（主キー）';
COMMENT ON COLUMN matchings.user_id IS 'ユーザーID（外部キー）';
COMMENT ON COLUMN matchings.notify_id IS '通知ID';
COMMENT ON COLUMN matchings.content IS 'マッチング内容・メッセージ';

COMMENT ON TABLE bookmarks IS 'ブックマークテーブル';
COMMENT ON COLUMN bookmarks.id IS 'ブックマークID（主キー）';
COMMENT ON COLUMN bookmarks.user_id IS 'ブックマークしたユーザーID';
COMMENT ON COLUMN bookmarks.bookmarked_user_id IS 'ブックマークされたユーザーID';