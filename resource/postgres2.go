package resource

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"lep/config"

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
		return db, nil
	}

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s",
		config.DB_USER, config.DB_PASS, config.DB_NAME, config.INSTANCE_UNIX_SOCKET)

	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	ctx := context.Background()
	sqlDB, err := dbPool.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sql.DB: %w", err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to open GORM connection: %w", err)
	}

	log.Printf("Database connection established successfully for environment: %s", config.ENV)
	return db, nil
}
