package tests

import (
	"time"
)

// TestData contains all test fixtures and sample data
type TestData struct {
	OrganizationID string
	ProjectID      string
	UserID         string
	CustomerID     string
	ProductID      string
	TableID        string
	OrderID        string
	ReservationID  string
	WaitlistID     string
}

// NewTestData creates a new set of test data with consistent UUIDs
func NewTestData() *TestData {
	return &TestData{
		OrganizationID: "123e4567-e89b-12d3-a456-426614174000",
		ProjectID:      "123e4567-e89b-12d3-a456-426614174001",
		UserID:         "123e4567-e89b-12d3-a456-426614174002",
		CustomerID:     "123e4567-e89b-12d3-a456-426614174003",
		ProductID:      "123e4567-e89b-12d3-a456-426614174004",
		TableID:        "123e4567-e89b-12d3-a456-426614174005",
		OrderID:        "123e4567-e89b-12d3-a456-426614174006",
		ReservationID:  "123e4567-e89b-12d3-a456-426614174007",
		WaitlistID:     "123e4567-e89b-12d3-a456-426614174008",
	}
}

// Sample data structures for testing

// SampleOrganization returns a sample organization for testing
func (td *TestData) SampleOrganization() map[string]interface{} {
	return map[string]interface{}{
		"name":        "LEP Restaurante Teste",
		"email":       "teste@lep-restaurante.com",
		"phone":       "+55 11 9999-8888",
		"address":     "Rua Augusta, 123 - São Paulo, SP",
		"website":     "https://lep-restaurante.com",
		"description": "Restaurante de teste para LEP System",
		"active":      true,
	}
}

// SampleUser returns a sample user for testing
func (td *TestData) SampleUser() map[string]interface{} {
	return map[string]interface{}{
		"organization_id": td.OrganizationID,
		"project_id":      td.ProjectID,
		"name":            "João Teste",
		"email":           "joao.teste@lep.com",
		"password":        "senha123",
		"role":            "waiter",
		"permissions":     []string{"orders", "tables", "customers"},
	}
}

// SampleCustomer returns a sample customer for testing
func (td *TestData) SampleCustomer() map[string]interface{} {
	return map[string]interface{}{
		"organization_id": td.OrganizationID,
		"project_id":      td.ProjectID,
		"name":            "Maria Silva",
		"email":           "maria.silva@email.com",
		"phone":           "+55 11 98765-4321",
		"birth_date":      "1990-05-15",
		"preferences":     "Vegetariana",
	}
}

// SampleProduct returns a sample product for testing
func (td *TestData) SampleProduct() map[string]interface{} {
	return map[string]interface{}{
		"organization_id":   td.OrganizationID,
		"project_id":        td.ProjectID,
		"name":              "Lasanha Bolonhesa",
		"description":       "Deliciosa lasanha com molho bolonhesa caseiro",
		"price":             42.90,
		"category":          "Pratos Principais",
		"prep_time_minutes": 25,
		"available":         true,
		"ingredients":       []string{"massa", "carne moída", "molho de tomate", "queijo"},
	}
}

// SampleTable returns a sample table for testing
func (td *TestData) SampleTable() map[string]interface{} {
	return map[string]interface{}{
		"organization_id": td.OrganizationID,
		"project_id":      td.ProjectID,
		"number":          1,
		"capacity":        4,
		"status":          "livre",
		"location":        "Salão Principal",
		"environment_id":  td.ProjectID, // Using project ID as environment for simplicity
	}
}

// SampleOrder returns a sample order for testing
func (td *TestData) SampleOrder() map[string]interface{} {
	return map[string]interface{}{
		"organization_id": td.OrganizationID,
		"project_id":      td.ProjectID,
		"customer_id":     td.CustomerID,
		"table_id":        td.TableID,
		"status":          "preparing",
		"total_amount":    42.90,
		"source":          "internal",
		"items": []map[string]interface{}{
			{
				"product_id": td.ProductID,
				"quantity":   1,
				"price":      42.90,
				"notes":      "Sem cebola",
			},
		},
	}
}

// SampleReservation returns a sample reservation for testing
func (td *TestData) SampleReservation() map[string]interface{} {
	futureDate := time.Now().Add(24 * time.Hour)
	return map[string]interface{}{
		"organization_id":  td.OrganizationID,
		"project_id":       td.ProjectID,
		"customer_id":      td.CustomerID,
		"table_id":         td.TableID,
		"datetime":         futureDate.Format(time.RFC3339),
		"party_size":       4,
		"status":           "confirmed",
		"special_requests": "Mesa próxima à janela",
	}
}

// SampleWaitlist returns a sample waitlist entry for testing
func (td *TestData) SampleWaitlist() map[string]interface{} {
	return map[string]interface{}{
		"organization_id":  td.OrganizationID,
		"project_id":       td.ProjectID,
		"customer_id":      td.CustomerID,
		"party_size":       2,
		"estimated_wait":   30,
		"status":           "waiting",
		"special_requests": "Mesa para duas pessoas",
	}
}

// SampleProject returns a sample project for testing
func (td *TestData) SampleProject() map[string]interface{} {
	return map[string]interface{}{
		"organization_id": td.OrganizationID,
		"name":            "Projeto Teste",
		"description":     "Projeto de teste para LEP System",
		"active":          true,
	}
}

// SampleSettings returns sample settings for testing
func (td *TestData) SampleSettings() map[string]interface{} {
	return map[string]interface{}{
		"organization_id":           td.OrganizationID,
		"project_id":                td.ProjectID,
		"reservation_advance_hours": 24,
		"max_party_size":            12,
		"default_reservation_time":  120,
		"allow_overbooking":         false,
	}
}

// SampleEnvironment returns a sample environment for testing
func (td *TestData) SampleEnvironment() map[string]interface{} {
	return map[string]interface{}{
		"organization_id": td.OrganizationID,
		"project_id":      td.ProjectID,
		"name":            "Salão Principal",
		"description":     "Área principal do restaurante",
		"capacity":        50,
		"active":          true,
	}
}

// SampleNotificationTemplate returns a sample notification template
func (td *TestData) SampleNotificationTemplate() map[string]interface{} {
	return map[string]interface{}{
		"organization_id": td.OrganizationID,
		"project_id":      td.ProjectID,
		"event_type":      "reservation_confirmed",
		"channel":         "sms",
		"template":        "Olá {{nome}}, sua reserva para {{data}} às {{hora}} foi confirmada!",
		"active":          true,
	}
}

// InvalidData returns various invalid data scenarios for testing
type InvalidData struct{}

// InvalidUUID returns an invalid UUID for testing
func (InvalidData) InvalidUUID() string {
	return "invalid-uuid"
}

// EmptyString returns empty string for required fields
func (InvalidData) EmptyString() string {
	return ""
}

// InvalidEmail returns invalid email for testing
func (InvalidData) InvalidEmail() string {
	return "invalid-email"
}

// InvalidPhone returns invalid phone for testing
func (InvalidData) InvalidPhone() string {
	return "123"
}

// NegativeNumber returns negative number for testing
func (InvalidData) NegativeNumber() float64 {
	return -1.0
}

// TooLongString returns a string that exceeds typical field limits
func (InvalidData) TooLongString() string {
	return string(make([]byte, 1000)) // 1000 character string
}

// GetInvalidData returns an instance of invalid data scenarios
func GetInvalidData() InvalidData {
	return InvalidData{}
}
