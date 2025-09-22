package tests

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

// RealIntegrationTestSuite tests actual business logic with real database
type RealIntegrationTestSuite struct {
	suite.Suite
	config   *TestConfig
	helper   *TestHelper
}

// SetupSuite runs once before all tests
func (suite *RealIntegrationTestSuite) SetupSuite() {
	// Initialize test configuration with real database and handlers
	suite.config = NewTestConfig()

	// Create helper with real router
	suite.helper = &TestHelper{
		Router:    suite.config.Router,
		OrgID:     suite.config.OrgID,
		ProjectID: suite.config.ProjectID,
		AuthToken: suite.config.AuthToken,
	}
}

// SetupTest runs before each test
func (suite *RealIntegrationTestSuite) SetupTest() {
	// Setup fresh test database for each test
	err := suite.config.SetupTestDatabase()
	suite.Require().NoError(err, "Failed to setup test database")
}

// TearDownTest runs after each test
func (suite *RealIntegrationTestSuite) TearDownTest() {
	// Clean up after each test
	err := suite.config.CleanupTestDatabase()
	if err != nil {
		suite.T().Logf("Warning: Failed to cleanup test database: %v", err)
	}
}

// TearDownSuite runs once after all tests
func (suite *RealIntegrationTestSuite) TearDownSuite() {
	// Final cleanup
	if suite.config != nil {
		err := suite.config.CleanupTestDatabase()
		if err != nil {
			suite.T().Logf("Warning: Failed to cleanup test database: %v", err)
		}
	}
}

// TestProductCRUD tests complete CRUD operations for products
func (suite *RealIntegrationTestSuite) TestProductCRUD() {
	// Test creating a product
	productData := map[string]interface{}{
		"name":              "Pizza Margherita",
		"description":       "Pizza clássica com molho de tomate e queijo",
		"price":             35.90,
		"category":          "Pizzas",
		"prep_time_minutes": 20,
		"available":         true,
		"ingredients":       []string{"massa", "molho de tomate", "queijo mozzarella", "manjericão"},
	}

	// Create product
	resp := suite.helper.MakeRequest(suite.T(), "POST", "/product", productData)
	suite.helper.AssertStatusCode(suite.T(), 201)

	var createResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &createResponse)
	suite.Require().NoError(err)

	productID, exists := createResponse["id"]
	suite.Require().True(exists, "Response should contain product ID")

	// Test reading the created product
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/product/"+productID.(string), nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	var product map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &product)
	suite.Require().NoError(err)

	suite.Equal("Pizza Margherita", product["name"])
	suite.Equal(35.90, product["price"])

	// Test updating the product
	updateData := map[string]interface{}{
		"price": 39.90,
	}

	resp = suite.helper.MakeRequest(suite.T(), "PUT", "/product/"+productID.(string), updateData)
	suite.helper.AssertStatusCode(suite.T(), 200)

	// Verify update
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/product/"+productID.(string), nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	err = json.Unmarshal(resp.Body.Bytes(), &product)
	suite.Require().NoError(err)
	suite.Equal(39.90, product["price"])

	// Test listing products
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/product", nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	var products []map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &products)
	suite.Require().NoError(err)
	suite.GreaterOrEqual(len(products), 1, "Should have at least one product")

	// Test deleting the product
	resp = suite.helper.MakeRequest(suite.T(), "DELETE", "/product/"+productID.(string), nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	// Verify deletion (should return 404)
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/product/"+productID.(string), nil)
	suite.helper.AssertStatusCode(suite.T(), 404)
}

// TestCustomerCRUD tests complete CRUD operations for customers
func (suite *RealIntegrationTestSuite) TestCustomerCRUD() {
	// Test creating a customer
	customerData := map[string]interface{}{
		"name":        "João Silva",
		"email":       "joao.silva@email.com",
		"phone":       "+55 11 98765-4321",
		"birth_date":  "1985-10-15",
		"preferences": "Vegetariano",
	}

	// Create customer
	resp := suite.helper.MakeRequest(suite.T(), "POST", "/customer", customerData)
	suite.helper.AssertStatusCode(suite.T(), 201)

	var createResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &createResponse)
	suite.Require().NoError(err)

	customerID, exists := createResponse["id"]
	suite.Require().True(exists, "Response should contain customer ID")

	// Test reading the created customer
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/customer/"+customerID.(string), nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	var customer map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &customer)
	suite.Require().NoError(err)

	suite.Equal("João Silva", customer["name"])
	suite.Equal("joao.silva@email.com", customer["email"])

	// Test listing customers
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/customer", nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	var customers []map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &customers)
	suite.Require().NoError(err)
	suite.GreaterOrEqual(len(customers), 1, "Should have at least one customer")
}

// TestTableManagement tests table operations
func (suite *RealIntegrationTestSuite) TestTableManagement() {
	// Test creating a table
	tableData := map[string]interface{}{
		"number":   10,
		"capacity": 4,
		"status":   "livre",
		"location": "Salão Principal",
	}

	// Create table
	resp := suite.helper.MakeRequest(suite.T(), "POST", "/table", tableData)
	suite.helper.AssertStatusCode(suite.T(), 201)

	var createResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &createResponse)
	suite.Require().NoError(err)

	tableID, exists := createResponse["id"]
	suite.Require().True(exists, "Response should contain table ID")

	// Test reading the created table
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/table/"+tableID.(string), nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	var table map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &table)
	suite.Require().NoError(err)

	suite.Equal(float64(10), table["number"]) // JSON numbers come as float64
	suite.Equal("livre", table["status"])

	// Test updating table status
	updateData := map[string]interface{}{
		"status": "ocupada",
	}

	resp = suite.helper.MakeRequest(suite.T(), "PUT", "/table/"+tableID.(string), updateData)
	suite.helper.AssertStatusCode(suite.T(), 200)

	// Verify status update
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/table/"+tableID.(string), nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	err = json.Unmarshal(resp.Body.Bytes(), &table)
	suite.Require().NoError(err)
	suite.Equal("ocupada", table["status"])
}

// TestDataValidation tests input validation
func (suite *RealIntegrationTestSuite) TestDataValidation() {
	// Test invalid product data
	invalidProduct := map[string]interface{}{
		"name":  "", // Empty name should fail
		"price": -10, // Negative price should fail
	}

	resp := suite.helper.MakeRequest(suite.T(), "POST", "/product", invalidProduct)
	suite.helper.AssertStatusCode(suite.T(), 400) // Bad request for validation error

	// Test invalid customer email
	invalidCustomer := map[string]interface{}{
		"name":  "Test Customer",
		"email": "invalid-email", // Invalid email format
	}

	resp = suite.helper.MakeRequest(suite.T(), "POST", "/customer", invalidCustomer)
	suite.helper.AssertStatusCode(suite.T(), 400) // Bad request for validation error
}

// TestOrganizationHeaderValidation tests multi-tenant header validation
func (suite *RealIntegrationTestSuite) TestOrganizationHeaderValidation() {
	// Create helper without headers
	helperNoHeaders := &TestHelper{
		Router: suite.config.Router,
		// No OrgID or ProjectID
	}

	// Test that requests without headers fail
	resp := helperNoHeaders.MakeRequest(suite.T(), "GET", "/product", nil)
	suite.helper.AssertStatusCode(suite.T(), 400) // Should fail without headers

	// Test with invalid headers
	helperInvalidHeaders := &TestHelper{
		Router:    suite.config.Router,
		OrgID:     "invalid-uuid",
		ProjectID: "invalid-uuid",
	}

	resp = helperInvalidHeaders.MakeRequest(suite.T(), "GET", "/product", nil)
	suite.helper.AssertStatusCode(suite.T(), 400) // Should fail with invalid UUIDs
}

// TestRealIntegrationTestSuite runs the real integration test suite
func TestRealIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RealIntegrationTestSuite))
}