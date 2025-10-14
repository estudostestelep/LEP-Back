package handler

import (
	"errors"
	"lep/repositories"
	"lep/repositories/models"
	"lep/resource/validation"

	"github.com/google/uuid"
)

type resourceUserProject struct {
	repo *repositories.DBconn
}

type IHandlerUserProject interface {
	AddUserToProject(userProj *models.UserProject) error
	RemoveUserFromProject(userId, projectId string) error
	UpdateUserProject(userProj *models.UserProject) error
	GetUserProjects(userId string) ([]models.UserProject, error)
	GetUserProjectsByOrganization(userId, orgId string) ([]models.UserProject, error)
	GetProjectUsers(projectId string) ([]models.UserProject, error)
	UserHasAccessToProject(userId, projectId string) (bool, error)
}

func (r *resourceUserProject) AddUserToProject(userProj *models.UserProject) error {
	// Validar dados
	if err := validation.CreateUserProjectValidation(userProj); err != nil {
		return err
	}

	// Verificar se já existe
	existing, _ := r.repo.UserProjects.GetByUserAndProject(userProj.UserId.String(), userProj.ProjectId.String())
	if existing != nil {
		return errors.New("usuário já vinculado a este projeto")
	}

	// Gerar ID se necessário
	if userProj.Id == uuid.Nil {
		userProj.Id = uuid.New()
	}

	// Criar relacionamento - foreign keys do banco garantem integridade referencial
	return r.repo.UserProjects.Create(userProj)
}

func (r *resourceUserProject) RemoveUserFromProject(userId, projectId string) error {
	userProj, err := r.repo.UserProjects.GetByUserAndProject(userId, projectId)
	if err != nil {
		return errors.New("relacionamento não encontrado")
	}

	return r.repo.UserProjects.Delete(userProj.Id.String())
}

func (r *resourceUserProject) UpdateUserProject(userProj *models.UserProject) error {
	// Validar dados
	if err := validation.UpdateUserProjectValidation(userProj); err != nil {
		return err
	}

	existing, err := r.repo.UserProjects.GetById(userProj.Id.String())
	if err != nil {
		return errors.New("relacionamento não encontrado")
	}

	// Manter campos imutáveis
	userProj.UserId = existing.UserId
	userProj.ProjectId = existing.ProjectId

	return r.repo.UserProjects.Update(userProj)
}

func (r *resourceUserProject) GetUserProjects(userId string) ([]models.UserProject, error) {
	return r.repo.UserProjects.ListByUser(userId)
}

func (r *resourceUserProject) GetUserProjectsByOrganization(userId, orgId string) ([]models.UserProject, error) {
	return r.repo.UserProjects.ListByUserAndOrganization(userId, orgId)
}

func (r *resourceUserProject) GetProjectUsers(projectId string) ([]models.UserProject, error) {
	return r.repo.UserProjects.ListByProject(projectId)
}

func (r *resourceUserProject) UserHasAccessToProject(userId, projectId string) (bool, error) {
	return r.repo.UserProjects.UserBelongsToProject(userId, projectId)
}

func NewSourceHandlerUserProject(repo *repositories.DBconn) IHandlerUserProject {
	return &resourceUserProject{repo: repo}
}
