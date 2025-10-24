package utils

import (
	"time"

	"github.com/google/uuid"
	"lep/repositories/models"
)

// Fattoria-specific IDs
var (
	FattoriaOrgID      = uuid.MustParse("223e4567-e89b-12d3-a456-426614174100")
	FattoriaProjectID  = uuid.MustParse("223e4567-e89b-12d3-a456-426614174101")
	FattoriaMenuID     = uuid.MustParse("223e4567-e89b-12d3-a456-426614174102")

	// Categories
	FattoriaCategoryPizzasID      = uuid.MustParse("223e4567-e89b-12d3-a456-426614174110")
	FattoriaCategoryBebidasID     = uuid.MustParse("223e4567-e89b-12d3-a456-426614174111")

	// Subcategories
	FattoriaSubcategoryEntradasID = uuid.MustParse("223e4567-e89b-12d3-a456-426614174120")
	FattoriaSubcategoryPizzasID   = uuid.MustParse("223e4567-e89b-12d3-a456-426614174121")
	FattoriaSubcategorySoftID     = uuid.MustParse("223e4567-e89b-12d3-a456-426614174122")
	FattoriaSubcategoryCervejasID = uuid.MustParse("223e4567-e89b-12d3-a456-426614174123")
	FattoriaSubcatCervArteID      = uuid.MustParse("223e4567-e89b-12d3-a456-426614174124")
	FattoriaSubcategoryCoqueisID  = uuid.MustParse("223e4567-e89b-12d3-a456-426614174125")

	// Tags
	FattoriaTagVegetarianaID = uuid.MustParse("223e4567-e89b-12d3-a456-426614174130")
	FattoriaTagVeganaID      = uuid.MustParse("223e4567-e89b-12d3-a456-426614174131")

	// Products
	FattoriaProductCrostiniID       = uuid.MustParse("223e4567-e89b-12d3-a456-426614174200")
	FattoriaProductMargueritaID     = uuid.MustParse("223e4567-e89b-12d3-a456-426614174201")
	FattoriaProductMarinaraID       = uuid.MustParse("223e4567-e89b-12d3-a456-426614174202")
	FattoriaProductParmaID          = uuid.MustParse("223e4567-e89b-12d3-a456-426614174203")
	FattoriaProductVeganaID         = uuid.MustParse("223e4567-e89b-12d3-a456-426614174204")
	FattoriaProductSucoID           = uuid.MustParse("223e4567-e89b-12d3-a456-426614174205")
	FattoriaProductBadenBadenID     = uuid.MustParse("223e4567-e89b-12d3-a456-426614174206")
	FattoriaProductSoniaZeID        = uuid.MustParse("223e4567-e89b-12d3-a456-426614174207")
	FattoriaProductHeinekeID        = uuid.MustParse("223e4567-e89b-12d3-a456-426614174208")

	// Environment
	FattoriaEnvironmentID = uuid.MustParse("223e4567-e89b-12d3-a456-426614174300")
)

// GenerateFattoriaData creates complete seed data for Fattoria restaurant
func GenerateFattoriaData() *SeedData {
	now := time.Now()

	return &SeedData{
		Organizations: []models.Organization{
			{
				Id:          FattoriaOrgID,
				Name:        "Fattoria Pizzeria",
				Email:       "contato@fattoria.com.br",
				Phone:       "+55 11 9999-9999",
				Address:     "Rua dos Italianos, 456 - São Paulo, SP",
				Website:     "https://fattoria.com.br",
				Description: "Autêntica pizzaria italiana com receitas tradicionais",
				Active:      true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},

		Projects: []models.Project{
			{
				Id:             FattoriaProjectID,
				OrganizationId: FattoriaOrgID,
				Name:           "Fattoria Pizzeria - Projeto Principal",
				Description:    "Gerenciamento operacional da Fattoria Pizzeria",
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Users: []models.User{
			{
				Id:          uuid.MustParse("223e4567-e89b-12d3-a456-426614174310"),
				Name:        "Admin Fattoria",
				Email:       "admin@fattoria.com.br",
				Password:    hashPassword("password"),
				Permissions: []string{"admin", "products", "orders", "reservations", "customers", "tables", "reports"},
				Active:      true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},

		UserOrganizations: []models.UserOrganization{
			{
				Id:             uuid.MustParse("223e4567-e89b-12d3-a456-426614174410"),
				UserId:         uuid.MustParse("223e4567-e89b-12d3-a456-426614174310"),
				OrganizationId: FattoriaOrgID,
				Role:           "owner",
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		UserProjects: []models.UserProject{
			{
				Id:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174510"),
				UserId:    uuid.MustParse("223e4567-e89b-12d3-a456-426614174310"),
				ProjectId: FattoriaProjectID,
				Role:      "admin",
				Active:    true,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},

		Environments: []models.Environment{
			{
				Id:             FattoriaEnvironmentID,
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				Name:           "Salão Principal",
				Description:    "Área principal da Fattoria Pizzeria",
				Capacity:       60,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Customers: []models.Customer{
			{
				Id:             uuid.MustParse("223e4567-e89b-12d3-a456-426614174610"),
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				Name:           "Cliente Fattoria 1",
				Email:          "cliente1@email.com",
				Phone:          "+55 11 98765-4321",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Menus: []models.Menu{
			{
				Id:             FattoriaMenuID,
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				Name:           "Cardápio Fattoria",
				Order:          1,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Categories: []models.Category{
			// Parent categories
			{
				Id:             FattoriaCategoryPizzasID,
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				MenuId:         FattoriaMenuID,
				Name:           "Pizzas",
				Order:          1,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				Id:             FattoriaCategoryBebidasID,
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				MenuId:         FattoriaMenuID,
				Name:           "Bebidas",
				Order:          2,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			// Subcategories (stored as categories)
			{
				Id:             FattoriaSubcategoryEntradasID,
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				MenuId:         FattoriaMenuID,
				Name:           "Entradas",
				Order:          1,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				Id:             FattoriaSubcategoryPizzasID,
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				MenuId:         FattoriaMenuID,
				Name:           "Pizzas",
				Order:          2,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				Id:             FattoriaSubcategorySoftID,
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				MenuId:         FattoriaMenuID,
				Name:           "Soft drinks",
				Order:          3,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				Id:             FattoriaSubcategoryCervejasID,
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				MenuId:         FattoriaMenuID,
				Name:           "Cervejas",
				Order:          4,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				Id:             FattoriaSubcatCervArteID,
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				MenuId:         FattoriaMenuID,
				Name:           "Cervejas artesanais",
				Order:          5,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				Id:             FattoriaSubcategoryCoqueisID,
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				MenuId:         FattoriaMenuID,
				Name:           "Coquetéis",
				Order:          6,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Tags: []models.Tag{
			{
				Id:             FattoriaTagVegetarianaID,
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				Name:           "Vegetariana",
				Color:          "#4CAF50",
				Description:    "Prato vegetariano",
				EntityType:     "product",
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				Id:             FattoriaTagVeganaID,
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				Name:           "Vegana",
				Color:          "#8BC34A",
				Description:    "Prato vegano",
				EntityType:     "product",
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Products: []models.Product{
			// Entradas
			{
				Id:              FattoriaProductCrostiniID,
				OrganizationId:  FattoriaOrgID,
				ProjectId:       FattoriaProjectID,
				Name:            "Crostini",
				Description:     "Massa fina levemente crocante com alecrim, parmesão e azeite",
				Type:            "prato",
				CategoryId:      &FattoriaSubcategoryEntradasID,
				PriceNormal:     30.00,
				Active:          true,
				Order:           1,
				PrepTimeMinutes: 15,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			// Pizzas
			{
				Id:              FattoriaProductMargueritaID,
				OrganizationId:  FattoriaOrgID,
				ProjectId:       FattoriaProjectID,
				Name:            "Marguerita",
				Description:     "Molho pomodoro, mussarela de búfala, manjericão fresco, azeite de oliva e orégano",
				Type:            "prato",
				CategoryId:      &FattoriaSubcategoryPizzasID,
				PriceNormal:     80.00,
				Active:          true,
				Order:           2,
				PrepTimeMinutes: 25,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			{
				Id:              FattoriaProductMarinaraID,
				OrganizationId:  FattoriaOrgID,
				ProjectId:       FattoriaProjectID,
				Name:            "Marinara",
				Description:     "Molho pomodoro, alho em lascas, azeite de oliva e orégano",
				Type:            "prato",
				CategoryId:      &FattoriaSubcategoryPizzasID,
				PriceNormal:     58.00,
				Active:          true,
				Order:           3,
				PrepTimeMinutes: 25,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			{
				Id:              FattoriaProductParmaID,
				OrganizationId:  FattoriaOrgID,
				ProjectId:       FattoriaProjectID,
				Name:            "Parma",
				Description:     "Molho pomodoro, mussarela, tomate seco, orégano, parmesão, presunto parma e rúcula",
				Type:            "prato",
				CategoryId:      &FattoriaSubcategoryPizzasID,
				PriceNormal:     109.00,
				Active:          true,
				Order:           4,
				PrepTimeMinutes: 25,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			{
				Id:              FattoriaProductVeganaID,
				OrganizationId:  FattoriaOrgID,
				ProjectId:       FattoriaProjectID,
				Name:            "Vegana",
				Description:     "Molho pomodoro, tomate confit, alho em lascas, azeitona preta, manjericão fresco e orégano",
				Type:            "prato",
				CategoryId:      &FattoriaSubcategoryPizzasID,
				PriceNormal:     60.00,
				Active:          true,
				Order:           5,
				PrepTimeMinutes: 25,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			// Bebidas - Soft drinks
			{
				Id:              FattoriaProductSucoID,
				OrganizationId:  FattoriaOrgID,
				ProjectId:       FattoriaProjectID,
				Name:            "Suco de caju integral",
				Description:     "Suco natural de caju integral",
				Type:            "bebida",
				CategoryId:      &FattoriaSubcategorySoftID,
				PriceNormal:     15.00,
				Active:          true,
				Order:           6,
				Volume:          func(i int) *int { return &i }(300),
				PrepTimeMinutes: 2,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			// Bebidas - Cervejas
			{
				Id:              FattoriaProductHeinekeID,
				OrganizationId:  FattoriaOrgID,
				ProjectId:       FattoriaProjectID,
				Name:            "Heineken s/ álcool",
				Description:     "Cerveja Heineken sem álcool",
				Type:            "bebida",
				CategoryId:      &FattoriaSubcategoryCervejasID,
				PriceNormal:     13.00,
				Active:          true,
				Order:           7,
				Volume:          func(i int) *int { return &i }(330),
				PrepTimeMinutes: 2,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			// Bebidas - Cervejas artesanais
			{
				Id:              FattoriaProductBadenBadenID,
				OrganizationId:  FattoriaOrgID,
				ProjectId:       FattoriaProjectID,
				Name:            "Baden Baden IPA",
				Description:     "Cerveja artesanal Baden Baden IPA",
				Type:            "bebida",
				CategoryId:      &FattoriaSubcatCervArteID,
				PriceNormal:     23.00,
				Active:          true,
				Order:           8,
				Volume:          func(i int) *int { return &i }(600),
				PrepTimeMinutes: 2,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			// Bebidas - Coquetéis
			{
				Id:              FattoriaProductSoniaZeID,
				OrganizationId:  FattoriaOrgID,
				ProjectId:       FattoriaProjectID,
				Name:            "Sônia e Zé",
				Description:     "Suco de limão siciliano, Ramazzotti, cachaça Dom Drinks, Monin de flor de sabugueiro e manjericão para decorar",
				Type:            "bebida",
				CategoryId:      &FattoriaSubcategoryCoqueisID,
				PriceNormal:     32.00,
				Active:          true,
				Order:           9,
				PrepTimeMinutes: 5,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
		},

		ProductTags: []models.ProductTag{
			// Marguerita - Vegetariana
			{
				Id:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174700"),
				ProductId: FattoriaProductMargueritaID,
				TagId:     FattoriaTagVegetarianaID,
				CreatedAt: now,
			},
			// Marinara - Vegana
			{
				Id:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174701"),
				ProductId: FattoriaProductMarinaraID,
				TagId:     FattoriaTagVeganaID,
				CreatedAt: now,
			},
			// Vegana - Vegana
			{
				Id:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174702"),
				ProductId: FattoriaProductVeganaID,
				TagId:     FattoriaTagVeganaID,
				CreatedAt: now,
			},
		},

		Tables: []models.Table{
			{
				Id:            uuid.MustParse("223e4567-e89b-12d3-a456-426614174800"),
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				EnvironmentId:  &FattoriaEnvironmentID,
				Number:        1,
				Capacity:      4,
				Status:        "livre",
				Location:      "Salão Principal - Entrada",
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			{
				Id:            uuid.MustParse("223e4567-e89b-12d3-a456-426614174801"),
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				EnvironmentId:  &FattoriaEnvironmentID,
				Number:        2,
				Capacity:      2,
				Status:        "livre",
				Location:      "Salão Principal - Janela",
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			{
				Id:            uuid.MustParse("223e4567-e89b-12d3-a456-426614174802"),
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				EnvironmentId:  &FattoriaEnvironmentID,
				Number:        3,
				Capacity:      6,
				Status:        "livre",
				Location:      "Salão Principal - Fundo",
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		},

		Orders:       []models.Order{},
		Reservations: []models.Reservation{},
		Waitlists:    []models.Waitlist{},

		Settings: []models.Settings{
			{
				Id:             uuid.MustParse("223e4567-e89b-12d3-a456-426614174900"),
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},

		Templates: []models.NotificationTemplate{
			{
				Id:             uuid.MustParse("223e4567-e89b-12d3-a456-426614174901"),
				OrganizationId: FattoriaOrgID,
				ProjectId:      FattoriaProjectID,
				Channel:        "sms",
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
	}
}
