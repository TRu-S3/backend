-- Fix file_metadata table schema
-- First, check current schema
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'file_metadata' 
ORDER BY ordinal_position;

-- Fix created_at and updated_at columns if they are not timestamp with timezone
DO $$
BEGIN
    -- Check if created_at is not timestamptz
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'file_metadata' 
        AND column_name = 'created_at' 
        AND data_type != 'timestamp with time zone'
    ) THEN
        -- Drop and recreate created_at column
        ALTER TABLE file_metadata DROP COLUMN IF EXISTS created_at;
        ALTER TABLE file_metadata ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
    
    -- Check if updated_at is not timestamptz
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'file_metadata' 
        AND column_name = 'updated_at' 
        AND data_type != 'timestamp with time zone'
    ) THEN
        -- Drop and recreate updated_at column
        ALTER TABLE file_metadata DROP COLUMN IF EXISTS updated_at;
        ALTER TABLE file_metadata ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
END $$;

-- Update existing records to have proper timestamps
UPDATE file_metadata SET 
    created_at = NOW(), 
    updated_at = NOW() 
WHERE created_at IS NULL OR updated_at IS NULL;

-- Show final schema
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'file_metadata' 
ORDER BY ordinal_position;