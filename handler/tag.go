package handler

import (
	"errors"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

type resourceTag struct {
	repo *repositories.DBconn
}

type IHandlerTag interface {
	GetTag(id string) (*models.Tag, error)
	CreateTag(tag *models.Tag) error
	UpdateTag(updatedTag *models.Tag) error
	DeleteTag(id string) error
	ListTags(orgId, projectId string) ([]models.Tag, error)
	ListActiveTags(orgId, projectId string) ([]models.Tag, error)
	GetTagsByEntityType(orgId, projectId, entityType string) ([]models.Tag, error)
}

func (r *resourceTag) GetTag(id string) (*models.Tag, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	resp, err := r.repo.Tags.GetTagById(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceTag) CreateTag(tag *models.Tag) error {
	// Validar campos obrigatórios
	if strings.TrimSpace(tag.Name) == "" {
		return errors.New("tag name is required")
	}

	if tag.OrganizationId == uuid.Nil {
		return errors.New("organization_id is required")
	}

	if tag.ProjectId == uuid.Nil {
		return errors.New("project_id is required")
	}

	// Normalizar nome (trim spaces)
	tag.Name = strings.TrimSpace(tag.Name)

	// Validar duplicata: verificar se já existe tag com mesmo nome E tipo no projeto
	existingTag, err := r.repo.Tags.GetTagByNameAndType(tag.OrganizationId, tag.ProjectId, tag.Name, tag.EntityType)
	if err != nil {
		return fmt.Errorf("erro ao verificar duplicata: %w", err)
	}

	if existingTag != nil {
		return fmt.Errorf("tag with name '%s' and type '%s' already exists in this project", tag.Name, tag.EntityType)
	}

	// Gerar ID e timestamps
	tag.Id = uuid.New()
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	// Criar tag
	err = r.repo.Tags.CreateTag(tag)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceTag) UpdateTag(updatedTag *models.Tag) error {
	// Validar campos obrigatórios
	if strings.TrimSpace(updatedTag.Name) == "" {
		return errors.New("tag name is required")
	}

	if updatedTag.OrganizationId == uuid.Nil {
		return errors.New("organization_id is required")
	}

	if updatedTag.ProjectId == uuid.Nil {
		return errors.New("project_id is required")
	}

	// Normalizar nome (trim spaces)
	updatedTag.Name = strings.TrimSpace(updatedTag.Name)

	// Validar duplicata: verificar se já existe outra tag com mesmo nome E tipo no projeto
	// (excluindo a tag atual sendo atualizada)
	existingTag, err := r.repo.Tags.GetTagByNameAndType(updatedTag.OrganizationId, updatedTag.ProjectId, updatedTag.Name, updatedTag.EntityType)
	if err != nil {
		return fmt.Errorf("erro ao verificar duplicata: %w", err)
	}

	// Se encontrou outra tag com mesmo nome+tipo, verificar se é a mesma tag
	if existingTag != nil && existingTag.Id != updatedTag.Id {
		return fmt.Errorf("tag with name '%s' and type '%s' already exists in this project", updatedTag.Name, updatedTag.EntityType)
	}

	updatedTag.UpdatedAt = time.Now()
	err = r.repo.Tags.UpdateTag(updatedTag)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceTag) DeleteTag(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	err = r.repo.Tags.SoftDeleteTag(uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceTag) ListTags(orgId, projectId string) ([]models.Tag, error) {
	orgUuid, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	resp, err := r.repo.Tags.GetTagList(orgUuid, projectUuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceTag) ListActiveTags(orgId, projectId string) ([]models.Tag, error) {
	orgUuid, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	resp, err := r.repo.Tags.GetActiveTagList(orgUuid, projectUuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceTag) GetTagsByEntityType(orgId, projectId, entityType string) ([]models.Tag, error) {
	orgUuid, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	resp, err := r.repo.Tags.GetTagsByEntityType(orgUuid, projectUuid, entityType)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewSourceHandlerTag(repo *repositories.DBconn) IHandlerTag {
	return &resourceTag{repo: repo}
}
