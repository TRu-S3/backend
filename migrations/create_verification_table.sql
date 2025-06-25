-- Create verification results table
CREATE TABLE IF NOT EXISTS verification_results (
    id SERIAL PRIMARY KEY,
    test_name VARCHAR(50),
    result_value INTEGER,
    details TEXT,
    verified_at TIMESTAMP DEFAULT NOW()
);

-- Clear previous results
DELETE FROM verification_results;

-- Insert table count verification
INSERT INTO verification_results (test_name, result_value, details)
SELECT 
    'required_tables_found',
    COUNT(*),
    'Tables: ' || STRING_AGG(table_name, ', ' ORDER BY table_name)
FROM information_schema.tables 
WHERE table_schema = 'public' 
  AND table_name IN ('users', 'tags', 'profiles', 'matchings', 'bookmarks', 'contests', 'file_metadata', 'hackathons', 'hackathon_participants');

-- Insert hackathon count
INSERT INTO verification_results (test_name, result_value, details)
SELECT 
    'hackathons_available',
    COUNT(*),
    'Hackathons: ' || STRING_AGG(name, ', ' ORDER BY name)
FROM hackathons;

-- Insert total tables count
INSERT INTO verification_results (test_name, result_value, details)
SELECT 
    'total_tables_in_public',
    COUNT(*),
    'All tables: ' || STRING_AGG(table_name, ', ' ORDER BY table_name)
FROM information_schema.tables 
WHERE table_schema = 'public';