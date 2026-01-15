package server

import (
	"lep/handler"
	"lep/repositories/models"
	"log"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ServerController struct {
	SourceUsers              IServerUsers
	SourceUserOrganization   IServerUserOrganization
	SourceUserProject        IServerUserProject
	SourceUserAccess         *UserAccessServer
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
}

func (h *ServerController) Inject(handler *handler.Handlers) {
	// Create default user with master_admin permission
	userId := uuid.New()
	var user models.User
	user.Id = userId
	user.Email = "pablo@lep.com"
	user.Password = "senha123"
	user.Name = "Pablo"
	user.Permissions = pq.StringArray{"master_admin"}
	// Para master_admin, não precisa passar orgId/projectId/roleId pois a função adiciona automaticamente a todas as orgs
	err := handler.HandlerUser.CreateUser(&user, "", "", "")
	if err != nil {
		log.Printf("❌ Error creating master_admin user: %v", err)
	} else {
		log.Printf("✅ Master admin user created: Pablo (pablo@lep.com)")

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

				// Associate user with organization
				userOrg := &models.UserOrganization{
					Id:             uuid.New(),
					UserId:         userId,
					OrganizationId: orgId,
					Role:           "owner",
				}
				errUserOrg := handler.HandlerUserOrganization.AddUserToOrganization(userOrg)
				if errUserOrg != nil {
					log.Printf("❌ Error assigning user to organization: %v", errUserOrg)
				} else {
					log.Printf("✅ User assigned to organization")
				}

				// Associate user with project
				userProject := &models.UserProject{
					Id:        uuid.New(),
					UserId:    userId,
					ProjectId: projectId,
					Role:      "admin",
				}
				errUserProj := handler.HandlerUserProject.AddUserToProject(userProject)
				if errUserProj != nil {
					log.Printf("❌ Error assigning user to project: %v", errUserProj)
				} else {
					log.Printf("✅ User assigned to project")
				}
			}
		}
	}

	h.SourceUsers = NewSourceServerUsers(handler)
	h.SourceUserOrganization = NewSourceServerUserOrganization(handler)
	h.SourceUserProject = NewSourceServerUserProject(handler)
	h.SourceUserAccess = NewUserAccessServer(handler.HandlerUserAccess)
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
	h.SourcePlanChangeRequest = NewPlanChangeRequestServer(handler.HandlerPlanChangeRequest)
	// AdminController is initialized separately with DB in resource/inject.go

	// Sidebar Config Server
	h.SourceSidebarConfig = NewSidebarConfigServer(handler.HandlerSidebarConfig)
}
