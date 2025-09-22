package tests

import (
	"os"
	"testing"
)

// SetupTestEnvironment configures environment variables for testing
func SetupTestEnvironment(t *testing.T) {
	// Set test-specific environment variables
	testVars := map[string]string{
		"GIN_MODE":     "test",
		"ENVIRONMENT":  "test",
		"DB_NAME":      "lep_test",
		"DB_USER":      "postgres",
		"DB_PASS":      "postgres",
		"DB_HOST":      "localhost",
		"DB_PORT":      "5432",
		"JWT_SECRET_PRIVATE_KEY": "test-private-key",
		"JWT_SECRET_PUBLIC_KEY":  "test-public-key",
		"ENABLE_CRON_JOBS":       "false",
	}

	// Set environment variables for testing
	for key, value := range testVars {
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	t.Cleanup(func() {
		// Optionally clean up environment variables after test
		// This is usually not necessary as tests run in isolated processes
	})
}

// EnsureTestDatabase checks if test database is available
func EnsureTestDatabase(t *testing.T) {
	// Skip tests if database is not available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests (SKIP_DB_TESTS=true)")
	}

	// Check if PostgreSQL is available
	if os.Getenv("DB_HOST") == "" {
		t.Skip("Skipping database tests (no DB_HOST configured)")
	}
}

// MockExternalServices disables external service calls during testing
func MockExternalServices() {
	// Disable external service calls
	os.Setenv("TWILIO_ACCOUNT_SID", "")
	os.Setenv("TWILIO_AUTH_TOKEN", "")
	os.Setenv("SMTP_HOST", "")
	os.Setenv("ENABLE_NOTIFICATIONS", "false")
}