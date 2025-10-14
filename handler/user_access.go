package handler

import (
	"errors"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAccessHandler interface {
	GetUserOrganizationsAndProjects(userId uuid.UUID) (*UserAccessResponse, error)
	UpdateUserOrganizationsAndProjects(userId uuid.UUID, req *UpdateUserAccessRequest) (*UpdateUserAccessResponse, error)
}

type userAccessHandler struct {
	db *gorm.DB
}

func NewUserAccessHandler(db interface{}) UserAccessHandler {
	return &userAccessHandler{db: db.(*gorm.DB)}
}

// Response structures
type UserAccessResponse struct {
	Organizations []models.UserOrganization `json:"organizations"`
	Projects      []models.UserProject      `json:"projects"`
}

type UpdateUserAccessRequest struct {
	OrganizationIds []uuid.UUID `json:"organization_ids" binding:"required"`
	ProjectIds      []uuid.UUID `json:"project_ids" binding:"required"`
}

type UpdateUserAccessResponse struct {
	Message              string `json:"message"`
	OrganizationsAdded   int    `json:"organizations_added"`
	OrganizationsRemoved int    `json:"organizations_removed"`
	ProjectsAdded        int    `json:"projects_added"`
	ProjectsRemoved      int    `json:"projects_removed"`
}

// GetUserOrganizationsAndProjects retorna todas as organizações e projetos vinculados a um usuário
func (h *userAccessHandler) GetUserOrganizationsAndProjects(userId uuid.UUID) (*UserAccessResponse, error) {
	// Verificar se usuário existe
	var user models.User
	if err := h.db.Where("id = ? AND deleted_at IS NULL", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	// Buscar organizações
	var organizations []models.UserOrganization
	if err := h.db.Where("user_id = ? AND deleted_at IS NULL", userId).Find(&organizations).Error; err != nil {
		return nil, err
	}

	// Buscar projetos
	var projects []models.UserProject
	if err := h.db.Where("user_id = ? AND deleted_at IS NULL", userId).Find(&projects).Error; err != nil {
		return nil, err
	}

	return &UserAccessResponse{
		Organizations: organizations,
		Projects:      projects,
	}, nil
}

// UpdateUserOrganizationsAndProjects atualiza os vínculos de organizações e projetos de um usuário
func (h *userAccessHandler) UpdateUserOrganizationsAndProjects(userId uuid.UUID, req *UpdateUserAccessRequest) (*UpdateUserAccessResponse, error) {
	// Verificar se usuário existe
	var user models.User
	if err := h.db.Where("id = ? AND deleted_at IS NULL", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	// Validar que todas as organizações existem
	if len(req.OrganizationIds) > 0 {
		var orgCount int64
		h.db.Model(&models.Organization{}).Where("id IN ? AND deleted_at IS NULL", req.OrganizationIds).Count(&orgCount)
		if int(orgCount) != len(req.OrganizationIds) {
			return nil, errors.New("uma ou mais organizações não existem")
		}
	}

	// Validar que todos os projetos existem
	if len(req.ProjectIds) > 0 {
		var projCount int64
		h.db.Model(&models.Project{}).Where("id IN ? AND deleted_at IS NULL", req.ProjectIds).Count(&projCount)
		if int(projCount) != len(req.ProjectIds) {
			return nil, errors.New("um ou mais projetos não existem")
		}
	}

	// ===== VALIDAÇÃO DE CONSISTÊNCIA: Projetos devem pertencer às Organizações =====
	if len(req.ProjectIds) > 0 && len(req.OrganizationIds) > 0 {
		// Buscar todos os projetos selecionados com suas organizações
		var projects []models.Project
		if err := h.db.Where("id IN ? AND deleted_at IS NULL", req.ProjectIds).Find(&projects).Error; err != nil {
			return nil, err
		}

		// Criar map das organizações selecionadas
		selectedOrgs := make(map[uuid.UUID]bool)
		for _, orgId := range req.OrganizationIds {
			selectedOrgs[orgId] = true
		}

		// Validar que cada projeto pertence a uma das organizações selecionadas
		for _, project := range projects {
			if !selectedOrgs[project.OrganizationId] {
				return nil, errors.New("um ou mais projetos não pertencem às organizações selecionadas")
			}
		}
	}

	response := &UpdateUserAccessResponse{
		Message: "Acessos atualizados com sucesso",
	}

	// Iniciar transação
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ===== PROCESSAR ORGANIZAÇÕES =====
	// Buscar vínculos atuais de organizações
	var currentOrgs []models.UserOrganization
	if err := tx.Where("user_id = ? AND deleted_at IS NULL", userId).Find(&currentOrgs).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Criar map de IDs atuais
	currentOrgIds := make(map[uuid.UUID]bool)
	for _, org := range currentOrgs {
		currentOrgIds[org.OrganizationId] = true
	}

	// Criar map de novos IDs
	newOrgIds := make(map[uuid.UUID]bool)
	for _, id := range req.OrganizationIds {
		newOrgIds[id] = true
	}

	// Remover vínculos que não estão na nova lista
	for _, org := range currentOrgs {
		if !newOrgIds[org.OrganizationId] {
			now := time.Now()
			if err := tx.Model(&models.UserOrganization{}).
				Where("id = ?", org.Id).
				Update("deleted_at", now).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			response.OrganizationsRemoved++
		}
	}

	// Adicionar novos vínculos
	now := time.Now()
	for _, orgId := range req.OrganizationIds {
		if !currentOrgIds[orgId] {
			newUserOrg := models.UserOrganization{
				Id:             uuid.New(),
				UserId:         userId,
				OrganizationId: orgId,
				Role:           "member", // Role padrão
				Active:         true,
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			if err := tx.Create(&newUserOrg).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			response.OrganizationsAdded++
		}
	}

	// ===== PROCESSAR PROJETOS =====
	// Buscar vínculos atuais de projetos
	var currentProjs []models.UserProject
	if err := tx.Where("user_id = ? AND deleted_at IS NULL", userId).Find(&currentProjs).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Criar map de IDs atuais
	currentProjIds := make(map[uuid.UUID]bool)
	for _, proj := range currentProjs {
		currentProjIds[proj.ProjectId] = true
	}

	// Criar map de novos IDs
	newProjIds := make(map[uuid.UUID]bool)
	for _, id := range req.ProjectIds {
		newProjIds[id] = true
	}

	// Remover vínculos que não estão na nova lista
	for _, proj := range currentProjs {
		if !newProjIds[proj.ProjectId] {
			now := time.Now()
			if err := tx.Model(&models.UserProject{}).
				Where("id = ?", proj.Id).
				Update("deleted_at", now).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			response.ProjectsRemoved++
		}
	}

	// Adicionar novos vínculos
	for _, projId := range req.ProjectIds {
		if !currentProjIds[projId] {
			newUserProj := models.UserProject{
				Id:        uuid.New(),
				UserId:    userId,
				ProjectId: projId,
				Role:      "member", // Role padrão
				Active:    true,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := tx.Create(&newUserProj).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			response.ProjectsAdded++
		}
	}

	// Commit da transação
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return response, nil
}
