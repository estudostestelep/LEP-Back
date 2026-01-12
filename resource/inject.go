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
	db, err := OpenConnDBPostgres()
	if err != nil {
		panic(fmt.Sprintf("Conexão com banco falhou: %v", err))
	}
	server.Start(db)
	Repository.InjectPostgres(db)
	Handlers.Inject(&Repository, db)
	ServersControllers.Inject(&Handlers)
	// Initialize AdminController with DB
	ServersControllers.SourceAdmin = &server.AdminController{DB: db}

	// Seed de roles e permissões
	if err := handler.SeedRolesAndPermissions(db); err != nil {
		fmt.Printf("⚠️ Erro no seed de roles: %v\n", err)
	}
}
