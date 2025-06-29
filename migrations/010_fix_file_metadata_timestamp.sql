-- file_metadataテーブルのcreated_atカラムの型変換問題を解決
-- 作成日: 2025-06-29

-- 既存のcreated_atカラムを削除（データがある場合は注意）
-- まず、データの確認
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'file_metadata' AND column_name = 'created_at';

-- データがある場合はバックアップを取る
-- CREATE TABLE file_metadata_backup AS SELECT * FROM file_metadata;

-- 既存のcreated_atカラムを削除
ALTER TABLE file_metadata DROP COLUMN IF EXISTS created_at;

-- 新しいcreated_atカラムを追加（timestamptz型）
ALTER TABLE file_metadata 
ADD COLUMN created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- updated_atカラムも同様に修正
ALTER TABLE file_metadata DROP COLUMN IF EXISTS updated_at;
ALTER TABLE file_metadata 
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_file_metadata_created_at ON file_metadata(created_at);
CREATE INDEX IF NOT EXISTS idx_file_metadata_updated_at ON file_metadata(updated_at);

-- 確認
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'file_metadata' 
AND column_name IN ('created_at', 'updated_at'); 