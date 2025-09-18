package server

import "lep/handler"

type ServerController struct {
	SourceUsers        IServerUsers
	SourceProducts     IServerProducts
	SourceAuth         IServerAuth
	SourceOrders       IOrderServer
	SourceTables       IServerTables
	SourceWaitlist     IServerWaitlist
	SourceReservation  IServerReservation
	SourceCustomer     IServerCustomer
}

func (h *ServerController) Inject(handler *handler.Handlers) {
	h.SourceUsers = NewSourceServerUsers(handler)
	h.SourceProducts = NewSourceServerProducts(handler)
	h.SourceAuth = NewSourceServerAuth(handler)
	h.SourceOrders = NewOrderServer(handler.HandlerOrder)
	h.SourceTables = NewSourceServerTables(handler)
	h.SourceWaitlist = NewSourceServerWaitlist(handler)
	h.SourceReservation = NewSourceServerReservation(handler)
	h.SourceCustomer = NewSourceServerCustomer(handler)
}
