package server

import "lep/handler"

type ServerController struct {
	SourcePurchasesIServerPurchases
	SourceUsers    IServerUsers
	SourceProducts IServerProducts
	SourceAuth     IServerAuth
}

func (h *ServerController) Inject(handler *handler.Handlers) {
	h.SourcePurchases = NewSourceServerPurchases(handler)
	h.SourceUsers = NewSourceServerUsers(handler)
	h.SourceProducts = NewSourceServerProducts(handler)
	h.SourceAuth = NewSourceServerAuth(handler)
}
