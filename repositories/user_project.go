package repositories

import (
	"lep/repositories/models"

	"gorm.io/gorm"
)

type resourceUserProject struct {
	db *gorm.DB
}

type IUserProjectRepository interface {
	// CRUD básico
	Create(userProj *models.UserProject) error
	GetById(id string) (*models.UserProject, error)
	Update(userProj *models.UserProject) error
	Delete(id string) error

	// Consultas específicas
	GetByUserAndProject(userId, projectId string) (*models.UserProject, error)
	ListByUser(userId string) ([]models.UserProject, error)
	ListByProject(projectId string) ([]models.UserProject, error)

	// Consultas avançadas
	ListByUserAndOrganization(userId, orgId string) ([]models.UserProject, error)

	// Verificações
	UserBelongsToProject(userId, projectId string) (bool, error)
}

func NewUserProjectRepository(db *gorm.DB) IUserProjectRepository {
	return &resourceUserProject{db: db}
}

func (r *resourceUserProject) Create(userProj *models.UserProject) error {
	return r.db.Create(userProj).Error
}

func (r *resourceUserProject) GetById(id string) (*models.UserProject, error) {
	var userProj models.UserProject
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&userProj).Error
	if err != nil {
		return nil, err
	}
	return &userProj, nil
}

func (r *resourceUserProject) Update(userProj *models.UserProject) error {
	return r.db.Save(userProj).Error
}

func (r *resourceUserProject) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.UserProject{}).Error
}

func (r *resourceUserProject) GetByUserAndProject(userId, projectId string) (*models.UserProject, error) {
	var userProj models.UserProject
	err := r.db.Where("user_id = ? AND project_id = ? AND deleted_at IS NULL", userId, projectId).First(&userProj).Error
	if err != nil {
		return nil, err
	}
	return &userProj, nil
}

func (r *resourceUserProject) ListByUser(userId string) ([]models.UserProject, error) {
	var userProjs []models.UserProject
	err := r.db.Where("user_id = ? AND deleted_at IS NULL AND active = true", userId).Find(&userProjs).Error
	return userProjs, err
}

func (r *resourceUserProject) ListByProject(projectId string) ([]models.UserProject, error) {
	var userProjs []models.UserProject
	err := r.db.Where("project_id = ? AND deleted_at IS NULL AND active = true", projectId).Find(&userProjs).Error
	return userProjs, err
}

func (r *resourceUserProject) ListByUserAndOrganization(userId, orgId string) ([]models.UserProject, error) {
	var userProjs []models.UserProject
	// Join com a tabela de projetos para filtrar por organização
	err := r.db.Table("user_projects").
		Joins("JOIN projects ON projects.id = user_projects.project_id").
		Where("user_projects.user_id = ? AND projects.organization_id = ? AND user_projects.deleted_at IS NULL AND user_projects.active = true", userId, orgId).
		Find(&userProjs).Error
	return userProjs, err
}

func (r *resourceUserProject) UserBelongsToProject(userId, projectId string) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserProject{}).
		Where("user_id = ? AND project_id = ? AND deleted_at IS NULL AND active = true", userId, projectId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
