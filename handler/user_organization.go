package handler

import (
	"errors"
	"lep/repositories"
	"lep/repositories/models"
	"lep/resource/validation"

	"github.com/google/uuid"
)

type resourceUserOrganization struct {
	repo *repositories.DBconn
}

type IHandlerUserOrganization interface {
	AddUserToOrganization(userOrg *models.UserOrganization) error
	RemoveUserFromOrganization(userId, orgId string) error
	UpdateUserOrganization(userOrg *models.UserOrganization) error
	GetUserOrganizations(userId string) ([]models.UserOrganization, error)
	GetOrganizationUsers(orgId string) ([]models.UserOrganization, error)
	UserHasAccessToOrganization(userId, orgId string) (bool, error)
}

func (r *resourceUserOrganization) AddUserToOrganization(userOrg *models.UserOrganization) error {
	// Validar dados
	if err := validation.CreateUserOrganizationValidation(userOrg); err != nil {
		return err
	}

	// Verificar se já existe
	existing, _ := r.repo.UserOrganizations.GetByUserAndOrganization(userOrg.UserId.String(), userOrg.OrganizationId.String())
	if existing != nil {
		return errors.New("usuário já vinculado a esta organização")
	}

	// Gerar ID se necessário
	if userOrg.Id == uuid.Nil {
		userOrg.Id = uuid.New()
	}

	// Criar relacionamento - foreign keys do banco garantem integridade referencial
	return r.repo.UserOrganizations.Create(userOrg)
}

func (r *resourceUserOrganization) RemoveUserFromOrganization(userId, orgId string) error {
	userOrg, err := r.repo.UserOrganizations.GetByUserAndOrganization(userId, orgId)
	if err != nil {
		return errors.New("relacionamento não encontrado")
	}

	return r.repo.UserOrganizations.Delete(userOrg.Id.String())
}

func (r *resourceUserOrganization) UpdateUserOrganization(userOrg *models.UserOrganization) error {
	// Validar dados
	if err := validation.UpdateUserOrganizationValidation(userOrg); err != nil {
		return err
	}

	existing, err := r.repo.UserOrganizations.GetById(userOrg.Id.String())
	if err != nil {
		return errors.New("relacionamento não encontrado")
	}

	// Manter campos imutáveis
	userOrg.UserId = existing.UserId
	userOrg.OrganizationId = existing.OrganizationId

	return r.repo.UserOrganizations.Update(userOrg)
}

func (r *resourceUserOrganization) GetUserOrganizations(userId string) ([]models.UserOrganization, error) {
	return r.repo.UserOrganizations.ListByUser(userId)
}

func (r *resourceUserOrganization) GetOrganizationUsers(orgId string) ([]models.UserOrganization, error) {
	return r.repo.UserOrganizations.ListByOrganization(orgId)
}

func (r *resourceUserOrganization) UserHasAccessToOrganization(userId, orgId string) (bool, error) {
	return r.repo.UserOrganizations.UserBelongsToOrganization(userId, orgId)
}

func NewSourceHandlerUserOrganization(repo *repositories.DBconn) IHandlerUserOrganization {
	return &resourceUserOrganization{repo: repo}
}
