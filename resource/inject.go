package resource

import (
	"fmt"
	"lep/handler"
	"lep/repositories"
	"lep/repositories/models"
	"lep/server"

	"github.com/google/uuid"
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
	// Set repository for NotificationServer (needed for trigger-scheduled endpoint)
	ServersControllers.SourceNotification.SetRepo(&Repository)

	// Seed de roles e permissões
	if err := handler.SeedRolesAndPermissions(db); err != nil {
		fmt.Printf("⚠️ Erro no seed de roles: %v\n", err)
	}

	// Seed da organização demo (criada automaticamente se não existir)
	if err := handler.SeedDemoOrganization(db); err != nil {
		fmt.Printf("⚠️ Erro no seed demo: %v\n", err)
	} else {
		fmt.Printf("Executado seed demo: %v\n", err)
	}

	// Atribuir role super_admin ao admin principal (pablo@lep.com)
	assignSuperAdminRole(&Handlers)
}

// assignSuperAdminRole atribui o role super_admin ao admin principal
func assignSuperAdminRole(handlers *handler.Handlers) {
	// Buscar admin pelo email
	admin, err := handlers.HandlerAdminUser.GetAdminByEmail("pablo@lep.com")
	if err != nil || admin == nil {
		fmt.Printf("⚠️ Admin pablo@lep.com não encontrado para atribuir role\n")
		return
	}

	// Verificar se já tem roles atribuídos
	existingRoles, _ := handlers.HandlerAdminUser.GetAdminRoles(admin.Id.String())
	if len(existingRoles) > 0 {
		fmt.Printf("✅ Admin pablo@lep.com já possui roles atribuídos\n")
		return
	}

	// Buscar role super_admin
	superAdminRole, err := handlers.HandlerRole.GetRoleByName("super_admin")
	if err != nil || superAdminRole == nil {
		fmt.Printf("⚠️ Role super_admin não encontrado\n")
		return
	}

	// Criar associação admin-role
	adminRole := &models.AdminRole{
		Id:             uuid.New(),
		AdminId:        admin.Id,
		RoleId:         superAdminRole.Id,
		OrganizationId: nil, // NULL = cargo admin global
		Active:         true,
	}

	if err := handlers.HandlerAdminUser.AssignRoleToAdmin(adminRole); err != nil {
		fmt.Printf("❌ Erro ao atribuir role super_admin: %v\n", err)
	} else {
		fmt.Printf("✅ Role super_admin atribuído ao admin Pablo (pablo@lep.com)\n")
	}
}
