package resource

import (
	"lep/handler"
	"lep/server"

	"lep/repositories"
)

var Repository repositories.DBconn
var Handlers handler.Handlers
var ServersControllers server.ServerController

func Inject() {
	db, err := OpenConnDBPostgres()
	if err != nil {
		panic("Conecxão falhou")
	}
	server.Start(db)
	Repository.InjectProstgres(db)
	Handlers.Inject(&Repository)
	ServersControllers.Inject(&Handlers)
}
