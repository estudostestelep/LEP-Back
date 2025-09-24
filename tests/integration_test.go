package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// APITestSuite is the main test suite for all API endpoints
type APITestSuite struct {
	suite.Suite
	config   *TestConfig
	helper   *TestHelper
	testData *TestData
}

// SetupSuite runs once before all tests
func (suite *APITestSuite) SetupSuite() {
	fmt.Println("\n=== LEP Backend API Test Suite ===")
	fmt.Printf("Starting comprehensive route testing at %s\n", time.Now().Format("2006-01-02 15:04:05"))

	// Initialize test configuration with real database and handlers
	suite.config = NewTestConfig()
	suite.testData = suite.config.TestData

	// Create helper with real router
	suite.helper = &TestHelper{
		Router:    suite.config.Router,
		OrgID:     suite.config.OrgID,
		ProjectID: suite.config.ProjectID,
		AuthToken: suite.config.AuthToken,
	}

	// Setup test database
	err := suite.config.SetupTestDatabase()
	if err != nil {
		suite.T().Fatalf("Failed to setup test database: %v", err)
	}
}

// TearDownSuite runs once after all tests
func (suite *APITestSuite) TearDownSuite() {
	fmt.Println("\n=== Test Suite Complete ===")
	fmt.Printf("Finished at %s\n", time.Now().Format("2006-01-02 15:04:05"))

	// Cleanup test database
	if suite.config != nil {
		err := suite.config.CleanupTestDatabase()
		if err != nil {
			fmt.Printf("Warning: Failed to cleanup test database: %v\n", err)
		}
	}
}


// Test Public Routes
func (suite *APITestSuite) TestPublicRoutes() {
	// Test ping endpoint
	resp := suite.helper.MakeRequest(suite.T(), "GET", "/ping", nil)
	suite.helper.AssertStatusCode(suite.T(), 200)
	suite.Equal("\"pong\"", resp.Body.String())

	// Test health endpoint
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/health", nil)
	suite.helper.AssertStatusCode(suite.T(), 200)
	response := suite.helper.AssertJSONResponse(suite.T(), "status")
	suite.Equal("healthy", response["status"])

	// Note: Login and user creation require real authentication setup
	// These will be tested in dedicated authentication tests
}

// Test Auth Routes
func (suite *APITestSuite) TestAuthRoutes() {
	// Note: Authentication tests require proper JWT token setup
	// For now, we test that routes exist and require authentication
	resp := suite.helper.MakeRequest(suite.T(), "POST", "/logout", nil)
	// Should require authentication (401 or 400)
	suite.True(resp.Code == 401 || resp.Code == 400, "Logout should require authentication")

	resp = suite.helper.MakeRequest(suite.T(), "POST", "/checkToken", nil)
	// Should require authentication (401 or 400)
	suite.True(resp.Code == 401 || resp.Code == 400, "CheckToken should require authentication")
}

// Test User Routes
func (suite *APITestSuite) TestUserRoutes() {
	// Test that user routes require authentication and headers
	resp := suite.helper.MakeRequest(suite.T(), "GET", "/user", nil)
	// Should require authentication/headers
	suite.True(resp.Code == 401 || resp.Code == 400, "User list should require authentication")

	// Test specific user endpoint
	resp = suite.helper.MakeRequest(suite.T(), "GET", "/user/"+suite.testData.UserID, nil)
	suite.True(resp.Code == 401 || resp.Code == 400, "Get user should require authentication")
}

// Test Product Routes
func (suite *APITestSuite) TestProductRoutes() {
	routes := []TestRoute{
		{Method: "GET", Path: "/product/" + suite.testData.ProductID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Get product by ID"},
		{Method: "GET", Path: "/product/purchase/" + suite.testData.ProductID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Get product by purchase"},
		{Method: "GET", Path: "/product", ExpectedStatus: 200, RequiresHeaders: true, Description: "List products"},
		{Method: "POST", Path: "/product", ExpectedStatus: 201, RequiresHeaders: true, TestBody: suite.testData.SampleProduct(), Description: "Create product"},
		{Method: "PUT", Path: "/product/" + suite.testData.ProductID, ExpectedStatus: 200, RequiresHeaders: true, TestBody: suite.testData.SampleProduct(), Description: "Update product"},
		{Method: "DELETE", Path: "/product/" + suite.testData.ProductID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Delete product"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Product Routes")
}

// Test Table Routes
func (suite *APITestSuite) TestTableRoutes() {
	routes := []TestRoute{
		{Method: "GET", Path: "/table/" + suite.testData.TableID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Get table by ID"},
		{Method: "GET", Path: "/table", ExpectedStatus: 200, RequiresHeaders: true, Description: "List tables"},
		{Method: "POST", Path: "/table", ExpectedStatus: 201, RequiresHeaders: true, TestBody: suite.testData.SampleTable(), Description: "Create table"},
		{Method: "PUT", Path: "/table/" + suite.testData.TableID, ExpectedStatus: 200, RequiresHeaders: true, TestBody: suite.testData.SampleTable(), Description: "Update table"},
		{Method: "DELETE", Path: "/table/" + suite.testData.TableID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Delete table"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Table Routes")
}

// Test Order Routes
func (suite *APITestSuite) TestOrderRoutes() {
	routes := []TestRoute{
		{Method: "GET", Path: "/order/" + suite.testData.OrderID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Get order by ID"},
		{Method: "GET", Path: "/order/" + suite.testData.OrderID + "/progress", ExpectedStatus: 200, RequiresHeaders: true, Description: "Get order progress"},
		{Method: "GET", Path: "/order", ExpectedStatus: 200, RequiresHeaders: true, Description: "List orders"},
		{Method: "POST", Path: "/order", ExpectedStatus: 201, RequiresHeaders: true, TestBody: suite.testData.SampleOrder(), Description: "Create order"},
		{Method: "PUT", Path: "/order/" + suite.testData.OrderID, ExpectedStatus: 200, RequiresHeaders: true, TestBody: suite.testData.SampleOrder(), Description: "Update order"},
		{Method: "PUT", Path: "/order/" + suite.testData.OrderID + "/status", ExpectedStatus: 200, RequiresHeaders: true, TestBody: map[string]string{"status": "ready"}, Description: "Update order status"},
		{Method: "DELETE", Path: "/order/" + suite.testData.OrderID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Delete order"},
		{Method: "GET", Path: "/kitchen/queue", ExpectedStatus: 200, RequiresHeaders: true, Description: "Get kitchen queue"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Order & Kitchen Routes")
}

// Test Reservation Routes
func (suite *APITestSuite) TestReservationRoutes() {
	routes := []TestRoute{
		{Method: "GET", Path: "/reservation/" + suite.testData.ReservationID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Get reservation by ID"},
		{Method: "GET", Path: "/reservation", ExpectedStatus: 200, RequiresHeaders: true, Description: "List reservations"},
		{Method: "POST", Path: "/reservation", ExpectedStatus: 201, RequiresHeaders: true, TestBody: suite.testData.SampleReservation(), Description: "Create reservation"},
		{Method: "PUT", Path: "/reservation/" + suite.testData.ReservationID, ExpectedStatus: 200, RequiresHeaders: true, TestBody: suite.testData.SampleReservation(), Description: "Update reservation"},
		{Method: "DELETE", Path: "/reservation/" + suite.testData.ReservationID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Delete reservation"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Reservation Routes")
}

// Test Customer Routes
func (suite *APITestSuite) TestCustomerRoutes() {
	routes := []TestRoute{
		{Method: "GET", Path: "/customer/" + suite.testData.CustomerID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Get customer by ID"},
		{Method: "GET", Path: "/customer", ExpectedStatus: 200, RequiresHeaders: true, Description: "List customers"},
		{Method: "POST", Path: "/customer", ExpectedStatus: 201, RequiresHeaders: true, TestBody: suite.testData.SampleCustomer(), Description: "Create customer"},
		{Method: "PUT", Path: "/customer/" + suite.testData.CustomerID, ExpectedStatus: 200, RequiresHeaders: true, TestBody: suite.testData.SampleCustomer(), Description: "Update customer"},
		{Method: "DELETE", Path: "/customer/" + suite.testData.CustomerID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Delete customer"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Customer Routes")
}

// Test Waitlist Routes
func (suite *APITestSuite) TestWaitlistRoutes() {
	routes := []TestRoute{
		{Method: "GET", Path: "/waitlist/" + suite.testData.WaitlistID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Get waitlist by ID"},
		{Method: "GET", Path: "/waitlist", ExpectedStatus: 200, RequiresHeaders: true, Description: "List waitlist"},
		{Method: "POST", Path: "/waitlist", ExpectedStatus: 201, RequiresHeaders: true, TestBody: suite.testData.SampleWaitlist(), Description: "Create waitlist entry"},
		{Method: "PUT", Path: "/waitlist/" + suite.testData.WaitlistID, ExpectedStatus: 200, RequiresHeaders: true, TestBody: suite.testData.SampleWaitlist(), Description: "Update waitlist entry"},
		{Method: "DELETE", Path: "/waitlist/" + suite.testData.WaitlistID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Delete waitlist entry"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Waitlist Routes")
}

// Test Project Routes
func (suite *APITestSuite) TestProjectRoutes() {
	routes := []TestRoute{
		{Method: "GET", Path: "/project/" + suite.testData.ProjectID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Get project by ID"},
		{Method: "GET", Path: "/project", ExpectedStatus: 200, RequiresHeaders: true, Description: "List projects"},
		{Method: "GET", Path: "/project/active", ExpectedStatus: 200, RequiresHeaders: true, Description: "List active projects"},
		{Method: "POST", Path: "/project", ExpectedStatus: 201, RequiresHeaders: true, TestBody: suite.testData.SampleProject(), Description: "Create project"},
		{Method: "PUT", Path: "/project/" + suite.testData.ProjectID, ExpectedStatus: 200, RequiresHeaders: true, TestBody: suite.testData.SampleProject(), Description: "Update project"},
		{Method: "DELETE", Path: "/project/" + suite.testData.ProjectID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Delete project"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Project Routes")
}

// Test Settings Routes
func (suite *APITestSuite) TestSettingsRoutes() {
	routes := []TestRoute{
		{Method: "GET", Path: "/settings", ExpectedStatus: 200, RequiresHeaders: true, Description: "Get settings"},
		{Method: "PUT", Path: "/settings", ExpectedStatus: 200, RequiresHeaders: true, TestBody: suite.testData.SampleSettings(), Description: "Update settings"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Settings Routes")
}

// Test Environment Routes
func (suite *APITestSuite) TestEnvironmentRoutes() {
	routes := []TestRoute{
		{Method: "GET", Path: "/environment/env-id", ExpectedStatus: 200, RequiresHeaders: true, Description: "Get environment by ID"},
		{Method: "GET", Path: "/environment", ExpectedStatus: 200, RequiresHeaders: true, Description: "List environments"},
		{Method: "GET", Path: "/environment/active", ExpectedStatus: 200, RequiresHeaders: true, Description: "List active environments"},
		{Method: "POST", Path: "/environment", ExpectedStatus: 201, RequiresHeaders: true, TestBody: suite.testData.SampleEnvironment(), Description: "Create environment"},
		{Method: "PUT", Path: "/environment/env-id", ExpectedStatus: 200, RequiresHeaders: true, TestBody: suite.testData.SampleEnvironment(), Description: "Update environment"},
		{Method: "DELETE", Path: "/environment/env-id", ExpectedStatus: 200, RequiresHeaders: true, Description: "Delete environment"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Environment Routes")
}

// Test Organization Routes
func (suite *APITestSuite) TestOrganizationRoutes() {
	routes := []TestRoute{
		{Method: "GET", Path: "/organization/" + suite.testData.OrganizationID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Get organization by ID"},
		{Method: "GET", Path: "/organization", ExpectedStatus: 200, RequiresHeaders: true, Description: "List organizations"},
		{Method: "GET", Path: "/organization/active", ExpectedStatus: 200, RequiresHeaders: true, Description: "List active organizations"},
		{Method: "GET", Path: "/organization/email", ExpectedStatus: 200, RequiresHeaders: true, Description: "Get organization by email"},
		{Method: "POST", Path: "/organization", ExpectedStatus: 201, RequiresHeaders: true, TestBody: suite.testData.SampleOrganization(), Description: "Create organization"},
		{Method: "PUT", Path: "/organization/" + suite.testData.OrganizationID, ExpectedStatus: 200, RequiresHeaders: true, TestBody: suite.testData.SampleOrganization(), Description: "Update organization"},
		{Method: "DELETE", Path: "/organization/" + suite.testData.OrganizationID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Soft delete organization"},
		{Method: "DELETE", Path: "/organization/" + suite.testData.OrganizationID + "/permanent", ExpectedStatus: 200, RequiresHeaders: true, Description: "Hard delete organization"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Organization Routes")
}

// Test Reports Routes
func (suite *APITestSuite) TestReportsRoutes() {
	routes := []TestRoute{
		{Method: "GET", Path: "/reports/occupancy", ExpectedStatus: 200, RequiresHeaders: true, Description: "Get occupancy report"},
		{Method: "GET", Path: "/reports/reservations", ExpectedStatus: 200, RequiresHeaders: true, Description: "Get reservations report"},
		{Method: "GET", Path: "/reports/waitlist", ExpectedStatus: 200, RequiresHeaders: true, Description: "Get waitlist report"},
		{Method: "GET", Path: "/reports/leads", ExpectedStatus: 200, RequiresHeaders: true, Description: "Get leads report"},
		{Method: "GET", Path: "/reports/export/occupancy", ExpectedStatus: 200, RequiresHeaders: true, Description: "Export occupancy report"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Reports Routes")
}

// Test Notification Routes
func (suite *APITestSuite) TestNotificationRoutes() {
	routes := []TestRoute{
		{Method: "POST", Path: "/notification/send", ExpectedStatus: 200, RequiresHeaders: true, TestBody: map[string]string{"message": "test"}, Description: "Send notification"},
		{Method: "POST", Path: "/notification/event", ExpectedStatus: 200, RequiresHeaders: true, TestBody: map[string]string{"event": "test"}, Description: "Process notification event"},
		{Method: "GET", Path: "/notification/logs/" + suite.testData.OrganizationID + "/" + suite.testData.ProjectID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Get notification logs"},
		{Method: "GET", Path: "/notification/templates/" + suite.testData.OrganizationID + "/" + suite.testData.ProjectID, ExpectedStatus: 200, RequiresHeaders: true, Description: "Get notification templates"},
		{Method: "POST", Path: "/notification/template", ExpectedStatus: 201, RequiresHeaders: true, TestBody: suite.testData.SampleNotificationTemplate(), Description: "Create notification template"},
		{Method: "PUT", Path: "/notification/template", ExpectedStatus: 200, RequiresHeaders: true, TestBody: suite.testData.SampleNotificationTemplate(), Description: "Update notification template"},
		{Method: "POST", Path: "/notification/config", ExpectedStatus: 200, RequiresHeaders: true, TestBody: map[string]string{"config": "test"}, Description: "Update notification config"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Notification Routes")
}

// Test Webhook Routes
func (suite *APITestSuite) TestWebhookRoutes() {
	routes := []TestRoute{
		{Method: "POST", Path: "/webhook/twilio/status", ExpectedStatus: 200, TestBody: map[string]string{"status": "delivered"}, Description: "Twilio status webhook"},
		{Method: "POST", Path: "/webhook/twilio/inbound/" + suite.testData.OrganizationID + "/" + suite.testData.ProjectID, ExpectedStatus: 200, TestBody: map[string]string{"message": "test"}, Description: "Twilio inbound webhook"},
	}

	suite.helper.BatchRouteTest(suite.T(), routes, "Webhook Routes")
}

// TestAPITestSuite runs the entire test suite
func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}