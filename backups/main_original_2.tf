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

# Reference existing resources created manually
data "google_artifact_registry_repository" "lep_repo" {
  repository_id = "lep-backend"
  location      = var.region
  project       = var.project_id
}

data "google_service_account" "lep_backend_sa" {
  account_id = "lep-backend-sa"
  project    = var.project_id
}

data "google_secret_manager_secret" "jwt_private_key" {
  secret_id = "jwt-private-key-dev"
  project   = var.project_id
}

data "google_secret_manager_secret" "jwt_public_key" {
  secret_id = "jwt-public-key-dev"
  project   = var.project_id
}

data "google_secret_manager_secret" "db_password" {
  secret_id = "db-password-dev"
  project   = var.project_id
}

# Create secret versions
resource "google_secret_manager_secret_version" "jwt_private_key_version" {
  secret      = data.google_secret_manager_secret.jwt_private_key.id
  secret_data = var.jwt_private_key
}

resource "google_secret_manager_secret_version" "jwt_public_key_version" {
  secret      = data.google_secret_manager_secret.jwt_public_key.id
  secret_data = var.jwt_public_key
}

resource "google_secret_manager_secret_version" "db_password_version" {
  secret      = data.google_secret_manager_secret.db_password.id
  secret_data = random_password.db_password.result
}

# Reference Cloud SQL instance (once it's ready)
data "google_sql_database_instance" "lep_postgres" {
  name    = "leps-postgres-dev"
  project = var.project_id
}

# Create database
resource "google_sql_database" "lep_database" {
  name     = var.database_name
  instance = data.google_sql_database_instance.lep_postgres.name
}

# Create database user
resource "google_sql_user" "lep_user" {
  name     = var.database_user
  instance = data.google_sql_database_instance.lep_postgres.name
  password = random_password.db_password.result
}

# IAM roles for the service account
resource "google_project_iam_member" "sa_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${data.google_service_account.lep_backend_sa.email}"
}

resource "google_project_iam_member" "sa_cloudsql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${data.google_service_account.lep_backend_sa.email}"
}

# IAM bindings for secrets
resource "google_secret_manager_secret_iam_member" "sa_access_jwt_private" {
  secret_id = data.google_secret_manager_secret.jwt_private_key.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_service_account.lep_backend_sa.email}"
}

resource "google_secret_manager_secret_iam_member" "sa_access_jwt_public" {
  secret_id = data.google_secret_manager_secret.jwt_public_key.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_service_account.lep_backend_sa.email}"
}

resource "google_secret_manager_secret_iam_member" "sa_access_db_password" {
  secret_id = data.google_secret_manager_secret.db_password.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_service_account.lep_backend_sa.email}"
}

# Cloud Run Service
resource "google_cloud_run_v2_service" "lep_backend" {
  name         = "${var.project_name}-backend-${var.environment}"
  location     = var.region
  ingress      = "INGRESS_TRAFFIC_ALL"
  launch_stage = "GA"

  template {
    service_account = data.google_service_account.lep_backend_sa.email

    scaling {
      min_instance_count = var.min_instances
      max_instance_count = var.max_instances
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [data.google_sql_database_instance.lep_postgres.connection_name]
      }
    }

    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/${data.google_artifact_registry_repository.lep_repo.repository_id}/lep-backend:latest"

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

      # JWT configuration
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