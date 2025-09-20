# LEP System - Scripts Directory

This directory contains automation scripts for the LEP System deployment and development.

## Available Scripts

### 🚀 Production Deployment

#### `setup.sh`
Complete environment setup script that prepares your system for development and deployment.
```bash
./scripts/setup.sh
```

#### `terraform-setup.sh`
Initializes and applies Terraform infrastructure configuration.
```bash
./scripts/terraform-setup.sh
```

#### `build-and-deploy.sh`
Builds Docker image and deploys to Google Cloud Run.
```bash
./scripts/build-and-deploy.sh
```

### 🛠️ Development Tools

#### `local-dev.sh`
Local development helper with multiple commands.
```bash
# Start the application
./scripts/local-dev.sh run

# Build the application
./scripts/local-dev.sh build

# Run tests
./scripts/local-dev.sh test

# Build and run in Docker
./scripts/local-dev.sh docker

# Generate JWT keys
./scripts/local-dev.sh generate

# Clean build artifacts
./scripts/local-dev.sh clean

# Health check
./scripts/local-dev.sh health
```

## Quick Start

1. **First-time setup:**
   ```bash
   ./scripts/setup.sh
   ```

2. **Deploy infrastructure:**
   ```bash
   ./scripts/terraform-setup.sh
   ```

3. **Deploy application:**
   ```bash
   ./scripts/build-and-deploy.sh
   ```

## Script Dependencies

Ensure you have the following tools installed:
- **gcloud CLI** - Google Cloud command-line tool
- **terraform** - Infrastructure as Code tool
- **docker** - Container platform
- **go** - Go programming language
- **git** - Version control system

## Configuration Files

Make sure these files are properly configured:
- `terraform.tfvars` - Terraform variables
- `.env` - Environment variables for local development

## Security Notes

- Scripts automatically configure `.gitignore` to exclude sensitive files
- JWT keys and terraform.tfvars are not committed to version control
- Use encrypted PEM files for JWT keys
- All scripts include proper error handling and validation

## Troubleshooting

### Authentication Issues
```bash
gcloud auth login
gcloud auth application-default login
gcloud auth configure-docker us-central1-docker.pkg.dev
```

### Permission Issues
```bash
chmod +x scripts/*.sh
```

### Build Issues
```bash
go mod tidy
go mod verify
```

For more detailed information, check the individual script files or the main project documentation.