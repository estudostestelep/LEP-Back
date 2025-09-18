package repositories

import (
	"lep/repositories/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceBannedLists struct {
	db *gorm.DB
}

type IBannedListsRepository interface {
	GetBannedList(id uuid.UUID) (*models.BannedLists, error)
	GetBannedAllList() (*[]models.BannedLists, error)
	CreateBannedList(bannedList *models.BannedLists) error
	UpdateBannedList(id uuid.UUID, updatedBannedList *models.BannedLists) error
	DeleteBannedList(id uuid.UUID) error
}

func NewConnBannedLists(db *gorm.DB) IBannedListsRepository {
	return &resourceBannedLists{db: db}
}

func (r *resourceBannedLists) GetBannedList(id uuid.UUID) (*models.BannedLists, error) {
	var bannedList models.BannedLists
	result := r.db.Where("banned_list_id = ?", id).First(&bannedList)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bannedList, nil
}

func (r *resourceBannedLists) GetBannedAllList() (*[]models.BannedLists, error) {
	var bannedList []models.BannedLists
	result := r.db.Find(&bannedList)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bannedList, nil
}

func (r *resourceBannedLists) CreateBannedList(bannedList *models.BannedLists) error {
	result := r.db.Create(bannedList)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *resourceBannedLists) UpdateBannedList(id uuid.UUID, updatedBannedList *models.BannedLists) error {
	result := r.db.Model(&models.BannedLists{}).Where("banned_list_id = ?", id).Updates(updatedBannedList)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *resourceBannedLists) DeleteBannedList(id uuid.UUID) error {
	result := r.db.Where("banned_list_id = ?", id).Delete(&models.BannedLists{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
