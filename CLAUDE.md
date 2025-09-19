# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

LEP System is a Go-based REST API for restaurant management, built with clean architecture and modular design. The system handles users, products, orders, tables, reservations, waitlists, and customers.

## Development Commands

### Build and Run
```bash
# Install dependencies
go mod tidy

# Run the application
go run main.go

# Build binary
go build -o lep-system .

# Run built binary
./lep-system
```

### Testing
```bash
# Test basic connectivity
curl http://localhost:8080/ping
# Expected response: "pong"

# Health check
curl http://localhost:8080/health
# Expected response: {"status":"healthy"}
```

### Docker
```bash
# Build Docker image
docker build -t lep-system .

# Run containerized
docker run -p 8080:8080 lep-system
```

## Architecture

### Clean Architecture Pattern
The codebase follows a 3-layer clean architecture:

1. **Handler Layer** (`handler/`) - Business logic and validation
2. **Server Layer** (`server/`) - HTTP controllers and request/response handling
3. **Repository Layer** (`repositories/`) - Database access via GORM

### Key Architectural Concepts

- **Dependency Injection**: All dependencies are injected via `resource/inject.go`
- **Organization/Project Isolation**: All entities require `organization_id` and `project_id` headers
- **Soft Delete**: All entities use `deleted_at` for logical deletion
- **JWT Authentication**: Token-based auth with blacklist/whitelist management
- **Audit Logging**: All operations are logged via `AuditLog` model
- **Event-Driven Notifications**: Automated notification system with support for SMS, WhatsApp, and Email
- **Cron Jobs**: Background tasks for 24h confirmation reminders and event processing
- **Template System**: Dynamic notification templates with variable substitution

### Directory Structure
- `config/` - Database and environment configuration
- `handler/` - Business logic interfaces and implementations
- `middleware/` - Authentication middleware
- `repositories/` - Database models and CRUD operations
- `resource/` - Global resource management and dependency injection
- `routes/` - Route organization and setup
- `server/` - HTTP controllers
- `utils/` - Shared utilities

## Database Models

Core entities in `repositories/models/PostgresLEP.go`:
- `User` - Employees/admins with roles and permissions
- `Customer` - Restaurant customers
- `Table` - Restaurant tables with capacity and availability
- `Product` - Menu items with pricing
- `Order` - Customer orders with items and status
- `Reservation` - Table reservations with datetime
- `Waitlist` - Queue management for occupied tables
- `AuditLog` - Operation tracking and audit trail
- `NotificationConfig` - Event notification configuration per project
- `NotificationTemplate` - Message templates with variables
- `NotificationLog` - Delivery tracking and status logging
- `NotificationEvent` - Event queue for async processing
- `NotificationInbound` - Inbound message handling

## Authentication & Authorization

### Required Headers (all endpoints except `/login` and `POST /user`)
```
X-Lpe-Organization-Id: <organization-uuid>
X-Lpe-Project-Id: <project-uuid>
Authorization: Bearer <jwt-token>
```

### Token Management
- JWT tokens use HS256 signing (24h expiration)
- Active tokens stored in `LoggedLists`
- Revoked tokens stored in `BannedLists`
- Passwords hashed with bcrypt

## Environment Variables

Required environment variables:
```bash
# Database
DB_USER=postgres_username
DB_PASS=postgres_password
DB_NAME=database_name
INSTANCE_UNIX_SOCKET=/path/to/socket  # For GCP Cloud SQL

# Authentication
JWT_SECRET_PRIVATE_KEY=jwt_private_key
JWT_SECRET_PUBLIC_KEY=jwt_public_key

# Twilio (SMS/WhatsApp)
TWILIO_ACCOUNT_SID=your_account_sid
TWILIO_AUTH_TOKEN=your_auth_token
TWILIO_PHONE_NUMBER=+1234567890

# SMTP (Email)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# Cron Jobs (optional)
ENABLE_CRON_JOBS=true
```

## API Endpoints

### Public Routes
- `POST /login` - User authentication
- `POST /user` - Create user (public signup)

### Protected Routes
All other routes require authentication headers:

- **Auth**: `/logout`, `/checkToken`
- **Users**: `/user/:id`, `/user/group/:id` (GET/PUT/DELETE)
- **Products**: `/product/*` (Full CRUD)
- **Tables**: `/table/*` (Full CRUD + list)
- **Orders**: `/order/*` (Full CRUD) + `/orders` (list)
- **Reservations**: `/reservation/*` (Full CRUD + list)
- **Waitlist**: `/waitlist/*` (Full CRUD + list)
- **Customers**: `/customer/*` (Full CRUD + list)
- **Notifications**: `/notification/*` (Config, templates, logs, webhooks)
- **Reports**: `/reports/*` (Occupancy, reservations, waitlist, lead reports)

## Infrastructure

### Local Development
The application runs on port 8080 by default with CORS enabled for all origins.

### GCP Deployment
Terraform configuration in `example.main.tf` provisions:
- Cloud Run for application hosting
- Cloud SQL (PostgreSQL) for database
- Secret Manager for JWT keys
- Automatic scaling up to 2 instances

## Code Conventions

- Use GORM for all database operations
- Follow clean architecture separation of concerns
- All entities must include `organization_id` and `project_id`
- Use UUID for all primary keys
- Implement soft delete pattern with `deleted_at`
- Log all operations via `AuditLog`
- Validate headers in server layer before processing
- Hash passwords with bcrypt in handler layer

## Notification System

### Event-Driven Architecture
The notification system is triggered by business events:

- **Reservation Events**: Creation, update, cancellation
- **Table Events**: Available tables for waitlist
- **Scheduled Events**: 24h confirmation reminders via cron jobs

### Supported Channels
- **SMS**: Via Twilio API
- **WhatsApp**: Via Twilio Business API
- **Email**: Via SMTP with configurable providers

### Key Components

#### EventService (`utils/event_service.go`)
Handles event creation and notification processing:
```go
// Trigger notification for reservation creation
func (e *EventService) TriggerReservationCreated(orgId, projectId uuid.UUID,
    reservation *models.Reservation, customer *models.Customer, table *models.Table) error
```

#### CronService (`utils/cron_service.go`)
Background jobs for automated notifications:
- **24h Confirmation**: Runs hourly, sends confirmation requests 24h before reservations
- **Pending Events**: Processes event queue every 5 minutes
- **Cleanup**: Daily log cleanup for old notifications

#### NotificationService (`utils/notification_service.go`)
Handles external API integration:
```go
type NotificationService struct {
    twilioClient *twilio.RestClient
    smtpConfig   SMTPConfig
}
```

### Template Variables
Dynamic variables available in templates:
- `{{nome}}` / `{{cliente}}` - Customer name
- `{{mesa}}` / `{{numero_mesa}}` - Table number
- `{{data}}` - Date (DD/MM/YYYY)
- `{{hora}}` - Time (HH:MM)
- `{{data_hora}}` - Full datetime
- `{{pessoas}}` - Party size
- `{{tempo_espera}}` - Estimated wait time
- `{{status}}` - Reservation status

### Webhook Integration
Support for bidirectional communication:
- **Inbound Messages**: Process customer replies via webhooks
- **Status Updates**: Track message delivery and read receipts
- **Event Triggers**: External systems can trigger notifications

## Reports System

### Available Reports

#### Occupancy Report
- Daily table utilization metrics
- Peak and average occupancy rates
- Best/worst performing days

#### Reservation Report
- Total reservations with status breakdown
- Cancellation and no-show rates
- Daily reservation metrics

#### Waitlist Report
- Waitlist conversion rates
- Average wait times
- Daily waitlist metrics

#### Lead Report (Future)
- Customer acquisition tracking
- Source attribution
- Conversion funnel metrics

### Export Features
- CSV export for all report types
- Date range filtering
- Project-specific metrics
