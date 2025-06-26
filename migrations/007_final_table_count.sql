-- 最終テーブル数確認
SELECT 
    'TOTAL_TABLES' as check_type,
    COUNT(*) as count
FROM information_schema.tables 
WHERE table_schema = 'public';

-- 必須テーブル数確認
SELECT 
    'REQUIRED_TABLES' as check_type,
    COUNT(*) as count
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name IN ('users', 'tags', 'profiles', 'matchings', 'bookmarks', 'contests', 'file_metadata', 'hackathons', 'hackathon_participants');

-- ハッカソンデータ確認
SELECT 
    'HACKATHON_COUNT' as check_type,
    COUNT(*) as count
FROM hackathons;