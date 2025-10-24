# LEP System - Changelog

## [Unreleased]

### 🚀 **Updated**
- **Go version**: Upgraded from Go 1.21.5 to Go 1.24 across all Docker images and configuration files
  - Updated `Dockerfile` (production build)
  - Updated `Dockerfile.dev` (development build)
  - Updated `Dockerfile.prod` (production optimized build)
  - Updated `go.mod` module version
  - Updated Cloud Build configuration (`cloudbuild.yaml`)
  - Updated documentation references

### 🔧 **Fixed**
- **Docker Compose environment**: Resolved JWT key parsing error in local development
  - Fixed multi-line environment variable format incompatibility
  - Updated `environments/local-dev.env` to use simple keys for development
  - Modified `docker-compose.yml` to define environment variables directly
  - Improved deployment script error handling

### 📚 **Added**
- **Interactive deployment script**: Cross-platform support for multiple environments
  - Linux/Mac: `./scripts/deploy-interactive.sh`
  - Windows: `./scripts/Deploy-Interactive.ps1`
  - Support for 4 environments: local-dev, gcp-dev, gcp-stage, gcp-prd
- **Comprehensive troubleshooting guide**: `TROUBLESHOOTING.md`
- **Deployment documentation**: Enhanced `DEPLOYMENT-GUIDE.md`

### 🏗️ **Infrastructure**
- **Environment-specific configurations**:
  - `environments/local-dev.env` - Local development setup
  - `environments/gcp-dev.tfvars` - GCP minimal setup
  - `environments/gcp-stage.tfvars` - GCP staging environment
  - `environments/gcp-prd.tfvars` - GCP production environment
- **Docker Compose stack**: Complete local development environment
  - PostgreSQL database with initialization scripts
  - Redis for caching
  - MailHog for email testing
  - PgAdmin for database management

### 🛡️ **Security**
- **Container security**: Multi-stage builds with minimal attack surface
- **Non-root user**: All containers run with dedicated `lepuser`
- **Secret management**: Enhanced JWT key handling for different environments

---

## Previous Changes

For earlier changes, please refer to git commit history.

---

## Migration Notes

### **From Go 1.21.5 to 1.22.5**

This update includes:
- **Performance improvements** in Go 1.22.5
- **Security patches** and bug fixes
- **Better module support** and dependency resolution

**Action Required:**
- No breaking changes for existing LEP System functionality
- Recommend rebuilding Docker images: `docker-compose build --no-cache`
- If developing locally, ensure Go 1.22+ is installed

**Compatibility:**
- ✅ All existing APIs remain compatible
- ✅ Database schema unchanged
- ✅ Environment variables unchanged
- ✅ Deployment procedures unchanged