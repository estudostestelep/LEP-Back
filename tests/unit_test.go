package tests

import (
	"os"
	"testing"
	"time"

	"lep/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDataGeneration tests the seed data generation
func TestDataGeneration(t *testing.T) {
	// Test that seed data generation works
	seedData := utils.GenerateCompleteData()

	// Verify organizations
	require.NotEmpty(t, seedData.Organizations, "Should have organizations")
	org := seedData.Organizations[0]
	assert.NotEmpty(t, org.Name, "Organization should have name")
	assert.NotEmpty(t, org.Email, "Organization should have email")
	assert.True(t, org.Active, "Organization should be active")

	// Verify projects
	require.NotEmpty(t, seedData.Projects, "Should have projects")
	project := seedData.Projects[0]
	assert.NotEmpty(t, project.Name, "Project should have name")
	assert.True(t, project.Active, "Project should be active")

	// Verify users
	require.NotEmpty(t, seedData.Users, "Should have users")
	user := seedData.Users[0]
	assert.NotEmpty(t, user.Name, "User should have name")
	assert.NotEmpty(t, user.Email, "User should have email")
	assert.NotEmpty(t, user.Role, "User should have role")

	// Verify products
	require.NotEmpty(t, seedData.Products, "Should have products")
	product := seedData.Products[0]
	assert.NotEmpty(t, product.Name, "Product should have name")
	assert.Greater(t, product.Price, 0.0, "Product should have positive price")
	assert.True(t, product.Available, "Product should be available")

	// Verify tables
	require.NotEmpty(t, seedData.Tables, "Should have tables")
	table := seedData.Tables[0]
	assert.Greater(t, table.Number, 0, "Table should have positive number")
	assert.Greater(t, table.Capacity, 0, "Table should have positive capacity")
	assert.NotEmpty(t, table.Status, "Table should have status")

	// Verify customers
	require.NotEmpty(t, seedData.Customers, "Should have customers")
	customer := seedData.Customers[0]
	assert.NotEmpty(t, customer.Name, "Customer should have name")
	assert.NotEmpty(t, customer.Email, "Customer should have email")

	// Verify templates
	require.NotEmpty(t, seedData.Templates, "Should have notification templates")
	template := seedData.Templates[0]
	assert.NotEmpty(t, template.EventType, "Template should have event type")
	assert.NotEmpty(t, template.Channel, "Template should have channel")
	assert.NotEmpty(t, template.Template, "Template should have content")
}

// TestTestDataStructure tests the test data structure
func TestTestDataStructure(t *testing.T) {
	testData := NewTestData()

	// Verify UUIDs are valid format
	assert.Len(t, testData.OrganizationID, 36, "Organization ID should be valid UUID")
	assert.Len(t, testData.ProjectID, 36, "Project ID should be valid UUID")
	assert.Len(t, testData.UserID, 36, "User ID should be valid UUID")

	// Test sample data methods
	sampleOrg := testData.SampleOrganization()
	assert.NotEmpty(t, sampleOrg["name"], "Sample organization should have name")
	assert.NotEmpty(t, sampleOrg["email"], "Sample organization should have email")

	sampleUser := testData.SampleUser()
	assert.Equal(t, testData.OrganizationID, sampleUser["organization_id"], "Sample user should have org ID")
	assert.Equal(t, testData.ProjectID, sampleUser["project_id"], "Sample user should have project ID")

	sampleProduct := testData.SampleProduct()
	assert.NotEmpty(t, sampleProduct["name"], "Sample product should have name")
	assert.Greater(t, sampleProduct["price"], 0.0, "Sample product should have positive price")

	sampleCustomer := testData.SampleCustomer()
	assert.NotEmpty(t, sampleCustomer["name"], "Sample customer should have name")
	assert.NotEmpty(t, sampleCustomer["email"], "Sample customer should have email")

	sampleTable := testData.SampleTable()
	assert.Greater(t, sampleTable["number"], 0, "Sample table should have positive number")
	assert.Greater(t, sampleTable["capacity"], 0, "Sample table should have positive capacity")

	sampleReservation := testData.SampleReservation()
	assert.NotEmpty(t, sampleReservation["datetime"], "Sample reservation should have datetime")
	assert.Greater(t, sampleReservation["party_size"], 0, "Sample reservation should have positive party size")
}

// TestEnvironmentConfiguration tests environment setup
func TestEnvironmentConfiguration(t *testing.T) {
	// Test environment setup
	SetupTestEnvironment(t)

	// Verify test environment variables are set
	assert.Equal(t, "test", os.Getenv("GIN_MODE"))
	assert.Equal(t, "test", os.Getenv("ENVIRONMENT"))
	assert.Equal(t, "lep_test", os.Getenv("DB_NAME"))
	assert.Equal(t, "false", os.Getenv("ENABLE_CRON_JOBS"))
}

// TestTimeHandling tests time-related functionality
func TestTimeHandling(t *testing.T) {
	testData := NewTestData()
	reservation := testData.SampleReservation()

	// Verify datetime is in future
	datetimeStr, ok := reservation["datetime"].(string)
	require.True(t, ok, "Datetime should be string")

	parsedTime, err := time.Parse(time.RFC3339, datetimeStr)
	require.NoError(t, err, "Datetime should be valid RFC3339")

	assert.True(t, parsedTime.After(time.Now()), "Reservation should be in future")
}

// TestValidationData tests invalid data scenarios
func TestValidationData(t *testing.T) {
	invalidData := GetInvalidData()

	// Test invalid UUID
	assert.Equal(t, "invalid-uuid", invalidData.InvalidUUID())

	// Test empty string
	assert.Equal(t, "", invalidData.EmptyString())

	// Test invalid email
	assert.Equal(t, "invalid-email", invalidData.InvalidEmail())

	// Test invalid phone
	assert.Equal(t, "123", invalidData.InvalidPhone())

	// Test negative number
	assert.Equal(t, -1.0, invalidData.NegativeNumber())

	// Test too long string
	assert.Equal(t, 1000, len(invalidData.TooLongString()))
}