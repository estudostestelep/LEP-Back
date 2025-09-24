-- LEP System - Local Development Database Initialization
-- This script sets up the initial database schema for local development

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create development user with appropriate permissions
-- (User is already created by POSTGRES_USER, this is just for reference)

-- Grant necessary permissions
GRANT ALL PRIVILEGES ON DATABASE lep_database TO lep_user;
GRANT ALL PRIVILEGES ON SCHEMA public TO lep_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO lep_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO lep_user;

-- Set default permissions for future objects
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO lep_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO lep_user;

-- Optional: Create some initial development data
-- (Your GORM will handle table creation, this is just for any manual setup needed)

-- Log the initialization
DO $$
BEGIN
    RAISE NOTICE 'LEP System local development database initialized successfully';
    RAISE NOTICE 'Database: lep_database';
    RAISE NOTICE 'User: lep_user';
    RAISE NOTICE 'Extensions: uuid-ossp, pgcrypto';
END $$;