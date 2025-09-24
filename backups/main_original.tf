terraform {
  required_version = ">= 1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.26.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.4"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# Generate random password for database
resource "random_password" "db_password" {
  length  = 32
  special = true
}

# APIs are already enabled manually, so we comment this out to avoid permission issues
# resource "google_project_service" "required_apis" {
#   for_each = toset([
#     "secretmanager.googleapis.com",
#     "sqladmin.googleapis.com",
#     "run.googleapis.com",
#     "cloudbuild.googleapis.com",
#     "artifactregistry.googleapis.com"
#   ])
#
#   project                    = var.project_id
#   service                    = each.value
#   disable_on_destroy         = false
#   disable_dependent_services = false
#
#   timeouts {
#     create = "30m"
#     update = "40m"
#   }
# }

# Create Artifact Registry repository for Docker images
resource "google_artifact_registry_repository" "lep_repo" {
  location      = var.region
  repository_id = "lep-backend"
  description   = "LEP Backend Docker repository"
  format        = "DOCKER"

  # depends_on = [google_project_service.required_apis] # APIs manually enabled
}

# Create Service Account for Cloud Run
resource "google_service_account" "lep_backend_sa" {
  account_id   = "lep-backend-sa"
  display_name = "LEP Backend Service Account"
  description  = "Service account for LEP Backend Cloud Run service"
}

# IAM roles for the service account
resource "google_project_iam_member" "sa_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.lep_backend_sa.email}"
}

resource "google_project_iam_member" "sa_cloudsql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.lep_backend_sa.email}"
}

# Cloud SQL Instance
resource "google_sql_database_instance" "lep_postgres" {
  name                = "${var.project_name}-postgres-${var.environment}"
  region              = var.region
  database_version    = "POSTGRES_15"
  deletion_protection = var.enable_deletion_protection

  settings {
    tier              = var.db_tier
    availability_type = var.db_availability_type
    disk_type         = "PD_SSD"
    disk_size         = var.db_disk_size

    backup_configuration {
      enabled                        = true
      start_time                     = "03:00"
      location                       = var.region
      point_in_time_recovery_enabled = true
    }

    ip_configuration {
      ipv4_enabled = false
    }

    database_flags {
      name  = "log_checkpoints"
      value = "on"
    }

    database_flags {
      name  = "log_connections"
      value = "on"
    }

    database_flags {
      name  = "log_disconnections"
      value = "on"
    }
  }

  # depends_on = [google_project_service.required_apis] # APIs manually enabled
}

# Database
resource "google_sql_database" "lep_database" {
  name     = var.database_name
  instance = google_sql_database_instance.lep_postgres.name
}

# Database user
resource "google_sql_user" "lep_user" {
  name     = var.database_user
  instance = google_sql_database_instance.lep_postgres.name
  password = random_password.db_password.result
}

# Secret Manager secrets
resource "google_secret_manager_secret" "jwt_private_key" {
  secret_id = "jwt-private-key-${var.environment}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "jwt_private_key_version" {
  secret      = google_secret_manager_secret.jwt_private_key.id
  secret_data = var.jwt_private_key
}

resource "google_secret_manager_secret" "jwt_public_key" {
  secret_id = "jwt-public-key-${var.environment}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "jwt_public_key_version" {
  secret      = google_secret_manager_secret.jwt_public_key.id
  secret_data = var.jwt_public_key
}

resource "google_secret_manager_secret" "db_password" {
  secret_id = "db-password-${var.environment}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "db_password_version" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = random_password.db_password.result
}

# Twilio secrets (if provided)
resource "google_secret_manager_secret" "twilio_account_sid" {
  count     = var.twilio_account_sid != "" ? 1 : 0
  secret_id = "twilio-account-sid-${var.environment}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "twilio_account_sid_version" {
  count       = var.twilio_account_sid != "" ? 1 : 0
  secret      = google_secret_manager_secret.twilio_account_sid[0].id
  secret_data = var.twilio_account_sid
}

resource "google_secret_manager_secret" "twilio_auth_token" {
  count     = var.twilio_auth_token != "" ? 1 : 0
  secret_id = "twilio-auth-token-${var.environment}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "twilio_auth_token_version" {
  count       = var.twilio_auth_token != "" ? 1 : 0
  secret      = google_secret_manager_secret.twilio_auth_token[0].id
  secret_data = var.twilio_auth_token
}

resource "google_secret_manager_secret" "twilio_phone_number" {
  count     = var.twilio_phone_number != "" ? 1 : 0
  secret_id = "twilio-phone-number-${var.environment}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "twilio_phone_number_version" {
  count       = var.twilio_phone_number != "" ? 1 : 0
  secret      = google_secret_manager_secret.twilio_phone_number[0].id
  secret_data = var.twilio_phone_number
}

# SMTP secrets (if provided)
resource "google_secret_manager_secret" "smtp_username" {
  count     = var.smtp_username != "" ? 1 : 0
  secret_id = "smtp-username-${var.environment}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "smtp_username_version" {
  count       = var.smtp_username != "" ? 1 : 0
  secret      = google_secret_manager_secret.smtp_username[0].id
  secret_data = var.smtp_username
}

resource "google_secret_manager_secret" "smtp_password" {
  count     = var.smtp_password != "" ? 1 : 0
  secret_id = "smtp-password-${var.environment}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "smtp_password_version" {
  count       = var.smtp_password != "" ? 1 : 0
  secret      = google_secret_manager_secret.smtp_password[0].id
  secret_data = var.smtp_password
}

# IAM bindings for secrets
resource "google_secret_manager_secret_iam_member" "sa_access_jwt_private" {
  secret_id = google_secret_manager_secret.jwt_private_key.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.lep_backend_sa.email}"
}

resource "google_secret_manager_secret_iam_member" "sa_access_jwt_public" {
  secret_id = google_secret_manager_secret.jwt_public_key.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.lep_backend_sa.email}"
}

resource "google_secret_manager_secret_iam_member" "sa_access_db_password" {
  secret_id = google_secret_manager_secret.db_password.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.lep_backend_sa.email}"
}

# Cloud Run Service
resource "google_cloud_run_v2_service" "lep_backend" {
  name         = "${var.project_name}-backend-${var.environment}"
  location     = var.region
  ingress      = "INGRESS_TRAFFIC_ALL"
  launch_stage = "GA"

  template {
    service_account = google_service_account.lep_backend_sa.email

    scaling {
      min_instance_count = var.min_instances
      max_instance_count = var.max_instances
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.lep_postgres.connection_name]
      }
    }

    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.lep_repo.repository_id}/lep-backend:latest"

      ports {
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = var.cpu_limit
          memory = var.memory_limit
        }
      }

      # Database configuration
      env {
        name  = "DB_USER"
        value = google_sql_user.lep_user.name
      }

      env {
        name = "DB_PASS"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.db_password.secret_id
            version = "latest"
          }
        }
      }

      env {
        name  = "DB_NAME"
        value = google_sql_database.lep_database.name
      }

      env {
        name  = "INSTANCE_UNIX_SOCKET"
        value = "/cloudsql/${google_sql_database_instance.lep_postgres.connection_name}"
      }

      # JWT configuration
      env {
        name = "JWT_SECRET_PRIVATE_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.jwt_private_key.secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "JWT_SECRET_PUBLIC_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.jwt_public_key.secret_id
            version = "latest"
          }
        }
      }

      # Twilio configuration (conditional)
      env {
        name = "TWILIO_ACCOUNT_SID"
        value_source {
          secret_key_ref {
            secret  = var.twilio_account_sid != "" ? google_secret_manager_secret.twilio_account_sid[0].secret_id : ""
            version = "latest"
          }
        }
      }

      env {
        name = "TWILIO_AUTH_TOKEN"
        value_source {
          secret_key_ref {
            secret  = var.twilio_auth_token != "" ? google_secret_manager_secret.twilio_auth_token[0].secret_id : ""
            version = "latest"
          }
        }
      }

      env {
        name = "TWILIO_PHONE_NUMBER"
        value_source {
          secret_key_ref {
            secret  = var.twilio_phone_number != "" ? google_secret_manager_secret.twilio_phone_number[0].secret_id : ""
            version = "latest"
          }
        }
      }

      # SMTP configuration
      env {
        name  = "SMTP_HOST"
        value = var.smtp_host
      }

      env {
        name  = "SMTP_PORT"
        value = tostring(var.smtp_port)
      }

      env {
        name = "SMTP_USERNAME"
        value_source {
          secret_key_ref {
            secret  = var.smtp_username != "" ? google_secret_manager_secret.smtp_username[0].secret_id : ""
            version = "latest"
          }
        }
      }

      env {
        name = "SMTP_PASSWORD"
        value_source {
          secret_key_ref {
            secret  = var.smtp_password != "" ? google_secret_manager_secret.smtp_password[0].secret_id : ""
            version = "latest"
          }
        }
      }

      # Application configuration
      env {
        name  = "ENABLE_CRON_JOBS"
        value = tostring(var.enable_cron_jobs)
      }

      env {
        name  = "PORT"
        value = "8080"
      }

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }

      startup_probe {
        http_get {
          path = "/health"
          port = 8080
        }
        initial_delay_seconds = 10
        timeout_seconds       = 5
        period_seconds        = 10
        failure_threshold     = 3
      }

      liveness_probe {
        http_get {
          path = "/health"
          port = 8080
        }
        initial_delay_seconds = 30
        timeout_seconds       = 5
        period_seconds        = 30
        failure_threshold     = 3
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }

  depends_on = [
    google_secret_manager_secret_version.jwt_private_key_version,
    google_secret_manager_secret_version.jwt_public_key_version,
    google_secret_manager_secret_version.db_password_version,
    google_sql_database.lep_database,
    google_sql_user.lep_user
  ]
}

# IAM policy to allow public access (adjust as needed)
resource "google_cloud_run_service_iam_member" "public_access" {
  location = google_cloud_run_v2_service.lep_backend.location
  service  = google_cloud_run_v2_service.lep_backend.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}