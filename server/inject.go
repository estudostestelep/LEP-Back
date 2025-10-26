package server

import "lep/handler"

type ServerController struct {
	SourceUsers             IServerUsers
	SourceUserOrganization  IServerUserOrganization
	SourceUserProject       IServerUserProject
	SourceUserAccess        *UserAccessServer
	SourceProducts          IServerProducts
	SourceAuth              IServerAuth
	SourceOrders            IOrderServer
	SourceOrganization      IServerOrganization
	SourceTables            IServerTables
	SourceWaitlist          IServerWaitlist
	SourceReservation       IServerReservation
	SourceCustomer          IServerCustomer
	SourceProject           IProjectServer
	SourceSettings          ISettingsServer
	SourceEnvironment       IEnvironmentServer
	SourceNotification      *NotificationServer
	SourceReports           IReportsServer
	SourcePublic            IServerPublic
	SourceUpload            *ResourceUpload  // Mudado para *ResourceUpload para permitir injeção
	SourceTag               IServerTag
	SourceMenu              IServerMenu
	SourceCategory          IServerCategory
	SourceSubcategory       IServerSubcategory
	SourceImageManagement   IServerImageManagement
	SourceAdmin             *AdminController
}

func (h *ServerController) Inject(handler *handler.Handlers) {
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
	// AdminController is initialized separately with DB in resource/inject.go
}
