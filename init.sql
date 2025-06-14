-- TRu-S3 Database Initialization
-- This file is automatically executed when the PostgreSQL container starts

-- Create database if not exists (this is handled by POSTGRES_DB env var)
-- CREATE DATABASE IF NOT EXISTS tru_s3;

-- Set timezone
SET timezone = 'UTC';

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Log initialization
DO $$
BEGIN
    RAISE NOTICE 'TRu-S3 Database initialized successfully';
END $$;
