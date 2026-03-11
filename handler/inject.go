package handler

import (
	"lep/repositories"
	"lep/service"
	"lep/utils"
)

type Handlers struct {
	HandlerProducts           IHandlerProducts
	HandlerAuth               IHandlerAuth
	HandlerOrder              IOrderHandler
	HandlerOrganization       IHandlerOrganization
	HandlerTables             IHandlerTables
	HandlerWaitlist           IHandlerWaitlist
	HandlerReservation        IHandlerReservation
	HandlerCustomer           IHandlerCustomer
	HandlerProject            IProjectHandler
	HandlerSettings           ISettingsHandler
	HandlerDisplaySettings    IDisplaySettingsHandler
	HandlerThemeCustomization IThemeCustomizationHandler
	HandlerEnvironment        IEnvironmentHandler
	HandlerNotification       *NotificationHandler
	HandlerReports            IReportsHandler
	HandlerTag                IHandlerTag
	HandlerMenu               IHandlerMenu
	HandlerCategory           IHandlerCategory
	HandlerSubcategory        IHandlerSubcategory
	HandlerImageManagement    IHandlerImageManagement
	HandlerOnboarding         IOnboardingHandler
	HandlerRole               *RoleHandler
	HandlerPlanChangeRequest  IPlanChangeRequestHandler
	HandlerLimits             *LimitHandler
	ImageManagementService    service.IImageManagementService // Service direto para o Upload server
	HandlerSidebarConfig      ISidebarConfigHandler
	HandlerAdminAuditLog      IAdminAuditLogHandler  // Logs de auditoria administrativa (read-only)
	HandlerClientAuditLog     IClientAuditLogHandler // Logs de auditoria de cliente (módulo opcional)
	// Novos handlers para Admin e Client separados
	HandlerAdminUser  IHandlerAdminUser  // Gestão de usuários admin
	HandlerClientUser IHandlerClientUser // Gestão de usuários cliente
	HandlerUserAccess IHandlerUserAccess // Gestão de acesso a organizações/projetos
	EventService      *utils.EventService
	// Staff Management System
	HandlerStaffAvailability IHandlerStaffAvailability // Disponibilidade de equipe
	HandlerStaffSchedule     IHandlerStaffSchedule     // Escalas de trabalho
	HandlerStaffAttendance   IHandlerStaffAttendance   // Presença e consumo
	HandlerStaffStock        IHandlerStaffStock        // Estoque operacional
	HandlerStaffCommission   IHandlerStaffCommission   // Comissões
	HandlerStaffDashboard    IHandlerStaffDashboard    // Dashboard e relatórios
}

func (h *Handlers) Inject(repo *repositories.DBconn) {
	// Role Handler precisa ser criado primeiro para ser usado pelo UserHandler
	h.HandlerRole = NewRoleHandler(repo.Roles, repo.Permissions, repo.Modules, repo.Plans, repo.Admins, repo.Clients)

	// Admin Audit Log Handler (precisa ser criado antes do UserHandler para injeção)
	h.HandlerAdminAuditLog = NewAdminAuditLogHandler(
		repo.AdminAuditLogs,
		repo.Organizations,
		repo.Projects,
	)

	// Injetar handler de auditoria no RoleHandler
	h.HandlerRole.SetAdminAuditHandler(h.HandlerAdminAuditLog)

	h.HandlerProducts = NewSourceHandlerProducts(repo)
	h.HandlerAuth = NewAuthHandler(repo)
	h.HandlerOrder = NewOrderHandler(repo.Orders, repo.Products, repo.KitchenQueue)
	h.HandlerOrganization = NewSourceHandlerOrganization(repo, repo.DB)
	h.HandlerTables = NewSourceHandlerTables(repo)
	h.HandlerWaitlist = NewSourceHandlerWaitlist(repo)
	h.HandlerReservation = NewSourceHandlerReservation(repo)
	h.HandlerCustomer = NewSourceHandlerCustomer(repo)
	h.HandlerProject = NewProjectHandler(repo.Projects, repo.Settings, repo.Notifications, repo.CascadeDelete)
	h.HandlerSettings = NewSettingsHandler(repo.Settings)
	h.HandlerDisplaySettings = NewDisplaySettingsHandler(repo.DisplaySettings)
	h.HandlerThemeCustomization = NewThemeCustomizationHandler(repo.ThemeCustomization)
	h.HandlerEnvironment = NewEnvironmentHandler(repo.Environments)
	h.HandlerNotification = NewNotificationHandler(
		repo.Notifications,
		repo.Projects,
		repo.Reservations,
		repo.Customers,
		repo.Tables,
		repo.Settings,
	)
	h.HandlerReports = NewReportsHandler(repo)
	h.HandlerTag = NewSourceHandlerTag(repo)
	h.HandlerMenu = NewSourceHandlerMenu(repo)
	h.HandlerCategory = NewSourceHandlerCategory(repo)
	h.HandlerSubcategory = NewSourceHandlerSubcategory(repo)
	h.HandlerOnboarding = NewOnboardingHandler(repo)

	// Image Management Service e Handler
	imageManagementSvc := service.NewImageManagementService(repo.FileReference, repo.EntityFileReference, "./uploads")
	h.HandlerImageManagement = NewHandlerImageManagement(imageManagementSvc)
	h.ImageManagementService = imageManagementSvc // Armazenar o service direto para o Upload server

	// Plan Change Request Handler
	h.HandlerPlanChangeRequest = NewPlanChangeRequestHandler(repo.PlanChangeRequests, h.HandlerRole)

	// Limits Handler - Verificação de limites de plano
	h.HandlerLimits = NewLimitHandler(
		repo.Plans,
		repo.Tables,
		repo.Roles,
		repo.Products,
		repo.Reservations,
		repo.Modules,
	)

	// Sidebar Config Handler
	h.HandlerSidebarConfig = NewSidebarConfigHandler(repo.SidebarConfig)

	// Client Audit Log Handler (módulo opcional)
	h.HandlerClientAuditLog = NewClientAuditLogHandler(repo.ClientAuditLogs)

	// Novos handlers para Admin e Client separados
	h.HandlerAdminUser = NewAdminUserHandler(repo)
	h.HandlerClientUser = NewClientUserHandler(repo)
	h.HandlerUserAccess = NewUserAccessHandler(repo)

	// EventService para disparo de notificações
	h.EventService = utils.NewEventService(repo.Notifications, repo.Projects, repo.Settings)

	// Staff Management System
	h.HandlerStaffAvailability = NewStaffAvailabilityHandler(repo)
	h.HandlerStaffSchedule = NewStaffScheduleHandler(repo)
	h.HandlerStaffAttendance = NewStaffAttendanceHandler(repo)
	h.HandlerStaffStock = NewStaffStockHandler(repo)
	h.HandlerStaffCommission = NewStaffCommissionHandler(repo)
	h.HandlerStaffDashboard = NewStaffDashboardHandler(repo)
}
