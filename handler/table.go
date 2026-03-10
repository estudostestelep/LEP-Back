package handler

import (
	"errors"
	"fmt"
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
	ListTables(orgId, projectId string, environmentId *string) ([]models.Table, error)
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
	// Verificar se já existe mesa com o mesmo número no projeto
	exists, err := r.repo.Tables.CheckTableNumberExists(table.OrganizationId, table.ProjectId, table.Number, nil)
	if err != nil {
		return fmt.Errorf("erro ao verificar duplicata: %w", err)
	}
	if exists {
		return errors.New("already_exists: table with this number already exists in this project")
	}

	table.Id = uuid.New()
	table.CreatedAt = time.Now()
	table.UpdatedAt = time.Now()
	err = r.repo.Tables.CreateTable(table)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceTables) UpdateTable(updatedTable *models.Table) error {
	// Verificar se já existe outra mesa com o mesmo número no projeto
	exists, err := r.repo.Tables.CheckTableNumberExists(updatedTable.OrganizationId, updatedTable.ProjectId, updatedTable.Number, &updatedTable.Id)
	if err != nil {
		return fmt.Errorf("erro ao verificar duplicata: %w", err)
	}
	if exists {
		return errors.New("already_exists: table with this number already exists in this project")
	}

	updatedTable.UpdatedAt = time.Now()
	err = r.repo.Tables.UpdateTable(updatedTable)
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

func (r *resourceTables) ListTables(orgId, projectId string, environmentId *string) ([]models.Table, error) {
	orgUuid, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	var envUuid *uuid.UUID
	if environmentId != nil && *environmentId != "" {
		parsed, parseErr := uuid.Parse(*environmentId)
		if parseErr != nil {
			return nil, parseErr
		}
		envUuid = &parsed
	}
	resp, err := r.repo.Tables.ListTables(orgUuid, projectUuid, envUuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewSourceHandlerTables(repo *repositories.DBconn) IHandlerTables {
	return &resourceTables{repo: repo}
}
