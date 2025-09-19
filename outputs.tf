# Cloud Run service outputs
output "service_url" {
  description = "URL of the Cloud Run service"
  value       = google_cloud_run_v2_service.lep_backend.uri
}

output "service_name" {
  description = "Name of the Cloud Run service"
  value       = google_cloud_run_v2_service.lep_backend.name
}

# Database outputs
output "database_connection_name" {
  description = "Cloud SQL instance connection name"
  value       = google_sql_database_instance.lep_postgres.connection_name
}

output "database_instance_name" {
  description = "Cloud SQL instance name"
  value       = google_sql_database_instance.lep_postgres.name
}

output "database_name" {
  description = "Database name"
  value       = google_sql_database.lep_database.name
}

output "database_user" {
  description = "Database user"
  value       = google_sql_user.lep_user.name
  sensitive   = true
}

# Service Account outputs
output "service_account_email" {
  description = "Email of the service account"
  value       = google_service_account.lep_backend_sa.email
}

# Artifact Registry outputs
output "docker_repository_url" {
  description = "Docker repository URL for pushing images"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.lep_repo.repository_id}"
}

output "docker_image_url" {
  description = "Full Docker image URL"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.lep_repo.repository_id}/lep-backend:latest"
}

# Secret Manager outputs
output "secrets_list" {
  description = "List of created secrets"
  value = {
    jwt_private_key = google_secret_manager_secret.jwt_private_key.secret_id
    jwt_public_key  = google_secret_manager_secret.jwt_public_key.secret_id
    db_password     = google_secret_manager_secret.db_password.secret_id
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
  value = "docker build -t ${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.lep_repo.repository_id}/lep-backend:latest . && docker push ${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.lep_repo.repository_id}/lep-backend:latest"
}

output "cloud_run_deploy_command" {
  description = "Command to deploy to Cloud Run"
  value = "gcloud run deploy ${google_cloud_run_v2_service.lep_backend.name} --image=${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.lep_repo.repository_id}/lep-backend:latest --region=${var.region} --platform=managed"
}