-- GCPデータベースの現在の状況確認とレポート

-- 1. 現在存在するテーブル一覧
SELECT 
    'Current Tables:' as section,
    '' as table_name,
    '' as status
UNION ALL
SELECT 
    '',
    table_name,
    table_type
FROM information_schema.tables 
WHERE table_schema = 'public' 
ORDER BY section DESC, table_name;

-- 2. 必須テーブルの存在チェック
SELECT 
    'Required Tables Status:' as section,
    '' as table_name,
    '' as status
UNION ALL
SELECT 
    '',
    'users' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as status
UNION ALL
SELECT 
    '',
    'tags' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'tags' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as status
UNION ALL
SELECT 
    '',
    'profiles' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'profiles' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as status
UNION ALL
SELECT 
    '',
    'matchings' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'matchings' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as status
UNION ALL
SELECT 
    '',
    'bookmarks' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'bookmarks' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as status
UNION ALL
SELECT 
    '',
    'contests' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'contests' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as status
UNION ALL
SELECT 
    '',
    'file_metadata' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'file_metadata' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as status
UNION ALL
SELECT 
    '',
    'hackathons' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'hackathons' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as status
UNION ALL
SELECT 
    '',
    'hackathon_participants' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'hackathon_participants' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as status
ORDER BY section DESC, table_name;

-- 3. ハッカソンデータの確認
SELECT 
    'Hackathon Data:' as section,
    '' as info,
    '' as value
UNION ALL
SELECT 
    '',
    'Total Hackathons',
    CAST(COUNT(*) as TEXT)
FROM hackathons
UNION ALL
SELECT 
    '',
    'Hackathon Names',
    string_agg(name, ', ')
FROM hackathons
ORDER BY section DESC;