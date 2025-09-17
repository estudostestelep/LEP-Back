terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.26.0"
    }
  }
}

provider "google" {
  project = "market-split"
  region  = "us-central1"
}

# Enable APIs
resource "google_project_service" "secretmanager_api" {
  project             = "market-split"
  service             = "secretmanager.googleapis.com"
  disable_on_destroy  = false
}

resource "google_project_service" "sqladmin_api" {
  project             = "market-split"
  service             = "sqladmin.googleapis.com"
  disable_on_destroy  = false
}

resource "google_project_service" "cloudrun_api" {
  project             = "market-split"
  service             = "run.googleapis.com"
  disable_on_destroy  = false
}

# Create Service Account
resource "google_service_account" "cloud_run_service_account" {
  account_id   = "cloud-run-service-account"
  display_name = "Cloud Run Service Account"
}

# Assign Secret Manager Access 
resource "google_project_iam_member" "service_account_secret_accessor" {
  project = "market-split"
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.cloud_run_service_account.email}"
  depends_on = [null_resource.wait_for_apis]
}

# Assign Cloud SQL Client Role to Service Account
resource "google_project_iam_member" "cloud_run_sql_access" {
  project = "market-split"
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.cloud_run_service_account.email}"
  depends_on = [null_resource.wait_for_apis]
}

# Ensure APIs are enabled before proceeding
resource "null_resource" "wait_for_apis" {
  provisioner "local-exec" {
    command = <<EOT
      gcloud services enable run.googleapis.com secretmanager.googleapis.com sqladmin.googleapis.com --project=market-split
      while ! gcloud services list --enabled --filter="run.googleapis.com" --project=market-split | grep -q run.googleapis.com; do
        echo "Waiting for Cloud Run API to be enabled..."
        sleep 30
      done

      while ! gcloud services list --enabled --filter="secretmanager.googleapis.com" --project=market-split | grep -q secretmanager.googleapis.com; do
        echo "Waiting for Secret Manager API to be enabled..."
        sleep 30
      done

      while ! gcloud services list --enabled --filter="sqladmin.googleapis.com" --project=market-split | grep -q sqladmin.googleapis.com; do
        echo "Waiting for Cloud SQL Admin API to be enabled..."
        sleep 30
      done
    EOT
  }

  depends_on = [
    google_project_service.secretmanager_api,
    google_project_service.sqladmin_api,
    google_project_service.cloudrun_api
  ]
}

# Cloud SQL Instance
resource "google_sql_database_instance" "postgres2" {
  name             = "postgres2"
  region           = "us-central1"
  database_version = "POSTGRES_15"
  root_password    = "12345"

  settings {
    tier = "db-f1-micro"
  }

  deletion_protection = false
  depends_on = [null_resource.wait_for_apis]
}

# Secret Manager Secrets
resource "google_secret_manager_secret" "JWT_SECRET_PRIVATE_KEY" {
  secret_id = "JWT_SECRET_PRIVATE_KEY"
  replication {
    auto {}
  }
  lifecycle {
    ignore_changes = [secret_id]
  }
}

resource "google_secret_manager_secret_version" "secret_private_key" {
  secret      = google_secret_manager_secret.JWT_SECRET_PRIVATE_KEY.id
  secret_data =  "JWT_SECRET_PRIVATE_KEY real key"
  }

resource "google_secret_manager_secret" "JWT_SECRET_PUBLIC_KEY" {
  secret_id = "JWT_SECRET_PUBLIC_KEY"
  replication {
    auto {}
  }
  lifecycle {
    ignore_changes = [secret_id]
  }
}

resource "google_secret_manager_secret_version" "secret_public_key" {
  secret      = google_secret_manager_secret.JWT_SECRET_PUBLIC_KEY.id
  secret_data = "JWT_SECRET_PUBLIC_KEY real key"
}

resource "google_secret_manager_secret_iam_member" "secret_access_private_key" {
  secret_id = google_secret_manager_secret.JWT_SECRET_PRIVATE_KEY.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_run_service_account.email}"
}

resource "google_secret_manager_secret_iam_member" "secret_access_public_key" {
  secret_id = google_secret_manager_secret.JWT_SECRET_PUBLIC_KEY.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_run_service_account.email}"
}

# Cloud Run Service
resource "google_cloud_run_v2_service" "my_service1" {
  name        = "cloudrun-service"
  location    = "us-central1"
  ingress     = "INGRESS_TRAFFIC_ALL"
  launch_stage = "BETA"

  template {
    service_account = google_service_account.cloud_run_service_account.email

    scaling {
      max_instance_count = 2
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.postgres2.connection_name]
      }
    }

    containers {
      image = "gcr.io/market-split/test-2"

      env {
        name  = "DB_USER"
        value = "real value postgres"
      }
      env {
        name  = "DB_PASS"
        value = "real value 12345"
      }
      env {
        name  = "DB_NAME"
        value = "real value postgres"
      }
      env {
        name  = "INSTANCE_UNIX_SOCKET"
        value = "real value '/cloudsql/project:region:instance'"
      }
      env {
        name = "JWT_SECRET_PRIVATE_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.JWT_SECRET_PRIVATE_KEY.secret_id
            version = "latest"
          }
        }
      }
      env {
        name = "JWT_SECRET_PUBLIC_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.JWT_SECRET_PUBLIC_KEY.secret_id
            version = "latest"
          }
        }
      }

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }

 depends_on = [
    google_secret_manager_secret_version.secret_private_key,
    google_secret_manager_secret_version.secret_public_key,
    google_secret_manager_secret_iam_member.secret_access_private_key,
    google_secret_manager_secret_iam_member.secret_access_public_key
  ]
}

# Cloud Run IAM Policy
resource "google_cloud_run_service_iam_member" "cloud_run_sql_iam_policy" {
  service  = google_cloud_run_v2_service.my_service1.name
  location = google_cloud_run_v2_service.my_service1.location
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.cloud_run_service_account.email}"
}
