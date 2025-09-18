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
}

func (h *Handlers) Inject(repo *repositories.DBconn) {
	h.HandlerUser = NewSourceHandlerUser(repo)
	h.HandlerProducts = NewSourceHandlerProducts(repo)
	h.HandlerAuth = NewAuthHandler(repo)
	h.HandlerOrder = NewOrderHandler(repo.Orders)
	h.HandlerTables = NewSourceHandlerTables(repo)
	h.HandlerWaitlist = NewSourceHandlerWaitlist(repo)
	h.HandlerReservation = NewSourceHandlerReservation(repo)
	h.HandlerCustomer = NewSourceHandlerCustomer(repo)
}
