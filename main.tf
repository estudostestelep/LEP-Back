# LEP System GCP Infrastructure
# Minimal version - Cloud Run service managed separately

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

# Generate random password for database
resource "random_password" "db_password" {
  length  = 32
  special = true
}

# Cloud Run service is managed separately via gcloud commands
# Service name: leps-backend-dev
# URL: https://leps-backend-dev-516622888070.us-central1.run.app

# Allow unauthenticated access to Cloud Run service
# Note: This will be applied manually since the service is managed outside Terraform
# resource "google_cloud_run_service_iam_binding" "default" {
#   location = "us-central1"
#   service  = "leps-backend-dev"
#   role     = "roles/run.invoker"
#   members = [
#     "allUsers"
#   ]
# }