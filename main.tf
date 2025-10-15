# LEP System - Complete GCP Infrastructure
# Este arquivo gerencia TODOS os recursos GCP necessários para o LEP System

terraform {
  required_version = ">= 1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# ==============================================================================
# ENABLE REQUIRED APIS
# ==============================================================================

resource "google_project_service" "services" {
  for_each = toset([
    "run.googleapis.com",
    "sql-component.googleapis.com",
    "sqladmin.googleapis.com",
    "storage.googleapis.com",
    "secretmanager.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com",
    "cloudbuild.googleapis.com",
    "compute.googleapis.com"
  ])

  project = var.project_id
  service = each.key

  disable_on_destroy = false
}

# ==============================================================================
# RANDOM PASSWORD FOR DATABASE
# ==============================================================================

resource "random_password" "db_password" {
  length  = 32
  special = true
}

# ==============================================================================
# CLOUD SQL - POSTGRESQL DATABASE
# ==============================================================================

resource "google_sql_database_instance" "main" {
  name             = "${var.project_name}-postgres-${var.environment}"
  database_version = "POSTGRES_15"
  region           = var.region
  project          = var.project_id

  settings {
    tier              = var.db_tier
    availability_type = var.db_availability_type
    disk_size         = var.db_disk_size
    disk_type         = "PD_SSD"

    backup_configuration {
      enabled                        = true
      start_time                     = "03:00"
      point_in_time_recovery_enabled = true
      transaction_log_retention_days = 7
      backup_retention_settings {
        retained_backups = 7
        retention_unit   = "COUNT"
      }
    }

    ip_configuration {
      ipv4_enabled    = false
      private_network = null
      require_ssl     = false
    }

    database_flags {
      name  = "max_connections"
      value = "100"
    }

    database_flags {
      name  = "shared_buffers"
      value = "32768" # 256MB in 8KB pages
    }

    maintenance_window {
      day  = 7 # Sunday
      hour = 3
    }

    insights_config {
      query_insights_enabled  = true
      query_plans_per_minute  = 5
      query_string_length     = 1024
      record_application_tags = true
    }
  }

  deletion_protection = var.enable_deletion_protection

  depends_on = [google_project_service.services]
}

resource "google_sql_database" "main" {
  name     = var.database_name
  instance = google_sql_database_instance.main.name
  project  = var.project_id

  charset   = "UTF8"
  collation = "en_US.UTF8"
}

resource "google_sql_user" "main" {
  name     = var.database_user
  instance = google_sql_database_instance.main.name
  password = random_password.db_password.result
  project  = var.project_id
}

# ==============================================================================
# GOOGLE CLOUD STORAGE - IMAGE BUCKETS
# ==============================================================================

resource "google_storage_bucket" "images" {
  name          = var.bucket_name != "" ? var.bucket_name : "${var.project_id}-${var.project_name}-images-${var.environment}"
  location      = var.region
  project       = var.project_id
  force_destroy = var.environment != "prod"

  uniform_bucket_level_access = true

  versioning {
    enabled = var.environment == "prod"
  }

  lifecycle_rule {
    condition {
      age = 90
    }
    action {
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }
  }

  cors {
    origin          = ["*"]
    method          = ["GET", "HEAD", "POST", "PUT", "DELETE"]
    response_header = ["*"]
    max_age_seconds = 3600
  }

  depends_on = [google_project_service.services]
}

# Make bucket publicly readable
resource "google_storage_bucket_iam_member" "images_public_read" {
  bucket = google_storage_bucket.images.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

# ==============================================================================
# SERVICE ACCOUNT
# ==============================================================================

resource "google_service_account" "backend" {
  account_id   = "${var.project_name}-backend-sa"
  display_name = "LEP Backend Service Account"
  description  = "Service account for LEP backend application running on Cloud Run"
  project      = var.project_id

  depends_on = [google_project_service.services]
}

# Cloud SQL Client permissions
resource "google_project_iam_member" "backend_sql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

# Storage Admin permissions
resource "google_project_iam_member" "backend_storage_admin" {
  project = var.project_id
  role    = "roles/storage.admin"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

# Secret Manager permissions
resource "google_project_iam_member" "backend_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

# Logging permissions
resource "google_project_iam_member" "backend_log_writer" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

# Monitoring permissions
resource "google_project_iam_member" "backend_monitoring_writer" {
  project = var.project_id
  role    = "roles/monitoring.metricWriter"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

# ==============================================================================
# SECRET MANAGER - SECRETS
# ==============================================================================

# Database Password Secret
resource "google_secret_manager_secret" "db_password" {
  secret_id = "db-password-${var.environment}"
  project   = var.project_id

  replication {
    auto {}
  }

  depends_on = [google_project_service.services]
}

resource "google_secret_manager_secret_version" "db_password" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = random_password.db_password.result
}

# JWT Private Key Secret
resource "google_secret_manager_secret" "jwt_private_key" {
  secret_id = "jwt-private-key-${var.environment}"
  project   = var.project_id

  replication {
    auto {}
  }

  depends_on = [google_project_service.services]
}

resource "google_secret_manager_secret_version" "jwt_private_key" {
  secret      = google_secret_manager_secret.jwt_private_key.id
  secret_data = var.jwt_private_key
}

# JWT Public Key Secret
resource "google_secret_manager_secret" "jwt_public_key" {
  secret_id = "jwt-public-key-${var.environment}"
  project   = var.project_id

  replication {
    auto {}
  }

  depends_on = [google_project_service.services]
}

resource "google_secret_manager_secret_version" "jwt_public_key" {
  secret      = google_secret_manager_secret.jwt_public_key.id
  secret_data = var.jwt_public_key
}

# Twilio Account SID (optional)
resource "google_secret_manager_secret" "twilio_account_sid" {
  count     = var.twilio_account_sid != "" ? 1 : 0
  secret_id = "twilio-account-sid-${var.environment}"
  project   = var.project_id

  replication {
    auto {}
  }

  depends_on = [google_project_service.services]
}

resource "google_secret_manager_secret_version" "twilio_account_sid" {
  count       = var.twilio_account_sid != "" ? 1 : 0
  secret      = google_secret_manager_secret.twilio_account_sid[0].id
  secret_data = var.twilio_account_sid
}

# Twilio Auth Token (optional)
resource "google_secret_manager_secret" "twilio_auth_token" {
  count     = var.twilio_auth_token != "" ? 1 : 0
  secret_id = "twilio-auth-token-${var.environment}"
  project   = var.project_id

  replication {
    auto {}
  }

  depends_on = [google_project_service.services]
}

resource "google_secret_manager_secret_version" "twilio_auth_token" {
  count       = var.twilio_auth_token != "" ? 1 : 0
  secret      = google_secret_manager_secret.twilio_auth_token[0].id
  secret_data = var.twilio_auth_token
}

# SMTP Password (optional)
resource "google_secret_manager_secret" "smtp_password" {
  count     = var.smtp_password != "" ? 1 : 0
  secret_id = "smtp-password-${var.environment}"
  project   = var.project_id

  replication {
    auto {}
  }

  depends_on = [google_project_service.services]
}

resource "google_secret_manager_secret_version" "smtp_password" {
  count       = var.smtp_password != "" ? 1 : 0
  secret      = google_secret_manager_secret.smtp_password[0].id
  secret_data = var.smtp_password
}

# ==============================================================================
# OUTPUTS
# ==============================================================================

output "database_instance_name" {
  description = "Cloud SQL instance name"
  value       = google_sql_database_instance.main.name
}

output "database_connection_name" {
  description = "Cloud SQL instance connection name for Cloud Run"
  value       = google_sql_database_instance.main.connection_name
}

output "database_user" {
  description = "Database username"
  value       = google_sql_user.main.name
}

output "database_name" {
  description = "Database name"
  value       = google_sql_database.main.name
}

output "storage_bucket_name" {
  description = "GCS bucket name for images"
  value       = google_storage_bucket.images.name
}

output "storage_bucket_url" {
  description = "Public URL for GCS bucket"
  value       = "https://storage.googleapis.com/${google_storage_bucket.images.name}"
}

output "service_account_email" {
  description = "Service account email for Cloud Run"
  value       = google_service_account.backend.email
}

output "db_password_secret_name" {
  description = "Secret Manager secret name for database password"
  value       = google_secret_manager_secret.db_password.secret_id
}

output "jwt_private_key_secret_name" {
  description = "Secret Manager secret name for JWT private key"
  value       = google_secret_manager_secret.jwt_private_key.secret_id
}

output "jwt_public_key_secret_name" {
  description = "Secret Manager secret name for JWT public key"
  value       = google_secret_manager_secret.jwt_public_key.secret_id
}

# ==============================================================================
# NOTES FOR CLOUD RUN DEPLOYMENT
# ==============================================================================

# Cloud Run service is deployed via gcloud command (see scripts/stage-deploy.sh)
# Example deployment command:
#
# gcloud run deploy lep-system \
#   --source . \
#   --region=us-central1 \
#   --platform=managed \
#   --allow-unauthenticated \
#   --memory=1Gi \
#   --cpu=2 \
#   --min-instances=0 \
#   --max-instances=10 \
#   --add-cloudsql-instances=PROJECT_ID:REGION:INSTANCE_NAME \
#   --service-account=SERVICE_ACCOUNT_EMAIL \
#   --set-env-vars="ENVIRONMENT=stage,STORAGE_TYPE=gcs,BUCKET_NAME=BUCKET_NAME,BASE_URL=https://storage.googleapis.com/BUCKET_NAME,DB_USER=DB_USER,DB_NAME=DB_NAME,INSTANCE_UNIX_SOCKET=/cloudsql/CONNECTION_NAME" \
#   --set-secrets="DB_PASS=db-password-ENV:latest,JWT_SECRET_PRIVATE_KEY=jwt-private-key-ENV:latest,JWT_SECRET_PUBLIC_KEY=jwt-public-key-ENV:latest"
