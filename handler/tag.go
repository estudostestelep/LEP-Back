package handler

import (
	"lep/repositories"
	"lep/repositories/models"
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
	tag.Id = uuid.New()
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()
	err := r.repo.Tags.CreateTag(tag)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceTag) UpdateTag(updatedTag *models.Tag) error {
	updatedTag.UpdatedAt = time.Now()
	err := r.repo.Tags.UpdateTag(updatedTag)
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
