-- Fix all table schemas to ensure compatibility with GORM
-- This script fixes common data type issues

-- Fix users table
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'created_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE users DROP COLUMN IF EXISTS created_at;
        ALTER TABLE users ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'updated_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE users DROP COLUMN IF EXISTS updated_at;
        ALTER TABLE users ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
END $$;

-- Fix profiles table
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'profiles' AND column_name = 'created_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE profiles DROP COLUMN IF EXISTS created_at;
        ALTER TABLE profiles ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'profiles' AND column_name = 'updated_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE profiles DROP COLUMN IF EXISTS updated_at;
        ALTER TABLE profiles ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
END $$;

-- Fix tags table
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'tags' AND column_name = 'created_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE tags DROP COLUMN IF EXISTS created_at;
        ALTER TABLE tags ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'tags' AND column_name = 'updated_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE tags DROP COLUMN IF EXISTS updated_at;
        ALTER TABLE tags ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
END $$;

-- Fix bookmarks table
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'bookmarks' AND column_name = 'created_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE bookmarks DROP COLUMN IF EXISTS created_at;
        ALTER TABLE bookmarks ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'bookmarks' AND column_name = 'updated_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE bookmarks DROP COLUMN IF EXISTS updated_at;
        ALTER TABLE bookmarks ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
END $$;

-- Fix matchings table
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'matchings' AND column_name = 'created_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE matchings DROP COLUMN IF EXISTS created_at;
        ALTER TABLE matchings ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'matchings' AND column_name = 'updated_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE matchings DROP COLUMN IF EXISTS updated_at;
        ALTER TABLE matchings ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
END $$;

-- Fix contests table
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'contests' AND column_name = 'created_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE contests DROP COLUMN IF EXISTS created_at;
        ALTER TABLE contests ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'contests' AND column_name = 'updated_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE contests DROP COLUMN IF EXISTS updated_at;
        ALTER TABLE contests ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'contests' AND column_name = 'application_deadline' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE contests DROP COLUMN IF EXISTS application_deadline;
        ALTER TABLE contests ADD COLUMN application_deadline TIMESTAMPTZ;
    END IF;
END $$;

-- Fix hackathons table
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'hackathons' AND column_name = 'created_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE hackathons DROP COLUMN IF EXISTS created_at;
        ALTER TABLE hackathons ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'hackathons' AND column_name = 'updated_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE hackathons DROP COLUMN IF EXISTS updated_at;
        ALTER TABLE hackathons ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'hackathons' AND column_name = 'start_date' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE hackathons DROP COLUMN IF EXISTS start_date;
        ALTER TABLE hackathons ADD COLUMN start_date TIMESTAMPTZ;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'hackathons' AND column_name = 'end_date' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE hackathons DROP COLUMN IF EXISTS end_date;
        ALTER TABLE hackathons ADD COLUMN end_date TIMESTAMPTZ;
    END IF;
END $$;

-- Fix hackathon_participants table
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'hackathon_participants' AND column_name = 'created_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE hackathon_participants DROP COLUMN IF EXISTS created_at;
        ALTER TABLE hackathon_participants ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'hackathon_participants' AND column_name = 'updated_at' AND data_type != 'timestamp with time zone') THEN
        ALTER TABLE hackathon_participants DROP COLUMN IF EXISTS updated_at;
        ALTER TABLE hackathon_participants ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
END $$;