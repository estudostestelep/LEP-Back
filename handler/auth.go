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
		fmt.Printf("DEBUG: Erro ao buscar user_organizations: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Encontrados %d user_organizations para user %s\n", len(userOrgs), userId)

	// Criar slice de resposta enriquecida
	result := make([]UserOrganizationWithName, 0, len(userOrgs))

	// Para cada organização, buscar o nome
	for _, userOrg := range userOrgs {
		fmt.Printf("DEBUG: Buscando organização %s\n", userOrg.OrganizationId.String())
		org, err := r.repo.Organizations.GetOrganizationById(userOrg.OrganizationId)
		if err != nil {
			// Se organização não for encontrada, pular
			fmt.Printf("DEBUG: Erro ao buscar organização %s: %v\n", userOrg.OrganizationId.String(), err)
			continue
		}

		fmt.Printf("DEBUG: Organização encontrada: %s\n", org.Name)
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

	fmt.Printf("DEBUG: Retornando %d organizações\n", len(result))
	return result, nil
}

// GetUserProjectsWithNames busca projetos do usuário com seus nomes
func (r *resourceAuth) GetUserProjectsWithNames(userId string) ([]UserProjectWithName, error) {
	// Buscar relacionamentos user-project
	userProjs, err := r.repo.UserProjects.ListByUser(userId)
	if err != nil {
		fmt.Printf("DEBUG: Erro ao buscar user_projects: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Encontrados %d user_projects para user %s\n", len(userProjs), userId)

	// Criar slice de resposta enriquecida
	result := make([]UserProjectWithName, 0, len(userProjs))

	// Para cada projeto, buscar o nome
	for _, userProj := range userProjs {
		fmt.Printf("DEBUG: Buscando projeto %s\n", userProj.ProjectId.String())
		proj, err := r.repo.Projects.GetProjectById(userProj.ProjectId)
		if err != nil {
			// Se projeto não for encontrado, pular
			fmt.Printf("DEBUG: Erro ao buscar projeto %s: %v\n", userProj.ProjectId.String(), err)
			continue
		}

		fmt.Printf("DEBUG: Projeto encontrado: %s (org: %s)\n", proj.Name, proj.OrganizationId.String())
		result = append(result, UserProjectWithName{
			Id:             userProj.Id,
			UserId:         userProj.UserId,
			ProjectId:      userProj.ProjectId,
			ProjectName:    proj.Name,
			OrganizationId: proj.OrganizationId, // ✅ NOVO: Incluir organization_id
			Role:           userProj.Role,
			Active:         userProj.Active,
			CreatedAt:      userProj.CreatedAt,
			UpdatedAt:      userProj.UpdatedAt,
			DeletedAt:      userProj.DeletedAt,
		})
	}

	fmt.Printf("DEBUG: Retornando %d projetos\n", len(result))
	return result, nil
}

func NewAuthHandler(repo *repositories.DBconn) IHandlerAuth {
	return &resourceAuth{repo: repo}
}
