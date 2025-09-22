package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// BusinessScenariosTestSuite tests real business scenarios
type BusinessScenariosTestSuite struct {
	suite.Suite
	config   *TestConfig
	helper   *TestHelper
}

// SetupSuite runs once before all tests
func (suite *BusinessScenariosTestSuite) SetupSuite() {
	suite.config = NewTestConfig()
	suite.helper = &TestHelper{
		Router:    suite.config.Router,
		OrgID:     suite.config.OrgID,
		ProjectID: suite.config.ProjectID,
		AuthToken: suite.config.AuthToken,
	}
}

// SetupTest runs before each test
func (suite *BusinessScenariosTestSuite) SetupTest() {
	err := suite.config.SetupTestDatabase()
	suite.Require().NoError(err)
}

// TearDownTest runs after each test
func (suite *BusinessScenariosTestSuite) TearDownTest() {
	err := suite.config.CleanupTestDatabase()
	if err != nil {
		suite.T().Logf("Warning: Failed to cleanup: %v", err)
	}
}

// TestCompleteReservationFlow tests the complete reservation workflow
func (suite *BusinessScenariosTestSuite) TestCompleteReservationFlow() {
	// Step 1: Create a customer
	customerData := map[string]interface{}{
		"name":  "Maria Silva",
		"email": "maria.silva@email.com",
		"phone": "+55 11 99999-8888",
	}

	resp := suite.helper.MakeRequest(suite.T(), "POST", "/customer", customerData)
	suite.helper.AssertStatusCode(suite.T(), 201)

	var customerResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &customerResponse)
	suite.Require().NoError(err)
	customerID := customerResponse["id"].(string)

	// Step 2: Create a table
	tableData := map[string]interface{}{
		"number":   5,
		"capacity": 4,
		"status":   "livre",
		"location": "Salão Principal",
	}

	resp = suite.helper.MakeRequest(suite.T(), "POST", "/table", tableData)
	suite.helper.AssertStatusCode(suite.T(), 201)

	var tableResponse map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &tableResponse)
	suite.Require().NoError(err)
	tableID := tableResponse["id"].(string)

	// Step 3: Create a reservation
	futureTime := time.Now().Add(2 * time.Hour)
	reservationData := map[string]interface{}{
		"customer_id":      customerID,
		"table_id":         tableID,
		"datetime":         futureTime.Format(time.RFC3339),
		"party_size":       4,
		"status":           "confirmed",
		"special_requests": "Mesa próxima à janela",
	}

	resp = suite.helper.MakeRequest(suite.T(), "POST", "/reservation", reservationData)
	suite.helper.AssertStatusCode(suite.T(), 201)

	var reservationResponse map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &reservationResponse)
	suite.Require().NoError(err)
	reservationID := reservationResponse["id"].(string)

	// Step 4: Verify reservation was created correctly
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/reservation/"+reservationID, nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	var reservation map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &reservation)
	suite.Require().NoError(err)

	suite.Equal("confirmed", reservation["status"])
	suite.Equal(float64(4), reservation["party_size"])

	// Step 5: Update reservation status
	updateData := map[string]interface{}{
		"status": "seated",
	}

	resp = suite.helper.MakeRequest(suite.T(), "PUT", "/reservation/"+reservationID, updateData)
	suite.helper.AssertStatusCode(suite.T(), 200)

	// Step 6: Verify status update
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/reservation/"+reservationID, nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	err = json.Unmarshal(resp.Body.Bytes(), &reservation)
	suite.Require().NoError(err)
	suite.Equal("seated", reservation["status"])

	// Step 7: List all reservations
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/reservation", nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	var reservations []map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &reservations)
	suite.Require().NoError(err)
	suite.GreaterOrEqual(len(reservations), 1)
}

// TestOrderCreationFlow tests creating an order with multiple items
func (suite *BusinessScenariosTestSuite) TestOrderCreationFlow() {
	// Step 1: Create products
	products := []map[string]interface{}{
		{
			"name":              "Pizza Margherita",
			"description":       "Pizza clássica",
			"price":             35.90,
			"category":          "Pizzas",
			"prep_time_minutes": 20,
			"available":         true,
		},
		{
			"name":              "Refrigerante",
			"description":       "Refrigerante gelado",
			"price":             5.50,
			"category":          "Bebidas",
			"prep_time_minutes": 1,
			"available":         true,
		},
	}

	var productIDs []string
	for _, productData := range products {
		resp := suite.helper.MakeRequest(suite.T(), "POST", "/product", productData)
		suite.helper.AssertStatusCode(suite.T(), 201)

		var productResponse map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &productResponse)
		suite.Require().NoError(err)
		productIDs = append(productIDs, productResponse["id"].(string))
	}

	// Step 2: Create customer
	customerData := map[string]interface{}{
		"name":  "João Santos",
		"email": "joao.santos@email.com",
		"phone": "+55 11 88888-7777",
	}

	resp := suite.helper.MakeRequest(suite.T(), "POST", "/customer", customerData)
	suite.helper.AssertStatusCode(suite.T(), 201)

	var customerResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &customerResponse)
	suite.Require().NoError(err)
	customerID := customerResponse["id"].(string)

	// Step 3: Create table
	tableData := map[string]interface{}{
		"number":   3,
		"capacity": 2,
		"status":   "livre",
		"location": "Área Externa",
	}

	resp = suite.helper.MakeRequest(suite.T(), "POST", "/table", tableData)
	suite.helper.AssertStatusCode(suite.T(), 201)

	var tableResponse map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &tableResponse)
	suite.Require().NoError(err)
	tableID := tableResponse["id"].(string)

	// Step 4: Create order with multiple items
	orderData := map[string]interface{}{
		"customer_id":   customerID,
		"table_id":      tableID,
		"status":        "pending",
		"total_amount":  41.40, // Pizza + Refrigerante
		"source":        "internal",
		"items": []map[string]interface{}{
			{
				"product_id": productIDs[0],
				"quantity":   1,
				"price":      35.90,
				"notes":      "Sem cebola",
			},
			{
				"product_id": productIDs[1],
				"quantity":   1,
				"price":      5.50,
				"notes":      "Com gelo",
			},
		},
	}

	resp = suite.helper.MakeRequest(suite.T(), "POST", "/order", orderData)
	suite.helper.AssertStatusCode(suite.T(), 201)

	var orderResponse map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &orderResponse)
	suite.Require().NoError(err)
	orderID := orderResponse["id"].(string)

	// Step 5: Update order status to "preparing"
	updateData := map[string]interface{}{
		"status": "preparing",
	}

	resp = suite.helper.MakeRequest(suite.T(), "PUT", "/order/"+orderID, updateData)
	suite.helper.AssertStatusCode(suite.T(), 200)

	// Step 6: Verify order status
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/order/"+orderID, nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	var order map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &order)
	suite.Require().NoError(err)
	suite.Equal("preparing", order["status"])
	suite.Equal(41.40, order["total_amount"])

	// Step 7: Complete order
	updateData = map[string]interface{}{
		"status": "ready",
	}

	resp = suite.helper.MakeRequest(suite.T(), "PUT", "/order/"+orderID, updateData)
	suite.helper.AssertStatusCode(suite.T(), 200)
}

// TestWaitlistManagement tests waitlist functionality
func (suite *BusinessScenariosTestSuite) TestWaitlistManagement() {
	// Step 1: Create customer
	customerData := map[string]interface{}{
		"name":  "Pedro Oliveira",
		"email": "pedro.oliveira@email.com",
		"phone": "+55 11 77777-6666",
	}

	resp := suite.helper.MakeRequest(suite.T(), "POST", "/customer", customerData)
	suite.helper.AssertStatusCode(suite.T(), 201)

	var customerResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &customerResponse)
	suite.Require().NoError(err)
	customerID := customerResponse["id"].(string)

	// Step 2: Add customer to waitlist
	waitlistData := map[string]interface{}{
		"customer_id":      customerID,
		"party_size":       3,
		"estimated_wait":   25,
		"status":           "waiting",
		"special_requests": "Mesa para família",
	}

	resp = suite.helper.MakeRequest(suite.T(), "POST", "/waitlist", waitlistData)
	suite.helper.AssertStatusCode(suite.T(), 201)

	var waitlistResponse map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &waitlistResponse)
	suite.Require().NoError(err)
	waitlistID := waitlistResponse["id"].(string)

	// Step 3: Update estimated wait time
	updateData := map[string]interface{}{
		"estimated_wait": 15,
	}

	resp = suite.helper.MakeRequest(suite.T(), "PUT", "/waitlist/"+waitlistID, updateData)
	suite.helper.AssertStatusCode(suite.T(), 200)

	// Step 4: Mark as seated
	updateData = map[string]interface{}{
		"status": "seated",
	}

	resp = suite.helper.MakeRequest(suite.T(), "PUT", "/waitlist/"+waitlistID, updateData)
	suite.helper.AssertStatusCode(suite.T(), 200)

	// Step 5: Verify final status
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/waitlist/"+waitlistID, nil)
	suite.helper.AssertStatusCode(suite.T(), 200)

	var waitlist map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &waitlist)
	suite.Require().NoError(err)
	suite.Equal("seated", waitlist["status"])
	suite.Equal(float64(15), waitlist["estimated_wait"])
}

// TestMultiTenantIsolation tests that data is isolated between organizations
func (suite *BusinessScenariosTestSuite) TestMultiTenantIsolation() {
	// Create data with current organization
	productData := map[string]interface{}{
		"name":        "Produto Org 1",
		"description": "Produto da primeira organização",
		"price":       25.50,
		"category":    "Teste",
		"available":   true,
	}

	resp := suite.helper.MakeRequest(suite.T(), "POST", "/product", productData)
	suite.helper.AssertStatusCode(suite.T(), 201)

	// Create helper with different organization ID
	differentOrgHelper := &TestHelper{
		Router:    suite.config.Router,
		OrgID:     "different-org-id",
		ProjectID: "different-project-id",
		AuthToken: suite.config.AuthToken,
	}

	// Try to access data with different organization - should not see the product
	resp = differentOrgHelper.MakeRequest(suite.T(), "GET", "/product", nil)
	// Should either return empty list or fail with 400 (invalid headers)
	suite.True(resp.Code == 200 || resp.Code == 400)

	if resp.Code == 200 {
		var products []map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &products)
		suite.Require().NoError(err)
		// Should not see products from other organization
		for _, product := range products {
			suite.NotEqual("Produto Org 1", product["name"])
		}
	}
}

// TestBusinessValidationRules tests business rules validation
func (suite *BusinessScenariosTestSuite) TestBusinessValidationRules() {
	// Test 1: Product with negative price should fail
	invalidProduct := map[string]interface{}{
		"name":     "Produto Inválido",
		"price":    -10.00,
		"category": "Teste",
	}

	resp := suite.helper.MakeRequest(suite.T(), "POST", "/product", invalidProduct)
	suite.helper.AssertStatusCode(suite.T(), 400)

	// Test 2: Table with zero capacity should fail
	invalidTable := map[string]interface{}{
		"number":   1,
		"capacity": 0,
		"status":   "livre",
	}

	resp = suite.helper.MakeRequest(suite.T(), "POST", "/table", invalidTable)
	suite.helper.AssertStatusCode(suite.T(), 400)

	// Test 3: Reservation in the past should fail
	pastTime := time.Now().Add(-2 * time.Hour)
	invalidReservation := map[string]interface{}{
		"datetime":   pastTime.Format(time.RFC3339),
		"party_size": 2,
		"status":     "confirmed",
	}

	resp = suite.helper.MakeRequest(suite.T(), "POST", "/reservation", invalidReservation)
	suite.helper.AssertStatusCode(suite.T(), 400)

	// Test 4: Customer with invalid email should fail
	invalidCustomer := map[string]interface{}{
		"name":  "Cliente Inválido",
		"email": "email-invalido",
		"phone": "+55 11 99999-9999",
	}

	resp = suite.helper.MakeRequest(suite.T(), "POST", "/customer", invalidCustomer)
	suite.helper.AssertStatusCode(suite.T(), 400)
}

// TestBusinessScenariosTestSuite runs the business scenarios test suite
func TestBusinessScenariosTestSuite(t *testing.T) {
	suite.Run(t, new(BusinessScenariosTestSuite))
}