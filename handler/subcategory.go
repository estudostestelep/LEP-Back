package handler

import (
	"errors"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceSubcategory struct {
	repo *repositories.DBconn
}

type IHandlerSubcategory interface {
	GetSubcategory(id string) (*models.Subcategory, error)
	ListSubcategories(orgId, projectId string) ([]models.Subcategory, error)
	GetSubcategoriesByCategory(categoryId string) ([]models.Subcategory, error)
	ListActiveSubcategories(orgId, projectId string) ([]models.Subcategory, error)
	CreateSubcategory(subcategory *models.Subcategory) error
	UpdateSubcategory(updatedSubcategory *models.Subcategory) error
	UpdateSubcategoryOrder(id string, order int) error
	UpdateSubcategoryStatus(id string, active bool) error
	DeleteSubcategory(id string) error
	AddCategoryToSubcategory(subcategoryId, categoryId string) error
	RemoveCategoryFromSubcategory(subcategoryId, categoryId string) error
	GetSubcategoryCategories(subcategoryId string) ([]models.Category, error)
}

func (r *resourceSubcategory) GetSubcategory(id string) (*models.Subcategory, error) {
	subId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.Subcategories.GetSubcategoryById(subId)
}

func (r *resourceSubcategory) ListSubcategories(orgId, projectId string) ([]models.Subcategory, error) {
	orgUuid, _ := uuid.Parse(orgId)
	projectUuid, _ := uuid.Parse(projectId)
	return r.repo.Subcategories.GetSubcategoryList(orgUuid, projectUuid)
}

func (r *resourceSubcategory) GetSubcategoriesByCategory(categoryId string) ([]models.Subcategory, error) {
	catUuid, err := uuid.Parse(categoryId)
	if err != nil {
		return nil, err
	}
	return r.repo.Subcategories.GetSubcategoriesByCategory(catUuid)
}

func (r *resourceSubcategory) ListActiveSubcategories(orgId, projectId string) ([]models.Subcategory, error) {
	orgUuid, _ := uuid.Parse(orgId)
	projectUuid, _ := uuid.Parse(projectId)
	return r.repo.Subcategories.GetActiveSubcategoryList(orgUuid, projectUuid)
}

func (r *resourceSubcategory) CreateSubcategory(subcategory *models.Subcategory) error {
	// Verificar se já existe subcategoria com o mesmo nome no projeto
	exists, err := r.repo.Subcategories.CheckSubcategoryNameExists(subcategory.OrganizationId, subcategory.ProjectId, subcategory.Name, nil)
	if err != nil {
		return fmt.Errorf("erro ao verificar duplicata: %w", err)
	}
	if exists {
		return errors.New("already_exists: subcategory with this name already exists in this project")
	}

	subcategory.Id = uuid.New()
	subcategory.CreatedAt = time.Now()
	subcategory.UpdatedAt = time.Now()
	return r.repo.Subcategories.CreateSubcategory(subcategory)
}

func (r *resourceSubcategory) UpdateSubcategory(updatedSubcategory *models.Subcategory) error {
	// Verificar se já existe outra subcategoria com o mesmo nome no projeto
	exists, err := r.repo.Subcategories.CheckSubcategoryNameExists(updatedSubcategory.OrganizationId, updatedSubcategory.ProjectId, updatedSubcategory.Name, &updatedSubcategory.Id)
	if err != nil {
		return fmt.Errorf("erro ao verificar duplicata: %w", err)
	}
	if exists {
		return errors.New("already_exists: subcategory with this name already exists in this project")
	}

	updatedSubcategory.UpdatedAt = time.Now()
	return r.repo.Subcategories.UpdateSubcategory(updatedSubcategory)
}

func (r *resourceSubcategory) UpdateSubcategoryOrder(id string, order int) error {
	subId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.Subcategories.UpdateSubcategoryOrder(subId, order)
}

func (r *resourceSubcategory) UpdateSubcategoryStatus(id string, active bool) error {
	subId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.Subcategories.UpdateSubcategoryStatus(subId, active)
}

func (r *resourceSubcategory) DeleteSubcategory(id string) error {
	subId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.Subcategories.SoftDeleteSubcategory(subId)
}

func (r *resourceSubcategory) AddCategoryToSubcategory(subcategoryId, categoryId string) error {
	subId, _ := uuid.Parse(subcategoryId)
	catId, _ := uuid.Parse(categoryId)
	return r.repo.Subcategories.AddCategoryToSubcategory(subId, catId)
}

func (r *resourceSubcategory) RemoveCategoryFromSubcategory(subcategoryId, categoryId string) error {
	subId, _ := uuid.Parse(subcategoryId)
	catId, _ := uuid.Parse(categoryId)
	return r.repo.Subcategories.RemoveCategoryFromSubcategory(subId, catId)
}

func (r *resourceSubcategory) GetSubcategoryCategories(subcategoryId string) ([]models.Category, error) {
	subId, err := uuid.Parse(subcategoryId)
	if err != nil {
		return nil, err
	}
	return r.repo.Subcategories.GetSubcategoryCategories(subId)
}

func NewSourceHandlerSubcategory(repo *repositories.DBconn) IHandlerSubcategory {
	return &resourceSubcategory{repo: repo}
}
