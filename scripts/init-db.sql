-- LEP System Database Initialization Script
-- This script runs automatically when the PostgreSQL container is first created

-- Ensure the database uses UTF-8 encoding
SET client_encoding = 'UTF8';

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Grant privileges to the application user
GRANT ALL PRIVILEGES ON DATABASE lep_database TO lep_user;

-- Log initialization complete
DO $$
BEGIN
    RAISE NOTICE 'LEP Database initialization complete!';
END $$;
