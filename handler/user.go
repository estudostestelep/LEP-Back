package handler

import (
	"errors"
	"fmt"
	"lep/constants"
	"lep/repositories"
	"lep/repositories/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type resourceUser struct {
	repo              *repositories.DBconn
	roleRepo          repositories.IRoleRepository
	roleHandler       *RoleHandler
	adminAuditHandler IAdminAuditLogHandler
}

type IHandlerUser interface {
	GetUser(id string) (*models.User, error)
	GetUserByGroup(id string) ([]models.User, error)
	ListUsers(orgId, projectId string) ([]models.User, error)
	ListUsersWithRoles(orgId, projectId string) ([]models.UserWithRole, error)
	CreateUser(user *models.User, orgId, projectId, roleId string) error
	UpdateUser(updatedUser *models.User) error
	UpdateLastAccess(userId string) error
	DeleteUser(id string) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserWithRelations(id string) (*models.UserWithRelations, error)

	// Métodos com contexto para auditoria (usados pelas rotas admin)
	CreateUserWithContext(ctx *RequestContext, user *models.User, orgId, projectId, roleId string) error
	UpdateUserWithContext(ctx *RequestContext, userId string, updatedUser *models.User) error
	DeleteUserWithContext(ctx *RequestContext, userId string) error
}

func (r *resourceUser) GetUser(id string) (*models.User, error) {
	resp, err := r.repo.User.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceUser) GetUserByGroup(id string) ([]models.User, error) {
	resp, err := r.repo.User.GetUsersByGroup(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceUser) ListUsers(orgId, projectId string) ([]models.User, error) {
	resp, err := r.repo.User.ListUsersByOrganizationAndProject(orgId, projectId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceUser) ListUsersWithRoles(orgId, projectId string) ([]models.UserWithRole, error) {
	return r.repo.User.ListUsersWithRoles(orgId, projectId)
}

func (r *resourceUser) CreateUser(user *models.User, orgId, projectId, roleId string) error {
	existingUser, _ := r.repo.User.GetUserByEmail(user.Email)

	if existingUser != nil {
		return errors.New("E-mail já cadastrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	// Gerar ID apenas se não foi fornecido
	if user.Id == uuid.Nil {
		user.Id = uuid.New()
	}

	err = r.repo.User.CreateUser(user)
	if err != nil {
		return err
	}

	// 🔑 REGRA DE NEGÓCIO: Se o novo usuário é um master admin, adicioná-lo a todas as organizações
	isMasterAdmin := constants.HasPermission(user.Permissions, constants.PermissionMasterAdmin)
	if isMasterAdmin {
		if err := r.addMasterAdminToAllOrganizations(user.Id); err != nil {
			// Log error but don't fail user creation
			fmt.Printf("Aviso: erro ao adicionar master admin a organizações: %v\n", err)
		}
	} else {
		// Vincular usuário à organização e projeto especificados
		if err := r.linkUserToOrgAndProject(user.Id, orgId, projectId); err != nil {
			fmt.Printf("Aviso: erro ao vincular usuário a org/projeto: %v\n", err)
		}
	}

	// 🔑 ATRIBUIR CARGO: Se roleId foi fornecido, atribuir o cargo ao usuário
	// Para roles admin, orgId pode ser vazio (cargo global da zona administrativa)
	if roleId != "" && r.roleRepo != nil {
		if err := r.assignRoleToUser(user.Id, roleId, orgId, projectId); err != nil {
			fmt.Printf("Aviso: erro ao atribuir cargo ao usuário: %v\n", err)
		}
	}

	return nil
}

func (r *resourceUser) UpdateUser(updatedUser *models.User) error {
	existingUser, err := r.repo.User.GetUserByEmail(updatedUser.Email)

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updatedUser.Password = string(hashedPassword)
	}

	if existingUser != nil && existingUser.Id != updatedUser.Id {
		return fmt.Errorf("E-mail já cadastrado")
	}

	err = r.repo.User.UpdateUser(updatedUser)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceUser) DeleteUser(id string) error {
	// Deletar user_roles associados antes de deletar o user (foreign key constraint)
	if err := r.roleRepo.DeleteUserRolesByUserId(id); err != nil {
		return fmt.Errorf("erro ao remover cargos do usuário: %w", err)
	}

	if err := r.repo.User.DeleteUser(id); err != nil {
		return err
	}
	return nil
}

func (r *resourceUser) UpdateLastAccess(userId string) error {
	return r.repo.User.UpdateLastAccess(userId)
}

func (r *resourceUser) GetUserByEmail(email string) (*models.User, error) {
	resp, err := r.repo.User.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceUser) GetUserWithRelations(id string) (*models.UserWithRelations, error) {
	resp, err := r.repo.User.GetUserWithRelations(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// linkUserToOrgAndProject vincula um usuário a uma organização e projeto via UserRole
// NOTA: Esta função requer que um roleId seja atribuído posteriormente via assignRoleToUser
// O acesso é gerenciado exclusivamente via user_roles no novo sistema
func (r *resourceUser) linkUserToOrgAndProject(userId uuid.UUID, orgId, projectId string) error {
	fmt.Printf("🔗 linkUserToOrgAndProject: userId=%s, orgId=%s, projectId=%s\n", userId, orgId, projectId)
	fmt.Printf("ℹ️ No novo sistema, o acesso é gerenciado via UserRole com assignRoleToUser\n")
	// O acesso agora é gerenciado exclusivamente via user_roles
	// Esta função é mantida por compatibilidade mas não faz nada
	// O roleId deve ser fornecido na criação do usuário
	return nil
}

// assignRoleToUser atribui um cargo a um usuário no sistema de permissões granulares
// Para roles de escopo "admin", orgId pode ser vazio (cargo global da zona administrativa)
// IMPORTANTE: Sincroniza permissão master_admin quando role é super_admin
func (r *resourceUser) assignRoleToUser(userId uuid.UUID, roleId, orgId, projectId string) error {
	fmt.Printf("🔑 assignRoleToUser: userId=%s, roleId=%s, orgId=%s, projectId=%s\n", userId, roleId, orgId, projectId)

	roleUUID, err := uuid.Parse(roleId)
	if err != nil {
		return fmt.Errorf("ID do cargo inválido: %v", err)
	}

	// Verificar se o cargo é super_admin para sincronizar permissão master_admin
	role, err := r.roleRepo.GetById(roleId)
	if err != nil {
		return fmt.Errorf("erro ao buscar cargo: %v", err)
	}

	// Se o cargo é super_admin, adicionar permissão master_admin ao usuário
	if role.Name == "super_admin" {
		user, err := r.repo.User.GetUserById(userId.String())
		if err != nil {
			return fmt.Errorf("erro ao buscar usuário: %v", err)
		}

		// Verificar se já tem a permissão
		hasMasterAdmin := false
		for _, perm := range user.Permissions {
			if perm == "master_admin" {
				hasMasterAdmin = true
				break
			}
		}

		// Adicionar permissão se não tiver
		if !hasMasterAdmin {
			user.Permissions = append(user.Permissions, "master_admin")
			if err := r.repo.User.UpdateUser(user); err != nil {
				return fmt.Errorf("erro ao atualizar permissões do usuário: %w", err)
			}
			fmt.Printf("✅ Permissão master_admin adicionada ao usuário %s\n", userId)
		}
	}

	userRole := &models.UserRole{
		Id:     uuid.New(),
		UserId: userId,
		RoleId: roleUUID,
		Active: true,
	}

	// Se orgId foi fornecido, adicionar ao contexto
	// Se vazio, OrganizationId fica nil (cargo admin global)
	if orgId != "" {
		orgUUID, err := uuid.Parse(orgId)
		if err != nil {
			return fmt.Errorf("ID da organização inválido: %v", err)
		}
		userRole.OrganizationId = &orgUUID
	}

	// Se projectId foi fornecido, adicionar ao contexto
	if projectId != "" {
		projUUID, err := uuid.Parse(projectId)
		if err == nil {
			userRole.ProjectId = &projUUID
		}
	}

	// Atribuir cargo diretamente (durante criação de usuário não usa validação de hierarquia)
	if err := r.roleRepo.AssignRoleToUser(userRole); err != nil {
		return fmt.Errorf("erro ao atribuir cargo: %v", err)
	}

	fmt.Printf("✅ Cargo atribuído com sucesso ao usuário %s\n", userId)
	return nil
}

// addMasterAdminToAllOrganizations adiciona um novo master admin a todas as organizações existentes
// REGRA DE NEGÓCIO: Master admins têm acesso global via permissão master_admin, não precisam de UserRole por org
// O sistema verifica a permissão master_admin diretamente no middleware de autorização
func (r *resourceUser) addMasterAdminToAllOrganizations(userId uuid.UUID) error {
	fmt.Printf("ℹ️ Master admins têm acesso global via permissão master_admin no array de permissions\n")
	fmt.Printf("ℹ️ Não é necessário criar UserRole para cada organização\n")
	// Master admins têm acesso global via a permissão "master_admin" no array de permissions do User
	// O middleware de autorização verifica essa permissão e concede acesso total
	// Não é mais necessário criar UserOrganization/UserProject para cada org
	return nil
}

// CreateUserWithContext cria um usuário e registra auditoria se o contexto for de admin
func (r *resourceUser) CreateUserWithContext(ctx *RequestContext, user *models.User, orgId, projectId, roleId string) error {
	// Executar criação normal
	if err := r.CreateUser(user, orgId, projectId, roleId); err != nil {
		return err
	}

	// Registrar auditoria se for Master Admin
	if ctx.IsMasterAdmin() && r.adminAuditHandler != nil {
		actor := &models.User{Id: ctx.UserId, Email: ctx.UserEmail}
		go func() {
			if err := r.adminAuditHandler.LogUserCreate(actor, user, ctx.OrganizationId, ctx.ProjectId, ctx.IsAdminZone, ctx.IpAddress, ctx.UserAgent); err != nil {
				fmt.Printf("⚠️ Erro ao registrar log de auditoria (CREATE): %v\n", err)
			}
		}()
	}

	return nil
}

// UpdateUserWithContext atualiza um usuário e registra auditoria se o contexto for de admin
func (r *resourceUser) UpdateUserWithContext(ctx *RequestContext, userId string, updatedUser *models.User) error {

	// Executar atualização normal
	if err := r.UpdateUser(updatedUser); err != nil {
		return err
	}

	// Registrar auditoria se for Master Admin
	if ctx.IsMasterAdmin() && r.adminAuditHandler != nil {
		// Capturar estado anterior do usuário ANTES do update (necessário para auditoria)
		oldUser, _ := r.GetUser(userId)

		actor := &models.User{Id: ctx.UserId, Email: ctx.UserEmail}
		go func() {
			// Se foi reset de senha, logar separadamente
			if updatedUser.Password != "" {
				if err := r.adminAuditHandler.LogPasswordReset(actor, updatedUser, ctx.OrganizationId, ctx.ProjectId, ctx.IsAdminZone, ctx.IpAddress, ctx.UserAgent); err != nil {
					fmt.Printf("⚠️ Erro ao registrar log de auditoria (RESET_PASSWORD): %v\n", err)
				}
			}

			// Logar alterações de outros campos
			if err := r.adminAuditHandler.LogUserUpdate(actor, oldUser, updatedUser, ctx.OrganizationId, ctx.ProjectId, ctx.IsAdminZone, ctx.IpAddress, ctx.UserAgent); err != nil {
				fmt.Printf("⚠️ Erro ao registrar log de auditoria (UPDATE): %v\n", err)
			}
		}()

		isMasterAdmin := constants.HasPermission(updatedUser.Permissions, constants.PermissionMasterAdmin)
		if isMasterAdmin {
			if err := r.addMasterAdminToAllOrganizations(updatedUser.Id); err != nil {
				// Log error but don't fail user creation
				fmt.Printf("Aviso: erro ao adicionar master admin a organizações: %v\n", err)
			}
		}
	}

	return nil
}

// DeleteUserWithContext exclui um usuário e registra auditoria se o contexto for de admin
func (r *resourceUser) DeleteUserWithContext(ctx *RequestContext, userId string) error {
	// Capturar dados do usuário ANTES da exclusão (necessário para auditoria)
	targetUser, err := r.GetUser(userId)
	if err != nil {
		return fmt.Errorf("usuário não encontrado: %w", err)
	}

	// Executar exclusão normal
	if err := r.DeleteUser(userId); err != nil {
		return err
	}

	// Registrar auditoria se for Master Admin
	if ctx.IsMasterAdmin() && r.adminAuditHandler != nil && targetUser != nil {
		actor := &models.User{Id: ctx.UserId, Email: ctx.UserEmail}
		go func() {
			if err := r.adminAuditHandler.LogUserDelete(actor, targetUser, ctx.OrganizationId, ctx.ProjectId, ctx.IsAdminZone, ctx.IpAddress, ctx.UserAgent); err != nil {
				fmt.Printf("⚠️ Erro ao registrar log de auditoria (DELETE): %v\n", err)
			}
		}()
	}

	return nil
}

func NewSourceHandlerUser(repo *repositories.DBconn, roleRepo repositories.IRoleRepository, roleHandler *RoleHandler, adminAuditHandler IAdminAuditLogHandler) IHandlerUser {
	return &resourceUser{repo: repo, roleRepo: roleRepo, roleHandler: roleHandler, adminAuditHandler: adminAuditHandler}
}
