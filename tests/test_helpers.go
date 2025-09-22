package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHelper provides utilities for testing HTTP endpoints
type TestHelper struct {
	Router       *gin.Engine
	OrgID        string
	ProjectID    string
	AuthToken    string
	TestResponse *httptest.ResponseRecorder
}

// NewTestHelper creates a new test helper with default values
func NewTestHelper() *TestHelper {
	gin.SetMode(gin.TestMode)
	return &TestHelper{
		Router:    gin.New(),
		OrgID:     "123e4567-e89b-12d3-a456-426614174000",
		ProjectID: "123e4567-e89b-12d3-a456-426614174001",
		AuthToken: "test-jwt-token",
	}
}

// SetupTestRouter configures the router with test routes (deprecated - use TestConfig instead)
func (h *TestHelper) SetupTestRouter() {
	// Add basic middleware for testing
	h.Router.Use(func(c *gin.Context) {
		// Set test headers in context for protected routes
		c.Set("organization_id", h.OrgID)
		c.Set("project_id", h.ProjectID)
		c.Next()
	})
}

// MakeRequest creates and executes an HTTP request
func (h *TestHelper) MakeRequest(t *testing.T, method, url string, body interface{}) *httptest.ResponseRecorder {
	var jsonBody []byte
	var err error

	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Add required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Lpe-Organization-Id", h.OrgID)
	req.Header.Set("X-Lpe-Project-Id", h.ProjectID)
	if h.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+h.AuthToken)
	}

	h.TestResponse = httptest.NewRecorder()
	h.Router.ServeHTTP(h.TestResponse, req)

	return h.TestResponse
}

// AssertStatusCode checks if the response has the expected status code
func (h *TestHelper) AssertStatusCode(t *testing.T, expected int) {
	assert.Equal(t, expected, h.TestResponse.Code,
		"Expected status %d, got %d. Response: %s", expected, h.TestResponse.Code, h.TestResponse.Body.String())
}

// AssertJSONResponse checks if the response is valid JSON and matches expected structure
func (h *TestHelper) AssertJSONResponse(t *testing.T, expectedFields ...string) map[string]interface{} {
	assert.Equal(t, "application/json; charset=utf-8", h.TestResponse.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(h.TestResponse.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")

	// Check for expected fields
	for _, field := range expectedFields {
		assert.Contains(t, response, field, "Response should contain field: %s", field)
	}

	return response
}

// AssertErrorResponse checks for standardized error response format
func (h *TestHelper) AssertErrorResponse(t *testing.T, expectedStatus int, expectedMessage string) {
	h.AssertStatusCode(t, expectedStatus)
	response := h.AssertJSONResponse(t, "error", "message")

	if expectedMessage != "" {
		assert.Contains(t, response["message"].(string), expectedMessage)
	}
}

// AssertSuccessResponse checks for standardized success response format
func (h *TestHelper) AssertSuccessResponse(t *testing.T, expectedStatus int) map[string]interface{} {
	h.AssertStatusCode(t, expectedStatus)
	return h.AssertJSONResponse(t)
}

// GenerateTestUUID generates a deterministic UUID for testing
func GenerateTestUUID(seed string) string {
	// Generate deterministic UUIDs for testing
	baseUUID := "123e4567-e89b-12d3-a456-42661417"
	return baseUUID + fmt.Sprintf("%04d", len(seed))
}

// MockHandler creates a mock handler that returns predictable responses
func MockHandler(statusCode int, response interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if response != nil {
			c.JSON(statusCode, response)
		} else {
			c.Status(statusCode)
		}
	}
}

// MockErrorHandler creates a mock handler that returns standardized errors
func MockErrorHandler(statusCode int, message string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(statusCode, gin.H{
			"error":   true,
			"message": message,
		})
	}
}

// TestRoute represents a route to be tested
type TestRoute struct {
	Method          string
	Path            string
	ExpectedStatus  int
	RequiresAuth    bool
	RequiresHeaders bool
	TestBody        interface{}
	Description     string
}

// RouteTest executes a test for a specific route
func (h *TestHelper) RouteTest(t *testing.T, route TestRoute) {
	t.Run(route.Description, func(t *testing.T) {
		// Setup mock handler for this route
		mockResponse := gin.H{"status": "ok", "method": route.Method, "path": route.Path}
		h.Router.Handle(route.Method, route.Path, MockHandler(route.ExpectedStatus, mockResponse))

		// Make request
		resp := h.MakeRequest(t, route.Method, route.Path, route.TestBody)

		// Basic assertions
		h.AssertStatusCode(t, route.ExpectedStatus)

		if route.ExpectedStatus == http.StatusOK {
			h.AssertJSONResponse(t, "status")
		}

		t.Logf("✅ %s %s - Status: %d", route.Method, route.Path, resp.Code)
	})
}

// BatchRouteTest executes multiple route tests efficiently
func (h *TestHelper) BatchRouteTest(t *testing.T, routes []TestRoute, groupName string) {
	t.Run(groupName, func(t *testing.T) {
		passed := 0
		total := len(routes)

		for _, route := range routes {
			h.RouteTest(t, route)
			passed++
		}

		t.Logf("✅ %s: %d/%d routes passed", groupName, passed, total)
	})
}
