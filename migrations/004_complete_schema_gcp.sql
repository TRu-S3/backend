-- GCP用完全スキーマ作成マイグレーション
-- 全ての必要なテーブルを安全に作成

-- 既存テーブルの制約を安全に削除してから再作成
-- DROP制約は存在しない場合エラーになるため、IF EXISTSを使用

-- users テーブル（既存の場合は構造を調整）
DO $$
BEGIN
    -- usersテーブルが存在しない場合のみ作成
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users' AND table_schema = 'public') THEN
        CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            gmail VARCHAR(255) UNIQUE NOT NULL,
            name VARCHAR(100) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            icon_url TEXT
        );
        CREATE INDEX IF NOT EXISTS idx_users_gmail ON users(gmail);
        CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
    END IF;
END
$$;

-- tags テーブル
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);

-- profiles テーブル
CREATE TABLE IF NOT EXISTS profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    bio TEXT,
    tag_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_profiles_tag_id ON profiles(tag_id);

-- matchings テーブル
CREATE TABLE IF NOT EXISTS matchings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    notify_id INTEGER,
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_matchings_user_id ON matchings(user_id);
CREATE INDEX IF NOT EXISTS idx_matchings_created_at ON matchings(created_at);
CREATE INDEX IF NOT EXISTS idx_matchings_notify_id ON matchings(notify_id);

-- bookmarks テーブル
CREATE TABLE IF NOT EXISTS bookmarks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    bookmarked_user_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_bookmarked_user_id ON bookmarks(bookmarked_user_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_created_at ON bookmarks(created_at);

-- contests テーブル
CREATE TABLE IF NOT EXISTS contests (
    id SERIAL PRIMARY KEY,
    backend_quota INTEGER NOT NULL DEFAULT 0,
    frontend_quota INTEGER NOT NULL DEFAULT 0,
    ai_quota INTEGER NOT NULL DEFAULT 0,
    application_deadline TIMESTAMP WITH TIME ZONE NOT NULL,
    purpose TEXT NOT NULL,
    message TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_contests_author_id ON contests(author_id);
CREATE INDEX IF NOT EXISTS idx_contests_application_deadline ON contests(application_deadline);
CREATE INDEX IF NOT EXISTS idx_contests_created_at ON contests(created_at);
CREATE INDEX IF NOT EXISTS idx_contests_updated_at ON contests(updated_at);

-- file_metadata テーブル
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
CREATE INDEX IF NOT EXISTS idx_file_metadata_name ON file_metadata(name);
CREATE INDEX IF NOT EXISTS idx_file_metadata_path ON file_metadata(path);
CREATE INDEX IF NOT EXISTS idx_file_metadata_created_at ON file_metadata(created_at);

-- hackathons テーブル（既に作成済みなのでスキップ）
-- hackathon_participants テーブル（既に作成済みなのでスキップ）

-- トリガー関数の作成（既に存在する場合はスキップ）
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- トリガーの作成（既に存在する場合はスキップ）
DO $$
BEGIN
    -- tags用トリガー
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_tags_updated_at') THEN
        CREATE TRIGGER update_tags_updated_at 
            BEFORE UPDATE ON tags 
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    
    -- profiles用トリガー
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_profiles_updated_at') THEN
        CREATE TRIGGER update_profiles_updated_at 
            BEFORE UPDATE ON profiles 
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    
    -- matchings用トリガー
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_matchings_updated_at') THEN
        CREATE TRIGGER update_matchings_updated_at 
            BEFORE UPDATE ON matchings 
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    
    -- contests用トリガー
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_contests_updated_at') THEN
        CREATE TRIGGER update_contests_updated_at 
            BEFORE UPDATE ON contests 
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    
    -- hackathons用トリガー
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_hackathons_updated_at') THEN
        CREATE TRIGGER update_hackathons_updated_at 
            BEFORE UPDATE ON hackathons 
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
END
$$;

-- 後で外部キー制約を追加（テーブル作成後）
DO $$
BEGIN
    -- profiles -> users
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'profiles_user_id_fkey' AND table_name = 'profiles'
    ) THEN
        ALTER TABLE profiles 
        ADD CONSTRAINT profiles_user_id_fkey 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
    
    -- profiles -> tags
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'profiles_tag_id_fkey' AND table_name = 'profiles'
    ) THEN
        ALTER TABLE profiles 
        ADD CONSTRAINT profiles_tag_id_fkey 
        FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE SET NULL;
    END IF;
    
    -- matchings -> users
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'matchings_user_id_fkey' AND table_name = 'matchings'
    ) THEN
        ALTER TABLE matchings 
        ADD CONSTRAINT matchings_user_id_fkey 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
    
    -- bookmarks -> users (user_id)
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'bookmarks_user_id_fkey' AND table_name = 'bookmarks'
    ) THEN
        ALTER TABLE bookmarks 
        ADD CONSTRAINT bookmarks_user_id_fkey 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
    
    -- bookmarks -> users (bookmarked_user_id)
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'bookmarks_bookmarked_user_id_fkey' AND table_name = 'bookmarks'
    ) THEN
        ALTER TABLE bookmarks 
        ADD CONSTRAINT bookmarks_bookmarked_user_id_fkey 
        FOREIGN KEY (bookmarked_user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
    
    -- contests -> users
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'contests_author_id_fkey' AND table_name = 'contests'
    ) THEN
        ALTER TABLE contests 
        ADD CONSTRAINT contests_author_id_fkey 
        FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
    
    -- hackathon_participants -> users
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'hackathon_participants_user_id_fkey' AND table_name = 'hackathon_participants'
    ) THEN
        ALTER TABLE hackathon_participants 
        ADD CONSTRAINT hackathon_participants_user_id_fkey 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
    
    -- hackathon_participants -> hackathons
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'hackathon_participants_hackathon_id_fkey' AND table_name = 'hackathon_participants'
    ) THEN
        ALTER TABLE hackathon_participants 
        ADD CONSTRAINT hackathon_participants_hackathon_id_fkey 
        FOREIGN KEY (hackathon_id) REFERENCES hackathons(id) ON DELETE CASCADE;
    END IF;
END
$$;

-- UNIQUE制約の追加
DO $$
BEGIN
    -- profiles.user_id のユニーク制約
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'profiles_user_id_unique' AND table_name = 'profiles'
    ) THEN
        ALTER TABLE profiles 
        ADD CONSTRAINT profiles_user_id_unique 
        UNIQUE (user_id);
    END IF;
    
    -- bookmarks重複防止制約
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'bookmarks_user_bookmarked_unique' AND table_name = 'bookmarks'
    ) THEN
        ALTER TABLE bookmarks 
        ADD CONSTRAINT bookmarks_user_bookmarked_unique 
        UNIQUE (user_id, bookmarked_user_id);
    END IF;
    
    -- hackathon_participants重複防止制約
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'hackathon_participants_unique' AND table_name = 'hackathon_participants'
    ) THEN
        ALTER TABLE hackathon_participants 
        ADD CONSTRAINT hackathon_participants_unique 
        UNIQUE (hackathon_id, user_id);
    END IF;
END
$$;

-- CHECK制約の追加
DO $$
BEGIN
    -- bookmarks自己参照防止
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints 
        WHERE constraint_name = 'bookmarks_no_self_bookmark'
    ) THEN
        ALTER TABLE bookmarks 
        ADD CONSTRAINT bookmarks_no_self_bookmark 
        CHECK (user_id != bookmarked_user_id);
    END IF;
    
    -- hackathons日付制約
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints 
        WHERE constraint_name = 'hackathons_date_logic'
    ) THEN
        ALTER TABLE hackathons 
        ADD CONSTRAINT hackathons_date_logic 
        CHECK (end_date > start_date AND registration_deadline <= start_date AND registration_start <= registration_deadline);
    END IF;
    
    -- hackathons参加者数制約
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints 
        WHERE constraint_name = 'hackathons_max_participants_check'
    ) THEN
        ALTER TABLE hackathons 
        ADD CONSTRAINT hackathons_max_participants_check 
        CHECK (max_participants >= 0);
    END IF;
    
    -- hackathonsステータス制約
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints 
        WHERE constraint_name = 'hackathons_status_check'
    ) THEN
        ALTER TABLE hackathons 
        ADD CONSTRAINT hackathons_status_check 
        CHECK (status IN ('upcoming', 'ongoing', 'completed', 'cancelled'));
    END IF;
    
    -- hackathon_participantsステータス制約
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints 
        WHERE constraint_name = 'hackathon_participants_status_check'
    ) THEN
        ALTER TABLE hackathon_participants 
        ADD CONSTRAINT hackathon_participants_status_check 
        CHECK (status IN ('registered', 'confirmed', 'cancelled', 'disqualified'));
    END IF;
END
$$;