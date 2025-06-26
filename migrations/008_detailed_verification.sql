-- GCPデータベースの詳細検証クエリ

-- 1. 全テーブル一覧の表示
SELECT 
    '=== ALL TABLES IN DATABASE ===' as section,
    '' as table_name,
    '' as table_type,
    '' as row_count;

SELECT 
    'TABLE' as section,
    table_name,
    table_type,
    '' as row_count
FROM information_schema.tables 
WHERE table_schema = 'public' 
ORDER BY table_name;

-- 2. 必須テーブルの存在確認
SELECT 
    '=== REQUIRED TABLES CHECK ===' as section,
    '' as table_name,
    '' as table_type,
    '' as row_count;

SELECT 
    'REQUIRED' as section,
    'users' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as table_type,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users' AND table_schema = 'public') 
         THEN (SELECT COUNT(*)::text FROM users) 
         ELSE 'N/A' END as row_count
UNION ALL
SELECT 
    'REQUIRED' as section,
    'tags' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'tags' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as table_type,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'tags' AND table_schema = 'public') 
         THEN (SELECT COUNT(*)::text FROM tags) 
         ELSE 'N/A' END as row_count
UNION ALL
SELECT 
    'REQUIRED' as section,
    'profiles' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'profiles' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as table_type,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'profiles' AND table_schema = 'public') 
         THEN (SELECT COUNT(*)::text FROM profiles) 
         ELSE 'N/A' END as row_count
UNION ALL
SELECT 
    'REQUIRED' as section,
    'matchings' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'matchings' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as table_type,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'matchings' AND table_schema = 'public') 
         THEN (SELECT COUNT(*)::text FROM matchings) 
         ELSE 'N/A' END as row_count
UNION ALL
SELECT 
    'REQUIRED' as section,
    'bookmarks' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'bookmarks' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as table_type,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'bookmarks' AND table_schema = 'public') 
         THEN (SELECT COUNT(*)::text FROM bookmarks) 
         ELSE 'N/A' END as row_count
UNION ALL
SELECT 
    'REQUIRED' as section,
    'contests' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'contests' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as table_type,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'contests' AND table_schema = 'public') 
         THEN (SELECT COUNT(*)::text FROM contests) 
         ELSE 'N/A' END as row_count
UNION ALL
SELECT 
    'REQUIRED' as section,
    'file_metadata' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'file_metadata' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as table_type,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'file_metadata' AND table_schema = 'public') 
         THEN (SELECT COUNT(*)::text FROM file_metadata) 
         ELSE 'N/A' END as row_count
UNION ALL
SELECT 
    'REQUIRED' as section,
    'hackathons' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'hackathons' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as table_type,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'hackathons' AND table_schema = 'public') 
         THEN (SELECT COUNT(*)::text FROM hackathons) 
         ELSE 'N/A' END as row_count
UNION ALL
SELECT 
    'REQUIRED' as section,
    'hackathon_participants' as table_name,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'hackathon_participants' AND table_schema = 'public') THEN '✅ EXISTS' ELSE '❌ MISSING' END as table_type,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'hackathon_participants' AND table_schema = 'public') 
         THEN (SELECT COUNT(*)::text FROM hackathon_participants) 
         ELSE 'N/A' END as row_count
ORDER BY section, table_name;

-- 3. ハッカソンデータの詳細確認
SELECT 
    '=== HACKATHON DATA DETAILS ===' as section,
    '' as name,
    '' as organizer,
    '' as status;

SELECT 
    'HACKATHON' as section,
    name,
    organizer,
    status
FROM hackathons 
ORDER BY created_at;

-- 4. 統計サマリー
SELECT 
    '=== SUMMARY STATISTICS ===' as section,
    '' as metric,
    '' as value,
    '' as details;

SELECT 
    'STATS' as section,
    'Total Tables' as metric,
    COUNT(*)::text as value,
    'in public schema' as details
FROM information_schema.tables 
WHERE table_schema = 'public'
UNION ALL
SELECT 
    'STATS' as section,
    'Required Tables Found' as metric,
    COUNT(*)::text as value,
    'out of 9 required' as details
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name IN ('users', 'tags', 'profiles', 'matchings', 'bookmarks', 'contests', 'file_metadata', 'hackathons', 'hackathon_participants')
UNION ALL
SELECT 
    'STATS' as section,
    'Hackathons Available' as metric,
    COUNT(*)::text as value,
    'sample hackathons' as details
FROM hackathons;