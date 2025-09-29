# Project configuration
variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "leps"
}

variable "region" {
  description = "GCP region for resources"
  type        = string
  default     = "us-central1"
}

variable "environment" {
  description = "Environment name (dev, stage, prod)"
  type        = string
  default     = "dev"

  validation {
    condition     = contains(["dev", "stage", "prod"], var.environment)
    error_message = "Environment must be one of: dev, stage, prod."
  }
}

# Database configuration
variable "database_name" {
  description = "Name of the PostgreSQL database"
  type        = string
  default     = "lep_database"
}

variable "database_user" {
  description = "Database user name"
  type        = string
  default     = "lep_user"
}

variable "db_tier" {
  description = "Database instance tier"
  type        = string
  default     = "db-f1-micro"

  validation {
    condition = contains([
      "db-f1-micro", "db-g1-small", "db-n1-standard-1", "db-n1-standard-2",
      "db-n1-standard-4", "db-n1-standard-8", "db-n1-standard-16"
    ], var.db_tier)
    error_message = "Database tier must be a valid Cloud SQL tier."
  }
}

variable "db_availability_type" {
  description = "Database availability type"
  type        = string
  default     = "ZONAL"

  validation {
    condition     = contains(["ZONAL", "REGIONAL"], var.db_availability_type)
    error_message = "Availability type must be ZONAL or REGIONAL."
  }
}

variable "db_disk_size" {
  description = "Database disk size in GB"
  type        = number
  default     = 20

  validation {
    condition     = var.db_disk_size >= 10 && var.db_disk_size <= 30720
    error_message = "Database disk size must be between 10 and 30720 GB."
  }
}

variable "enable_deletion_protection" {
  description = "Enable deletion protection for database"
  type        = bool
  default     = true
}

# JWT configuration
variable "jwt_private_key" {
  description = "JWT private key for token signing"
  type        = string
  sensitive   = true
}

variable "jwt_public_key" {
  description = "JWT public key for token verification"
  type        = string
  sensitive   = true
}

# Cloud Run configuration
variable "min_instances" {
  description = "Minimum number of Cloud Run instances"
  type        = number
  default     = 0

  validation {
    condition     = var.min_instances >= 0 && var.min_instances <= 1000
    error_message = "Min instances must be between 0 and 1000."
  }
}

variable "max_instances" {
  description = "Maximum number of Cloud Run instances"
  type        = number
  default     = 100

  validation {
    condition     = var.max_instances >= 1 && var.max_instances <= 1000
    error_message = "Max instances must be between 1 and 1000."
  }
}

variable "cpu_limit" {
  description = "CPU limit for Cloud Run container"
  type        = string
  default     = "1"

  validation {
    condition = contains([
      "0.08", "0.17", "0.25", "0.33", "0.5", "0.58", "0.67", "0.75", "0.83", "1", "2", "4", "6", "8"
    ], var.cpu_limit)
    error_message = "CPU limit must be a valid Cloud Run CPU allocation."
  }
}

variable "memory_limit" {
  description = "Memory limit for Cloud Run container"
  type        = string
  default     = "512Mi"

  validation {
    condition = can(regex("^[0-9]+(Mi|Gi)$", var.memory_limit))
    error_message = "Memory limit must be in format like '512Mi' or '2Gi'."
  }
}

# Twilio configuration (optional)
variable "twilio_account_sid" {
  description = "Twilio Account SID for SMS/WhatsApp"
  type        = string
  default     = ""
  sensitive   = true
}

variable "twilio_auth_token" {
  description = "Twilio Auth Token"
  type        = string
  default     = ""
  sensitive   = true
}

variable "twilio_phone_number" {
  description = "Twilio phone number"
  type        = string
  default     = ""
  sensitive   = true
}

# SMTP configuration (optional)
variable "smtp_host" {
  description = "SMTP server host"
  type        = string
  default     = "smtp.gmail.com"
}

variable "smtp_port" {
  description = "SMTP server port"
  type        = number
  default     = 587

  validation {
    condition     = var.smtp_port > 0 && var.smtp_port <= 65535
    error_message = "SMTP port must be between 1 and 65535."
  }
}

variable "smtp_username" {
  description = "SMTP username"
  type        = string
  default     = ""
  sensitive   = true
}

variable "smtp_password" {
  description = "SMTP password"
  type        = string
  default     = ""
  sensitive   = true
}

# Application configuration
variable "enable_cron_jobs" {
  description = "Enable cron jobs for notifications"
  type        = bool
  default     = true
}

# Domain configuration (optional)
variable "domain_name" {
  description = "Custom domain name for the service"
  type        = string
  default     = ""
}

variable "enable_custom_domain" {
  description = "Enable custom domain mapping"
  type        = bool
  default     = false
}

# Bucket configuration for image storage
variable "bucket_name" {
  description = "GCS bucket name for image storage"
  type        = string
  default     = ""

  validation {
    condition = can(regex("^[a-z0-9._-]*$", var.bucket_name)) || var.bucket_name == ""
    error_message = "Bucket name must contain only lowercase letters, numbers, hyphens, underscores, and periods."
  }
}

variable "bucket_cache_control" {
  description = "Cache control header for uploaded files"
  type        = string
  default     = "public, max-age=3600"

  validation {
    condition = can(regex("^(public|private)(, max-age=[0-9]+)?(, no-cache|, no-store|, must-revalidate)*$", var.bucket_cache_control))
    error_message = "Cache control must be a valid HTTP cache control directive."
  }
}

variable "bucket_timeout" {
  description = "Timeout in seconds for GCS operations"
  type        = number
  default     = 30

  validation {
    condition     = var.bucket_timeout >= 5 && var.bucket_timeout <= 300
    error_message = "Bucket timeout must be between 5 and 300 seconds."
  }
}