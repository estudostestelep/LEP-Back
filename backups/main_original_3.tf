# LEP System - Terraform Hybrid Configuration
# Uses data sources for manually created resources (gcloud) + manages Cloud Run

terraform {
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

# Import existing resources created via gcloud bootstrap
data "google_service_account" "lep_backend_sa" {
  account_id = "lep-backend-sa"
  project    = var.project_id
}

data "google_artifact_registry_repository" "lep_repo" {
  repository_id = "lep-backend"
  location      = var.region
  project       = var.project_id
}

data "google_secret_manager_secret" "jwt_private_key" {
  secret_id = "jwt-private-key-${var.environment}"
  project   = var.project_id
}

data "google_secret_manager_secret" "jwt_public_key" {
  secret_id = "jwt-public-key-${var.environment}"
  project   = var.project_id
}

data "google_secret_manager_secret" "db_password" {
  secret_id = "db-password-${var.environment}"
  project   = var.project_id
}

data "google_sql_database_instance" "lep_postgres" {
  name    = "leps-postgres-${var.environment}"
  project = var.project_id
}

# Generate random password for database (if not exists)
resource "random_password" "db_password" {
  length  = 32
  special = true
}

# Store generated password in secret manager
resource "google_secret_manager_secret_version" "db_password_version" {
  secret      = data.google_secret_manager_secret.db_password.id
  secret_data = random_password.db_password.result
}

# Create database on existing instance
resource "google_sql_database" "lep_database" {
  name     = "lep_database"
  instance = data.google_sql_database_instance.lep_postgres.name
  project  = var.project_id
}

# Create database user
resource "google_sql_user" "lep_user" {
  name     = "lep_user"
  instance = data.google_sql_database_instance.lep_postgres.name
  password = random_password.db_password.result
  project  = var.project_id
}

# Cloud Run service (this is what Terraform manages)
resource "google_cloud_run_v2_service" "lep_backend" {
  name         = "${var.project_name}-backend-${var.environment}"
  location     = var.region
  ingress      = "INGRESS_TRAFFIC_ALL"
  launch_stage = "GA"

  template {
    service_account = data.google_service_account.lep_backend_sa.email

    scaling {
      min_instance_count = 0
      max_instance_count = 2
    }

    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/${data.google_artifact_registry_repository.lep_repo.repository_id}/lep-backend:latest"

      ports {
        container_port = 8080
      }

      env {
        name  = "DB_USER"
        value = google_sql_user.lep_user.name
      }

      env {
        name = "DB_PASS"
        value_source {
          secret_key_ref {
            secret  = data.google_secret_manager_secret.db_password.secret_id
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
        value = "/cloudsql/${data.google_sql_database_instance.lep_postgres.connection_name}"
      }

      env {
        name = "JWT_SECRET_PRIVATE_KEY"
        value_source {
          secret_key_ref {
            secret  = data.google_secret_manager_secret.jwt_private_key.secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "JWT_SECRET_PUBLIC_KEY"
        value_source {
          secret_key_ref {
            secret  = data.google_secret_manager_secret.jwt_public_key.secret_id
            version = "latest"
          }
        }
      }

      env {
        name  = "PORT"
        value = "8080"
      }

      env {
        name  = "ENVIRONMENT"
        value = var.environment
      }

      resources {
        limits = {
          cpu    = "1000m"
          memory = "512Mi"
        }
        cpu_idle          = true
        startup_cpu_boost = false
      }

      startup_probe {
        http_get {
          path = "/ping"
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
        initial_delay_seconds = 15
        timeout_seconds       = 5
        period_seconds        = 30
        failure_threshold     = 3
      }
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [data.google_sql_database_instance.lep_postgres.connection_name]
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }

  depends_on = [
    google_sql_database.lep_database,
    google_sql_user.lep_user
  ]
}

# Allow unauthenticated access to Cloud Run service
resource "google_cloud_run_service_iam_binding" "default" {
  location = google_cloud_run_v2_service.lep_backend.location
  service  = google_cloud_run_v2_service.lep_backend.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}