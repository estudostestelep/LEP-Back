package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SidebarConfigRepository struct {
	db *gorm.DB
}

type ISidebarConfigRepository interface {
	GetGlobal() (*models.SidebarConfig, error)
	Create(config *models.SidebarConfig) error
	Update(config *models.SidebarConfig) error
	GetOrCreate(defaultItemsJSON string) (*models.SidebarConfig, error)
}

func NewSidebarConfigRepository(db *gorm.DB) ISidebarConfigRepository {
	return &SidebarConfigRepository{db: db}
}

// GetGlobal busca a configuração global da sidebar
func (r *SidebarConfigRepository) GetGlobal() (*models.SidebarConfig, error) {
	var config models.SidebarConfig
	err := r.db.First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Create cria uma nova configuração de sidebar
func (r *SidebarConfigRepository) Create(config *models.SidebarConfig) error {
	return r.db.Create(config).Error
}

// Update atualiza uma configuração de sidebar existente
func (r *SidebarConfigRepository) Update(config *models.SidebarConfig) error {
	config.UpdatedAt = time.Now()
	return r.db.Save(config).Error
}

// GetOrCreate busca ou cria configuração global padrão
func (r *SidebarConfigRepository) GetOrCreate(defaultItemsJSON string) (*models.SidebarConfig, error) {
	// Tenta buscar configuração existente
	config, err := r.GetGlobal()
	if err == nil {
		return config, nil
	}

	// Se não encontrou, cria configuração padrão
	if err == gorm.ErrRecordNotFound {
		newConfig := &models.SidebarConfig{
			Id:          uuid.New(),
			ItemConfigs: defaultItemsJSON,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err = r.Create(newConfig)
		if err != nil {
			return nil, err
		}
		return newConfig, nil
	}

	return nil, err
}
