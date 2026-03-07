package handler

import (
	"fmt"
	"lep/repositories"
	"lep/repositories/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedStaffModules popula os módulos, permissões e configurações relacionados a gestão de equipe
// Deve ser chamado após SeedRolesAndPermissions
func SeedStaffModules(db *gorm.DB) error {
	fmt.Println("🌱 Iniciando seed de módulos Staff...")

	moduleRepo := repositories.NewModuleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)
	planRepo := repositories.NewPlanRepository(db)

	// 1. Criar Módulos Staff
	modules := createStaffModules()
	for _, m := range modules {
		existing, _ := moduleRepo.GetByCodeName(m.Code)
		if existing == nil {
			if err := moduleRepo.Create(&m); err != nil {
				fmt.Printf("⚠️ Erro ao criar módulo %s: %v\n", m.Code, err)
			} else {
				fmt.Printf("✅ Módulo Staff criado: %s\n", m.Name)
			}
		}
	}

	// 2. Criar Permissões Staff
	permissions := createStaffPermissions(moduleRepo)
	for _, p := range permissions {
		existing, _ := permissionRepo.GetByCodeName(p.Code)
		if existing == nil {
			if err := permissionRepo.Create(&p); err != nil {
				fmt.Printf("⚠️ Erro ao criar permissão %s: %v\n", p.Code, err)
			} else {
				fmt.Printf("✅ Permissão Staff criada: %s\n", p.DisplayName)
			}
		}
	}

	// 3. Adicionar módulos Staff aos planos
	configureStaffModulesInPlans(planRepo, moduleRepo)

	fmt.Println("🌱 Seed de módulos Staff concluído!")
	return nil
}

func createStaffModules() []models.Module {
	return []models.Module{
		{
			Id:           uuid.New(),
			Code:         "client_staff_availability",
			Name:         "Disponibilidade de Equipe",
			Description:  "Gerenciamento de disponibilidade semanal dos funcionários",
			Icon:         "calendar-check",
			Scope:        "client",
			DisplayOrder: 20,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_staff_schedule",
			Name:         "Escalas",
			Description:  "Gerenciamento de escalas de trabalho",
			Icon:         "users",
			Scope:        "client",
			DisplayOrder: 21,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_staff_attendance",
			Name:         "Presença",
			Description:  "Registro de presença, transporte e consumo",
			Icon:         "clipboard-list",
			Scope:        "client",
			DisplayOrder: 22,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_staff_stock",
			Name:         "Estoque Operacional",
			Description:  "Controle de estoque com mínimo/máximo e lista de compras",
			Icon:         "package",
			Scope:        "client",
			DisplayOrder: 23,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_staff_commission",
			Name:         "Comissões",
			Description:  "Registro de comissões e relatório de pagamentos",
			Icon:         "dollar-sign",
			Scope:        "client",
			DisplayOrder: 24,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_staff_dashboard",
			Name:         "Dashboard Operacional",
			Description:  "Dashboard de faturamento com gráficos e filtros",
			Icon:         "bar-chart-3",
			Scope:        "client",
			DisplayOrder: 25,
			IsFree:       false,
			Active:       true,
		},
		{
			Id:           uuid.New(),
			Code:         "client_staff_reports",
			Name:         "Relatórios de Equipe",
			Description:  "Relatórios detalhados de pagamentos por colaborador",
			Icon:         "file-text",
			Scope:        "client",
			DisplayOrder: 26,
			IsFree:       false,
			Active:       true,
		},
	}
}

func createStaffPermissions(moduleRepo repositories.IModuleRepository) []models.Permission {
	var permissions []models.Permission

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

	// Disponibilidade - CRUD
	addCRUDPermissions("client_staff_availability", "Disponibilidade")

	// Escalas - CRUD + envio de emails
	addCRUDPermissions("client_staff_schedule", "Escalas")
	if moduleExists("client_staff_schedule") {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), Code: "client_staff_schedule_send_email", Module: "client_staff_schedule", Action: "send", DisplayName: "Enviar Emails de Escala", Description: "Pode enviar emails com a escala", Active: true},
		)
	}

	// Presença - CRUD
	addCRUDPermissions("client_staff_attendance", "Presença")

	// Estoque - CRUD + gerar lista de compras
	addCRUDPermissions("client_staff_stock", "Estoque Operacional")
	if moduleExists("client_staff_stock") {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), Code: "client_staff_stock_generate_list", Module: "client_staff_stock", Action: "export", DisplayName: "Gerar Lista de Compras", Description: "Pode gerar lista de compras em PDF", Active: true},
		)
	}

	// Comissões - CRUD + gerar relatório
	addCRUDPermissions("client_staff_commission", "Comissões")
	if moduleExists("client_staff_commission") {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), Code: "client_staff_commission_report", Module: "client_staff_commission", Action: "export", DisplayName: "Gerar Relatório de Pagamentos", Description: "Pode gerar relatório de pagamentos", Active: true},
		)
	}

	// Dashboard - view + export
	if moduleExists("client_staff_dashboard") {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), Code: "client_staff_dashboard_view", Module: "client_staff_dashboard", Action: "view", DisplayName: "Visualizar Dashboard Operacional", Description: "Pode visualizar dashboard de faturamento", Active: true},
			models.Permission{Id: uuid.New(), Code: "client_staff_dashboard_export", Module: "client_staff_dashboard", Action: "export", DisplayName: "Exportar Dashboard", Description: "Pode exportar dados do dashboard", Active: true},
			models.Permission{Id: uuid.New(), Code: "client_staff_dashboard_import", Module: "client_staff_dashboard", Action: "import", DisplayName: "Importar Dados", Description: "Pode importar dados de vendas via CSV", Active: true},
		)
	}

	// Relatórios de Equipe - view + export
	if moduleExists("client_staff_reports") {
		permissions = append(permissions,
			models.Permission{Id: uuid.New(), Code: "client_staff_reports_view", Module: "client_staff_reports", Action: "view", DisplayName: "Visualizar Relatórios de Equipe", Description: "Pode visualizar relatórios de equipe", Active: true},
			models.Permission{Id: uuid.New(), Code: "client_staff_reports_export", Module: "client_staff_reports", Action: "export", DisplayName: "Exportar Relatórios de Equipe", Description: "Pode exportar relatórios de equipe", Active: true},
		)
	}

	return permissions
}

// configureStaffModulesInPlans adiciona os módulos staff aos planos professional e enterprise
func configureStaffModulesInPlans(planRepo repositories.IPlanRepository, moduleRepo repositories.IModuleRepository) {
	fmt.Println("📦 Adicionando módulos Staff aos planos...")

	// Módulos staff para adicionar aos planos maiores
	staffModules := []string{
		"client_staff_availability",
		"client_staff_schedule",
		"client_staff_attendance",
		"client_staff_stock",
		"client_staff_commission",
		"client_staff_dashboard",
		"client_staff_reports",
	}

	// Planos que terão os módulos staff
	plansWithStaffModules := []string{"demo", "professional", "enterprise"}

	for _, planCode := range plansWithStaffModules {
		plan, err := planRepo.GetByCode(planCode)
		if err != nil || plan == nil {
			continue
		}

		for _, modCode := range staffModules {
			mod, err := moduleRepo.GetByCodeName(modCode)
			if err != nil || mod == nil {
				continue
			}

			err = planRepo.AddModuleToPlan(plan.Id.String(), mod.Id.String())
			if err != nil {
				fmt.Printf("⚠️ Erro ao adicionar módulo %s ao plano %s: %v\n", modCode, planCode, err)
			}
		}
		fmt.Printf("✅ Módulos Staff adicionados ao plano: %s\n", plan.Name)
	}
}
