# 🚀 LEP System - Backend API

Complete REST API for LEP restaurant management system built with Go, featuring multi-tenant architecture, JWT authentication, and comprehensive restaurant operations management.

## 📊 What This Backend Does

- ✅ **Authentication & Authorization** - JWT tokens with whitelist/blacklist system
- ✅ **Multi-tenant Isolation** - Organization and Project isolation with security validation
- ✅ **User Management** - CRUD operations for users with role-based access control
- ✅ **Restaurant Operations** - Tables, customers, reservations, waitlist management
- ✅ **Product Management** - Menu items with categories, pricing, and image uploads
- ✅ **Order Management** - Complete order lifecycle with prep time and kitchen queue
- ✅ **Notification System** - SMS, WhatsApp, and Email with templates and webhooks
- ✅ **Analytics & Reports** - Occupancy, reservation, waitlist, and lead reports
- ✅ **Image Management** - Deduplication system with automatic reference tracking
- ✅ **Cron Jobs** - Scheduled confirmations and event processing
- ✅ **Cloud Ready** - GCP deployment with Terraform, Cloud Run, and Cloud SQL

## 🛠️ Technology Stack

- **Go 1.24.0** - Programming language
- **Gin Web Framework** - HTTP routing and middleware
- **GORM** - PostgreSQL ORM
- **JWT** - Token-based authentication
- **Twilio** - SMS/WhatsApp notifications
- **SMTP** - Email notifications
- **GCS** - Image storage (production)
- **Terraform** - Infrastructure as Code for GCP

## 🚀 Quick Start

### Prerequisites

- **Go 1.24.0+**
- **PostgreSQL 15+**
- **Git**

### Installation & Running

```bash
# Clone repository
git clone <repository-url>
cd LEP-Back

# Install dependencies
go mod tidy

# Configure environment
cp .env.example .env
# Edit .env with your database credentials and API keys

# Run migrations (if needed)
go run cmd/migrate/main.go

# Start backend server
go run .
# Server starts on http://localhost:8080
```

### Verify Installation

```bash
# Check server health
curl http://localhost:8080/ping
# Expected: "pong"

curl http://localhost:8080/health
# Expected: {"status": "ok"}
```

## 📁 Project Structure

```
LEP-Back/
├── handler/           # Business logic layer
├── server/            # HTTP controllers
├── repositories/      # Data access layer
│   └── models/       # Database entities
├── routes/           # Route definitions
├── middleware/       # HTTP middleware
├── resource/         # Global resource injection
├── config/           # Configuration
├── utils/            # Utilities
├── cmd/              # Command-line tools
└── main.go          # Entry point
```

## 🔐 Authentication & Multi-tenant

### Headers (Required for Protected Routes)

```bash
X-Lpe-Organization-Id: <organization-uuid>
X-Lpe-Project-Id: <project-uuid>
Authorization: Bearer <jwt-token>
```

### Public Routes (No Headers Required)

```bash
POST   /login              # User login
POST   /user               # Create user
GET    /ping               # Health check
GET    /health             # Health status
POST   /webhook/*          # Webhook endpoints
```

### Protected Routes (Headers Required)

All other routes require:
1. Valid JWT token in `Authorization` header
2. Valid organization ID in `X-Lpe-Organization-Id`
3. Valid project ID in `X-Lpe-Project-Id`
4. User must have access to specified organization/project

## 📡 Main API Endpoints

### Authentication
```bash
POST   /login              # Login (email + password)
POST   /logout             # Logout user
POST   /checkToken         # Validate JWT token
```

### Users
```bash
GET    /user/:id                           # Get user
POST   /user                               # Create user (public)
PUT    /user/:id                           # Update user
DELETE /user/:id                           # Soft delete user
GET    /user/:id/organizations-projects    # Get user access (Admin)
POST   /user/:id/organizations-projects    # Update user access (Admin)
```

### Products
```bash
GET    /product/:id           # Get product
POST   /product              # Create product
PUT    /product/:id          # Update product
DELETE /product/:id          # Soft delete product
POST   /upload/product/image # Upload product image
```

### Orders
```bash
GET    /order/:id   # Get order
GET    /orders      # List orders
POST   /order       # Create order
PUT    /order/:id   # Update order
DELETE /order/:id   # Soft delete order
GET    /kitchen/queue  # Kitchen queue
```

### Tables & Reservations
```bash
GET    /table/:id       # Get table
GET    /table           # List tables
POST   /table           # Create table
PUT    /table/:id       # Update table
DELETE /table/:id       # Soft delete table

GET    /reservation/:id # Get reservation
GET    /reservation     # List reservations
POST   /reservation     # Create reservation
PUT    /reservation/:id # Update reservation
DELETE /reservation/:id # Cancel reservation
```

### Waitlist & Customers
```bash
GET    /waitlist/:id    # Get waitlist entry
GET    /waitlist        # List waitlist
POST   /waitlist        # Create waitlist entry
PUT    /waitlist/:id    # Update waitlist entry
DELETE /waitlist/:id    # Remove from waitlist

GET    /customer/:id    # Get customer
GET    /customer        # List customers
POST   /customer        # Create customer
PUT    /customer/:id    # Update customer
DELETE /customer/:id    # Soft delete customer
```

### Notifications
```bash
POST   /notification/config              # Create/update notification config
GET    /notification/config/:event       # Get config for event
POST   /notification/template            # Create notification template
PUT    /notification/template/:id        # Update template
GET    /notification/templates           # List templates
POST   /notification/send               # Send manual notification
GET    /notification/logs               # List notification logs
POST   /notification/webhook/twilio/status    # Twilio status webhook
POST   /notification/webhook/twilio/inbound   # Twilio inbound webhook
```

### Reports
```bash
GET    /reports/occupancy        # Table occupancy report
GET    /reports/reservations     # Reservation statistics
GET    /reports/waitlist         # Waitlist metrics
GET    /reports/leads            # Customer leads
GET    /reports/export/csv       # Export as CSV
```

## ⚙️ Environment Variables

```bash
# Database
DB_USER=postgres
DB_PASS=password
DB_NAME=lep_database
INSTANCE_UNIX_SOCKET=/path/to/socket  # For GCP Cloud SQL

# Authentication
JWT_SECRET_PRIVATE_KEY=your_private_key
JWT_SECRET_PUBLIC_KEY=your_public_key

# Twilio (SMS/WhatsApp)
TWILIO_ACCOUNT_SID=your_sid
TWILIO_AUTH_TOKEN=your_token
TWILIO_PHONE_NUMBER=+55119999999

# SMTP (Email)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# Optional Features
ENABLE_CRON_JOBS=true
GIN_MODE=debug  # or release
```

## 🐳 Docker & Cloud Deployment

### Build Docker Image

```bash
docker build -t lep-backend:stage .
docker tag lep-backend:stage gcr.io/your-project/lep-backend:stage
docker push gcr.io/your-project/lep-backend:stage
```

### Deploy to Google Cloud Run

```bash
gcloud run deploy lep-system \
  --image=gcr.io/your-project/lep-backend:stage \
  --region=us-central1 \
  --allow-unauthenticated \
  --set-env-vars ENVIRONMENT=staging,PORT=8080
```

### Terraform Infrastructure

```bash
# Initialize and deploy
terraform init
terraform plan
terraform apply
```

This creates:
- Cloud SQL PostgreSQL database
- Cloud Storage for images
- Cloud Run service
- IAM roles and service accounts

## 🌱 Database Seeding

### Standard Seed

```bash
bash ./scripts/run_seed.sh
# Creates: 1 organization, 6 users, 12 products, 8 tables, 4 orders, etc.
```

### Seeding via API (Remote)

```bash
# For staging/production without database access
bash ./scripts/run_seed_remote.sh --url https://your-api.com --verbose
```

## 📋 Master Admin Access

To grant Master Admin permissions (for testing or administration):

```bash
# Set user as Master Admin
go run cmd/create-master-admins/main.go

# Grants:
# - master_admin permission
# - Access to all organizations
# - Access to all projects
# - Ability to manage user permissions
```

## 📝 Build & Testing

```bash
# Build binary
go build -o lep-system .

# Run tests
bash scripts/run_tests.sh

# Format code
go fmt ./...

# Check code quality
go vet ./...
```

## 📖 Additional Documentation

- **[LATEST_UPDATES.txt](LATEST_UPDATES.txt)** - Recent changes and security fixes
- **[CLAUDE.md](CLAUDE.md)** - Development guidelines
- **[CHANGELOG.md](CHANGELOG.md)** - Complete changelog
- **[STRUCTURE.md](STRUCTURE.md)** - Detailed project structure
- **[routes/routes.md](routes/routes.md)** - Complete API routes documentation

## ✅ Build Status

- ✅ Backend: Compiles successfully (0 errors)
- ✅ Routes: 150+ endpoints registered
- ✅ Security: Multi-tenant validation active
- ✅ Tests: Ready for execution
- ✅ Production: Ready for deployment

---

**Version**: 1.0
**Status**: ✅ Production Ready
**Last Updated**: 2025-11-09
