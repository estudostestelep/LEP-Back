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

# Random suffix for unique resource names
resource "random_id" "suffix" {
  byte_length = 4
}

locals {
  resource_suffix = random_id.suffix.hex
  service_name    = "${var.app_name}-${var.environment}"
  db_instance_name = "${var.app_name}-db-${var.environment}-${local.resource_suffix}"
}

# Enable required APIs
resource "google_project_service" "apis" {
  for_each = toset(var.enable_apis)

  project = var.project_id
  service = each.value

  disable_dependent_services = false
  disable_on_destroy        = false
}

# Service Account for Cloud Run
resource "google_service_account" "cloud_run" {
  account_id   = "${var.app_name}-cloudrun-${var.environment}"
  display_name = "Cloud Run Service Account for ${var.app_name} ${var.environment}"
  description  = "Service account for Cloud Run service with access to Cloud SQL and Secret Manager"
}

# IAM roles for Service Account
resource "google_project_iam_member" "secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.cloud_run.email}"

  depends_on = [google_project_service.apis]
}

resource "google_project_iam_member" "cloudsql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.cloud_run.email}"

  depends_on = [google_project_service.apis]
}

# Cloud SQL Instance
resource "google_sql_database_instance" "postgres" {
  name             = local.db_instance_name
  region           = var.region
  database_version = var.database_version

  deletion_protection = var.deletion_protection

  settings {
    tier              = var.database_tier
    availability_type = "ZONAL"
    disk_type         = "PD_SSD"
    disk_size         = 10
    disk_autoresize   = true

    backup_configuration {
      enabled                        = true
      start_time                    = "03:00"
      point_in_time_recovery_enabled = true
      location                      = var.region
      transaction_log_retention_days = 7
      backup_retention_settings {
        retained_backups = 7
        retention_unit   = "COUNT"
      }
    }

    ip_configuration {
      ipv4_enabled = false
      require_ssl  = true
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

  depends_on = [google_project_service.apis]
}

# Database
resource "google_sql_database" "database" {
  name     = var.database_name
  instance = google_sql_database_instance.postgres.name
}

# Database User
resource "google_sql_user" "user" {
  name     = var.database_user
  instance = google_sql_database_instance.postgres.name
  password = var.database_password
}

# Secret Manager Secrets
resource "google_secret_manager_secret" "jwt_private_key" {
  secret_id = "jwt-secret-private-key-${var.environment}"

  labels = {
    environment = var.environment
    app         = var.app_name
  }

  replication {
    auto {}
  }

  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret_version" "jwt_private_key" {
  secret      = google_secret_manager_secret.jwt_private_key.id
  secret_data = var.jwt_secret_private_key
}

resource "google_secret_manager_secret" "jwt_public_key" {
  secret_id = "jwt-secret-public-key-${var.environment}"

  labels = {
    environment = var.environment
    app         = var.app_name
  }

  replication {
    auto {}
  }

  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret_version" "jwt_public_key" {
  secret      = google_secret_manager_secret.jwt_public_key.id
  secret_data = var.jwt_secret_public_key
}

resource "google_secret_manager_secret" "db_password" {
  secret_id = "db-password-${var.environment}"

  labels = {
    environment = var.environment
    app         = var.app_name
  }

  replication {
    auto {}
  }

  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret_version" "db_password" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = var.database_password
}

# IAM permissions for secrets
resource "google_secret_manager_secret_iam_member" "jwt_private_key_access" {
  secret_id = google_secret_manager_secret.jwt_private_key.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_run.email}"
}

resource "google_secret_manager_secret_iam_member" "jwt_public_key_access" {
  secret_id = google_secret_manager_secret.jwt_public_key.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_run.email}"
}

resource "google_secret_manager_secret_iam_member" "db_password_access" {
  secret_id = google_secret_manager_secret.db_password.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_run.email}"
}

# Cloud Run Service
resource "google_cloud_run_v2_service" "app" {
  name     = local.service_name
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  labels = {
    environment = var.environment
    app         = var.app_name
  }

  template {
    service_account = google_service_account.cloud_run.email

    labels = {
      environment = var.environment
      app         = var.app_name
    }

    scaling {
      min_instance_count = 0
      max_instance_count = var.max_instances
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.postgres.connection_name]
      }
    }

    containers {
      image = var.container_image

      ports {
        container_port = 8080
        name          = "http1"
      }

      resources {
        limits = {
          cpu    = "1000m"
          memory = "512Mi"
        }
        cpu_idle = true
        startup_cpu_boost = false
      }

      # Environment variables
      env {
        name  = "DB_USER"
        value = var.database_user
      }

      env {
        name  = "DB_NAME"
        value = var.database_name
      }

      env {
        name  = "INSTANCE_UNIX_SOCKET"
        value = "/cloudsql/${google_sql_database_instance.postgres.connection_name}"
      }

      env {
        name  = "PORT"
        value = "8080"
      }

      # Secret environment variables
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
        timeout_seconds      = 1
        period_seconds       = 10
        failure_threshold    = 3
      }

      liveness_probe {
        http_get {
          path = "/health"
          port = 8080
        }
        initial_delay_seconds = 30
        timeout_seconds      = 1
        period_seconds       = 30
        failure_threshold    = 3
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }

  depends_on = [
    google_secret_manager_secret_version.jwt_private_key,
    google_secret_manager_secret_version.jwt_public_key,
    google_secret_manager_secret_version.db_password,
    google_sql_database.database,
    google_sql_user.user
  ]
}

# Allow unauthenticated access to Cloud Run
resource "google_cloud_run_service_iam_member" "public_access" {
  service  = google_cloud_run_v2_service.app.name
  location = google_cloud_run_v2_service.app.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}