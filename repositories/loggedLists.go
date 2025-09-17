package repositories

import (
	"lep/repositories/models"

	"gorm.io/gorm"
)

type resourceLoggedLists struct {
	db *gorm.DB
}

type ILoggedListsRepository interface {
	GetLoggedToken(token string) (*models.LoggedLists, error)
	GetLoggedAllList() (*[]models.LoggedLists, error)
	CreateLoggedList(loggedList *models.LoggedLists) error
	UpdateLoggedList(id int, updatedLoggedList *models.LoggedLists) error
	DeleteLoggedList(token string) error
}

func NewConnLoggedLists(db *gorm.DB) ILoggedListsRepository {
	return &resourceLoggedLists{db: db}
}

func (r *resourceLoggedLists) GetLoggedToken(token string) (*models.LoggedLists, error) {
	var loggedList models.LoggedLists
	result := r.db.Where("token = ?", token).First(&loggedList)
	if result.Error != nil {
		return nil, result.Error
	}
	return &loggedList, nil
}

func (r *resourceLoggedLists) GetLoggedAllList() (*[]models.LoggedLists, error) {
	var loggedList []models.LoggedLists
	result := r.db.Find(&loggedList)
	if result.Error != nil {
		return nil, result.Error
	}
	return &loggedList, nil
}

func (r *resourceLoggedLists) CreateLoggedList(loggedList *models.LoggedLists) error {
	result := r.db.Create(loggedList)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *resourceLoggedLists) UpdateLoggedList(id int, updatedLoggedList *models.LoggedLists) error {
	result := r.db.Model(&models.LoggedLists{}).Where("logged_list_id = ?", id).Updates(updatedLoggedList)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *resourceLoggedLists) DeleteLoggedList(token string) error {
	result := r.db.Where("token = ?", token).Delete(&models.LoggedLists{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
