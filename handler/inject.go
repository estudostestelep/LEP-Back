package handler

import (
	"lep/repositories"
	"lep/service"

	"gorm.io/gorm"
)

type Handlers struct {
	HandlerUser               IHandlerUser
	HandlerUserOrganization   IHandlerUserOrganization
	HandlerUserProject        IHandlerUserProject
	HandlerUserAccess         UserAccessHandler
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
	ImageManagementService    service.IImageManagementService // Service direto para o Upload server
}

func (h *Handlers) Inject(repo *repositories.DBconn, db interface{}) {
	h.HandlerUser = NewSourceHandlerUser(repo)
	h.HandlerUserOrganization = NewSourceHandlerUserOrganization(repo)
	h.HandlerUserProject = NewSourceHandlerUserProject(repo)
	h.HandlerUserAccess = NewUserAccessHandler(db)
	h.HandlerProducts = NewSourceHandlerProducts(repo)
	h.HandlerAuth = NewAuthHandler(repo)
	h.HandlerOrder = NewOrderHandler(repo.Orders, repo.Products, repo.KitchenQueue)
	h.HandlerOrganization = NewSourceHandlerOrganization(repo)
	h.HandlerTables = NewSourceHandlerTables(repo)
	h.HandlerWaitlist = NewSourceHandlerWaitlist(repo)
	h.HandlerReservation = NewSourceHandlerReservation(repo)
	h.HandlerCustomer = NewSourceHandlerCustomer(repo)
	h.HandlerProject = NewProjectHandler(repo.Projects, repo.Settings, repo.Notifications)
	h.HandlerSettings = NewSettingsHandler(repo.Settings)
	h.HandlerDisplaySettings = NewDisplaySettingsHandler(repo.DisplaySettings)
	h.HandlerThemeCustomization = NewThemeCustomizationHandler(repo.ThemeCustomization)
	h.HandlerEnvironment = NewEnvironmentHandler(repo.Environments)
	h.HandlerNotification = NewNotificationHandler(repo.Notifications, repo.Projects)
	h.HandlerReports = NewReportsHandler(repo)
	h.HandlerTag = NewSourceHandlerTag(repo)
	h.HandlerMenu = NewSourceHandlerMenu(repo)
	h.HandlerCategory = NewSourceHandlerCategory(repo)
	h.HandlerSubcategory = NewSourceHandlerSubcategory(repo)
	h.HandlerOnboarding = NewOnboardingHandler(repo)

	// Image Management Service e Handler
	// Nota: db é *gorm.DB, necessário para os novos repositories
	gormDB := db.(*gorm.DB)
	fileRefRepo := repositories.NewFileReferenceRepository(gormDB)
	entityFileRefRepo := repositories.NewEntityFileReferenceRepository(gormDB)
	imageManagementSvc := service.NewImageManagementService(fileRefRepo, entityFileRefRepo, "./uploads")
	h.HandlerImageManagement = NewHandlerImageManagement(imageManagementSvc)
	h.ImageManagementService = imageManagementSvc // Armazenar o service direto para o Upload server

	// Role & Permission Handler
	h.HandlerRole = NewRoleHandler(repo.Roles, repo.Permissions, repo.Modules, repo.Packages)
}
