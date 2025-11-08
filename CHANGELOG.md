# LEP System Changelog

All notable changes to the LEP System backend will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

#### Menu Features
- **Menu Name Uniqueness Validation**: Implemented case-insensitive validation to prevent duplicate menu names within a project
  - Repository layer: `CheckMenuNameExists()` with LOWER() comparison
  - Handler layer: Custom error handling with proper validation
  - Server layer: HTTP 409 Conflict responses
  - Support for update operations (menu can keep its own name)
  - Comprehensive test suite: `test_menu_uniqueness.sh`
  - Documentation: `MENU_UNIQUENESS_IMPLEMENTATION.md`

#### Theme Customization Features
- **Extended Theme Customization**: Expanded from 7 to 15 fields
  - New numeric configuration fields:
    - `disabled_opacity` (0.0-1.0): Opacity for disabled states
    - `shadow_intensity` (0.0-2.0): Shadow and depth effects
  - Frontend UI with range sliders and number inputs
  - Complete backend validation
  - Seed data for all 15 fields (generic + Fattoria)
  - Comprehensive test coverage: `test_theme_customization.sh`

### Fixed
- Theme customization route ordering (fixed 404 errors on specific routes)
- Partial update support for theme customization
- Default value consistency across model, handler, and repository layers
- Soft delete compliance for theme customization queries

### Changed
- Menu validation workflow: Added business-layer uniqueness checks
- Theme customization update logic: Support for partial field updates
- Error response handling: Introduced HTTP 409 Conflict for resource conflicts

### Deprecated
- N/A

### Removed
- N/A

### Security
- Menu name validation prevents data integrity issues at application level
- Case-insensitive comparison handles potential security edge cases

---

## [1.0.0] - 2024-11-08

### Added

#### Core Features
- **REST API** with Gin web framework
- **Multi-tenant Architecture**: Organization and Project isolation
- **Authentication**: JWT-based with token blacklist/whitelist
- **Clean Architecture**: 3-layer pattern (Handler → Server → Repository)
- **GORM ORM**: PostgreSQL database integration
- **Soft Delete Pattern**: Logical deletion with `deleted_at` field
- **Audit Logging**: Complete operation tracking

#### Entity Management
- **Users**: Employee/admin management with roles and permissions
- **Customers**: Customer relationship management
- **Tables**: Restaurant table capacity and availability tracking
- **Products/Menu**: Item management with categories and pricing
- **Orders**: Order creation, status tracking, and preparation times
- **Reservations**: Table booking with datetime and party size
- **Waitlist**: Queue management for occupied tables
- **Settings**: Project-level configuration
- **Environment**: Physical space definitions (dining areas)

#### Notification System
- **Multi-channel Support**: SMS, WhatsApp, Email
- **Event-driven Architecture**: Automatic triggers for business events
- **Template System**: Dynamic templates with variable substitution
- **Cron Jobs**: Scheduled tasks (24h confirmation reminders)
- **Webhook Support**: Inbound message handling and status updates
- **Notification Logging**: Complete delivery tracking

#### Kitchen Management
- **Kitchen Queue**: Real-time order tracking
- **Prep Time Management**: Estimated and actual preparation times
- **Order Status Workflow**: Complete order lifecycle management

#### Reports
- **Occupancy Reports**: Table utilization metrics
- **Reservation Reports**: Booking patterns and analytics
- **Waitlist Reports**: Queue performance metrics
- **Lead Reports**: Customer acquisition tracking

#### Infrastructure
- **Docker Support**: Containerization for development and production
- **Terraform Configuration**: GCP deployment with Cloud Run + Cloud SQL
- **Seed System**: Data population with Fattoria Italian restaurant theme
- **GCP Integration**: Cloud Storage, Secret Manager, Cloud SQL

### Infrastructure & DevOps
- Cloud Run deployment configuration
- Cloud SQL PostgreSQL integration
- GCS (Google Cloud Storage) image storage
- JWT key management via Secret Manager
- Terraform infrastructure as code
- Docker and docker-compose support

### Documentation
- Comprehensive API documentation
- Setup and installation guides
- Deployment guides (local, staging, production)
- GCP authentication setup
- Fattoria demo theme documentation
- Quick start guides

### Scripts
- `run_seed.sh`: Standard seed data population
- `run_tests.sh`: Automated test execution
- `run_seed_fattoria.sh`: Fattoria theme seeding
- `master-interactive.sh`: Interactive menu-driven operations
- `dev-local.sh`: Local development environment
- `stage-local.sh`: Staging environment setup
- `stage-deploy.sh`: Staging deployment

---

## Development Notes

### Current Focus Areas
- Menu system enhancements (uniqueness validation, smart selection)
- Theme customization extensibility
- Performance optimization for multi-tenant queries
- Webhook integration improvements

### Known Limitations
- Validation is at application level (not database constraints)
- Soft delete queries require explicit `deleted_at IS NULL` checks
- Multi-tenant isolation depends on header validation

### Future Roadmap
- Database-level uniqueness constraints (with migration strategy)
- Advanced analytics and reporting
- Loyalty program integration
- POS system integration
- Mobile app support

---

## Migration Guide

### From Previous Versions
See `docs/` folder for version-specific migration guides:
- `docs/QUICKSTART.md`: Quick start guide
- `docs/SETUP.md`: Detailed setup instructions
- `docs/ENVIRONMENT_SETUP.md`: Environment configuration

### Database Migrations
Run migrations before starting the server:
```bash
go run cmd/migrate/main.go
```

---

## Contributing

When adding new features:
1. Update this CHANGELOG
2. Follow the [Keep a Changelog](https://keepachangelog.com/) format
3. Add entries under [Unreleased] section
4. Create pull requests with clear descriptions
5. Include tests and documentation

---

## Support

For issues and questions:
- Check `docs/` for detailed documentation
- Review test files for usage examples
- Consult `CLAUDE.md` for development guidelines

---

## License

All rights reserved. See project configuration for details.
