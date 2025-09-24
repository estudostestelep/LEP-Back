# LEP System - Terraform Minimal Configuration
# Only manages Cloud Run service with hardcoded references

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

# Generate random password for database
resource "random_password" "db_password" {
  length  = 32
  special = true
}

# Cloud Run service - hardcoded references to avoid permission issues
resource "google_cloud_run_v2_service" "lep_backend" {
  name         = "${var.project_name}-backend-${var.environment}"
  location     = var.region
  ingress      = "INGRESS_TRAFFIC_ALL"
  launch_stage = "GA"

  template {
    # Hardcoded service account created by bootstrap
    service_account = "lep-backend-sa@${var.project_id}.iam.gserviceaccount.com"

    scaling {
      min_instance_count = 0
      max_instance_count = 2
    }

    containers {
      # Hardcoded image path from bootstrap artifact registry
      image = "${var.region}-docker.pkg.dev/${var.project_id}/lep-backend/lep-backend:latest"

      ports {
        container_port = 8080
      }

      env {
        name  = "DB_USER"
        value = "lep_user"
      }

      env {
        name = "DB_PASS"
        value_source {
          secret_key_ref {
            secret  = "db-password-${var.environment}"
            version = "latest"
          }
        }
      }

      env {
        name  = "DB_NAME"
        value = "lep_database"
      }

      env {
        name  = "INSTANCE_UNIX_SOCKET"
        value = "/cloudsql/${var.project_id}:${var.region}:leps-postgres-${var.environment}"
      }

      env {
        name = "JWT_SECRET_PRIVATE_KEY"
        value_source {
          secret_key_ref {
            secret  = "jwt-private-key-${var.environment}"
            version = "latest"
          }
        }
      }

      env {
        name = "JWT_SECRET_PUBLIC_KEY"
        value_source {
          secret_key_ref {
            secret  = "jwt-public-key-${var.environment}"
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

    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = ["${var.project_id}:${var.region}:leps-postgres-${var.environment}"]
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
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