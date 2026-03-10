package handler

import (
	"encoding/json"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// AdminAuditLogHandler - Handler para logs de auditoria administrativa
type AdminAuditLogHandler struct {
	repo    repositories.IAdminAuditLogRepository
	orgRepo repositories.IOrganizationRepository
	projRepo repositories.IProjectRepository
}

// IAdminAuditLogHandler - Interface do handler
type IAdminAuditLogHandler interface {
	// Métodos públicos (API)
	ListLogs(filters models.AdminAuditLogFilters) (*models.AdminAuditLogPaginatedResponse, error)
	GetLogById(id uuid.UUID) (*models.AdminAuditLog, error)
	DeleteOlderThan(days int) (int64, error)

	// Métodos internos (interceptação de operações de usuário)
	LogUserCreate(actor *models.User, target *models.User, orgId, projectId *uuid.UUID, isAdminZone bool, ipAddress, userAgent string) error
	LogUserUpdate(actor *models.User, oldUser, newUser *models.User, orgId, projectId *uuid.UUID, isAdminZone bool, ipAddress, userAgent string) error
	LogUserDelete(actor *models.User, target *models.User, orgId, projectId *uuid.UUID, isAdminZone bool, ipAddress, userAgent string) error
	LogPasswordReset(actor *models.User, target *models.User, orgId, projectId *uuid.UUID, isAdminZone bool, ipAddress, userAgent string) error

	// Métodos genéricos (para qualquer tipo de entidade)
	LogGenericAction(params AuditLogParams) error
	LogRoleAssignment(actorId uuid.UUID, actorEmail string, targetUserId uuid.UUID, targetUserEmail string, roleId uuid.UUID, roleName string, orgId, projectId *uuid.UUID, ipAddress, userAgent string) error
	LogRoleRemoval(actorId uuid.UUID, actorEmail string, targetUserId uuid.UUID, targetUserEmail string, roleId uuid.UUID, roleName string, orgId, projectId *uuid.UUID, ipAddress, userAgent string) error
}

// NewAdminAuditLogHandler - Construtor do handler
func NewAdminAuditLogHandler(
	repo repositories.IAdminAuditLogRepository,
	orgRepo repositories.IOrganizationRepository,
	projRepo repositories.IProjectRepository,
) IAdminAuditLogHandler {
	return &AdminAuditLogHandler{
		repo:    repo,
		orgRepo: orgRepo,
		projRepo: projRepo,
	}
}

// ListLogs - Lista logs com filtros e paginação
func (h *AdminAuditLogHandler) ListLogs(filters models.AdminAuditLogFilters) (*models.AdminAuditLogPaginatedResponse, error) {
	return h.repo.ListWithFilters(filters)
}

// GetLogById - Busca um log específico pelo ID
func (h *AdminAuditLogHandler) GetLogById(id uuid.UUID) (*models.AdminAuditLog, error) {
	return h.repo.GetById(id)
}

// DeleteOlderThan - Deleta logs mais antigos que X dias
func (h *AdminAuditLogHandler) DeleteOlderThan(days int) (int64, error) {
	if days < 1 {
		return 0, fmt.Errorf("days must be at least 1")
	}
	return h.repo.DeleteOlderThan(days)
}

// LogUserCreate - Registra log de criação de usuário
func (h *AdminAuditLogHandler) LogUserCreate(actor *models.User, target *models.User, orgId, projectId *uuid.UUID, isAdminZone bool, ipAddress, userAgent string) error {
	// Preparar dados do novo usuário (sem senha)
	newValues := h.sanitizeUserForLog(target)
	newValuesJSON, _ := json.Marshal(newValues)

	// Obter nomes de org/projeto se disponíveis
	orgName, projName := h.getOrgProjectNames(orgId, projectId)

	log := &models.AdminAuditLog{
		Id:               uuid.New(),
		ActorId:          actor.Id,
		ActorEmail:       actor.Email,
		TargetId:         target.Id,
		TargetEmail:      target.Email,
		Action:           models.AdminAuditActionCreate,
		EntityType:       models.AdminAuditEntityUser,
		OrganizationId:   orgId,
		OrganizationName: orgName,
		ProjectId:        projectId,
		ProjectName:      projName,
		IsAdminZone:      isAdminZone,
		OldValues:        nil,
		NewValues:        newValuesJSON,
		ChangedFields:    pq.StringArray{"*"}, // Todos os campos são novos
		IpAddress:        ipAddress,
		UserAgent:        userAgent,
		CreatedAt:        time.Now(),
	}

	if err := h.repo.Create(log); err != nil {
		fmt.Printf("❌ Erro ao criar log de auditoria (CREATE): %v\n", err)
		return err
	}

	fmt.Printf("✅ Log de auditoria registrado: %s criou usuário %s\n", actor.Email, target.Email)
	return nil
}

// LogUserUpdate - Registra log de atualização de usuário
func (h *AdminAuditLogHandler) LogUserUpdate(actor *models.User, oldUser, newUser *models.User, orgId, projectId *uuid.UUID, isAdminZone bool, ipAddress, userAgent string) error {
	// Detectar campos alterados
	changedFields := h.detectChangedFields(oldUser, newUser)

	// Se nenhum campo foi alterado, não logar
	if len(changedFields) == 0 {
		return nil
	}

	// Preparar dados antigos e novos (sem senha)
	oldValues := h.sanitizeUserForLog(oldUser)
	newValues := h.sanitizeUserForLog(newUser)

	oldValuesJSON, _ := json.Marshal(oldValues)
	newValuesJSON, _ := json.Marshal(newValues)

	// Obter nomes de org/projeto se disponíveis
	orgName, projName := h.getOrgProjectNames(orgId, projectId)

	log := &models.AdminAuditLog{
		Id:               uuid.New(),
		ActorId:          actor.Id,
		ActorEmail:       actor.Email,
		TargetId:         newUser.Id,
		TargetEmail:      newUser.Email,
		Action:           models.AdminAuditActionUpdate,
		EntityType:       models.AdminAuditEntityUser,
		OrganizationId:   orgId,
		OrganizationName: orgName,
		ProjectId:        projectId,
		ProjectName:      projName,
		IsAdminZone:      isAdminZone,
		OldValues:        oldValuesJSON,
		NewValues:        newValuesJSON,
		ChangedFields:    changedFields,
		IpAddress:        ipAddress,
		UserAgent:        userAgent,
		CreatedAt:        time.Now(),
	}

	if err := h.repo.Create(log); err != nil {
		fmt.Printf("❌ Erro ao criar log de auditoria (UPDATE): %v\n", err)
		return err
	}

	fmt.Printf("✅ Log de auditoria registrado: %s atualizou usuário %s (campos: %v)\n", actor.Email, newUser.Email, changedFields)
	return nil
}

// LogUserDelete - Registra log de exclusão de usuário
func (h *AdminAuditLogHandler) LogUserDelete(actor *models.User, target *models.User, orgId, projectId *uuid.UUID, isAdminZone bool, ipAddress, userAgent string) error {
	// Preparar dados do usuário excluído (sem senha)
	oldValues := h.sanitizeUserForLog(target)
	oldValuesJSON, _ := json.Marshal(oldValues)

	// Obter nomes de org/projeto se disponíveis
	orgName, projName := h.getOrgProjectNames(orgId, projectId)

	log := &models.AdminAuditLog{
		Id:               uuid.New(),
		ActorId:          actor.Id,
		ActorEmail:       actor.Email,
		TargetId:         target.Id,
		TargetEmail:      target.Email,
		Action:           models.AdminAuditActionDelete,
		EntityType:       models.AdminAuditEntityUser,
		OrganizationId:   orgId,
		OrganizationName: orgName,
		ProjectId:        projectId,
		ProjectName:      projName,
		IsAdminZone:      isAdminZone,
		OldValues:        oldValuesJSON,
		NewValues:        nil,
		ChangedFields:    pq.StringArray{"*"}, // Todos os campos foram removidos
		IpAddress:        ipAddress,
		UserAgent:        userAgent,
		CreatedAt:        time.Now(),
	}

	if err := h.repo.Create(log); err != nil {
		fmt.Printf("❌ Erro ao criar log de auditoria (DELETE): %v\n", err)
		return err
	}

	fmt.Printf("✅ Log de auditoria registrado: %s excluiu usuário %s\n", actor.Email, target.Email)
	return nil
}

// LogPasswordReset - Registra log de reset de senha
func (h *AdminAuditLogHandler) LogPasswordReset(actor *models.User, target *models.User, orgId, projectId *uuid.UUID, isAdminZone bool, ipAddress, userAgent string) error {
	// Obter nomes de org/projeto se disponíveis
	orgName, projName := h.getOrgProjectNames(orgId, projectId)

	log := &models.AdminAuditLog{
		Id:               uuid.New(),
		ActorId:          actor.Id,
		ActorEmail:       actor.Email,
		TargetId:         target.Id,
		TargetEmail:      target.Email,
		Action:           models.AdminAuditActionResetPassword,
		EntityType:       models.AdminAuditEntityUser,
		OrganizationId:   orgId,
		OrganizationName: orgName,
		ProjectId:        projectId,
		ProjectName:      projName,
		IsAdminZone:      isAdminZone,
		OldValues:        nil, // Não logamos senhas
		NewValues:        nil, // Não logamos senhas
		ChangedFields:    pq.StringArray{"password"},
		IpAddress:        ipAddress,
		UserAgent:        userAgent,
		CreatedAt:        time.Now(),
	}

	if err := h.repo.Create(log); err != nil {
		fmt.Printf("❌ Erro ao criar log de auditoria (RESET_PASSWORD): %v\n", err)
		return err
	}

	fmt.Printf("✅ Log de auditoria registrado: %s resetou senha do usuário %s\n", actor.Email, target.Email)
	return nil
}

// detectChangedFields - Detecta quais campos foram alterados entre duas versões do usuário
func (h *AdminAuditLogHandler) detectChangedFields(old, new *models.User) pq.StringArray {
	var changed []string

	if old.Name != new.Name {
		changed = append(changed, "name")
	}
	if old.Email != new.Email {
		changed = append(changed, "email")
	}
	if old.Active != new.Active {
		changed = append(changed, "active")
	}
	// Verificar se password foi alterado (se new.Password não está vazio e é diferente)
	if new.Password != "" && new.Password != old.Password {
		changed = append(changed, "password")
	}
	// Comparar permissions
	if !reflect.DeepEqual([]string(old.Permissions), []string(new.Permissions)) {
		changed = append(changed, "permissions")
	}

	return changed
}

// sanitizeUserForLog - Remove dados sensíveis do usuário para logging
func (h *AdminAuditLogHandler) sanitizeUserForLog(user *models.User) map[string]interface{} {
	if user == nil {
		return nil
	}

	return map[string]interface{}{
		"id":            user.Id.String(),
		"name":          user.Name,
		"email":         user.Email,
		"active":        user.Active,
		"permissions":   user.Permissions,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
		"last_access_at": user.LastAccessAt,
	}
}

// getOrgProjectNames - Obtém nomes de organização e projeto
func (h *AdminAuditLogHandler) getOrgProjectNames(orgId, projectId *uuid.UUID) (string, string) {
	var orgName, projName string

	if orgId != nil && h.orgRepo != nil {
		org, err := h.orgRepo.GetOrganizationById(*orgId)
		if err == nil && org != nil {
			orgName = org.Name
		}
	}

	if projectId != nil && h.projRepo != nil {
		proj, err := h.projRepo.GetProjectById(*projectId)
		if err == nil && proj != nil {
			projName = proj.Name
		}
	}

	return orgName, projName
}

// ==================== Métodos Genéricos para Auditoria ====================

// AuditLogParams - Parâmetros para log de auditoria genérico
type AuditLogParams struct {
	ActorId       uuid.UUID
	ActorEmail    string
	TargetId      uuid.UUID
	TargetEmail   string // Opcional, use "" se não aplicável
	Action        string // Use constantes AdminAuditAction*
	EntityType    string // Use constantes AdminAuditEntity*
	OrgId         *uuid.UUID
	ProjectId     *uuid.UUID
	IsAdminZone   bool
	OldValues     interface{} // Será convertido para JSON
	NewValues     interface{} // Será convertido para JSON
	ChangedFields []string
	IpAddress     string
	UserAgent     string
}

// LogGenericAction - Registra um log de auditoria genérico para qualquer tipo de ação
func (h *AdminAuditLogHandler) LogGenericAction(params AuditLogParams) error {
	var oldValuesJSON, newValuesJSON []byte
	var err error

	// Converter OldValues para JSON se não for nil
	if params.OldValues != nil {
		oldValuesJSON, err = json.Marshal(params.OldValues)
		if err != nil {
			fmt.Printf("⚠️ Erro ao serializar OldValues: %v\n", err)
			oldValuesJSON = nil
		}
	}

	// Converter NewValues para JSON se não for nil
	if params.NewValues != nil {
		newValuesJSON, err = json.Marshal(params.NewValues)
		if err != nil {
			fmt.Printf("⚠️ Erro ao serializar NewValues: %v\n", err)
			newValuesJSON = nil
		}
	}

	// Obter nomes de org/projeto
	orgName, projName := h.getOrgProjectNames(params.OrgId, params.ProjectId)

	log := &models.AdminAuditLog{
		Id:               uuid.New(),
		ActorId:          params.ActorId,
		ActorEmail:       params.ActorEmail,
		TargetId:         params.TargetId,
		TargetEmail:      params.TargetEmail,
		Action:           params.Action,
		EntityType:       params.EntityType,
		OrganizationId:   params.OrgId,
		OrganizationName: orgName,
		ProjectId:        params.ProjectId,
		ProjectName:      projName,
		IsAdminZone:      params.IsAdminZone,
		OldValues:        oldValuesJSON,
		NewValues:        newValuesJSON,
		ChangedFields:    params.ChangedFields,
		IpAddress:        params.IpAddress,
		UserAgent:        params.UserAgent,
		CreatedAt:        time.Now(),
	}

	if err := h.repo.Create(log); err != nil {
		fmt.Printf("❌ Erro ao criar log de auditoria (%s/%s): %v\n", params.Action, params.EntityType, err)
		return err
	}

	fmt.Printf("✅ Log de auditoria registrado: %s executou %s em %s (target: %s)\n",
		params.ActorEmail, params.Action, params.EntityType, params.TargetId)
	return nil
}

// LogRoleAssignment - Registra log de atribuição de cargo a usuário
func (h *AdminAuditLogHandler) LogRoleAssignment(
	actorId uuid.UUID,
	actorEmail string,
	targetUserId uuid.UUID,
	targetUserEmail string,
	roleId uuid.UUID,
	roleName string,
	orgId, projectId *uuid.UUID,
	ipAddress, userAgent string,
) error {
	return h.LogGenericAction(AuditLogParams{
		ActorId:       actorId,
		ActorEmail:    actorEmail,
		TargetId:      targetUserId,
		TargetEmail:   targetUserEmail,
		Action:        models.AdminAuditActionAssign,
		EntityType:    models.AdminAuditEntityUserRole,
		OrgId:         orgId,
		ProjectId:     projectId,
		IsAdminZone:   true,
		OldValues:     nil,
		NewValues:     map[string]interface{}{"role_id": roleId.String(), "role_name": roleName},
		ChangedFields: []string{"role_id", "role_name"},
		IpAddress:     ipAddress,
		UserAgent:     userAgent,
	})
}

// LogRoleRemoval - Registra log de remoção de cargo de usuário
func (h *AdminAuditLogHandler) LogRoleRemoval(
	actorId uuid.UUID,
	actorEmail string,
	targetUserId uuid.UUID,
	targetUserEmail string,
	roleId uuid.UUID,
	roleName string,
	orgId, projectId *uuid.UUID,
	ipAddress, userAgent string,
) error {
	return h.LogGenericAction(AuditLogParams{
		ActorId:       actorId,
		ActorEmail:    actorEmail,
		TargetId:      targetUserId,
		TargetEmail:   targetUserEmail,
		Action:        models.AdminAuditActionRemove,
		EntityType:    models.AdminAuditEntityUserRole,
		OrgId:         orgId,
		ProjectId:     projectId,
		IsAdminZone:   true,
		OldValues:     map[string]interface{}{"role_id": roleId.String(), "role_name": roleName},
		NewValues:     nil,
		ChangedFields: []string{"role_id", "role_name"},
		IpAddress:     ipAddress,
		UserAgent:     userAgent,
	})
}
