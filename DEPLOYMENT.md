# LEP System - Deployment Guide

## 🎯 Project Overview

LEP System is now fully configured for deployment to Google Cloud Platform with the following setup:

- **Project ID**: `leps-472702`
- **Project Name**: `leps`
- **Region**: `us-central1`
- **Environment**: `dev` (configurable)

## 📁 Current Configuration Status

### ✅ Completed Configurations

1. **Terraform Infrastructure** - Complete GCP setup
   - Cloud Run service
   - Cloud SQL PostgreSQL database
   - Artifact Registry for Docker images
   - Secret Manager for sensitive data
   - Service accounts with proper IAM

2. **JWT Authentication** - Ready for use
   - ✅ Private key: `jwt_private_key.pem` (encrypted)
   - ✅ Public key: `jwt_public_key.pem`
   - ✅ Keys configured in `terraform.tfvars`

3. **Docker Configuration** - Production-ready
   - ✅ Main `Dockerfile` (production-optimized, multi-stage, scratch-based)
   - ✅ `Dockerfile.dev` (development with hot reloading)
   - ✅ `Dockerfile.prod` (alternative production build)

4. **Security Setup** - Enhanced protection
   - ✅ Comprehensive `.gitignore` for sensitive files
   - ✅ Non-root user in containers
   - ✅ Minimal attack surface (scratch-based images)
   - ✅ Encrypted JWT keys

5. **Automation Scripts** - Ready to use
   - ✅ `scripts/setup.sh` - Complete environment setup
   - ✅ `scripts/terraform-setup.sh` - Infrastructure deployment
   - ✅ `scripts/build-and-deploy.sh` - Application deployment
   - ✅ `scripts/local-dev.sh` - Development tools

## 🚀 Deployment Process

### Step 1: Environment Setup
```bash
./scripts/setup.sh
```
This script will:
- Verify all required tools are installed
- Configure Google Cloud authentication
- Set up project structure
- Validate Go dependencies
- Configure Docker for GCP

### Step 2: Deploy Infrastructure
```bash
./scripts/terraform-setup.sh
```
This script will:
- Enable required GCP APIs
- Initialize Terraform
- Create all infrastructure resources
- Store secrets in Secret Manager

### Step 3: Deploy Application
```bash
./scripts/build-and-deploy.sh
```
This script will:
- Build optimized Docker image
- Push to Artifact Registry
- Deploy to Cloud Run
- Perform health checks

## 🛠️ Development Workflow

### Local Development
```bash
# Start local development server
./scripts/local-dev.sh run

# Build application
./scripts/local-dev.sh build

# Run tests
./scripts/local-dev.sh test

# Run in Docker
./scripts/local-dev.sh docker
```

### Environment Files

#### `.env` (Local Development)
```bash
DB_USER=postgres
DB_PASS=your_password
DB_NAME=lep_database
JWT_SECRET_PRIVATE_KEY=your_private_key
JWT_SECRET_PUBLIC_KEY=your_public_key
# ... other configuration
```

#### `terraform.tfvars` (Infrastructure)
✅ Already configured with:
- Project ID: `leps-472702`
- JWT keys from PEM files
- Default resource configurations

## 🔐 Security Features

### Container Security
- **Multi-stage builds** reduce image size and attack surface
- **Non-root user** (`lepuser`) for container execution
- **Scratch-based final image** with minimal dependencies
- **Distroless approach** with only necessary runtime components

### Secrets Management
- **JWT keys** stored in Google Secret Manager
- **Database passwords** auto-generated and secured
- **Twilio/SMTP credentials** optionally stored in Secret Manager
- **Environment isolation** between dev/staging/prod

### Access Control
- **Service accounts** with minimal required permissions
- **IAM policies** following least-privilege principle
- **Private database** with no public IP
- **VPC connectivity** for secure communication

## 📊 Infrastructure Components

### Compute
- **Cloud Run**: Serverless container platform
  - Auto-scaling (0-10 instances)
  - CPU: 1 vCPU, Memory: 512Mi
  - Health checks enabled

### Database
- **Cloud SQL PostgreSQL 15**
  - Tier: db-f1-micro (upgradeable)
  - Disk: 20GB SSD
  - Automated backups enabled
  - Point-in-time recovery

### Storage & Registry
- **Artifact Registry**: Docker image storage
- **Secret Manager**: Secure credential storage

### Networking
- **Private IP**: Database not publicly accessible
- **Cloud SQL Proxy**: Secure database connections
- **HTTPS**: Automatic SSL/TLS termination

## 🌍 Environment Management

### Current Setup
- **Environment**: `dev`
- **Deletion Protection**: Disabled (safe for development)
- **High Availability**: ZONAL (cost-optimized)

### Production Recommendations
Update `terraform.tfvars` for production:
```bash
environment = "prod"
enable_deletion_protection = true
db_availability_type = "REGIONAL"
db_tier = "db-n1-standard-1"
min_instances = 1
```

## 📋 Post-Deployment Checklist

### Immediate Tasks
- [ ] Run `./scripts/setup.sh` to verify environment
- [ ] Execute `./scripts/terraform-setup.sh` to deploy infrastructure
- [ ] Deploy application with `./scripts/build-and-deploy.sh`
- [ ] Test health endpoint: `curl https://your-service-url/health`

### Optional Configurations
- [ ] Configure Twilio credentials for SMS/WhatsApp
- [ ] Set up SMTP credentials for email notifications
- [ ] Configure custom domain (if needed)
- [ ] Set up monitoring and alerting

### Production Readiness
- [ ] Update resource limits for expected load
- [ ] Enable regional database for high availability
- [ ] Configure backup retention policies
- [ ] Set up CI/CD pipeline
- [ ] Implement monitoring dashboards

## 🆘 Troubleshooting

### Common Issues

**Authentication Errors**
```bash
gcloud auth login
gcloud auth application-default login
```

**Docker Permission Errors**
```bash
gcloud auth configure-docker us-central1-docker.pkg.dev
```

**Terraform State Issues**
```bash
terraform init -reconfigure
```

**Build Failures**
```bash
go mod tidy
go mod verify
```

### Support Resources
- **Scripts Documentation**: `scripts/README.md`
- **Terraform Outputs**: `terraform output` for deployment info
- **Health Checks**: Use `/health` and `/ping` endpoints
- **Logs**: Check Cloud Run logs in GCP Console

## 🎉 Ready for Deployment!

Your LEP System is now fully configured and ready for deployment. All scripts are executable and all configuration files are properly set up with your GCP project details.

**Next Steps**: Run `./scripts/setup.sh` to begin the deployment process.