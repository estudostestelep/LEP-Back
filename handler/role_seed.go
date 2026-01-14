package handler

import (
	"fmt"
	"lep/repositories"
	"lep/repositories/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedRolesAndPermissions popula os dados iniciais de módulos, permissões e cargos
func SeedRolesAndPermissions(db *gorm.DB) error {
	fmt.Println("🌱 Iniciando seed de roles e permissões...")

	// Criar repositórios
	moduleRepo := repositories.NewModuleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	packageRepo := repositories.NewPackageRepository(db)

	// 1. Criar Módulos
	modules := createModules()
	for _, m := range modules {
		existing, _ := moduleRepo.GetByCodeName(m.CodeName)
		if existing == nil {
			if err := moduleRepo.Create(&m); err != nil {
				fmt.Printf("⚠️ Erro ao criar módulo %s: %v\n", m.CodeName, err)
			} else {
				fmt.Printf("✅ Módulo criado: %s\n", m.DisplayName)
			}
		}
	}

	// 2. Criar Permissões
	permissions := createPermissions(moduleRepo)
	for _, p := range permissions {
		existing, _ := permissionRepo.GetByCodeName(p.CodeName)
		if existing == nil {
			if err := permissionRepo.Create(&p); err != nil {
				fmt.Printf("⚠️ Erro ao criar permissão %s: %v\n", p.CodeName, err)
			} else {
				fmt.Printf("✅ Permissão criada: %s\n", p.DisplayName)
			}
		}
	}

	// 3. Criar Cargos do Sistema
	roles := createSystemRoles()
	for _, r := range roles {
		existing, _ := roleRepo.GetByName(r.Name)
		if existing == nil {
			if err := roleRepo.Create(&r); err != nil {
				fmt.Printf("⚠️ Erro ao criar cargo %s: %v\n", r.Name, err)
			} else {
				fmt.Printf("✅ Cargo criado: %s (nível %d)\n", r.DisplayName, r.HierarchyLevel)
			}
		}
	}

	// 4. Criar Pacotes
	packages := createPackages()
	for _, pkg := range packages {
		existing, _ := packageRepo.GetByCodeName(pkg.CodeName)
		if existing == nil {
			if err := packageRepo.Create(&pkg); err != nil {
				fmt.Printf("⚠️ Erro ao criar pacote %s: %v\n", pkg.CodeName, err)
			} else {
				fmt.Printf("✅ Pacote criado: %s\n", pkg.DisplayName)
			}
		}
	}

	// 5. Configurar permissões padrão para cargos
	configureDefaultPermissions(roleRepo, permissionRepo)

	// 6. Configurar limites e módulos dos pacotes
	configurePackageLimitsAndModules(packageRepo, moduleRepo)

	fmt.Println("🌱 Seed de roles e permissões concluído!")
	return nil
}

func createModules() []models.Module {
	return []models.Module{
		// Módulos Admin
		{
			Id:           uuid.New(),
			CodeName:     "admin_organizations",
			DisplayName:  "Organizações",
			Description:  "Gerenciamento de organizações do sistema",
			Icon:         "building",
			Scope:        "admin",
			DisplayOrder: 1,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "admin_users",
			DisplayName:  "Usuários Admin",
			Description:  "Gerenciamento de usuários administradores",
			Icon:         "users-cog",
			Scope:        "admin",
			DisplayOrder: 2,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "admin_packages",
			DisplayName:  "Pacotes e Planos",
			Description:  "Gerenciamento de pacotes e assinaturas",
			Icon:         "package",
			Scope:        "admin",
			DisplayOrder: 3,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "admin_reports",
			DisplayName:  "Relatórios Globais",
			Description:  "Relatórios e métricas do sistema",
			Icon:         "chart-bar",
			Scope:        "admin",
			DisplayOrder: 4,
			IsFree:       false,
			Active:       true,
		},

		// Módulos Cliente
		{
			Id:           uuid.New(),
			CodeName:     "client_users",
			DisplayName:  "Usuários",
			Description:  "Gerenciamento de usuários da organização",
			Icon:         "users",
			Scope:        "client",
			DisplayOrder: 1,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "client_tables",
			DisplayName:  "Mesas",
			Description:  "Gerenciamento de mesas do estabelecimento",
			Icon:         "table",
			Scope:        "client",
			DisplayOrder: 2,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "client_customers",
			DisplayName:  "Clientes",
			Description:  "Gerenciamento de clientes",
			Icon:         "user-check",
			Scope:        "client",
			DisplayOrder: 3,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "client_menu",
			DisplayName:  "Cardápio",
			Description:  "Gerenciamento do cardápio digital",
			Icon:         "book-open",
			Scope:        "client",
			DisplayOrder: 4,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "client_products",
			DisplayName:  "Produtos",
			Description:  "Gerenciamento de produtos",
			Icon:         "package",
			Scope:        "client",
			DisplayOrder: 5,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "client_orders",
			DisplayName:  "Pedidos",
			Description:  "Gerenciamento de pedidos",
			Icon:         "shopping-cart",
			Scope:        "client",
			DisplayOrder: 6,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "client_reservations",
			DisplayName:  "Reservas",
			Description:  "Gerenciamento de reservas",
			Icon:         "calendar",
			Scope:        "client",
			DisplayOrder: 7,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "client_waitlist",
			DisplayName:  "Fila de Espera",
			Description:  "Gerenciamento da fila de espera",
			Icon:         "clock",
			Scope:        "client",
			DisplayOrder: 8,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "client_reports",
			DisplayName:  "Relatórios",
			Description:  "Relatórios e estatísticas",
			Icon:         "chart-bar",
			Scope:        "client",
			DisplayOrder: 9,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "client_settings",
			DisplayName:  "Configurações",
			Description:  "Configurações do projeto",
			Icon:         "settings",
			Scope:        "client",
			DisplayOrder: 10,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "client_notifications",
			DisplayName:  "Notificações",
			Description:  "Configuração de notificações",
			Icon:         "bell",
			Scope:        "client",
			DisplayOrder: 11,
			IsFree:       false,
			Active:       true,
		},
	}
}

func createPermissions(moduleRepo repositories.IModuleRepository) []models.Permission {
	var permissions []models.Permission

	// Helper para obter module ID
	getModuleId := func(codeName string) uuid.UUID {
		module, _ := moduleRepo.GetByCodeName(codeName)
		if module != nil {
			return module.Id
		}
		return uuid.Nil
	}

	// Permissões Admin - Organizações
	adminOrgModuleId := getModuleId("admin_organizations")
	if adminOrgModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "admin_organizations_view", DisplayName: "Visualizar Organizações", Description: "Pode visualizar lista de organizações", ModuleId: adminOrgModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "admin_organizations_manage", DisplayName: "Gerenciar Organizações", Description: "Pode criar, editar e excluir organizações", ModuleId: adminOrgModuleId, Active: true},
		)
	}

	// Permissões Admin - Usuários
	adminUsersModuleId := getModuleId("admin_users")
	if adminUsersModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "admin_users_view", DisplayName: "Visualizar Usuários Admin", Description: "Pode visualizar usuários administradores", ModuleId: adminUsersModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "admin_users_manage", DisplayName: "Gerenciar Usuários Admin", Description: "Pode gerenciar usuários administradores", ModuleId: adminUsersModuleId, Active: true},
		)
	}

	// Permissões Admin - Pacotes
	adminPkgModuleId := getModuleId("admin_packages")
	if adminPkgModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "admin_packages_view", DisplayName: "Visualizar Pacotes", Description: "Pode visualizar pacotes e planos", ModuleId: adminPkgModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "admin_packages_manage", DisplayName: "Gerenciar Pacotes", Description: "Pode criar e editar pacotes", ModuleId: adminPkgModuleId, Active: true},
		)
	}

	// Permissões Cliente - Usuários
	clientUsersModuleId := getModuleId("client_users")
	if clientUsersModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "client_users_view", DisplayName: "Visualizar Usuários", Description: "Pode visualizar lista de usuários", ModuleId: clientUsersModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "client_users_manage", DisplayName: "Gerenciar Usuários", Description: "Pode criar, editar e excluir usuários", ModuleId: clientUsersModuleId, Active: true},
		)
	}

	// Permissões Cliente - Mesas
	clientTablesModuleId := getModuleId("client_tables")
	if clientTablesModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "client_tables_view", DisplayName: "Visualizar Mesas", Description: "Pode visualizar mesas", ModuleId: clientTablesModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "client_tables_manage", DisplayName: "Gerenciar Mesas", Description: "Pode criar, editar e excluir mesas", ModuleId: clientTablesModuleId, Active: true},
		)
	}

	// Permissões Cliente - Clientes
	clientCustomersModuleId := getModuleId("client_customers")
	if clientCustomersModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "client_customers_view", DisplayName: "Visualizar Clientes", Description: "Pode visualizar clientes", ModuleId: clientCustomersModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "client_customers_manage", DisplayName: "Gerenciar Clientes", Description: "Pode criar, editar e excluir clientes", ModuleId: clientCustomersModuleId, Active: true},
		)
	}

	// Permissões Cliente - Cardápio
	clientMenuModuleId := getModuleId("client_menu")
	if clientMenuModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "client_menu_view", DisplayName: "Visualizar Cardápio", Description: "Pode visualizar o cardápio", ModuleId: clientMenuModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "client_menu_manage", DisplayName: "Gerenciar Cardápio", Description: "Pode editar menus, categorias e produtos", ModuleId: clientMenuModuleId, Active: true},
		)
	}

	// Permissões Cliente - Produtos
	clientProductsModuleId := getModuleId("client_products")
	if clientProductsModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "client_products_view", DisplayName: "Visualizar Produtos", Description: "Pode visualizar produtos", ModuleId: clientProductsModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "client_products_manage", DisplayName: "Gerenciar Produtos", Description: "Pode criar, editar e excluir produtos", ModuleId: clientProductsModuleId, Active: true},
		)
	}

	// Permissões Cliente - Pedidos
	clientOrdersModuleId := getModuleId("client_orders")
	if clientOrdersModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "client_orders_view", DisplayName: "Visualizar Pedidos", Description: "Pode visualizar pedidos", ModuleId: clientOrdersModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "client_orders_manage", DisplayName: "Gerenciar Pedidos", Description: "Pode criar, editar e atualizar status de pedidos", ModuleId: clientOrdersModuleId, Active: true},
		)
	}

	// Permissões Cliente - Reservas
	clientReservationsModuleId := getModuleId("client_reservations")
	if clientReservationsModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "client_reservations_view", DisplayName: "Visualizar Reservas", Description: "Pode visualizar reservas", ModuleId: clientReservationsModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "client_reservations_manage", DisplayName: "Gerenciar Reservas", Description: "Pode criar, editar e cancelar reservas", ModuleId: clientReservationsModuleId, Active: true},
		)
	}

	// Permissões Cliente - Fila de Espera
	clientWaitlistModuleId := getModuleId("client_waitlist")
	if clientWaitlistModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "client_waitlist_view", DisplayName: "Visualizar Fila de Espera", Description: "Pode visualizar fila de espera", ModuleId: clientWaitlistModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "client_waitlist_manage", DisplayName: "Gerenciar Fila de Espera", Description: "Pode adicionar e remover da fila", ModuleId: clientWaitlistModuleId, Active: true},
		)
	}

	// Permissões Cliente - Relatórios
	clientReportsModuleId := getModuleId("client_reports")
	if clientReportsModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "client_reports_view", DisplayName: "Visualizar Relatórios", Description: "Pode visualizar relatórios", ModuleId: clientReportsModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "client_reports_export", DisplayName: "Exportar Relatórios", Description: "Pode exportar relatórios", ModuleId: clientReportsModuleId, Active: true},
		)
	}

	// Permissões Cliente - Configurações
	clientSettingsModuleId := getModuleId("client_settings")
	if clientSettingsModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "client_settings_view", DisplayName: "Visualizar Configurações", Description: "Pode visualizar configurações", ModuleId: clientSettingsModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "client_settings_manage", DisplayName: "Gerenciar Configurações", Description: "Pode alterar configurações do projeto", ModuleId: clientSettingsModuleId, Active: true},
		)
	}

	// Permissões Cliente - Notificações
	clientNotificationsModuleId := getModuleId("client_notifications")
	if clientNotificationsModuleId != uuid.Nil {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), CodeName: "client_notifications_view", DisplayName: "Visualizar Notificações", Description: "Pode visualizar logs de notificações", ModuleId: clientNotificationsModuleId, Active: true},
			models.Permission{Id: uuid.New(), CodeName: "client_notifications_manage", DisplayName: "Gerenciar Notificações", Description: "Pode configurar templates e enviar notificações", ModuleId: clientNotificationsModuleId, Active: true},
		)
	}

	return permissions
}

func createSystemRoles() []models.Role {
	return []models.Role{
		// Cargos Admin (sistema)
		{
			Id:             uuid.New(),
			Name:           "super_admin",
			DisplayName:    "Super Administrador",
			Description:    "Acesso total ao sistema",
			Scope:          "admin",
			HierarchyLevel: 10,
			IsSystem:       true,
			Active:         true,
		},
		{
			Id:             uuid.New(),
			Name:           "admin_support",
			DisplayName:    "Suporte Técnico",
			Description:    "Acesso para suporte técnico",
			Scope:          "admin",
			HierarchyLevel: 8,
			IsSystem:       true,
			Active:         true,
		},
		{
			Id:             uuid.New(),
			Name:           "admin_sales",
			DisplayName:    "Comercial",
			Description:    "Acesso para equipe comercial",
			Scope:          "admin",
			HierarchyLevel: 6,
			IsSystem:       true,
			Active:         true,
		},

		// Cargos Cliente (organização)
		{
			Id:             uuid.New(),
			Name:           "owner",
			DisplayName:    "Proprietário",
			Description:    "Dono da organização com acesso total",
			Scope:          "client",
			HierarchyLevel: 10,
			IsSystem:       true,
			Active:         true,
		},
		{
			Id:             uuid.New(),
			Name:           "manager",
			DisplayName:    "Gerente",
			Description:    "Gerente com amplo acesso",
			Scope:          "client",
			HierarchyLevel: 8,
			IsSystem:       true,
			Active:         true,
		},
		{
			Id:             uuid.New(),
			Name:           "supervisor",
			DisplayName:    "Supervisor",
			Description:    "Supervisor de equipe",
			Scope:          "client",
			HierarchyLevel: 6,
			IsSystem:       true,
			Active:         true,
		},
		{
			Id:             uuid.New(),
			Name:           "attendant",
			DisplayName:    "Atendente",
			Description:    "Atendente com acesso operacional",
			Scope:          "client",
			HierarchyLevel: 4,
			IsSystem:       true,
			Active:         true,
		},
		{
			Id:             uuid.New(),
			Name:           "waiter",
			DisplayName:    "Garçom",
			Description:    "Garçom com acesso a pedidos e mesas",
			Scope:          "client",
			HierarchyLevel: 3,
			IsSystem:       true,
			Active:         true,
		},
		{
			Id:             uuid.New(),
			Name:           "kitchen",
			DisplayName:    "Cozinha",
			Description:    "Acesso apenas à fila de pedidos",
			Scope:          "client",
			HierarchyLevel: 2,
			IsSystem:       true,
			Active:         true,
		},
		{
			Id:             uuid.New(),
			Name:           "viewer",
			DisplayName:    "Visualizador",
			Description:    "Apenas visualização",
			Scope:          "client",
			HierarchyLevel: 1,
			IsSystem:       true,
			Active:         true,
		},
	}
}

func createPackages() []models.Package {
	return []models.Package{
		{
			Id:           uuid.New(),
			CodeName:     "free",
			DisplayName:  "Gratuito",
			Description:  "Plano gratuito com funcionalidades básicas",
			PriceMonthly: 0,
			PriceYearly:  0,
			IsPublic:     true,
			DisplayOrder: 1,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "starter",
			DisplayName:  "Starter",
			Description:  "Ideal para pequenos estabelecimentos",
			PriceMonthly: 99.90,
			PriceYearly:  999.00,
			IsPublic:     true,
			DisplayOrder: 2,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "professional",
			DisplayName:  "Profissional",
			Description:  "Para estabelecimentos em crescimento",
			PriceMonthly: 199.90,
			PriceYearly:  1999.00,
			IsPublic:     true,
			DisplayOrder: 3,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			CodeName:     "enterprise",
			DisplayName:  "Enterprise",
			Description:  "Solução completa para grandes operações",
			PriceMonthly: 499.90,
			PriceYearly:  4999.00,
			IsPublic:     true,
			DisplayOrder: 4,
			Active:       true,
		},
	}
}

func configureDefaultPermissions(roleRepo repositories.IRoleRepository, permissionRepo repositories.IPermissionRepository) {
	fmt.Println("🔧 Configurando permissões padrão para cargos...")

	// Buscar todos os cargos e permissões
	permissions, _ := permissionRepo.List()

	// Mapear permissões por código
	permByCode := make(map[string]models.Permission)
	for _, p := range permissions {
		permByCode[p.CodeName] = p
	}

	// Configurar permissões para cada cargo
	roleConfigs := map[string]map[string]int{
		"owner": {
			// Nível 2 = CRUD completo para tudo
			"client_users_view": 2, "client_users_manage": 2,
			"client_tables_view": 2, "client_tables_manage": 2,
			"client_customers_view": 2, "client_customers_manage": 2,
			"client_menu_view": 2, "client_menu_manage": 2,
			"client_products_view": 2, "client_products_manage": 2,
			"client_orders_view": 2, "client_orders_manage": 2,
			"client_reservations_view": 2, "client_reservations_manage": 2,
			"client_waitlist_view": 2, "client_waitlist_manage": 2,
			"client_reports_view": 2, "client_reports_export": 2,
			"client_settings_view": 2, "client_settings_manage": 2,
			"client_notifications_view": 2, "client_notifications_manage": 2,
		},
		"manager": {
			"client_users_view": 2, "client_users_manage": 2,
			"client_tables_view": 2, "client_tables_manage": 2,
			"client_customers_view": 2, "client_customers_manage": 2,
			"client_menu_view": 2, "client_menu_manage": 2,
			"client_products_view": 2, "client_products_manage": 2,
			"client_orders_view": 2, "client_orders_manage": 2,
			"client_reservations_view": 2, "client_reservations_manage": 2,
			"client_waitlist_view": 2, "client_waitlist_manage": 2,
			"client_reports_view": 2, "client_reports_export": 1,
			"client_settings_view": 1, "client_settings_manage": 0,
			"client_notifications_view": 1, "client_notifications_manage": 0,
		},
		"supervisor": {
			"client_users_view": 1, "client_users_manage": 0,
			"client_tables_view": 2, "client_tables_manage": 1,
			"client_customers_view": 2, "client_customers_manage": 1,
			"client_menu_view": 1, "client_menu_manage": 0,
			"client_products_view": 1, "client_products_manage": 0,
			"client_orders_view": 2, "client_orders_manage": 2,
			"client_reservations_view": 2, "client_reservations_manage": 2,
			"client_waitlist_view": 2, "client_waitlist_manage": 2,
			"client_reports_view": 1, "client_reports_export": 0,
		},
		"attendant": {
			"client_tables_view": 1, "client_tables_manage": 0,
			"client_customers_view": 2, "client_customers_manage": 1,
			"client_orders_view": 2, "client_orders_manage": 2,
			"client_reservations_view": 2, "client_reservations_manage": 2,
			"client_waitlist_view": 2, "client_waitlist_manage": 2,
		},
		"waiter": {
			"client_tables_view": 1,
			"client_menu_view": 1,
			"client_products_view": 1,
			"client_orders_view": 2, "client_orders_manage": 2,
		},
		"kitchen": {
			"client_orders_view": 1, "client_orders_manage": 1,
			"client_products_view": 1,
		},
		"viewer": {
			"client_tables_view": 1,
			"client_customers_view": 1,
			"client_menu_view": 1,
			"client_products_view": 1,
			"client_orders_view": 1,
			"client_reservations_view": 1,
			"client_waitlist_view": 1,
		},
	}

	for roleName, permLevels := range roleConfigs {
		role, err := roleRepo.GetByName(roleName)
		if err != nil || role == nil {
			continue
		}

		for permCode, level := range permLevels {
			perm, exists := permByCode[permCode]
			if !exists {
				continue
			}

			err := roleRepo.SetPermissionLevel(role.Id.String(), perm.Id.String(), level)
			if err != nil {
				fmt.Printf("⚠️ Erro ao configurar %s para %s: %v\n", permCode, roleName, err)
			}
		}
		fmt.Printf("✅ Permissões configuradas para: %s\n", role.DisplayName)
	}
}

// configurePackageLimitsAndModules configura os limites e módulos para cada pacote
func configurePackageLimitsAndModules(packageRepo repositories.IPackageRepository, moduleRepo repositories.IModuleRepository) {
	fmt.Println("📦 Configurando limites e módulos dos pacotes...")

	// Definição de limites por pacote
	// -1 = ilimitado, 0 = desabilitado
	packageLimits := map[string]map[string]int{
		"free": {
			"users":                3,
			"tables":               10,
			"products":             50,
			"reservations_per_day": 0, // Desabilitado no plano gratuito
		},
		"starter": {
			"users":                10,
			"tables":               30,
			"products":             200,
			"reservations_per_day": 20,
		},
		"professional": {
			"users":                50,
			"tables":               100,
			"products":             1000,
			"reservations_per_day": 100,
		},
		"enterprise": {
			"users":                -1, // Ilimitado
			"tables":               -1,
			"products":             -1,
			"reservations_per_day": -1,
		},
	}

	// Definição de módulos por pacote
	// Módulos gratuitos (IsFree=true) são incluídos em todos os pacotes
	packageModules := map[string][]string{
		"free": {
			// Apenas módulos gratuitos
			"client_users",
			"client_tables",
			"client_customers",
			"client_menu",
			"client_products",
			"client_orders",
			"client_settings",
		},
		"starter": {
			// Gratuitos + reservas e fila
			"client_users",
			"client_tables",
			"client_customers",
			"client_menu",
			"client_products",
			"client_orders",
			"client_settings",
			"client_reservations",
			"client_waitlist",
			"client_reports",
		},
		"professional": {
			// Starter + notificações
			"client_users",
			"client_tables",
			"client_customers",
			"client_menu",
			"client_products",
			"client_orders",
			"client_settings",
			"client_reservations",
			"client_waitlist",
			"client_reports",
			"client_notifications",
		},
		"enterprise": {
			// Todos os módulos
			"client_users",
			"client_tables",
			"client_customers",
			"client_menu",
			"client_products",
			"client_orders",
			"client_settings",
			"client_reservations",
			"client_waitlist",
			"client_reports",
			"client_notifications",
		},
	}

	// Aplicar limites para cada pacote
	for pkgCode, limits := range packageLimits {
		pkg, err := packageRepo.GetByCodeName(pkgCode)
		if err != nil || pkg == nil {
			fmt.Printf("⚠️ Pacote %s não encontrado\n", pkgCode)
			continue
		}

		for limitType, limitValue := range limits {
			err := packageRepo.SetPackageLimit(pkg.Id.String(), limitType, limitValue)
			if err != nil {
				fmt.Printf("⚠️ Erro ao definir limite %s para %s: %v\n", limitType, pkgCode, err)
			}
		}
		fmt.Printf("✅ Limites configurados para: %s\n", pkg.DisplayName)
	}

	// Aplicar módulos para cada pacote
	for pkgCode, modules := range packageModules {
		pkg, err := packageRepo.GetByCodeName(pkgCode)
		if err != nil || pkg == nil {
			continue
		}

		for _, modCode := range modules {
			mod, err := moduleRepo.GetByCodeName(modCode)
			if err != nil || mod == nil {
				continue
			}

			err = packageRepo.AddModuleToPackage(pkg.Id.String(), mod.Id.String())
			if err != nil {
				fmt.Printf("⚠️ Erro ao adicionar módulo %s ao pacote %s: %v\n", modCode, pkgCode, err)
			}
		}
		fmt.Printf("✅ Módulos configurados para: %s\n", pkg.DisplayName)
	}
}
