package handler

import (
	"fmt"
	"lep/constants"
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type resourceUser struct {
	repo *repositories.DBconn
}

type IHandlerUser interface {
	GetUser(id string) (*models.User, error)
	GetUserByGroup(id string) ([]models.User, error)
	ListUsers(orgId, projectId string) ([]models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(updatedUser *models.User) error
	DeleteUser(id string) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserWithRelations(id string) (*models.UserWithRelations, error)
}

func (r *resourceUser) GetUser(id string) (*models.User, error) {
	resp, err := r.repo.User.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceUser) GetUserByGroup(id string) ([]models.User, error) {
	resp, err := r.repo.User.GetUsersByGroup(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceUser) ListUsers(orgId, projectId string) ([]models.User, error) {
	resp, err := r.repo.User.ListUsersByOrganizationAndProject(orgId, projectId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceUser) CreateUser(user *models.User) error {
	existingUser, _ := r.repo.User.GetUserByEmail(user.Email)

	if existingUser != nil {
		return errors.New("E-mail já cadastrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	// Gerar ID apenas se não foi fornecido
	if user.Id == uuid.Nil {
		user.Id = uuid.New()
	}

	err = r.repo.User.CreateUser(user)
	if err != nil {
		return err
	}

	// 🔑 REGRA DE NEGÓCIO: Se o novo usuário é um master admin, adicioná-lo a todas as organizações
	isMasterAdmin := constants.HasPermission(user.Permissions, constants.PermissionMasterAdmin)
	if isMasterAdmin {
		if err := r.addMasterAdminToAllOrganizations(user.Id); err != nil {
			// Log error but don't fail user creation
			fmt.Printf("Aviso: erro ao adicionar master admin a organizações: %v\n", err)
		}
	}

	return nil
}

func (r *resourceUser) UpdateUser(updatedUser *models.User) error {
	existingUser, err := r.repo.User.GetUserByEmail(updatedUser.Email)

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updatedUser.Password = string(hashedPassword)
	}

	if existingUser != nil && existingUser.Id != updatedUser.Id {
		return fmt.Errorf("E-mail já cadastrado")
	}

	err = r.repo.User.UpdateUser(updatedUser)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceUser) DeleteUser(id string) error {
	err := r.repo.User.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceUser) GetUserByEmail(email string) (*models.User, error) {
	resp, err := r.repo.User.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceUser) GetUserWithRelations(id string) (*models.UserWithRelations, error) {
	resp, err := r.repo.User.GetUserWithRelations(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// addMasterAdminToAllOrganizations adiciona um novo master admin a todas as organizações existentes
// REGRA DE NEGÓCIO: Master admins devem ter acesso automático a todas as orgs
func (r *resourceUser) addMasterAdminToAllOrganizations(userId uuid.UUID) error {
	// Buscar todas as organizações ativas
	orgs, err := r.repo.Organizations.ListActiveOrganizations()
	if err != nil {
		return fmt.Errorf("erro ao buscar organizações: %v", err)
	}

	now := time.Now()

	// Adicionar master admin a cada organização e seus projetos
	for _, org := range orgs {
		// Criar relacionamento usuário-organização (se não existir)
		userOrg := &models.UserOrganization{
			Id:             uuid.New(),
			UserId:         userId,
			OrganizationId: org.Id,
			Role:           "admin", // Master admins são admins da organização
			Active:         true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		// Ignorar erro se já existe (idempotente)
		_ = r.repo.UserOrganizations.Create(userOrg)

		// Buscar todos os projetos da organização
		projects, err := r.repo.Projects.GetProjectByOrganization(org.Id)
		if err != nil {
			// Log error but continue
			fmt.Printf("Aviso: erro ao buscar projetos da org %s: %v\n", org.Id, err)
			continue
		}

		// Adicionar master admin a cada projeto
		for _, proj := range projects {
			userProj := &models.UserProject{
				Id:        uuid.New(),
				UserId:    userId,
				ProjectId: proj.Id,
				Role:      "admin",
				Active:    true,
				CreatedAt: now,
				UpdatedAt: now,
			}

			// Ignorar erro se já existe (idempotente)
			_ = r.repo.UserProjects.Create(userProj)
		}
	}

	return nil
}

func NewSourceHandlerUser(repo *repositories.DBconn) IHandlerUser {
	return &resourceUser{repo: repo}
}
