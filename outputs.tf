# Cloud Run service outputs
output "service_url" {
  description = "URL of the Cloud Run service"
  value       = google_cloud_run_v2_service.lep_backend.uri
}

output "service_name" {
  description = "Name of the Cloud Run service"
  value       = google_cloud_run_v2_service.lep_backend.name
}

# Database outputs (hardcoded from bootstrap)
output "database_connection_name" {
  description = "Cloud SQL instance connection name"
  value       = "${var.project_id}:${var.region}:leps-postgres-${var.environment}"
}

output "database_instance_name" {
  description = "Cloud SQL instance name"
  value       = "leps-postgres-${var.environment}"
}

output "database_name" {
  description = "Database name"
  value       = "lep_database"
}

output "database_user" {
  description = "Database user"
  value       = "lep_user"
  sensitive   = true
}

# Service Account outputs (hardcoded from bootstrap)
output "service_account_email" {
  description = "Email of the service account"
  value       = "lep-backend-sa@${var.project_id}.iam.gserviceaccount.com"
}

# Artifact Registry outputs (hardcoded from bootstrap)
output "docker_repository_url" {
  description = "Docker repository URL for pushing images"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/lep-backend"
}

output "docker_image_url" {
  description = "Full Docker image URL"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/lep-backend/lep-backend:latest"
}

# Secret Manager outputs (hardcoded from bootstrap)
output "secrets_list" {
  description = "List of created secrets"
  value = {
    jwt_private_key = "jwt-private-key-${var.environment}"
    jwt_public_key  = "jwt-public-key-${var.environment}"
    db_password     = "db-password-${var.environment}"
  }
}

# Project information
output "project_id" {
  description = "GCP project ID"
  value       = var.project_id
}

output "region" {
  description = "GCP region"
  value       = var.region
}

output "environment" {
  description = "Environment name"
  value       = var.environment
}

# Deployment commands
output "docker_build_command" {
  description = "Command to build and push Docker image"
  value = "docker build -t ${var.region}-docker.pkg.dev/${var.project_id}/lep-backend/lep-backend:latest . && docker push ${var.region}-docker.pkg.dev/${var.project_id}/lep-backend/lep-backend:latest"
}

output "cloud_run_deploy_command" {
  description = "Command to deploy to Cloud Run"
  value = "gcloud run deploy ${google_cloud_run_v2_service.lep_backend.name} --image=${var.region}-docker.pkg.dev/${var.project_id}/lep-backend/lep-backend:latest --region=${var.region} --platform=managed"
}