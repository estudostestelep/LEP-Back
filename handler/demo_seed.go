package handler

import (
	"errors"
	"fmt"
	"log"
	"time"

	"lep/repositories/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// isIgnorableError verifica se o erro é esperado (duplicate key, not found)
// Esses erros não devem abortar a transação, apenas logar um aviso
func isIgnorableError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return true
		case "23503": // foreign_key_violation
			return true
		}
	}
	return false
}

// SeedDemoOrganization cria a organização demo com todos os dados de exemplo
// Esta função é chamada no startup da aplicação e verifica se já existe
func SeedDemoOrganization(db *gorm.DB) error {
	// Verificar se a organização demo já existe
	var existingOrg models.Organization
	if err := db.Where("slug = ?", "demo").First(&existingOrg).Error; err == nil {
		log.Println("✅ Organização demo já existe, pulando seed")
		return nil
	}

	log.Println("🌱 Iniciando seed da organização demo...")

	// IDs que serão usados em múltiplas entidades
	orgId := uuid.New()
	projectId := uuid.New()
	envSalaoId := uuid.New()
	envVarandaId := uuid.New()
	menuId := uuid.New()
	now := time.Now()

	// Hash da senha padrão "password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	// Executar tudo em uma transação
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Criar Organização
		org := &models.Organization{
			Id:          orgId,
			Name:        "LEP Demo",
			Slug:        "demo",
			Email:       "contato@lep-demo.com",
			Phone:       "(11) 99999-9999",
			Address:     "Rua Exemplo, 123 - São Paulo, SP",
			Description: "Organização de demonstração do sistema LEP",
			Active:      true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := tx.Create(org).Error; err != nil {
			return fmt.Errorf("erro ao criar organização: %w", err)
		}
		log.Println("  ✓ Organização criada")

		// 2. Atribuir plano "demo"
		var demoPlan models.Plan
		if err := tx.Where("code = ?", "demo").First(&demoPlan).Error; err != nil {
			log.Printf("  ⚠️ Plano demo não encontrado, usando free")
			if err := tx.Where("code = ?", "free").First(&demoPlan).Error; err != nil {
				log.Printf("  ⚠️ Nenhum plano encontrado")
			}
		}
		if demoPlan.Id != uuid.Nil {
			tx.SavePoint("sp_orgplan")
			orgPlan := &models.OrganizationPlan{
				Id:             uuid.New(),
				OrganizationId: orgId,
				PlanId:         demoPlan.Id,
				BillingCycle:   "monthly",
				Active:         true,
				StartsAt:       &now,
			}
			if err := tx.Create(orgPlan).Error; err != nil {
				if isIgnorableError(err) {
					log.Printf("  ⚠️ Plano já atribuído, pulando")
					tx.RollbackTo("sp_orgplan")
				} else {
					return fmt.Errorf("erro crítico ao atribuir plano: %w", err)
				}
			} else {
				log.Println("  ✓ Plano demo atribuído")
			}
		}

		// 3. Criar Projeto
		project := &models.Project{
			Id:             projectId,
			OrganizationId: orgId,
			Name:           "Restaurante Demo",
			Description:    "Restaurante de demonstração",
			Slug:           "restaurante-demo",
			IsDefault:      true,
			TimeZone:       "America/Sao_Paulo",
			Active:         true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := tx.Create(project).Error; err != nil {
			return fmt.Errorf("erro ao criar projeto: %w", err)
		}
		log.Println("  ✓ Projeto criado")

		// 4. Criar Settings
		tx.SavePoint("sp_settings")
		settings := &models.Settings{
			Id:                         uuid.New(),
			OrganizationId:             orgId,
			ProjectId:                  projectId,
			MinAdvanceHours:            2,
			MaxAdvanceDays:             30,
			NotifyReservationCreate:    true,
			NotifyReservationUpdate:    true,
			NotifyReservationCancel:    true,
			NotifyTableAvailable:       true,
			NotifyConfirmation24h:      true,
			ConfirmationHoursBefore:    24,
			DefaultNotificationChannel: "sms",
			EnableSms:                  true,
			CreatedAt:                  now,
			UpdatedAt:                  now,
		}
		if err := tx.Create(settings).Error; err != nil {
			if isIgnorableError(err) {
				log.Printf("  ⚠️ Settings já existe, pulando")
				tx.RollbackTo("sp_settings")
			} else {
				return fmt.Errorf("erro crítico ao criar settings: %w", err)
			}
		} else {
			log.Println("  ✓ Settings criado")
		}

		// 5. Criar Usuários (Clients)
		// Buscar roles
		var ownerRole, managerRole, waiterRole models.Role
		tx.Where("name = ?", "owner").First(&ownerRole)
		tx.Where("name = ?", "manager").First(&managerRole)
		tx.Where("name = ?", "waiter").First(&waiterRole)

		users := []struct {
			name  string
			email string
			role  models.Role
		}{
			{"Administrador Demo", "admin@lep-demo.com", ownerRole},
			{"Gerente Demo", "gerente@lep-demo.com", managerRole},
			{"Garçom Demo", "garcom@lep-demo.com", waiterRole},
		}

		usersCreated := 0
		for i, u := range users {
			spName := fmt.Sprintf("sp_user_%d", i)
			tx.SavePoint(spName)

			clientId := uuid.New()
			client := &models.Client{
				Id:        clientId,
				Name:      u.name,
				Email:     u.email,
				Password:  string(hashedPassword),
				OrgId:     orgId,
				ProjIds:   pq.StringArray{projectId.String()},
				Active:    true,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := tx.Create(client).Error; err != nil {
				if isIgnorableError(err) {
					log.Printf("  ⚠️ Usuário %s já existe, pulando", u.email)
					tx.RollbackTo(spName)
					continue
				}
				return fmt.Errorf("erro crítico ao criar usuário %s: %w", u.email, err)
			}

			// Atribuir role
			if u.role.Id != uuid.Nil {
				spRole := fmt.Sprintf("sp_role_%d", i)
				tx.SavePoint(spRole)
				clientRole := &models.ClientRole{
					ClientId:       clientId,
					RoleId:         u.role.Id,
					OrganizationId: orgId,
					ProjectId:      &projectId,
					Active:         true,
					CreatedAt:      now,
					UpdatedAt:      now,
				}
				if err := tx.Omit("Id").Create(clientRole).Error; err != nil {
					if isIgnorableError(err) {
						log.Printf("  ⚠️ Role para %s já existe, pulando", u.email)
						tx.RollbackTo(spRole)
					} else {
						return fmt.Errorf("erro crítico ao atribuir role para %s: %w", u.email, err)
					}
				}
			}
			usersCreated++
		}
		log.Printf("  ✓ %d usuários criados", usersCreated)

		// 6. Criar Ambientes
		environments := []models.Environment{
			{
				Id:             envSalaoId,
				OrganizationId: orgId,
				ProjectId:      projectId,
				Name:           "Salão Principal",
				Description:    "Ambiente interno climatizado",
				Capacity:       60,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				Id:             envVarandaId,
				OrganizationId: orgId,
				ProjectId:      projectId,
				Name:           "Varanda",
				Description:    "Área externa com vista",
				Capacity:       30,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		}
		envsCreated := 0
		for i, env := range environments {
			spName := fmt.Sprintf("sp_env_%d", i)
			tx.SavePoint(spName)
			if err := tx.Create(&env).Error; err != nil {
				if isIgnorableError(err) {
					log.Printf("  ⚠️ Ambiente %s já existe, pulando", env.Name)
					tx.RollbackTo(spName)
					continue
				}
				return fmt.Errorf("erro crítico ao criar ambiente %s: %w", env.Name, err)
			}
			envsCreated++
		}
		log.Printf("  ✓ %d ambientes criados", envsCreated)

		// 7. Criar Mesas
		tablesCreated := 0
		// Mesas do Salão (1-6)
		for i := 1; i <= 6; i++ {
			spName := fmt.Sprintf("sp_table_%d", i)
			tx.SavePoint(spName)
			table := &models.Table{
				Id:             uuid.New(),
				OrganizationId: orgId,
				ProjectId:      projectId,
				EnvironmentId:  &envSalaoId,
				Number:         i,
				Capacity:       4,
				Location:       "Salão Principal",
				Status:         "livre",
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			if err := tx.Create(table).Error; err != nil {
				if isIgnorableError(err) {
					log.Printf("  ⚠️ Mesa %d já existe, pulando", i)
					tx.RollbackTo(spName)
					continue
				}
				return fmt.Errorf("erro crítico ao criar mesa %d: %w", i, err)
			}
			tablesCreated++
		}
		// Mesas da Varanda (7-10)
		for i := 7; i <= 10; i++ {
			spName := fmt.Sprintf("sp_table_%d", i)
			tx.SavePoint(spName)
			table := &models.Table{
				Id:             uuid.New(),
				OrganizationId: orgId,
				ProjectId:      projectId,
				EnvironmentId:  &envVarandaId,
				Number:         i,
				Capacity:       6,
				Location:       "Varanda",
				Status:         "livre",
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			if err := tx.Create(table).Error; err != nil {
				if isIgnorableError(err) {
					log.Printf("  ⚠️ Mesa %d já existe, pulando", i)
					tx.RollbackTo(spName)
					continue
				}
				return fmt.Errorf("erro crítico ao criar mesa %d: %w", i, err)
			}
			tablesCreated++
		}
		log.Printf("  ✓ %d mesas criadas", tablesCreated)

		// 8. Criar Menu
		tx.SavePoint("sp_menu")
		menu := &models.Menu{
			Id:             menuId,
			OrganizationId: orgId,
			ProjectId:      projectId,
			Name:           "Cardápio Principal",
			Order:          1,
			Active:         true,
			Priority:       1,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := tx.Create(menu).Error; err != nil {
			if isIgnorableError(err) {
				log.Printf("  ⚠️ Menu já existe, pulando")
				tx.RollbackTo("sp_menu")
			} else {
				return fmt.Errorf("erro crítico ao criar menu: %w", err)
			}
		} else {
			log.Println("  ✓ Menu criado")
		}

		// 9. Criar Categorias
		categoryIds := make(map[string]uuid.UUID)
		categories := []struct {
			name  string
			order int
		}{
			{"Entradas", 1},
			{"Pratos Principais", 2},
			{"Bebidas", 3},
			{"Sobremesas", 4},
		}
		catsCreated := 0
		for i, cat := range categories {
			spName := fmt.Sprintf("sp_cat_%d", i)
			tx.SavePoint(spName)
			catId := uuid.New()
			categoryIds[cat.name] = catId
			category := &models.Category{
				Id:             catId,
				OrganizationId: orgId,
				ProjectId:      projectId,
				MenuId:         menuId,
				Name:           cat.name,
				Order:          cat.order,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			if err := tx.Create(category).Error; err != nil {
				if isIgnorableError(err) {
					log.Printf("  ⚠️ Categoria %s já existe, pulando", cat.name)
					tx.RollbackTo(spName)
					continue
				}
				return fmt.Errorf("erro crítico ao criar categoria %s: %w", cat.name, err)
			}
			catsCreated++
		}
		log.Printf("  ✓ %d categorias criadas", catsCreated)

		// 10. Criar Produtos
		products := []struct {
			name        string
			description string
			price       float64
			category    string
			prodType    string
		}{
			// Entradas
			{"Bruschetta", "Fatias de pão italiano com tomate e manjericão", 28.90, "Entradas", "prato"},
			{"Carpaccio", "Finas fatias de carne com alcaparras e parmesão", 42.90, "Entradas", "prato"},
			{"Salada Caesar", "Alface romana, croutons, parmesão e molho Caesar", 32.90, "Entradas", "prato"},

			// Pratos Principais
			{"Filé Mignon ao Molho Madeira", "Filé mignon grelhado com molho madeira", 89.90, "Pratos Principais", "prato"},
			{"Salmão Grelhado", "Salmão com ervas finas e legumes", 79.90, "Pratos Principais", "prato"},
			{"Risoto de Cogumelos", "Risoto cremoso com mix de cogumelos", 62.90, "Pratos Principais", "prato"},
			{"Massa ao Pesto", "Fettuccine ao pesto genovese", 48.90, "Pratos Principais", "prato"},

			// Bebidas
			{"Suco Natural", "Laranja, limão, abacaxi ou maracujá", 12.90, "Bebidas", "bebida"},
			{"Refrigerante", "Coca-Cola, Guaraná ou Sprite (350ml)", 8.90, "Bebidas", "bebida"},
			{"Água Mineral", "Com ou sem gás (500ml)", 6.90, "Bebidas", "bebida"},
			{"Cerveja Artesanal", "IPA, Pilsen ou Weiss (500ml)", 24.90, "Bebidas", "bebida"},

			// Sobremesas
			{"Petit Gâteau", "Bolo de chocolate com sorvete", 32.90, "Sobremesas", "prato"},
			{"Tiramisù", "Clássico italiano", 28.90, "Sobremesas", "prato"},
			{"Cheesecake", "Com calda de frutas vermelhas", 26.90, "Sobremesas", "prato"},
		}
		prodsCreated := 0
		for i, p := range products {
			spName := fmt.Sprintf("sp_prod_%d", i)
			tx.SavePoint(spName)
			catId := categoryIds[p.category]
			product := &models.Product{
				Id:             uuid.New(),
				OrganizationId: orgId,
				ProjectId:      projectId,
				Name:           p.name,
				Description:    p.description,
				Type:           p.prodType,
				CategoryId:     &catId,
				PriceNormal:    p.price,
				Order:          i + 1,
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			if err := tx.Create(product).Error; err != nil {
				if isIgnorableError(err) {
					log.Printf("  ⚠️ Produto %s já existe, pulando", p.name)
					tx.RollbackTo(spName)
					continue
				}
				return fmt.Errorf("erro crítico ao criar produto %s: %w", p.name, err)
			}
			prodsCreated++
		}
		log.Printf("  ✓ %d produtos criados", prodsCreated)

		// 11. Criar Clientes exemplo
		customers := []struct {
			name  string
			email string
			phone string
		}{
			{"Maria Silva", "maria@exemplo.com", "(11) 98765-4321"},
			{"João Santos", "joao@exemplo.com", "(11) 91234-5678"},
			{"Ana Costa", "ana@exemplo.com", "(11) 99876-5432"},
		}
		custsCreated := 0
		for i, c := range customers {
			spName := fmt.Sprintf("sp_cust_%d", i)
			tx.SavePoint(spName)
			customer := &models.Customer{
				Id:             uuid.New(),
				OrganizationId: orgId,
				ProjectId:      projectId,
				Name:           c.name,
				Email:          c.email,
				Phone:          c.phone,
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			if err := tx.Create(customer).Error; err != nil {
				if isIgnorableError(err) {
					log.Printf("  ⚠️ Cliente %s já existe, pulando", c.name)
					tx.RollbackTo(spName)
					continue
				}
				return fmt.Errorf("erro crítico ao criar cliente %s: %w", c.name, err)
			}
			custsCreated++
		}
		log.Printf("  ✓ %d clientes criados", custsCreated)

		log.Println("✅ Seed da organização demo concluído!")
		log.Println("")
		log.Println("📋 Credenciais de login:")
		log.Println("   Admin:   admin@lep-demo.com / password")
		log.Println("   Gerente: gerente@lep-demo.com / password")
		log.Println("   Garçom:  garcom@lep-demo.com / password")

		return nil
	})
}
