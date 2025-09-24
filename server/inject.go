package server

import "lep/handler"

type ServerController struct {
	SourceUsers        IServerUsers
	SourceProducts     IServerProducts
	SourceAuth         IServerAuth
	SourceOrders       IOrderServer
	SourceOrganization IServerOrganization
	SourceTables       IServerTables
	SourceWaitlist     IServerWaitlist
	SourceReservation  IServerReservation
	SourceCustomer     IServerCustomer
	SourceProject      IProjectServer
	SourceSettings     ISettingsServer
	SourceEnvironment  IEnvironmentServer
	SourceNotification *NotificationServer
	SourceReports      IReportsServer
}

func (h *ServerController) Inject(handler *handler.Handlers) {
	h.SourceUsers = NewSourceServerUsers(handler)
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
}
