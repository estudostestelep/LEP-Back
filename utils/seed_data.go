package utils

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"lep/repositories/models"
)

// SeedData contains all sample data for database seeding
type SeedData struct {
	Organizations []models.Organization
	Projects      []models.Project
	Users         []models.User
	Customers     []models.Customer
	Products      []models.Product
	Tables        []models.Table
	Orders        []models.Order
	Reservations  []models.Reservation
	Waitlists     []models.Waitlist
	Environments  []models.Environment
	Settings      []models.Settings
	Templates     []models.NotificationTemplate
}

// Base IDs for consistent relationships
var (
	SampleOrgID     = uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	SampleProjectID = uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
	SampleUserID1   = uuid.MustParse("123e4567-e89b-12d3-a456-426614174002")
	SampleUserID2   = uuid.MustParse("123e4567-e89b-12d3-a456-426614174003")
	SampleUserID3   = uuid.MustParse("123e4567-e89b-12d3-a456-426614174004")
)

// GenerateCompleteData creates a complete set of realistic sample data
func GenerateCompleteData() *SeedData {
	now := time.Now()

	return &SeedData{
		Organizations: []models.Organization{
			{
				Id:          SampleOrgID,
				Name:        "LEP Restaurante Demo",
				Email:       "teste331@gmail.com",
				Phone:       "+55 11 9999-8888",
				Address:     "Rua das Flores, 123 - São Paulo, SP",
				Website:     "https://lep-demo.com",
				Description: "Restaurante demo para demonstração do sistema LEP",
				Active:      true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},

		Projects: []models.Project{
			{
				Id:             SampleProjectID,
				OrganizationId: SampleOrgID,
				Name:           "Projeto Principal",
				Description:    "Projeto principal do restaurante LEP Demo",
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Users: []models.User{
			{
				Id:             SampleUserID1,
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				Name:           "Admin LEP",
				Email:          "teste221@gmail.com",
				Password:       "password", // password
				Role:           "admin",
				Permissions:    pq.StringArray{"admin"},
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				Id:             SampleUserID2,
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				Name:           "Garçom João",
				Email:          "garcom1@gmail.com",
				Password:       "password", // password
				Role:           "waiter",
				Permissions:    pq.StringArray{"orders", "tables", "customers"},
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				Id:             SampleUserID3,
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				Name:           "Gerente Maria",
				Email:          "gerente1@gmail.com",
				Password:       "password", // password
				Role:           "manager",
				Permissions:    pq.StringArray{"orders", "tables", "customers", "products", "reports"},
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Environments: []models.Environment{
			{
				Id:             uuid.MustParse("923e4567-e89b-12d3-a456-426614174001"),
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				Name:           "Salão Principal",
				Description:    "Área principal do restaurante",
				Capacity:       50,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Customers: []models.Customer{
			{
				Id:             uuid.MustParse("323e4567-e89b-12d3-a456-426614174001"),
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				Name:           "Maria Silva",
				Email:          "maria.silva1@email.com",
				Phone:          "+55 11 98765-4321",
				BirthDate:      "1990-05-15",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				Id:             uuid.MustParse("323e4567-e89b-12d3-a456-426614174002"),
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				Name:           "João Santos",
				Email:          "joao.santos1@email.com",
				Phone:          "+55 11 87654-3210",
				BirthDate:      "1985-11-22",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Products: []models.Product{
			{
				Id:              uuid.MustParse("423e4567-e89b-12d3-a456-426614174001"),
				OrganizationId:  SampleOrgID,
				ProjectId:       SampleProjectID,
				Name:            "Pizza Margherita",
				Description:     "Pizza clássica com molho de tomate, mussarela e manjericão",
				Price:           35.90,
				Category:        "Pizzas",
				Available:       true,
				PrepTimeMinutes: 20,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			{
				Id:              uuid.MustParse("423e4567-e89b-12d3-a456-426614174002"),
				OrganizationId:  SampleOrgID,
				ProjectId:       SampleProjectID,
				Name:            "Hambúrguer Clássico",
				Description:     "Hambúrguer com carne bovina, alface, tomate e queijo",
				Price:           28.50,
				Category:        "Hambúrgueres",
				Available:       true,
				PrepTimeMinutes: 15,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			{
				Id:              uuid.MustParse("423e4567-e89b-12d3-a456-426614174003"),
				OrganizationId:  SampleOrgID,
				ProjectId:       SampleProjectID,
				Name:            "Refrigerante Cola",
				Description:     "Refrigerante de cola gelado",
				Price:           5.50,
				Category:        "Bebidas",
				Available:       true,
				PrepTimeMinutes: 1,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
		},

		Tables: []models.Table{
			{
				Id:             uuid.MustParse("523e4567-e89b-12d3-a456-426614174001"),
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				EnvironmentId:  &[]uuid.UUID{uuid.MustParse("923e4567-e89b-12d3-a456-426614174001")}[0],
				Number:         1,
				Capacity:       4,
				Status:         "livre",
				Location:       "Salão Principal - Centro",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				Id:             uuid.MustParse("523e4567-e89b-12d3-a456-426614174002"),
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				EnvironmentId:  &[]uuid.UUID{uuid.MustParse("923e4567-e89b-12d3-a456-426614174001")}[0],
				Number:         2,
				Capacity:       2,
				Status:         "ocupada",
				Location:       "Salão Principal - Janela",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Orders: []models.Order{
			{
				Id:             uuid.MustParse("623e4567-e89b-12d3-a456-426614174001"),
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				CustomerId:     &[]uuid.UUID{uuid.MustParse("323e4567-e89b-12d3-a456-426614174001")}[0],
				TableId:        &[]uuid.UUID{uuid.MustParse("523e4567-e89b-12d3-a456-426614174002")}[0],
				Items: models.OrderItems{
					{
						ProductId: uuid.MustParse("423e4567-e89b-12d3-a456-426614174001"), // Pizza Margherita
						Quantity:  1,
						Price:     35.90,
						Notes:     "Massa fina",
					},
					{
						ProductId: uuid.MustParse("423e4567-e89b-12d3-a456-426614174003"), // Refrigerante
						Quantity:  2,
						Price:     5.50,
						Notes:     "Com gelo",
					},
				},
				TableNumber: func(i int) *int { return &i }(1),
				Source:      "internal",
				Status:      "preparing",
				TotalAmount: 46.90,
				CreatedAt:   now.Add(-30 * time.Minute),
				UpdatedAt:   now.Add(-15 * time.Minute),
			},
		},

		Reservations: []models.Reservation{
			{
				Id:             uuid.MustParse("723e4567-e89b-12d3-a456-426614174001"),
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				CustomerId:     uuid.MustParse("323e4567-e89b-12d3-a456-426614174001"),
				TableId:        uuid.MustParse("523e4567-e89b-12d3-a456-426614174001"),
				Datetime:       now.Add(2 * time.Hour).Format("2006-01-02 15:04:05"),
				PartySize:      4,
				Status:         "confirmed",
				Note:           "Mesa próxima à janela",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Waitlists: []models.Waitlist{
			{
				Id:             uuid.MustParse("823e4567-e89b-12d3-a456-426614174001"),
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				CustomerId:     uuid.MustParse("323e4567-e89b-12d3-a456-426614174002"),
				Status:         "waiting",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Settings: []models.Settings{
			{
				Id:             uuid.MustParse("a23e4567-e89b-12d3-a456-426614174001"),
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Templates: []models.NotificationTemplate{
			{
				Id:             uuid.MustParse("b23e4567-e89b-12d3-a456-426614174001"),
				OrganizationId: SampleOrgID,
				ProjectId:      SampleProjectID,
				Channel:        "sms",
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
	}
}
