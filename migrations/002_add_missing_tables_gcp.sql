-- GCP用不足テーブル追加マイグレーション

-- file_metadata テーブル（GORMで自動作成されないため手動追加）
CREATE TABLE IF NOT EXISTS file_metadata (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(500) NOT NULL,
    size BIGINT NOT NULL,
    content_type VARCHAR(100),
    checksum VARCHAR(64),
    tags TEXT,
    created_at BIGINT,
    updated_at BIGINT
);

-- file_metadata用インデックス
CREATE INDEX IF NOT EXISTS idx_file_metadata_name ON file_metadata(name);
CREATE INDEX IF NOT EXISTS idx_file_metadata_path ON file_metadata(path);
CREATE INDEX IF NOT EXISTS idx_file_metadata_created_at ON file_metadata(created_at);

-- hackathon_participantsテーブルに外部キー制約を追加
-- （既にテーブルは作成済みなので、制約のみ追加）
DO $$
BEGIN
    -- usersテーブルへの外部キー制約を追加（存在しない場合のみ）
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'hackathon_participants_user_id_fkey'
    ) THEN
        ALTER TABLE hackathon_participants 
        ADD CONSTRAINT hackathon_participants_user_id_fkey 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
    
    -- hackathonsテーブルへの外部キー制約を追加（存在しない場合のみ）
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'hackathon_participants_hackathon_id_fkey'
    ) THEN
        ALTER TABLE hackathon_participants 
        ADD CONSTRAINT hackathon_participants_hackathon_id_fkey 
        FOREIGN KEY (hackathon_id) REFERENCES hackathons(id) ON DELETE CASCADE;
    END IF;
END
$$;