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

Required for database connection:
```bash
DB_USER=postgres_username
DB_PASS=postgres_password
DB_NAME=database_name
INSTANCE_UNIX_SOCKET=/path/to/socket  # For GCP Cloud SQL
JWT_SECRET_PRIVATE_KEY=jwt_private_key
JWT_SECRET_PUBLIC_KEY=jwt_public_key
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
