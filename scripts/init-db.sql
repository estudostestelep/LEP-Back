-- LEP System - Database Initialization Script
-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE lep_database TO lep_user;
ALTER ROLE lep_user SET search_path TO public;
