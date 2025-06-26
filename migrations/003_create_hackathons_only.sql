-- ハッカソンテーブルのみ作成マイグレーション
-- 作成日: 2025-06-25

-- Hackathons テーブル
CREATE TABLE IF NOT EXISTS hackathons (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,                    -- ハッカソン名
    description TEXT,                              -- 説明
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,  -- 開始日時
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,    -- 終了日時
    registration_start TIMESTAMP WITH TIME ZONE NOT NULL, -- 参加登録開始日時
    registration_deadline TIMESTAMP WITH TIME ZONE NOT NULL, -- 参加登録締切日時
    max_participants INTEGER DEFAULT 0,           -- 最大参加者数（0は無制限）
    location VARCHAR(255),                         -- 開催場所（オンラインの場合はURL）
    organizer VARCHAR(255) NOT NULL,               -- 主催者
    contact_email VARCHAR(255),                    -- 連絡先メールアドレス
    prize_info TEXT,                               -- 賞品・賞金情報
    rules TEXT,                                    -- ルール・規則
    tech_stack TEXT,                               -- 推奨技術スタック（JSON形式）
    status VARCHAR(50) DEFAULT 'upcoming',        -- ステータス（upcoming, ongoing, completed, cancelled）
    is_public BOOLEAN DEFAULT true,                -- 公開/非公開
    banner_url TEXT,                               -- バナー画像URL
    website_url TEXT,                              -- 公式ウェブサイトURL
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- 制約
    CHECK (end_date > start_date),
    CHECK (registration_deadline <= start_date),
    CHECK (registration_start <= registration_deadline),
    CHECK (max_participants >= 0),
    CHECK (status IN ('upcoming', 'ongoing', 'completed', 'cancelled'))
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_hackathons_name ON hackathons(name);
CREATE INDEX IF NOT EXISTS idx_hackathons_start_date ON hackathons(start_date);
CREATE INDEX IF NOT EXISTS idx_hackathons_end_date ON hackathons(end_date);
CREATE INDEX IF NOT EXISTS idx_hackathons_registration_deadline ON hackathons(registration_deadline);
CREATE INDEX IF NOT EXISTS idx_hackathons_status ON hackathons(status);
CREATE INDEX IF NOT EXISTS idx_hackathons_is_public ON hackathons(is_public);
CREATE INDEX IF NOT EXISTS idx_hackathons_created_at ON hackathons(created_at);
CREATE INDEX IF NOT EXISTS idx_hackathons_organizer ON hackathons(organizer);

-- ハッカソン参加者管理テーブル
CREATE TABLE IF NOT EXISTS hackathon_participants (
    id SERIAL PRIMARY KEY,
    hackathon_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    team_name VARCHAR(255),                        -- チーム名
    role VARCHAR(100),                             -- 役割（developer, designer, pm, etc.）
    registration_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'registered',      -- ステータス（registered, confirmed, cancelled, disqualified）
    notes TEXT,                                    -- 参加時のメモ・自己紹介
    
    -- 制約
    UNIQUE(hackathon_id, user_id), -- 同じハッカソンに重複参加防止
    CHECK (status IN ('registered', 'confirmed', 'cancelled', 'disqualified'))
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_hackathon_participants_hackathon_id ON hackathon_participants(hackathon_id);
CREATE INDEX IF NOT EXISTS idx_hackathon_participants_user_id ON hackathon_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_hackathon_participants_team_name ON hackathon_participants(team_name);
CREATE INDEX IF NOT EXISTS idx_hackathon_participants_status ON hackathon_participants(status);
CREATE INDEX IF NOT EXISTS idx_hackathon_participants_registration_date ON hackathon_participants(registration_date);

-- サンプルデータ挿入
INSERT INTO hackathons (
    name, 
    description, 
    start_date, 
    end_date, 
    registration_start, 
    registration_deadline,
    max_participants,
    location,
    organizer,
    contact_email,
    prize_info,
    rules,
    tech_stack,
    status,
    website_url
) VALUES 
(
    'AI Hackathon with Google Cloud',
    'Google CloudのAI/ML技術を活用した革新的なアプリケーションを開発するハッカソンです。Vertex AI、Gemini API、AutoMLなどの最新技術を使って、社会課題解決に取り組みましょう。',
    '2024-08-15 09:00:00+09',
    '2024-08-17 18:00:00+09',
    '2024-07-01 00:00:00+09',
    '2024-08-10 23:59:59+09',
    100,
    'Google Cloud Tokyo オフィス',
    'Google Cloud Japan',
    'hackathon-ai@googlecloud.com',
    '優勝: Google Cloud クレジット $5,000 + Google Pixel、準優勝: Google Cloud クレジット $2,000、3位: Google Cloud クレジット $1,000',
    '・Google Cloud技術の使用必須\n・チーム人数は2-5人\n・オリジナル作品であること\n・48時間以内での開発',
    '["Google Cloud Vertex AI", "Gemini API", "AutoML", "BigQuery ML", "TensorFlow", "Python", "JavaScript", "React", "Node.js"]',
    'upcoming',
    'https://cloud.google.com/events/hackathon-ai'
),
(
    'Agent Development Kit Hackathon with Google Cloud',
    'Google CloudのAgent Development Kitを使用して、次世代のAIエージェントを開発するハッカソンです。自然言語処理、対話AI、マルチモーダルAIの最前線技術を体験できます。',
    '2024-09-20 10:00:00+09',
    '2024-09-22 17:00:00+09',
    '2024-08-01 00:00:00+09',
    '2024-09-15 23:59:59+09',
    80,
    'ハイブリッド開催（東京・オンライン）',
    'Google Cloud Japan',
    'hackathon-adk@googlecloud.com',
    '優勝: Google I/O 2025 参加券 + Google Cloud クレジット $3,000、準優勝: Google Cloud クレジット $1,500、3位: Google Cloud クレジット $800',
    '・Agent Development Kitの使用必須\n・チーム人数は1-4人\n・デモ動画の提出必須\n・エージェントの実用性を重視',
    '["Google Cloud Agent Development Kit", "Dialogflow CX", "Vertex AI", "Natural Language AI", "Speech-to-Text", "Text-to-Speech", "Python", "TypeScript", "Flutter", "Firebase"]',
    'upcoming',
    'https://cloud.google.com/events/hackathon-agent-dev-kit'
);

-- コメント追加
COMMENT ON TABLE hackathons IS 'ハッカソン情報テーブル';
COMMENT ON TABLE hackathon_participants IS 'ハッカソン参加者テーブル';