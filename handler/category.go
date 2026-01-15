package handler

import (
	"errors"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceCategory struct {
	repo *repositories.DBconn
}

type IHandlerCategory interface {
	GetCategory(id string) (*models.Category, error)
	ListCategories(orgId, projectId string) ([]models.Category, error)
	GetCategoriesByMenu(menuId string) ([]models.Category, error)
	ListActiveCategories(orgId, projectId string) ([]models.Category, error)
	CreateCategory(category *models.Category) error
	UpdateCategory(updatedCategory *models.Category) error
	UpdateCategoryOrder(id string, order int) error
	UpdateCategoryStatus(id string, active bool) error
	DeleteCategory(id string) error
}

func (r *resourceCategory) GetCategory(id string) (*models.Category, error) {
	catId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.Categories.GetCategoryById(catId)
}

func (r *resourceCategory) ListCategories(orgId, projectId string) ([]models.Category, error) {
	orgUuid, _ := uuid.Parse(orgId)
	projectUuid, _ := uuid.Parse(projectId)
	return r.repo.Categories.GetCategoryList(orgUuid, projectUuid)
}

func (r *resourceCategory) GetCategoriesByMenu(menuId string) ([]models.Category, error) {
	menuUuid, err := uuid.Parse(menuId)
	if err != nil {
		return nil, err
	}
	return r.repo.Categories.GetCategoriesByMenu(menuUuid)
}

func (r *resourceCategory) ListActiveCategories(orgId, projectId string) ([]models.Category, error) {
	orgUuid, _ := uuid.Parse(orgId)
	projectUuid, _ := uuid.Parse(projectId)
	return r.repo.Categories.GetActiveCategoryList(orgUuid, projectUuid)
}

func (r *resourceCategory) CreateCategory(category *models.Category) error {
	// Verificar se já existe categoria com o mesmo nome no projeto
	exists, err := r.repo.Categories.CheckCategoryNameExists(category.OrganizationId, category.ProjectId, category.Name, nil)
	if err != nil {
		return fmt.Errorf("erro ao verificar duplicata: %w", err)
	}
	if exists {
		return errors.New("already_exists: category with this name already exists in this project")
	}

	category.Id = uuid.New()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	return r.repo.Categories.CreateCategory(category)
}

func (r *resourceCategory) UpdateCategory(updatedCategory *models.Category) error {
	// Verificar se já existe outra categoria com o mesmo nome no projeto (excluindo a atual)
	exists, err := r.repo.Categories.CheckCategoryNameExists(updatedCategory.OrganizationId, updatedCategory.ProjectId, updatedCategory.Name, &updatedCategory.Id)
	if err != nil {
		return fmt.Errorf("erro ao verificar duplicata: %w", err)
	}
	if exists {
		return errors.New("already_exists: category with this name already exists in this project")
	}

	updatedCategory.UpdatedAt = time.Now()
	return r.repo.Categories.UpdateCategory(updatedCategory)
}

func (r *resourceCategory) UpdateCategoryOrder(id string, order int) error {
	catId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.Categories.UpdateCategoryOrder(catId, order)
}

func (r *resourceCategory) UpdateCategoryStatus(id string, active bool) error {
	catId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.Categories.UpdateCategoryStatus(catId, active)
}

func (r *resourceCategory) DeleteCategory(id string) error {
	catId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.Categories.SoftDeleteCategory(catId)
}

func NewSourceHandlerCategory(repo *repositories.DBconn) IHandlerCategory {
	return &resourceCategory{repo: repo}
}
