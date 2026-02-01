package server

import (
	"lep/handler"
	"lep/repositories/models"
	"log"

	"github.com/google/uuid"
)

type ServerController struct {
	SourceProducts           IServerProducts
	SourceAuth               IServerAuth
	SourceOrders             IOrderServer
	SourceOrganization       IServerOrganization
	SourceTables             IServerTables
	SourceWaitlist           IServerWaitlist
	SourceReservation        IServerReservation
	SourceCustomer           IServerCustomer
	SourceProject            IProjectServer
	SourceSettings           ISettingsServer
	SourceDisplaySettings    IDisplaySettingsServer
	SourceThemeCustomization IThemeCustomizationServer
	SourceEnvironment        IEnvironmentServer
	SourceNotification       *NotificationServer
	SourceReports            IReportsServer
	SourcePublic             IServerPublic
	SourceUpload             *ResourceUpload // Mudado para *ResourceUpload para permitir injeção
	SourceTag                IServerTag
	SourceMenu               IServerMenu
	SourceCategory           IServerCategory
	SourceSubcategory        IServerSubcategory
	SourceImageManagement    IServerImageManagement
	SourceOnboarding         IOnboardingServer
	SourceRole               *RoleServer
	SourcePlanChangeRequest  IPlanChangeRequestServer
	SourceAdmin              *AdminController
	SourceSidebarConfig      ISidebarConfigServer
	SourceAccessLog          IAccessLogServer
	SourceAdminAuditLog      IAdminAuditLogServer  // Logs de auditoria admin (read-only)
	SourceClientAuditLog     IClientAuditLogServer // Logs de auditoria de cliente (módulo opcional)
	// Novos servers para autenticação separada
	SourceAuthAdmin  IServerAuthAdmin  // Login de administradores
	SourceAuthClient IServerAuthClient // Login de clientes
	SourceTenant     IServerTenant     // Resolver tenant por slug
	// CRUD de admins e clients (novo sistema de usuários)
	SourceAdminUsers  IServerAdminUsers  // CRUD de admins (tabela admins)
	SourceClientUsers IServerClientUsers // CRUD de clients (tabela clients)
}

func (h *ServerController) Inject(handler *handler.Handlers) {

	// Nota: Permissões são gerenciadas via roles, não diretamente no admin
	admin := &models.Admin{
		Name:     "Pablo",
		Email:    "pablo@lep.com",
		Password: "senha123",
		Active:   true,
	}
	errAdmin := handler.HandlerAdminUser.CreateAdmin(admin)
	if errAdmin != nil {
		log.Printf("❌ Error creating admin: %v", errAdmin)
	} else {
		log.Printf("✅ Admin created: Pablo (pablo@lep.com)")
	}

	// Create default organization
	orgId := uuid.New()
	org := &models.Organization{
		Id:    orgId,
		Name:  "Default Organization",
		Email: "admin@default.com",
		Phone: "+5511999999999",
	}
	errOrg := handler.HandlerOrganization.CreateOrganization(org)
	if errOrg != nil {
		log.Printf("❌ Error creating default organization: %v", errOrg)
	} else {
		log.Printf("✅ Default organization created: Default Organization")

		// Create default project
		projectId := uuid.New()
		project := &models.Project{
			Id:             projectId,
			Name:           "Default Project",
			OrganizationId: orgId,
		}
		errProj := handler.HandlerProject.CreateProject(project)
		if errProj != nil {
			log.Printf("❌ Error creating default project: %v", errProj)
		} else {
			log.Printf("✅ Default project created: Default Project")

			// Associar usuário à organização e projeto via UserRole (novo sistema)
			// O master_admin já tem acesso global, mas podemos atribuir um role específico se necessário
			log.Printf("✅ User has master_admin permission - global access granted")
		}
	}

	h.SourceProducts = NewSourceServerProducts(handler)
	h.SourceAuth = NewSourceServerAuth(handler)
	h.SourceOrders = NewOrderServer(handler.HandlerOrder)
	h.SourceOrganization = NewSourceServerOrganization(handler)
	h.SourceTables = NewSourceServerTables(handler)
	h.SourceWaitlist = NewSourceServerWaitlist(handler)
	h.SourceReservation = NewSourceServerReservation(handler)
	h.SourceCustomer = NewSourceServerCustomer(handler)
	h.SourceProject = NewProjectServer(handler.HandlerProject)
	h.SourceSettings = NewSettingsServer(handler.HandlerSettings)
	h.SourceDisplaySettings = NewDisplaySettingsServer(handler.HandlerDisplaySettings)
	h.SourceThemeCustomization = NewThemeCustomizationServer(handler.HandlerThemeCustomization)
	h.SourceEnvironment = NewEnvironmentServer(handler.HandlerEnvironment)
	h.SourceNotification = NewNotificationServer(handler.HandlerNotification)
	h.SourceReports = NewReportsServer(handler.HandlerReports)
	h.SourcePublic = NewSourceServerPublic(handler)

	// Criar Upload controller e injetar ImageManagement service
	uploadServer := NewSourceServerUpload()
	h.SourceImageManagement = NewServerImageManagement(handler.HandlerImageManagement)
	// Injetar o service direto (não o handler)
	uploadServer.SetImageManagementService(handler.ImageManagementService)
	h.SourceUpload = uploadServer

	h.SourceTag = NewSourceServerTag(handler)
	h.SourceMenu = NewSourceServerMenu(handler)
	h.SourceCategory = NewSourceServerCategory(handler)
	h.SourceSubcategory = NewSourceServerSubcategory(handler)
	h.SourceOnboarding = NewOnboardingServer(handler.HandlerOnboarding)
	h.SourceRole = NewRoleServer(handler.HandlerRole)
	h.SourceRole.SetLimitHandler(handler.HandlerLimits) // Injetar handler de limites
	// Nota: auditoria agora é gerenciada pelo RoleHandler internamente
	h.SourcePlanChangeRequest = NewPlanChangeRequestServer(handler.HandlerPlanChangeRequest)
	// AdminController is initialized separately with DB in resource/inject.go

	// Sidebar Config Server
	h.SourceSidebarConfig = NewSidebarConfigServer(handler.HandlerSidebarConfig)

	// Access Log Server
	h.SourceAccessLog = NewAccessLogController(handler)

	// Admin Audit Log Server (read-only)
	h.SourceAdminAuditLog = NewAdminAuditLogServer(handler.HandlerAdminAuditLog)

	// Client Audit Log Server (módulo opcional)
	h.SourceClientAuditLog = NewClientAuditLogServer(handler.HandlerClientAuditLog)

	// Novos servers para autenticação separada Admin/Client
	h.SourceAuthAdmin = NewSourceServerAuthAdmin(handler)
	h.SourceAuthClient = NewSourceServerAuthClient(handler)
	h.SourceTenant = NewSourceServerTenant(handler)

	// CRUD de admins e clients (novo sistema de usuários)
	h.SourceAdminUsers = NewSourceServerAdminUsers(handler)
	h.SourceClientUsers = NewSourceServerClientUsers(handler)
}
