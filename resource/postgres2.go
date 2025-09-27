package resource

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"lep/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenConnDBPostgres2() (*gorm.DB, error) {
	// For local development (traditional host:port connection)
	if config.IsLocalDev() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Sao_Paulo",
			config.DB_HOST, config.DB_USER, config.DB_PASS, config.DB_NAME, config.DB_PORT, config.DB_SSL_MODE)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		// Configure connection pool for local development
		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}
		configureConnectionPool(sqlDB)
		return db, nil
	}

	// For GCP Cloud SQL (unix socket connection)
	// Fallback to TCP if INSTANCE_UNIX_SOCKET is not available (Windows)
	if config.INSTANCE_UNIX_SOCKET == "" || config.DB_HOST != "" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Sao_Paulo",
			config.DB_HOST, config.DB_USER, config.DB_PASS, config.DB_NAME, config.DB_PORT, config.DB_SSL_MODE)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		// Configure connection pool for TCP connections
		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}
		configureConnectionPool(sqlDB)
		return db, nil
	}

	dsn := fmt.Sprintf("user=%s password=%s database=%s host=%s sslmode=disable",
		config.DB_USER, config.DB_PASS, config.DB_NAME, config.INSTANCE_UNIX_SOCKET)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database via unix socket: %w", err)
	}

	// Configure connection pool for Unix socket connections
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	configureConnectionPool(sqlDB)

	log.Printf("Database connection established successfully for environment: %s", config.ENV)
	return db, nil
}

// configureConnectionPool sets optimal connection pool settings for different environments
func configureConnectionPool(db *sql.DB) {
	// Maximum number of open connections to the database
	// For Cloud Run: Lower values to avoid exceeding database connection limits
	if config.IsGCP() {
		// Very conservative for db-f1-micro (max 10 total connections)
		db.SetMaxOpenConns(2)  // Extremely conservative for f1-micro
		db.SetMaxIdleConns(1)  // Minimal idle connections
	} else {
		// Local development
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
	}

	// Maximum amount of time a connection may be reused
	// Cloud SQL connections should be refreshed regularly
	if config.IsGCP() {
		db.SetConnMaxLifetime(30 * time.Minute) // Shorter lifetime for f1-micro
		db.SetConnMaxIdleTime(10 * time.Minute) // Reduced idle time for f1-micro
	} else {
		db.SetConnMaxLifetime(30 * time.Minute)
		db.SetConnMaxIdleTime(5 * time.Minute)
	}

	if config.IsGCP() {
		log.Printf("Database connection pool configured (GCP/f1-micro): MaxOpen=%d, MaxIdle=%d, MaxLifetime=%s, MaxIdleTime=%s",
			2, 1, 30*time.Minute, 10*time.Minute)
	} else {
		log.Printf("Database connection pool configured (Local): MaxOpen=%d, MaxIdle=%d, MaxLifetime=%s, MaxIdleTime=%s",
			10, 5, 30*time.Minute, 5*time.Minute)
	}
}

// CheckConnectionHealth verifica se a conexão está saudável e tenta reconectar se necessário
func CheckConnectionHealth(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Tenta fazer um ping na conexão
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Printf("Database connection unhealthy, attempting to reconnect: %v", err)
		return fmt.Errorf("database connection unhealthy: %w", err)
	}

	return nil
}
