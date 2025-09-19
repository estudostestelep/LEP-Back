variable "project_id" {
  description = "The GCP project ID"
  type        = string
  default     = "market-split"
}

variable "region" {
  description = "The GCP region for resources"
  type        = string
  default     = "us-central1"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "app_name" {
  description = "Application name"
  type        = string
  default     = "lep-system"
}

variable "database_tier" {
  description = "Cloud SQL instance tier"
  type        = string
  default     = "db-f1-micro"
}

variable "database_version" {
  description = "PostgreSQL version"
  type        = string
  default     = "POSTGRES_15"
}

variable "max_instances" {
  description = "Maximum number of Cloud Run instances"
  type        = number
  default     = 2
}

variable "container_image" {
  description = "Container image for Cloud Run"
  type        = string
  default     = "gcr.io/market-split/lep-system"
}

variable "database_user" {
  description = "Database username"
  type        = string
  default     = "lep_user"
}

variable "database_name" {
  description = "Database name"
  type        = string
  default     = "lep_db"
}

variable "jwt_secret_private_key" {
  description = "JWT private key for token signing"
  type        = string
  sensitive   = true
}

variable "jwt_secret_public_key" {
  description = "JWT public key for token verification"
  type        = string
  sensitive   = true
}

variable "database_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "deletion_protection" {
  description = "Enable deletion protection for database"
  type        = bool
  default     = true
}

variable "enable_apis" {
  description = "List of APIs to enable"
  type        = list(string)
  default = [
    "secretmanager.googleapis.com",
    "sqladmin.googleapis.com",
    "run.googleapis.com",
    "cloudbuild.googleapis.com"
  ]
}