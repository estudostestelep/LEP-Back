package resource

import (
	"fmt"
	"lep/handler"
	"lep/server"

	"lep/repositories"
)

var Repository repositories.DBconn
var Handlers handler.Handlers
var ServersControllers server.ServerController

func Inject() {
	db, err := OpenConnDBPostgres2()
	if err != nil {
		panic(fmt.Sprintf("Conexão com banco falhou: %v", err))
	}
	server.Start(db)
	Repository.InjectProstgres(db)
	Handlers.Inject(&Repository)
	ServersControllers.Inject(&Handlers)
}
