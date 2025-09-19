package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Interface para TableRepository
type ITableRepository interface {
	CreateTable(table *models.Table) error
	GetById(id uuid.UUID) (*models.Table, error)
	GetTableById(id uuid.UUID) (*models.Table, error)
	ListTables(OrganizationId, projectId uuid.UUID) ([]models.Table, error)
	ListTablesByProject(OrganizationId, projectId uuid.UUID) ([]models.Table, error)
	GetTablesByProject(orgId, projectId uuid.UUID) ([]models.Table, error)
	UpdateTable(table *models.Table) error
	SoftDeleteTable(id uuid.UUID) error
}

type TableRepository struct {
	db *gorm.DB
}

func NewConnTable(db *gorm.DB) ITableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) CreateTable(table *models.Table) error {
	return r.db.Create(table).Error
}

func (r *TableRepository) GetById(id uuid.UUID) (*models.Table, error) {
	var table models.Table
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&table).Error
	if err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *TableRepository) GetTableById(id uuid.UUID) (*models.Table, error) {
	return r.GetById(id)
}

func (r *TableRepository) ListTables(OrganizationId, projectId uuid.UUID) ([]models.Table, error) {
	var tables []models.Table
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", OrganizationId, projectId).Find(&tables).Error
	return tables, err
}

func (r *TableRepository) ListTablesByProject(OrganizationId, projectId uuid.UUID) ([]models.Table, error) {
	return r.ListTables(OrganizationId, projectId)
}

func (r *TableRepository) UpdateTable(table *models.Table) error {
	// Garante que deleted_at não seja alterado
	return r.db.Model(table).Omit("deleted_at").Save(table).Error
}

func (r *TableRepository) SoftDeleteTable(id uuid.UUID) error {
	return r.db.Model(&models.Table{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now()).Error
}

func (r *TableRepository) GetTablesByProject(orgId, projectId uuid.UUID) ([]models.Table, error) {
	var tables []models.Table
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", orgId, projectId).
		Order("number ASC").Find(&tables).Error
	return tables, err
}
