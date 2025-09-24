package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

var (
	// Environment configuration
	ENV  = getEnvironment()
	PORT = getPort()

	// JWT configuration
	JWT_SECRET_PRIVATE_KEY = os.Getenv("JWT_SECRET_PRIVATE_KEY")
	JWT_SECRET_PUBLIC_KEY  = os.Getenv("JWT_SECRET_PUBLIC_KEY")

	// Database configuration
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_USER     = os.Getenv("DB_USER")
	DB_PASS     = os.Getenv("DB_PASS")
	DB_NAME     = os.Getenv("DB_NAME")
	DB_SSL_MODE = getDBSSLMode()

	// For GCP Cloud SQL
	INSTANCE_UNIX_SOCKET = os.Getenv("INSTANCE_UNIX_SOCKET")

	// Application configuration
	ENABLE_CRON_JOBS = getEnableCronJobs()
	GIN_MODE         = getGinMode()
	LOG_LEVEL        = getLogLevel()
)

// getEnvironment returns the current environment or defaults to "dev"
func getEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "dev"
	}

	// Validate environment
	validEnvs := []string{"local-dev", "dev", "staging", "prod"}
	for _, validEnv := range validEnvs {
		if env == validEnv {
			return env
		}
	}

	log.Printf("Warning: Invalid environment '%s', defaulting to 'dev'", env)
	return "dev"
}

// getPort returns the port number or defaults to 8080
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// getDBSSLMode returns SSL mode for database connections
func getDBSSLMode() string {
	sslMode := os.Getenv("DB_SSL_MODE")
	if sslMode == "" {
		// Default based on environment
		if ENV == "local-dev" {
			return "disable"
		}
		return "require"
	}
	return sslMode
}

// getEnableCronJobs returns whether cron jobs should be enabled
func getEnableCronJobs() bool {
	cronJobs := os.Getenv("ENABLE_CRON_JOBS")
	if cronJobs == "" {
		// Default based on environment
		return ENV == "staging" || ENV == "prod"
	}

	enabled, err := strconv.ParseBool(cronJobs)
	if err != nil {
		log.Printf("Warning: Invalid ENABLE_CRON_JOBS value '%s', defaulting to false", cronJobs)
		return false
	}
	return enabled
}

// getGinMode returns the Gin framework mode
func getGinMode() string {
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		// Default based on environment
		if ENV == "local-dev" {
			return "debug"
		} else if ENV == "dev" {
			return "debug"
		}
		return "release"
	}
	return mode
}

// getLogLevel returns the logging level
func getLogLevel() string {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		// Default based on environment
		if ENV == "local-dev" || ENV == "dev" {
			return "debug"
		}
		return "info"
	}
	return strings.ToLower(level)
}

// IsLocalDev returns true if running in local development mode
func IsLocalDev() bool {
	return ENV == "local-dev"
}

// IsDev returns true if running in any development mode
func IsDev() bool {
	return ENV == "local-dev" || ENV == "dev"
}

// IsStaging returns true if running in staging mode
func IsStaging() bool {
	return ENV == "staging"
}

// IsProd returns true if running in production mode
func IsProd() bool {
	return ENV == "prod"
}

// IsGCP returns true if running on Google Cloud Platform
func IsGCP() bool {
	return ENV == "dev" || ENV == "staging" || ENV == "prod"
}
