# LEP System - GCP Staging Environment
# Production-like configuration without Twilio (SMTP only)

# Project configuration
project_id   = "leps-472702"
project_name = "leps"
region       = "us-central1"
environment  = "stage"

# Database configuration - standardized credentials
database_name              = "lep_database"
database_user              = "lep_user"
db_tier                    = "db-n1-standard-1"  # Production-ready tier
db_availability_type       = "REGIONAL"          # High availability
db_disk_size               = 50                   # Larger disk
enable_deletion_protection = true                # Protect against accidental deletion

# JWT configuration - standardized keys for dev/stage compatibility
jwt_private_key = "dev-simple-private-key-for-testing-only"
jwt_public_key  = "dev-simple-public-key-for-testing-only"

# Cloud Run configuration - PRODUCTION-LIKE
min_instances = 1          # Always keep one instance warm
max_instances = 20         # Higher capacity for load testing
cpu_limit     = "2"        # More CPU power
memory_limit  = "1Gi"      # More memory

# Twilio configuration - DISABLED for staging
twilio_account_sid  = ""
twilio_auth_token   = ""
twilio_phone_number = ""

# SMTP configuration - ENABLED for email testing
smtp_host     = "smtp.gmail.com"
smtp_port     = 587
smtp_username = "your-staging-email@gmail.com"    # UPDATE THIS
smtp_password = "your-app-password"                # UPDATE THIS

# Application configuration
enable_cron_jobs = true    # Enable for staging testing

# Bucket configuration for staging
bucket_name          = "leps-472702-lep-images-stage"
bucket_cache_control = "public, max-age=7200"
bucket_timeout       = 60

# Custom domain - optional for staging
domain_name          = "staging-api.yourdomain.com"  # UPDATE THIS if needed
enable_custom_domain = false                         # Set to true if domain is configured