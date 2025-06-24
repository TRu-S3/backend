-- コンテスト管理テーブル作成マイグレーション
-- 作成日: 2025-06-23

-- Contests テーブル
CREATE TABLE IF NOT EXISTS contests (
    id SERIAL PRIMARY KEY,
    backend_quota INTEGER NOT NULL DEFAULT 0,     -- バックエンド募集人数
    frontend_quota INTEGER NOT NULL DEFAULT 0,    -- フロントエンド募集人数
    ai_quota INTEGER NOT NULL DEFAULT 0,          -- AI募集人数
    application_deadline TIMESTAMP WITH TIME ZONE NOT NULL, -- 募集掲載期日
    purpose TEXT NOT NULL,                        -- 目的（表示される文字列）
    message TEXT NOT NULL,                        -- メッセージ（表示される文字列）
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- 投稿者
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- 投稿日時
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP  -- 更新日時
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_contests_author_id ON contests(author_id);
CREATE INDEX IF NOT EXISTS idx_contests_application_deadline ON contests(application_deadline);
CREATE INDEX IF NOT EXISTS idx_contests_created_at ON contests(created_at);
CREATE INDEX IF NOT EXISTS idx_contests_updated_at ON contests(updated_at);

-- updated_at 自動更新トリガー
CREATE TRIGGER update_contests_updated_at 
    BEFORE UPDATE ON contests 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- コメント追加
COMMENT ON TABLE contests IS 'コンテスト募集情報テーブル';
COMMENT ON COLUMN contests.id IS 'コンテストID（主キー）';
COMMENT ON COLUMN contests.backend_quota IS 'バックエンド募集人数';
COMMENT ON COLUMN contests.frontend_quota IS 'フロントエンド募集人数';
COMMENT ON COLUMN contests.ai_quota IS 'AI募集人数';
COMMENT ON COLUMN contests.application_deadline IS '募集掲載期日';
COMMENT ON COLUMN contests.purpose IS '目的（表示される文字列）';
COMMENT ON COLUMN contests.message IS 'メッセージ（表示される文字列）';
COMMENT ON COLUMN contests.author_id IS '投稿者ID（外部キー）';
COMMENT ON COLUMN contests.created_at IS '投稿日時';
COMMENT ON COLUMN contests.updated_at IS '更新日時';