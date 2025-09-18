package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceTables struct {
	repo *repositories.DBconn
}

type IHandlerTables interface {
	GetTable(id string) (*models.Table, error)
	CreateTable(table *models.Table) error
	UpdateTable(updatedTable *models.Table) error
	DeleteTable(id string) error
	ListTables(orgId, projectId string) ([]models.Table, error)
}

func (r *resourceTables) GetTable(id string) (*models.Table, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	resp, err := r.repo.Tables.GetById(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceTables) CreateTable(table *models.Table) error {
	table.Id = uuid.New()
	table.CreatedAt = time.Now()
	table.UpdatedAt = time.Now()
	err := r.repo.Tables.CreateTable(table)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceTables) UpdateTable(updatedTable *models.Table) error {
	updatedTable.UpdatedAt = time.Now()
	err := r.repo.Tables.UpdateTable(updatedTable)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceTables) DeleteTable(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	err = r.repo.Tables.SoftDeleteTable(uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceTables) ListTables(orgId, projectId string) ([]models.Table, error) {
	orgUuid, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	resp, err := r.repo.Tables.ListTables(orgUuid, projectUuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewSourceHandlerTables(repo *repositories.DBconn) IHandlerTables {
	return &resourceTables{repo: repo}
}
