-- Quick verification query
SELECT 'TABLES_CHECK' as test_type, COUNT(*) as result 
FROM information_schema.tables 
WHERE table_schema = 'public' 
  AND table_name IN ('users', 'tags', 'profiles', 'matchings', 'bookmarks', 'contests', 'file_metadata', 'hackathons', 'hackathon_participants');

SELECT 'HACKATHONS_CHECK' as test_type, COUNT(*) as result FROM hackathons;