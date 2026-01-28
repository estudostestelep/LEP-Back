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

type IHandlerAuth interface {
	PostToken(user *models.User, token string) error
	Logout(token string) error
	VerificationToken(token string) (*models.User, error)
	GetUserOrganizationsWithNames(userId string) ([]UserOrganizationWithName, error)
	GetUserProjectsWithNames(userId string) ([]UserProjectWithName, error)
	RecordAccessLog(userId, ip, userAgent string) error
	GetAccessLogs(userId string, page, perPage int) (*models.AccessLogPaginatedResponse, error)
	// Admin Roles - para acesso à área administrativa
	GetUserAdminRoles(userId string) ([]UserAdminRoleInfo, error)
	UserHasAdminScopeRole(userId string) (bool, error)
}

func (r *resourceAuth) PostToken(user *models.User, token string) error {
	loggedList := &models.LoggedLists{
		LoggedListId: uuid.New(),
		Token:        token,
		UserEmail:    user.Email,
		UserId:       user.Id,
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

func (r *resourceAuth) VerificationToken(token string) (*models.User, error) {

	logged, err := r.repo.LoggedLists.GetLoggedToken(token)
	if err != nil {
		return nil, err
	}

	if logged == nil {
		return nil, errors.New("Not found")
	}

	user, err := r.repo.User.GetUserById(logged.UserId.String())
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserOrganizationsWithNames busca organizações do usuário com seus nomes via user_roles
func (r *resourceAuth) GetUserOrganizationsWithNames(userId string) ([]UserOrganizationWithName, error) {
	// Buscar todos os user_roles do usuário
	userRoles, err := r.repo.Roles.GetAllUserRoles(userId)
	if err != nil {
		return nil, err
	}

	// Agrupar por organização (evitar duplicatas)
	orgMap := make(map[uuid.UUID]UserOrganizationWithName)

	for _, ur := range userRoles {
		if ur.OrganizationId == nil {
			continue
		}

		// Se já temos essa organização, pular
		if _, exists := orgMap[*ur.OrganizationId]; exists {
			continue
		}

		// Buscar nome da organização
		org, err := r.repo.Organizations.GetOrganizationById(*ur.OrganizationId)
		if err != nil {
			continue
		}

		roleName := ""
		if ur.Role != nil {
			roleName = ur.Role.DisplayName
		}

		orgMap[*ur.OrganizationId] = UserOrganizationWithName{
			Id:               ur.Id,
			UserId:           ur.UserId,
			OrganizationId:   *ur.OrganizationId,
			OrganizationName: org.Name,
			Role:             roleName,
			Active:           ur.Active,
			CreatedAt:        ur.CreatedAt,
			UpdatedAt:        ur.UpdatedAt,
			DeletedAt:        ur.DeletedAt,
		}
	}

	// Converter map para slice
	result := make([]UserOrganizationWithName, 0, len(orgMap))
	for _, org := range orgMap {
		result = append(result, org)
	}

	return result, nil
}

// GetUserProjectsWithNames busca projetos do usuário com seus nomes via user_roles
func (r *resourceAuth) GetUserProjectsWithNames(userId string) ([]UserProjectWithName, error) {
	// Buscar todos os user_roles do usuário
	userRoles, err := r.repo.Roles.GetAllUserRoles(userId)
	if err != nil {
		return nil, err
	}

	// Agrupar por projeto (evitar duplicatas)
	projMap := make(map[uuid.UUID]UserProjectWithName)

	for _, ur := range userRoles {
		if ur.ProjectId == nil {
			continue
		}

		// Se já temos esse projeto, pular
		if _, exists := projMap[*ur.ProjectId]; exists {
			continue
		}

		// Buscar nome do projeto
		proj, err := r.repo.Projects.GetProjectById(*ur.ProjectId)
		if err != nil {
			continue
		}

		// Buscar nome da organização
		orgName := ""
		if org, err := r.repo.Organizations.GetOrganizationById(proj.OrganizationId); err == nil && org != nil {
			orgName = org.Name
		}

		roleName := ""
		if ur.Role != nil {
			roleName = ur.Role.DisplayName
		}

		projMap[*ur.ProjectId] = UserProjectWithName{
			Id:               ur.Id,
			UserId:           ur.UserId,
			ProjectId:        *ur.ProjectId,
			ProjectName:      proj.Name,
			OrganizationId:   proj.OrganizationId,
			OrganizationName: orgName,
			Role:             roleName,
			Active:           ur.Active,
			CreatedAt:        ur.CreatedAt,
			UpdatedAt:        ur.UpdatedAt,
			DeletedAt:        ur.DeletedAt,
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

// GetUserAdminRoles busca todos os cargos de escopo "admin" atribuídos ao usuário
// Esses são cargos com organization_id IS NULL (globais/administrativos)
func (r *resourceAuth) GetUserAdminRoles(userId string) ([]UserAdminRoleInfo, error) {
	// Buscar user_roles onde organization_id IS NULL (cargos admin globais)
	// e o role tem scope = 'admin'
	fmt.Printf("[DEBUG] GetUserAdminRoles - Buscando cargos admin para userId: %s\n", userId)

	userRoles, err := r.repo.Roles.GetUserRoles(userId, "", "admin")
	if err != nil {
		fmt.Printf("[DEBUG] GetUserAdminRoles - Erro ao buscar: %v\n", err)
		return nil, err
	}

	fmt.Printf("[DEBUG] GetUserAdminRoles - Encontrados %d user_roles (organization_id IS NULL + scope='admin')\n", len(userRoles))

	result := make([]UserAdminRoleInfo, 0, len(userRoles))
	for _, ur := range userRoles {
		if ur.Role == nil {
			fmt.Printf("[DEBUG] GetUserAdminRoles - UserRole %s tem Role nil, pulando\n", ur.Id)
			continue
		}

		fmt.Printf("[DEBUG] GetUserAdminRoles - UserRole encontrado: id=%s, role_name=%s, scope=%s, active=%v\n",
			ur.Id, ur.Role.Name, ur.Role.Scope, ur.Active)

		result = append(result, UserAdminRoleInfo{
			Id:              ur.Id,
			RoleId:          ur.RoleId,
			RoleName:        ur.Role.Name,
			RoleDisplayName: ur.Role.DisplayName,
			Scope:           ur.Role.Scope,
			HierarchyLevel:  ur.Role.HierarchyLevel,
			Active:          ur.Active,
		})
	}

	fmt.Printf("[DEBUG] GetUserAdminRoles - Retornando %d cargos admin\n", len(result))
	return result, nil
}

// UserHasAdminScopeRole verifica se o usuário possui pelo menos um cargo de escopo "admin"
func (r *resourceAuth) UserHasAdminScopeRole(userId string) (bool, error) {
	adminRoles, err := r.GetUserAdminRoles(userId)
	if err != nil {
		return false, err
	}

	return len(adminRoles) > 0, nil
}

func NewAuthHandler(repo *repositories.DBconn) IHandlerAuth {
	return &resourceAuth{repo: repo}
}
