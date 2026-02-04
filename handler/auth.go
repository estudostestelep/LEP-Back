package handler

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"lep/repositories"
	"lep/repositories/models"
	"time"
)

type resourceAuth struct {
	repo *repositories.DBconn
}

// AuthUser representa um usuário autenticado (Admin ou Client)
type AuthUser struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	UserType string    `json:"user_type"` // "admin" ou "client"
	Active   bool      `json:"active"`
}

type IHandlerAuth interface {
	PostTokenForUser(userId uuid.UUID, email, token string) error
	Logout(token string) error
	VerificationToken(token string) (*AuthUser, error)
	GetClientOrganizationsWithNames(clientId string) ([]UserOrganizationWithName, error)
	GetClientProjectsWithNames(clientId string) ([]UserProjectWithName, error)
	RecordAccessLog(userId, ip, userAgent string) error
	GetAccessLogs(userId string, page, perPage int) (*models.AccessLogPaginatedResponse, error)
	// Admin Roles - para acesso à área administrativa
	GetAdminRolesInfo(adminId string) ([]UserAdminRoleInfo, error)
	AdminHasAdminScopeRole(adminId string) (bool, error)
}

func (r *resourceAuth) PostTokenForUser(userId uuid.UUID, email, token string) error {
	loggedList := &models.LoggedLists{
		LoggedListId: uuid.New(),
		Token:        token,
		UserEmail:    email,
		UserId:       userId,
	}

	if err := r.repo.LoggedLists.CreateLoggedList(loggedList); err != nil {
		if err := r.repo.LoggedLists.DeleteLoggedList(loggedList.Token); err != nil {
			return fmt.Errorf("falha ao remover token existente: %v", err)
		}

		if err := r.repo.LoggedLists.CreateLoggedList(loggedList); err != nil {
			return fmt.Errorf("falha ao criar registro na LoggedLists: %v", err)
		}
	}

	return nil
}

func (r *resourceAuth) Logout(token string) error {
	bannedList := &models.BannedLists{
		Token: token,
	}

	if err := r.repo.BannedLists.CreateBannedList(bannedList); err != nil {
		return fmt.Errorf("falha ao criar registro na BannedLists: %v", err)
	}

	if err := r.repo.LoggedLists.DeleteLoggedList(token); err != nil {
		return fmt.Errorf("falha ao remover da LoggedLists: %v", err)
	}

	r.cleanupExpiredTokens()

	return nil
}

func (r *resourceAuth) cleanupExpiredTokens() {
	cutoffTime := time.Now().AddDate(0, 0, -7)

	resp, err := r.repo.BannedLists.GetBannedAllList()
	if err != nil || resp == nil {
		return
	}

	for _, item := range *resp {
		if item.CreatedAt.Before(cutoffTime) {
			r.repo.BannedLists.DeleteBannedList(item.BannedListId)
		}
	}
}

func (r *resourceAuth) VerificationToken(token string) (*AuthUser, error) {
	logged, err := r.repo.LoggedLists.GetLoggedToken(token)
	if err != nil {
		return nil, err
	}

	if logged == nil {
		return nil, errors.New("Not found")
	}

	// Tentar buscar como Admin primeiro
	admin, err := r.repo.Admins.GetAdminById(logged.UserId.String())
	if err == nil && admin != nil {
		return &AuthUser{
			Id:       admin.Id,
			Name:     admin.Name,
			Email:    admin.Email,
			UserType: "admin",
			Active:   admin.Active,
		}, nil
	}

	// Se não for Admin, tentar buscar como Client
	client, err := r.repo.Clients.GetClientById(logged.UserId.String())
	if err == nil && client != nil {
		return &AuthUser{
			Id:       client.Id,
			Name:     client.Name,
			Email:    client.Email,
			UserType: "client",
			Active:   client.Active,
		}, nil
	}

	return nil, errors.New("usuário não encontrado")
}

// GetClientOrganizationsWithNames busca organizações do cliente com seus nomes via client_roles
func (r *resourceAuth) GetClientOrganizationsWithNames(clientId string) ([]UserOrganizationWithName, error) {
	// Buscar todos os client_roles do cliente
	clientRoles, err := r.repo.Roles.GetClientRoles(clientId, "")
	if err != nil {
		return nil, err
	}

	// Agrupar por organização (evitar duplicatas)
	orgMap := make(map[uuid.UUID]UserOrganizationWithName)

	for _, cr := range clientRoles {
		// Se já temos essa organização, pular
		if _, exists := orgMap[cr.OrganizationId]; exists {
			continue
		}

		// Buscar nome da organização
		org, err := r.repo.Organizations.GetOrganizationById(cr.OrganizationId)
		if err != nil {
			continue
		}

		roleName := ""
		if cr.Role != nil {
			roleName = cr.Role.DisplayName
		}

		orgMap[cr.OrganizationId] = UserOrganizationWithName{
			Id:               cr.Id,
			UserId:           cr.ClientId,
			OrganizationId:   cr.OrganizationId,
			OrganizationName: org.Name,
			Role:             roleName,
			Active:           cr.Active,
			CreatedAt:        cr.CreatedAt,
			UpdatedAt:        cr.UpdatedAt,
			DeletedAt:        cr.DeletedAt,
		}
	}

	// Converter map para slice
	result := make([]UserOrganizationWithName, 0, len(orgMap))
	for _, org := range orgMap {
		result = append(result, org)
	}

	return result, nil
}

// GetClientProjectsWithNames busca projetos do cliente com seus nomes via client_roles
func (r *resourceAuth) GetClientProjectsWithNames(clientId string) ([]UserProjectWithName, error) {
	// Buscar todos os client_roles do cliente
	clientRoles, err := r.repo.Roles.GetClientRoles(clientId, "")
	if err != nil {
		return nil, err
	}

	// Agrupar por projeto (evitar duplicatas)
	projMap := make(map[uuid.UUID]UserProjectWithName)

	for _, cr := range clientRoles {
		if cr.ProjectId == nil {
			continue
		}

		// Se já temos esse projeto, pular
		if _, exists := projMap[*cr.ProjectId]; exists {
			continue
		}

		// Buscar nome do projeto
		proj, err := r.repo.Projects.GetProjectById(*cr.ProjectId)
		if err != nil {
			continue
		}

		// Buscar nome da organização
		orgName := ""
		if org, err := r.repo.Organizations.GetOrganizationById(proj.OrganizationId); err == nil && org != nil {
			orgName = org.Name
		}

		roleName := ""
		if cr.Role != nil {
			roleName = cr.Role.DisplayName
		}

		projMap[*cr.ProjectId] = UserProjectWithName{
			Id:               cr.Id,
			UserId:           cr.ClientId,
			ProjectId:        *cr.ProjectId,
			ProjectName:      proj.Name,
			OrganizationId:   proj.OrganizationId,
			OrganizationName: orgName,
			Role:             roleName,
			Active:           cr.Active,
			CreatedAt:        cr.CreatedAt,
			UpdatedAt:        cr.UpdatedAt,
			DeletedAt:        cr.DeletedAt,
		}
	}

	// Converter map para slice
	result := make([]UserProjectWithName, 0, len(projMap))
	for _, proj := range projMap {
		result = append(result, proj)
	}

	return result, nil
}

// RecordAccessLog registra um log de acesso do usuário
func (r *resourceAuth) RecordAccessLog(userId, ip, userAgent string) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return fmt.Errorf("ID de usuário inválido: %v", err)
	}

	accessLog := &models.AccessLog{
		Id:        uuid.New(),
		UserId:    userUUID,
		IP:        ip,
		UserAgent: userAgent,
		Location:  "", // TODO: Implementar geolocalização por IP
		LoginAt:   time.Now(),
		CreatedAt: time.Now(),
	}

	return r.repo.AccessLogs.Create(accessLog)
}

// GetAccessLogs retorna logs de acesso paginados de um usuário
func (r *resourceAuth) GetAccessLogs(userId string, page, perPage int) (*models.AccessLogPaginatedResponse, error) {
	return r.repo.AccessLogs.GetByUserId(userId, page, perPage)
}

// GetAdminRolesInfo busca todos os cargos de escopo "admin" atribuídos ao admin
func (r *resourceAuth) GetAdminRolesInfo(adminId string) ([]UserAdminRoleInfo, error) {
	fmt.Printf("[DEBUG] GetAdminRolesInfo - Buscando cargos admin para adminId: %s\n", adminId)

	adminRoles, err := r.repo.Roles.GetAdminRoles(adminId)
	if err != nil {
		fmt.Printf("[DEBUG] GetAdminRolesInfo - Erro ao buscar: %v\n", err)
		return nil, err
	}

	fmt.Printf("[DEBUG] GetAdminRolesInfo - Encontrados %d admin_roles\n", len(adminRoles))

	result := make([]UserAdminRoleInfo, 0, len(adminRoles))
	for _, ar := range adminRoles {
		if ar.Role == nil {
			fmt.Printf("[DEBUG] GetAdminRolesInfo - AdminRole %s tem Role nil, pulando\n", ar.Id)
			continue
		}

		fmt.Printf("[DEBUG] GetAdminRolesInfo - AdminRole encontrado: id=%s, role_name=%s, scope=%s, active=%v\n",
			ar.Id, ar.Role.Name, ar.Role.Scope, ar.Active)

		result = append(result, UserAdminRoleInfo{
			Id:              ar.Id,
			RoleId:          ar.RoleId,
			RoleName:        ar.Role.Name,
			RoleDisplayName: ar.Role.DisplayName,
			Scope:           ar.Role.Scope,
			HierarchyLevel:  ar.Role.HierarchyLevel,
			Active:          ar.Active,
		})
	}

	fmt.Printf("[DEBUG] GetAdminRolesInfo - Retornando %d cargos admin\n", len(result))
	return result, nil
}

// AdminHasAdminScopeRole verifica se o admin possui pelo menos um cargo de escopo "admin"
func (r *resourceAuth) AdminHasAdminScopeRole(adminId string) (bool, error) {
	adminRoles, err := r.GetAdminRolesInfo(adminId)
	if err != nil {
		return false, err
	}

	return len(adminRoles) > 0, nil
}

func NewAuthHandler(repo *repositories.DBconn) IHandlerAuth {
	return &resourceAuth{repo: repo}
}
