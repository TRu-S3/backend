-- コンテストテーブルに title, description, start_time, end_time カラムを追加
-- 作成日: 2025-06-29

-- 不足しているカラムを追加
ALTER TABLE contests 
ADD COLUMN IF NOT EXISTS title TEXT,
ADD COLUMN IF NOT EXISTS description TEXT,
ADD COLUMN IF NOT EXISTS start_time TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS end_time TIMESTAMP WITH TIME ZONE;

-- インデックス作成（必要に応じて）
CREATE INDEX IF NOT EXISTS idx_contests_title ON contests(title);
CREATE INDEX IF NOT EXISTS idx_contests_start_time ON contests(start_time);
CREATE INDEX IF NOT EXISTS idx_contests_end_time ON contests(end_time);

-- コメント追加
COMMENT ON COLUMN contests.title IS はハッカソンタイトル';
COMMENT ON COLUMN contests.description IS 'ハッカソン説明';
COMMENT ON COLUMN contests.start_time IS 'ハッカソン開始日時';
COMMENT ON COLUMN contests.end_time IS 'ハッカソン終了日時';

-- 既存データの確認
SELECT 
    id,
    purpose,
    message,
    title,
    description,
    start_time,
    end_time
FROM contests 
LIMIT 5; 