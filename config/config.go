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
	DISABLE_AUTO_SEED = getDisableAutoSeed()

	// Storage configuration
	STORAGE_TYPE        = getStorageType()
	STORAGE_BUCKET_NAME = os.Getenv("STORAGE_BUCKET_NAME")
	BASE_URL            = getBaseURL()

	// Bucket configuration
	BUCKET_NAME          = os.Getenv("BUCKET_NAME")
	BUCKET_CACHE_CONTROL = os.Getenv("BUCKET_CACHE_CONTROL")
	BUCKET_TIMEOUT, _    = strconv.Atoi(os.Getenv("BUCKET_TIMEOUT"))
)

// getEnvironment returns the current environment or defaults to "dev"
func getEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "dev"
	}

	// Validate environment - simplified to dev/stage/prod
	validEnvs := []string{"dev", "stage", "prod"}
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
		if ENV == "dev" {
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
		return ENV == "stage" || ENV == "prod"
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
		if ENV == "dev" {
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
		if ENV == "dev" {
			return "debug"
		}
		return "info"
	}
	return strings.ToLower(level)
}

// IsDev returns true if running in development mode
func IsDev() bool {
	return ENV == "dev"
}

// IsStage returns true if running in staging mode
func IsStage() bool {
	return ENV == "stage"
}

// IsProd returns true if running in production mode
func IsProd() bool {
	return ENV == "prod"
}

// IsGCP returns true if running on Google Cloud Platform
func IsGCP() bool {
	return ENV == "stage" || ENV == "prod"
}

// getStorageType returns the storage type based on environment
func getStorageType() string {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		// Default based on environment
		if IsDev() {
			return "local"
		}
		return "gcs"
	}
	return storageType
}

// getBaseURL returns the base URL for storage
func getBaseURL() string {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		// Default based on environment and storage type
		storageType := os.Getenv("STORAGE_TYPE")
		bucketName := os.Getenv("STORAGE_BUCKET_NAME")

		if (storageType == "gcs" || (storageType == "" && !IsDev())) && bucketName != "" {
			return "https://storage.googleapis.com/" + bucketName
		}
		// Default for local development
		return "http://localhost:8080"
	}
	return baseURL
}

// IsLocalStorage returns true if using local storage
func IsLocalStorage() bool {
	return STORAGE_TYPE == "local"
}

// IsGCSStorage returns true if using Google Cloud Storage
func IsGCSStorage() bool {
	return STORAGE_TYPE == "gcs"
}

// getDisableAutoSeed returns whether auto-seeding should be disabled
func getDisableAutoSeed() bool {
	autoSeed := os.Getenv("DISABLE_AUTO_SEED")
	if autoSeed == "" {
		// Default: allow auto-seed in all environments
		return false
	}

	disabled, err := strconv.ParseBool(autoSeed)
	if err != nil {
		log.Printf("Warning: Invalid DISABLE_AUTO_SEED value '%s', defaulting to false", autoSeed)
		return false
	}
	return disabled
}

// IsAutoSeedEnabled returns true if auto-seeding is enabled
func IsAutoSeedEnabled() bool {
	return !DISABLE_AUTO_SEED
}
