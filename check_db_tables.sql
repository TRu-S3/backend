-- Check if tables exist
SELECT table_name, table_type
FROM information_schema.tables 
WHERE table_schema = 'public'
ORDER BY table_name;

-- Check table counts if tables exist
SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'tags' as table_name, COUNT(*) as count FROM tags
UNION ALL
SELECT 'profiles' as table_name, COUNT(*) as count FROM profiles
UNION ALL
SELECT 'bookmarks' as table_name, COUNT(*) as count FROM bookmarks
UNION ALL
SELECT 'matchings' as table_name, COUNT(*) as count FROM matchings
UNION ALL
SELECT 'file_metadata' as table_name, COUNT(*) as count FROM file_metadata
UNION ALL
SELECT 'contests' as table_name, COUNT(*) as count FROM contests
UNION ALL
SELECT 'hackathons' as table_name, COUNT(*) as count FROM hackathons;