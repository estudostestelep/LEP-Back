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
	planRepo := repositories.NewPlanRepository(db)

	// 1. Criar Módulos
	modules := createModules()
	for _, m := range modules {
		existing, _ := moduleRepo.GetByCodeName(m.Code)
		if existing == nil {
			if err := moduleRepo.Create(&m); err != nil {
				fmt.Printf("⚠️ Erro ao criar módulo %s: %v\n", m.Code, err)
			} else {
				fmt.Printf("✅ Módulo criado: %s\n", m.Name)
			}
		}
	}

	// 2. Criar Permissões
	permissions := createPermissions(moduleRepo)
	for _, p := range permissions {
		existing, _ := permissionRepo.GetByCodeName(p.Code)
		if existing == nil {
			if err := permissionRepo.Create(&p); err != nil {
				fmt.Printf("⚠️ Erro ao criar permissão %s: %v\n", p.Code, err)
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

	// 4. Criar Planos
	plans := createPlans()
	for _, plan := range plans {
		existing, _ := planRepo.GetByCode(plan.Code)
		if existing == nil {
			if err := planRepo.Create(&plan); err != nil {
				fmt.Printf("⚠️ Erro ao criar plano %s: %v\n", plan.Code, err)
			} else {
				fmt.Printf("✅ Plano criado: %s\n", plan.Name)
			}
		}
	}

	// 5. Configurar permissões padrão para cargos
	configureDefaultPermissions(roleRepo, permissionRepo)

	// 6. Configurar limites e módulos dos planos
	configurePlanLimitsAndModules(planRepo, moduleRepo)

	fmt.Println("🌱 Seed de roles e permissões concluído!")
	return nil
}

func createModules() []models.Module {
	return []models.Module{
		// Módulos Admin
		{
			Id:           uuid.New(),
			Code:         "admin_organizations",
			Name:         "Organizações",
			Description:  "Gerenciamento de organizações do sistema",
			Icon:         "building",
			Scope:        "admin",
			DisplayOrder: 1,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "admin_users",
			Name:         "Usuários Admin",
			Description:  "Gerenciamento de usuários administradores",
			Icon:         "users-cog",
			Scope:        "admin",
			DisplayOrder: 2,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "admin_packages",
			Name:         "Pacotes e Planos",
			Description:  "Gerenciamento de pacotes e assinaturas",
			Icon:         "package",
			Scope:        "admin",
			DisplayOrder: 3,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "admin_reports",
			Name:         "Relatórios Globais",
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
			Code:         "client_users",
			Name:         "Usuários",
			Description:  "Gerenciamento de usuários da organização",
			Icon:         "users",
			Scope:        "client",
			DisplayOrder: 1,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_tables",
			Name:         "Mesas",
			Description:  "Gerenciamento de mesas do estabelecimento",
			Icon:         "table",
			Scope:        "client",
			DisplayOrder: 2,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_customers",
			Name:         "Clientes",
			Description:  "Gerenciamento de clientes",
			Icon:         "user-check",
			Scope:        "client",
			DisplayOrder: 3,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_menu",
			Name:         "Cardápio",
			Description:  "Gerenciamento do cardápio digital",
			Icon:         "book-open",
			Scope:        "client",
			DisplayOrder: 4,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_products",
			Name:         "Produtos",
			Description:  "Gerenciamento de produtos",
			Icon:         "package",
			Scope:        "client",
			DisplayOrder: 5,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_orders",
			Name:         "Pedidos",
			Description:  "Gerenciamento de pedidos",
			Icon:         "shopping-cart",
			Scope:        "client",
			DisplayOrder: 6,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_reservations",
			Name:         "Reservas",
			Description:  "Gerenciamento de reservas",
			Icon:         "calendar",
			Scope:        "client",
			DisplayOrder: 7,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_waitlist",
			Name:         "Fila de Espera",
			Description:  "Gerenciamento da fila de espera",
			Icon:         "clock",
			Scope:        "client",
			DisplayOrder: 8,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_reports",
			Name:         "Relatórios",
			Description:  "Relatórios e estatísticas",
			Icon:         "chart-bar",
			Scope:        "client",
			DisplayOrder: 9,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_settings",
			Name:         "Configurações",
			Description:  "Configurações do projeto",
			Icon:         "settings",
			Scope:        "client",
			DisplayOrder: 10,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_notifications",
			Name:         "Notificações",
			Description:  "Configuração de notificações",
			Icon:         "bell",
			Scope:        "client",
			DisplayOrder: 11,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_tags",
			Name:         "Tags",
			Description:  "Gerenciamento de tags e etiquetas",
			Icon:         "tag",
			Scope:        "client",
			DisplayOrder: 12,
			IsFree:       true,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_audit_logs",
			Name:         "Logs de Auditoria",
			Description:  "Visualização de logs de auditoria das operações",
			Icon:         "history",
			Scope:        "client",
			DisplayOrder: 13,
			IsFree:       false,
			Active:       true,
		},
	}
}

func createPermissions(moduleRepo repositories.IModuleRepository) []models.Permission {
	var permissions []models.Permission

	// Helper para verificar se módulo existe
	moduleExists := func(codeName string) bool {
		module, _ := moduleRepo.GetByCodeName(codeName)
		return module != nil
	}

	// Helper para criar permissões CRUD padrão para um módulo
	addCRUDPermissions := func(moduleCode, displayName string) {
		if !moduleExists(moduleCode) {
			return
		}
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), Code: moduleCode + "_view", Module: moduleCode, Action: "view", DisplayName: "Visualizar " + displayName, Description: "Pode visualizar " + displayName, Active: true},
			models.Permission{Id: uuid.New(), Code: moduleCode + "_create", Module: moduleCode, Action: "create", DisplayName: "Criar " + displayName, Description: "Pode criar " + displayName, Active: true},
			models.Permission{Id: uuid.New(), Code: moduleCode + "_edit", Module: moduleCode, Action: "edit", DisplayName: "Editar " + displayName, Description: "Pode editar " + displayName, Active: true},
			models.Permission{Id: uuid.New(), Code: moduleCode + "_delete", Module: moduleCode, Action: "delete", DisplayName: "Excluir " + displayName, Description: "Pode excluir " + displayName, Active: true},
		)
	}

	// Permissões Admin - Organizações (CRUD)
	addCRUDPermissions("admin_organizations", "Organizações")

	// Permissões Admin - Usuários (CRUD)
	addCRUDPermissions("admin_users", "Usuários Admin")

	// Permissões Admin - Pacotes (CRUD)
	addCRUDPermissions("admin_packages", "Pacotes")

	// Permissões Admin - Relatórios (apenas view e export)
	if moduleExists("admin_reports") {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), Code: "admin_reports_view", Module: "admin_reports", Action: "view", DisplayName: "Visualizar Relatórios Globais", Description: "Pode visualizar relatórios globais", Active: true},
			models.Permission{Id: uuid.New(), Code: "admin_reports_export", Module: "admin_reports", Action: "export", DisplayName: "Exportar Relatórios Globais", Description: "Pode exportar relatórios globais", Active: true},
		)
	}

	// Permissões Cliente - Usuários (CRUD)
	addCRUDPermissions("client_users", "Usuários")

	// Permissões Cliente - Mesas (CRUD)
	addCRUDPermissions("client_tables", "Mesas")

	// Permissões Cliente - Clientes (CRUD)
	addCRUDPermissions("client_customers", "Clientes")

	// Permissões Cliente - Cardápio (CRUD)
	addCRUDPermissions("client_menu", "Cardápio")

	// Permissões Cliente - Produtos (CRUD)
	addCRUDPermissions("client_products", "Produtos")

	// Permissões Cliente - Pedidos (CRUD)
	addCRUDPermissions("client_orders", "Pedidos")

	// Permissões Cliente - Reservas (CRUD)
	addCRUDPermissions("client_reservations", "Reservas")

	// Permissões Cliente - Fila de Espera (CRUD)
	addCRUDPermissions("client_waitlist", "Fila de Espera")

	// Permissões Cliente - Relatórios (apenas view e export)
	if moduleExists("client_reports") {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), Code: "client_reports_view", Module: "client_reports", Action: "view", DisplayName: "Visualizar Relatórios", Description: "Pode visualizar relatórios", Active: true},
			models.Permission{Id: uuid.New(), Code: "client_reports_export", Module: "client_reports", Action: "export", DisplayName: "Exportar Relatórios", Description: "Pode exportar relatórios", Active: true},
		)
	}

	// Permissões Cliente - Configurações (view e edit, sem create/delete)
	if moduleExists("client_settings") {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), Code: "client_settings_view", Module: "client_settings", Action: "view", DisplayName: "Visualizar Configurações", Description: "Pode visualizar configurações", Active: true},
			models.Permission{Id: uuid.New(), Code: "client_settings_edit", Module: "client_settings", Action: "edit", DisplayName: "Editar Configurações", Description: "Pode alterar configurações do projeto", Active: true},
		)
	}

	// Permissões Cliente - Notificações (CRUD para templates + enviar)
	if moduleExists("client_notifications") {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), Code: "client_notifications_view", Module: "client_notifications", Action: "view", DisplayName: "Visualizar Notificações", Description: "Pode visualizar logs de notificações", Active: true},
			models.Permission{Id: uuid.New(), Code: "client_notifications_create", Module: "client_notifications", Action: "create", DisplayName: "Criar Templates", Description: "Pode criar templates de notificação", Active: true},
			models.Permission{Id: uuid.New(), Code: "client_notifications_edit", Module: "client_notifications", Action: "edit", DisplayName: "Editar Templates", Description: "Pode editar templates de notificação", Active: true},
			models.Permission{Id: uuid.New(), Code: "client_notifications_delete", Module: "client_notifications", Action: "delete", DisplayName: "Excluir Templates", Description: "Pode excluir templates de notificação", Active: true},
			models.Permission{Id: uuid.New(), Code: "client_notifications_send", Module: "client_notifications", Action: "send", DisplayName: "Enviar Notificações", Description: "Pode enviar notificações manualmente", Active: true},
		)
	}

	// Permissões Cliente - Tags (CRUD)
	addCRUDPermissions("client_tags", "Tags")

	// Permissões Cliente - Logs de Auditoria (apenas view e configure)
	if moduleExists("client_audit_logs") {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), Code: "client_audit_logs_view", Module: "client_audit_logs", Action: "view", DisplayName: "Visualizar Logs de Auditoria", Description: "Pode visualizar logs de auditoria", Active: true},
			models.Permission{Id: uuid.New(), Code: "client_audit_logs_configure", Module: "client_audit_logs", Action: "configure", DisplayName: "Configurar Auditoria", Description: "Pode configurar o módulo de auditoria", Active: true},
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

func createPlans() []models.Plan {
	return []models.Plan{
		{
			Id:           uuid.New(),
			Code:         "demo",
			Name:         "Demo",
			Description:  "Plano de demonstração com todos os módulos",
			PriceMonthly: 0,
			PriceYearly:  0,
			IsPublic:     false,
			DisplayOrder: 0,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "free",
			Name:         "Gratuito",
			Description:  "Plano gratuito com funcionalidades básicas",
			PriceMonthly: 0,
			PriceYearly:  0,
			IsPublic:     true,
			DisplayOrder: 1,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "starter",
			Name:         "Starter",
			Description:  "Ideal para pequenos estabelecimentos",
			PriceMonthly: 99.90,
			PriceYearly:  999.00,
			IsPublic:     true,
			DisplayOrder: 2,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "professional",
			Name:         "Profissional",
			Description:  "Para estabelecimentos em crescimento",
			PriceMonthly: 199.90,
			PriceYearly:  1999.00,
			IsPublic:     true,
			DisplayOrder: 3,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "enterprise",
			Name:         "Enterprise",
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
		permByCode[p.Code] = p
	}

	// Configurar permissões para cada cargo (1 = habilitado, 0 = desabilitado)
	// Agora cada permissão é granular: _view, _create, _edit, _delete
	roleConfigs := map[string]map[string]int{
		"owner": {
			// CRUD completo para tudo
			"client_users_view": 1, "client_users_create": 1, "client_users_edit": 1, "client_users_delete": 1,
			"client_tables_view": 1, "client_tables_create": 1, "client_tables_edit": 1, "client_tables_delete": 1,
			"client_customers_view": 1, "client_customers_create": 1, "client_customers_edit": 1, "client_customers_delete": 1,
			"client_menu_view": 1, "client_menu_create": 1, "client_menu_edit": 1, "client_menu_delete": 1,
			"client_products_view": 1, "client_products_create": 1, "client_products_edit": 1, "client_products_delete": 1,
			"client_orders_view": 1, "client_orders_create": 1, "client_orders_edit": 1, "client_orders_delete": 1,
			"client_reservations_view": 1, "client_reservations_create": 1, "client_reservations_edit": 1, "client_reservations_delete": 1,
			"client_waitlist_view": 1, "client_waitlist_create": 1, "client_waitlist_edit": 1, "client_waitlist_delete": 1,
			"client_reports_view": 1, "client_reports_export": 1,
			"client_settings_view": 1, "client_settings_edit": 1,
			"client_notifications_view": 1, "client_notifications_create": 1, "client_notifications_edit": 1, "client_notifications_delete": 1, "client_notifications_send": 1,
			"client_tags_view": 1, "client_tags_create": 1, "client_tags_edit": 1, "client_tags_delete": 1,
			"client_audit_logs_view": 1, "client_audit_logs_configure": 1,
		},
		"manager": {
			// CRUD completo para a maioria, mas sem algumas configurações
			"client_users_view": 1, "client_users_create": 1, "client_users_edit": 1, "client_users_delete": 1,
			"client_tables_view": 1, "client_tables_create": 1, "client_tables_edit": 1, "client_tables_delete": 1,
			"client_customers_view": 1, "client_customers_create": 1, "client_customers_edit": 1, "client_customers_delete": 1,
			"client_menu_view": 1, "client_menu_create": 1, "client_menu_edit": 1, "client_menu_delete": 1,
			"client_products_view": 1, "client_products_create": 1, "client_products_edit": 1, "client_products_delete": 1,
			"client_orders_view": 1, "client_orders_create": 1, "client_orders_edit": 1, "client_orders_delete": 1,
			"client_reservations_view": 1, "client_reservations_create": 1, "client_reservations_edit": 1, "client_reservations_delete": 1,
			"client_waitlist_view": 1, "client_waitlist_create": 1, "client_waitlist_edit": 1, "client_waitlist_delete": 1,
			"client_reports_view": 1, "client_reports_export": 1,
			"client_settings_view": 1,
			"client_notifications_view": 1,
			"client_tags_view": 1, "client_tags_create": 1, "client_tags_edit": 1, "client_tags_delete": 1,
			"client_audit_logs_view": 1,
		},
		"supervisor": {
			// View e algumas edições, sem delete na maioria
			"client_users_view": 1,
			"client_tables_view": 1, "client_tables_create": 1, "client_tables_edit": 1,
			"client_customers_view": 1, "client_customers_create": 1, "client_customers_edit": 1,
			"client_menu_view": 1,
			"client_products_view": 1,
			"client_orders_view": 1, "client_orders_create": 1, "client_orders_edit": 1, "client_orders_delete": 1,
			"client_reservations_view": 1, "client_reservations_create": 1, "client_reservations_edit": 1, "client_reservations_delete": 1,
			"client_waitlist_view": 1, "client_waitlist_create": 1, "client_waitlist_edit": 1, "client_waitlist_delete": 1,
			"client_reports_view": 1,
			"client_tags_view": 1,
		},
		"attendant": {
			// Operacional: pedidos, reservas, fila
			"client_tables_view": 1,
			"client_customers_view": 1, "client_customers_create": 1, "client_customers_edit": 1,
			"client_orders_view": 1, "client_orders_create": 1, "client_orders_edit": 1,
			"client_reservations_view": 1, "client_reservations_create": 1, "client_reservations_edit": 1,
			"client_waitlist_view": 1, "client_waitlist_create": 1, "client_waitlist_edit": 1, "client_waitlist_delete": 1,
		},
		"waiter": {
			// Garçom: visualiza mesas, menu, produtos e gerencia pedidos
			"client_tables_view": 1,
			"client_menu_view": 1,
			"client_products_view": 1,
			"client_orders_view": 1, "client_orders_create": 1, "client_orders_edit": 1,
		},
		"kitchen": {
			// Cozinha: apenas pedidos e produtos (visualização)
			"client_orders_view": 1, "client_orders_edit": 1,
			"client_products_view": 1,
		},
		"viewer": {
			// Apenas visualização
			"client_tables_view": 1,
			"client_customers_view": 1,
			"client_menu_view": 1,
			"client_products_view": 1,
			"client_orders_view": 1,
			"client_reservations_view": 1,
			"client_waitlist_view": 1,
			"client_tags_view": 1,
		},
		// ==================== Roles Admin ====================
		"super_admin": {
			// CRUD completo em todos os módulos admin
			"admin_organizations_view": 2, "admin_organizations_create": 2, "admin_organizations_edit": 2, "admin_organizations_delete": 2,
			"admin_users_view": 2, "admin_users_create": 2, "admin_users_edit": 2, "admin_users_delete": 2,
			"admin_packages_view": 2, "admin_packages_create": 2, "admin_packages_edit": 2, "admin_packages_delete": 2,
			"admin_reports_view": 2, "admin_reports_export": 2,
		},
		"admin_support": {
			// View + Edit em orgs e users, View em packages, View + Export em reports
			"admin_organizations_view": 2, "admin_organizations_edit": 2,
			"admin_users_view": 2, "admin_users_edit": 2,
			"admin_packages_view": 1,
			"admin_reports_view": 2, "admin_reports_export": 2,
		},
		"admin_sales": {
			// View em orgs, View + Create em packages, View em reports
			"admin_organizations_view": 1,
			"admin_packages_view": 1, "admin_packages_create": 1,
			"admin_reports_view": 1,
		},
	}

	for roleName, permLevels := range roleConfigs {
		role, err := roleRepo.GetByName(roleName)
		if err != nil || role == nil {
			continue
		}

		for permCode, level := range permLevels {
			// Se level > 0, atribuir a permissão (sistema agora é binário)
			if level <= 0 {
				continue
			}

			perm, exists := permByCode[permCode]
			if !exists {
				continue
			}

			err := roleRepo.AddPermissionToRole(role.Id.String(), perm.Id.String())
			if err != nil {
				fmt.Printf("⚠️ Erro ao configurar %s para %s: %v\n", permCode, roleName, err)
			}
		}
		fmt.Printf("✅ Permissões configuradas para: %s\n", role.DisplayName)
	}
}

// configurePlanLimitsAndModules configura os limites e módulos para cada plano
func configurePlanLimitsAndModules(planRepo repositories.IPlanRepository, moduleRepo repositories.IModuleRepository) {
	fmt.Println("📦 Configurando limites e módulos dos planos...")

	// Definição de limites por pacote
	// -1 = ilimitado, 0 = desabilitado
	packageLimits := map[string]map[string]int{
		"demo": {
			"users":                10,
			"tables":               10,
			"products":             10,
			"reservations_per_day": 10,
			"audit_logs_limit":     100,
			"audit_logs_retention": 7,
		},
		"free": {
			"users":                  3,
			"tables":                 10,
			"products":               50,
			"reservations_per_day":   0,  // Desabilitado no plano gratuito
			"audit_logs_limit":       0,  // Desabilitado no plano gratuito
			"audit_logs_retention":   0,  // Desabilitado no plano gratuito
		},
		"starter": {
			"users":                  10,
			"tables":                 30,
			"products":               200,
			"reservations_per_day":   20,
			"audit_logs_limit":       0,  // Desabilitado no plano starter
			"audit_logs_retention":   0,  // Desabilitado no plano starter
		},
		"professional": {
			"users":                  50,
			"tables":                 100,
			"products":               1000,
			"reservations_per_day":   100,
			"audit_logs_limit":       10000, // 10.000 logs
			"audit_logs_retention":   90,    // 90 dias
		},
		"enterprise": {
			"users":                  -1, // Ilimitado
			"tables":                 -1,
			"products":               -1,
			"reservations_per_day":   -1,
			"audit_logs_limit":       -1, // Ilimitado
			"audit_logs_retention":   365, // 1 ano
		},
	}

	// Definição de módulos por pacote
	// Módulos gratuitos (IsFree=true) são incluídos em todos os pacotes
	packageModules := map[string][]string{
		"demo": {
			// Todos os módulos habilitados
			"client_users",
			"client_tables",
			"client_customers",
			"client_menu",
			"client_products",
			"client_orders",
			"client_settings",
			"client_tags",
			"client_reservations",
			"client_waitlist",
			"client_reports",
			"client_notifications",
			"client_audit_logs",
		},
		"free": {
			// Apenas módulos gratuitos
			"client_users",
			"client_tables",
			"client_customers",
			"client_menu",
			"client_products",
			"client_orders",
			"client_settings",
			"client_tags",
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
			"client_tags",
			"client_reservations",
			"client_waitlist",
			"client_reports",
		},
		"professional": {
			// Starter + notificações + audit logs
			"client_users",
			"client_tables",
			"client_customers",
			"client_menu",
			"client_products",
			"client_orders",
			"client_settings",
			"client_tags",
			"client_reservations",
			"client_waitlist",
			"client_reports",
			"client_notifications",
			"client_audit_logs",
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
			"client_tags",
			"client_reservations",
			"client_waitlist",
			"client_reports",
			"client_notifications",
			"client_audit_logs",
		},
	}

	// Aplicar limites para cada plano
	for planCode, limits := range packageLimits {
		plan, err := planRepo.GetByCode(planCode)
		if err != nil || plan == nil {
			fmt.Printf("⚠️ Plano %s não encontrado\n", planCode)
			continue
		}

		for limitType, limitValue := range limits {
			err := planRepo.SetPlanLimit(plan.Id.String(), limitType, limitValue)
			if err != nil {
				fmt.Printf("⚠️ Erro ao definir limite %s para %s: %v\n", limitType, planCode, err)
			}
		}
		fmt.Printf("✅ Limites configurados para: %s\n", plan.Name)
	}

	// Aplicar módulos para cada plano
	for planCode, modules := range packageModules {
		plan, err := planRepo.GetByCode(planCode)
		if err != nil || plan == nil {
			continue
		}

		for _, modCode := range modules {
			mod, err := moduleRepo.GetByCodeName(modCode)
			if err != nil || mod == nil {
				continue
			}

			err = planRepo.AddModuleToPlan(plan.Id.String(), mod.Id.String())
			if err != nil {
				fmt.Printf("⚠️ Erro ao adicionar módulo %s ao plano %s: %v\n", modCode, planCode, err)
			}
		}
		fmt.Printf("✅ Módulos configurados para: %s\n", plan.Name)
	}
}
