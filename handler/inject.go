package handler

import "lep/repositories"

type Handlers struct {
	HandlerUser         IHandlerUser
	HandlerProducts     IHandlerProducts
	HandlerAuth         IHandlerAuth
	HandlerOrder        IOrderHandler
	HandlerTables       IHandlerTables
	HandlerWaitlist     IHandlerWaitlist
	HandlerReservation  IHandlerReservation
	HandlerCustomer     IHandlerCustomer
	HandlerProject      IProjectHandler
	HandlerSettings     ISettingsHandler
	HandlerEnvironment  IEnvironmentHandler
	HandlerNotification *NotificationHandler
}

func (h *Handlers) Inject(repo *repositories.DBconn) {
	h.HandlerUser = NewSourceHandlerUser(repo)
	h.HandlerProducts = NewSourceHandlerProducts(repo)
	h.HandlerAuth = NewAuthHandler(repo)
	h.HandlerOrder = NewOrderHandler(repo.Orders, repo.Products, repo.KitchenQueue)
	h.HandlerTables = NewSourceHandlerTables(repo)
	h.HandlerWaitlist = NewSourceHandlerWaitlist(repo)
	h.HandlerReservation = NewSourceHandlerReservation(repo)
	h.HandlerCustomer = NewSourceHandlerCustomer(repo)
	h.HandlerProject = NewProjectHandler(repo.Projects, repo.Settings, repo.Notifications)
	h.HandlerSettings = NewSettingsHandler(repo.Settings)
	h.HandlerEnvironment = NewEnvironmentHandler(repo.Environments)
	h.HandlerNotification = NewNotificationHandler(repo.Notifications, repo.Projects)
}
