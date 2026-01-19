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

// GetUserOrganizationsWithNames busca organizações do usuário com seus nomes
func (r *resourceAuth) GetUserOrganizationsWithNames(userId string) ([]UserOrganizationWithName, error) {
	// Buscar relacionamentos user-organization
	userOrgs, err := r.repo.UserOrganizations.ListByUser(userId)
	if err != nil {
		return nil, err
	}

	// Criar slice de resposta enriquecida
	result := make([]UserOrganizationWithName, 0, len(userOrgs))

	// Para cada organização, buscar o nome
	for _, userOrg := range userOrgs {
		org, err := r.repo.Organizations.GetOrganizationById(userOrg.OrganizationId)
		if err != nil {
			// Se organização não for encontrada, pular
			continue
		}
		result = append(result, UserOrganizationWithName{
			Id:               userOrg.Id,
			UserId:           userOrg.UserId,
			OrganizationId:   userOrg.OrganizationId,
			OrganizationName: org.Name,
			Role:             userOrg.Role,
			Active:           userOrg.Active,
			CreatedAt:        userOrg.CreatedAt,
			UpdatedAt:        userOrg.UpdatedAt,
			DeletedAt:        userOrg.DeletedAt,
		})
	}

	return result, nil
}

// GetUserProjectsWithNames busca projetos do usuário com seus nomes
func (r *resourceAuth) GetUserProjectsWithNames(userId string) ([]UserProjectWithName, error) {
	// Buscar relacionamentos user-project
	userProjs, err := r.repo.UserProjects.ListByUser(userId)
	if err != nil {
		return nil, err
	}

	// Criar slice de resposta enriquecida
	result := make([]UserProjectWithName, 0, len(userProjs))

	// Para cada projeto, buscar o nome e também o nome da organização
	for _, userProj := range userProjs {
		proj, err := r.repo.Projects.GetProjectById(userProj.ProjectId)
		if err != nil {
			// Se projeto não for encontrado, pular
			continue
		}

		// Buscar nome da organização
		orgName := ""
		if org, err := r.repo.Organizations.GetOrganizationById(proj.OrganizationId); err == nil && org != nil {
			orgName = org.Name
		}

		result = append(result, UserProjectWithName{
			Id:               userProj.Id,
			UserId:           userProj.UserId,
			ProjectId:        userProj.ProjectId,
			ProjectName:      proj.Name,
			OrganizationId:   proj.OrganizationId,
			OrganizationName: orgName, // Nome da organização pai
			Role:             userProj.Role,
			Active:           userProj.Active,
			CreatedAt:        userProj.CreatedAt,
			UpdatedAt:        userProj.UpdatedAt,
			DeletedAt:        userProj.DeletedAt,
		})
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

func NewAuthHandler(repo *repositories.DBconn) IHandlerAuth {
	return &resourceAuth{repo: repo}
}
