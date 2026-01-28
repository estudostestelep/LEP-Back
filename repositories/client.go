package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type resourceClient struct {
	db *gorm.DB
}

type IClientRepository interface {
	GetClientById(id string) (*models.Client, error)
	GetClientByEmail(email string) (*models.Client, error)
	GetClientByEmailAndOrg(email string, orgId string) (*models.Client, error)
	ListClients() ([]models.Client, error)
	ListClientsByOrganization(orgId string) ([]models.Client, error)
	ListClientsByProject(orgId string, projectId string) ([]models.Client, error)
	CreateClient(client *models.Client) error
	UpdateClient(client *models.Client) error
	UpdateLastAccess(clientId string) error
	SoftDeleteClient(id string) error
	DeleteClient(id string) error
	ClientEmailExistsInOrg(email string, orgId string) (bool, error)
	AddProjectToClient(clientId string, projectId string) error
	RemoveProjectFromClient(clientId string, projectId string) error
	GetClientWithOrganization(id string) (*models.ClientWithOrganization, error)
}

func NewClientRepository(db *gorm.DB) IClientRepository {
	return &resourceClient{db: db}
}

func (r *resourceClient) GetClientById(id string) (*models.Client, error) {
	var client models.Client
	err := r.db.First(&client, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *resourceClient) GetClientByEmail(email string) (*models.Client, error) {
	var client models.Client
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *resourceClient) GetClientByEmailAndOrg(email string, orgId string) (*models.Client, error) {
	var client models.Client
	err := r.db.Where("email = ? AND org_id = ? AND deleted_at IS NULL", email, orgId).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *resourceClient) ListClients() ([]models.Client, error) {
	var clients []models.Client
	err := r.db.Where("deleted_at IS NULL").Order("created_at DESC").Find(&clients).Error
	return clients, err
}

func (r *resourceClient) ListClientsByOrganization(orgId string) ([]models.Client, error) {
	var clients []models.Client
	err := r.db.Where("org_id = ? AND deleted_at IS NULL", orgId).Order("created_at DESC").Find(&clients).Error
	return clients, err
}

func (r *resourceClient) ListClientsByProject(orgId string, projectId string) ([]models.Client, error) {
	var clients []models.Client
	err := r.db.Where("org_id = ? AND ? = ANY(proj_ids) AND deleted_at IS NULL", orgId, projectId).
		Order("created_at DESC").Find(&clients).Error
	return clients, err
}

func (r *resourceClient) CreateClient(client *models.Client) error {
	return r.db.Create(client).Error
}

func (r *resourceClient) UpdateClient(client *models.Client) error {
	// Se o password estiver vazio, ignora o campo para não sobrescrever
	if client.Password == "" {
		return r.db.Omit("Password").Save(client).Error
	}
	return r.db.Save(client).Error
}

func (r *resourceClient) UpdateLastAccess(clientId string) error {
	return r.db.Model(&models.Client{}).Where("id = ?", clientId).Update("last_access_at", time.Now()).Error
}

func (r *resourceClient) SoftDeleteClient(id string) error {
	return r.db.Model(&models.Client{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *resourceClient) DeleteClient(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Client{}).Error
}

func (r *resourceClient) ClientEmailExistsInOrg(email string, orgId string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Client{}).Where("email = ? AND org_id = ? AND deleted_at IS NULL", email, orgId).Count(&count).Error
	return count > 0, err
}

func (r *resourceClient) AddProjectToClient(clientId string, projectId string) error {
	return r.db.Model(&models.Client{}).
		Where("id = ?", clientId).
		Update("proj_ids", gorm.Expr("array_append(proj_ids, ?)", projectId)).Error
}

func (r *resourceClient) RemoveProjectFromClient(clientId string, projectId string) error {
	return r.db.Model(&models.Client{}).
		Where("id = ?", clientId).
		Update("proj_ids", gorm.Expr("array_remove(proj_ids, ?)", projectId)).Error
}

func (r *resourceClient) GetClientWithOrganization(id string) (*models.ClientWithOrganization, error) {
	var result models.ClientWithOrganization

	err := r.db.Table("clients").
		Select("clients.*, organizations.name as organization_name, organizations.slug as organization_slug").
		Joins("LEFT JOIN organizations ON clients.org_id = organizations.id").
		Where("clients.id = ? AND clients.deleted_at IS NULL", id).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetClientProjectsInfo busca informações dos projetos do cliente
func (r *resourceClient) GetClientProjectsInfo(client *models.Client) ([]models.ClientProjectInfo, error) {
	if len(client.ProjIds) == 0 {
		return []models.ClientProjectInfo{}, nil
	}

	// Converter pq.StringArray para slice de strings
	projectIds := []string(client.ProjIds)

	var projects []struct {
		Id   string `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}

	err := r.db.Table("projects").
		Select("id, name").
		Where("id IN ? AND deleted_at IS NULL", pq.Array(projectIds)).
		Find(&projects).Error

	if err != nil {
		return nil, err
	}

	result := make([]models.ClientProjectInfo, len(projects))
	for i, p := range projects {
		result[i] = models.ClientProjectInfo{
			ProjectName: p.Name,
		}
	}

	return result, nil
}
