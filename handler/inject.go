package handler

import "lep/repositories"

type Handlers struct {
	HandlerPurchasesIHandlerPurchases
	HandlerUser     IHandlerUser
	HandlerProducts IHandlerProducts
	HandlerAuth     IHandlerAuth
}

func (h *Handlers) Inject(repo *repositories.DBconn) {
	h.HandlerPurchases = NewSourceHandlerPurchases(repo)
	h.HandlerUser = NewSourceHandlerUser(repo)
	h.HandlerProducts = NewSourceHandlerProducts(repo)
	h.HandlerAuth = NewAuthHandler(repo)
}
