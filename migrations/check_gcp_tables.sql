-- GCPデータベースのテーブル一覧を取得
-- このクエリで現在存在するテーブルを確認

-- 1. 全テーブルの一覧
SELECT 
    table_name,
    table_type,
    CASE 
        WHEN table_name IN ('users', 'tags', 'profiles', 'matchings', 'bookmarks', 'contests', 'file_metadata', 'hackathons', 'hackathon_participants') 
        THEN 'REQUIRED'
        ELSE 'OTHER'
    END as status
FROM information_schema.tables 
WHERE table_schema = 'public' 
ORDER BY 
    CASE 
        WHEN table_name IN ('users', 'tags', 'profiles', 'matchings', 'bookmarks', 'contests', 'file_metadata', 'hackathons', 'hackathon_participants') 
        THEN 1 
        ELSE 2 
    END,
    table_name;

-- 2. 必須テーブルの存在確認
SELECT 
    'users' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users' AND table_schema = 'public') THEN 'EXISTS' ELSE 'MISSING' END as status
UNION ALL
SELECT 
    'tags' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'tags' AND table_schema = 'public') THEN 'EXISTS' ELSE 'MISSING' END as status
UNION ALL
SELECT 
    'profiles' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'profiles' AND table_schema = 'public') THEN 'EXISTS' ELSE 'MISSING' END as status
UNION ALL
SELECT 
    'matchings' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'matchings' AND table_schema = 'public') THEN 'EXISTS' ELSE 'MISSING' END as status
UNION ALL
SELECT 
    'bookmarks' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'bookmarks' AND table_schema = 'public') THEN 'EXISTS' ELSE 'MISSING' END as status
UNION ALL
SELECT 
    'contests' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'contests' AND table_schema = 'public') THEN 'EXISTS' ELSE 'MISSING' END as status
UNION ALL
SELECT 
    'file_metadata' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'file_metadata' AND table_schema = 'public') THEN 'EXISTS' ELSE 'MISSING' END as status
UNION ALL
SELECT 
    'hackathons' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'hackathons' AND table_schema = 'public') THEN 'EXISTS' ELSE 'MISSING' END as status
UNION ALL
SELECT 
    'hackathon_participants' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'hackathon_participants' AND table_schema = 'public') THEN 'EXISTS' ELSE 'MISSING' END as status
ORDER BY table_name;