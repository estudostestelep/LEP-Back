package handler

import (
	"errors"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type resourceAdminUser struct {
	repo *repositories.DBconn
}

type IHandlerAdminUser interface {
	GetAdminById(id string) (*models.Admin, error)
	GetAdminByEmail(email string) (*models.Admin, error)
	ListAdmins() ([]models.Admin, error)
	CreateAdmin(admin *models.Admin) error
	UpdateAdmin(admin *models.Admin) error
	DeleteAdmin(id string) error
	UpdateLastAccess(adminId string) error
	ValidateAdminCredentials(email, password string) (*models.Admin, error)
	ValidatePermissions(permissions []string) error
	ListAvailablePermissions() ([]models.Permission, error)
	// Role management
	GetRole(roleId string) (*models.Role, error)
	GetPermissionsFromRole(roleId string) ([]string, error)
	GetUserMaxHierarchyLevel(userId string) (int, error)
	AssignRoleToAdmin(adminRole *models.AdminRole) error
	GetAdminRoles(adminId string) ([]models.AdminRole, error)
}

func (r *resourceAdminUser) GetAdminById(id string) (*models.Admin, error) {
	return r.repo.Admins.GetAdminById(id)
}

func (r *resourceAdminUser) GetAdminByEmail(email string) (*models.Admin, error) {
	return r.repo.Admins.GetAdminByEmail(email)
}

func (r *resourceAdminUser) ListAdmins() ([]models.Admin, error) {
	return r.repo.Admins.ListAdmins()
}

func (r *resourceAdminUser) CreateAdmin(admin *models.Admin) error {
	// Verificar se email já existe
	exists, err := r.repo.Admins.AdminEmailExists(admin.Email)
	if err != nil {
		return fmt.Errorf("erro ao verificar email: %v", err)
	}
	if exists {
		return errors.New("email já cadastrado")
	}

	// Gerar UUID se não fornecido
	if admin.Id == uuid.Nil {
		admin.Id = uuid.New()
	}

	// Hash da senha
	if admin.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("erro ao gerar hash da senha: %v", err)
		}
		admin.Password = string(hashedPassword)
	}

	return r.repo.Admins.CreateAdmin(admin)
}

func (r *resourceAdminUser) UpdateAdmin(admin *models.Admin) error {
	// Se a senha foi fornecida, fazer hash
	if admin.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("erro ao gerar hash da senha: %v", err)
		}
		admin.Password = string(hashedPassword)
	}

	return r.repo.Admins.UpdateAdmin(admin)
}

func (r *resourceAdminUser) DeleteAdmin(id string) error {
	return r.repo.Admins.DeleteAdmin(id) // Hard delete - remove permanentemente
}

func (r *resourceAdminUser) UpdateLastAccess(adminId string) error {
	return r.repo.Admins.UpdateLastAccess(adminId)
}

func (r *resourceAdminUser) ValidateAdminCredentials(email, password string) (*models.Admin, error) {
	admin, err := r.repo.Admins.GetAdminByEmail(email)
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	if admin == nil || !admin.IsActive() {
		return nil, errors.New("credenciais inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	return admin, nil
}

// ValidatePermissions valida se as permissões existem no banco
// master_admin é sempre válido (permissão especial de admin)
func (r *resourceAdminUser) ValidatePermissions(permissions []string) error {
	if len(permissions) == 0 {
		return errors.New("pelo menos uma permissão é obrigatória")
	}

	// Buscar permissões que não são master_admin
	var permissionsToValidate []string
	for _, p := range permissions {
		if p == "" {
			return errors.New("permissão não pode ser vazia")
		}
		// master_admin é sempre válido (permissão especial)
		if p != "master_admin" {
			permissionsToValidate = append(permissionsToValidate, p)
		}
	}

	// Se só tem master_admin, não precisa validar no banco
	if len(permissionsToValidate) == 0 {
		return nil
	}

	// Buscar permissões no banco
	dbPermissions, err := r.repo.Permissions.GetByCodeNames(permissionsToValidate)
	if err != nil {
		return fmt.Errorf("erro ao buscar permissões: %v", err)
	}

	// Criar mapa de permissões encontradas
	found := make(map[string]bool)
	for _, p := range dbPermissions {
		found[p.Code] = true
	}

	// Verificar se todas as permissões existem
	for _, p := range permissionsToValidate {
		if !found[p] {
			return fmt.Errorf("permissão inválida: %s", p)
		}
	}

	return nil
}

// ListAvailablePermissions lista todas as permissões disponíveis no sistema
func (r *resourceAdminUser) ListAvailablePermissions() ([]models.Permission, error) {
	return r.repo.Permissions.List()
}

// GetRole busca role por ID
func (r *resourceAdminUser) GetRole(roleId string) (*models.Role, error) {
	return r.repo.Roles.GetById(roleId)
}

// GetPermissionsFromRole extrai permissões de um role
func (r *resourceAdminUser) GetPermissionsFromRole(roleId string) ([]string, error) {
	codes, err := r.repo.Roles.GetRolePermissionCodes(roleId)
	if err != nil {
		return nil, err
	}

	// Se role tem HierarchyLevel >= 10, adicionar indicação de master_admin
	role, _ := r.repo.Roles.GetById(roleId)
	if role != nil && role.HierarchyLevel >= 10 {
		codes = append(codes, "master_admin")
	}

	return codes, nil
}

// GetUserMaxHierarchyLevel retorna o nível máximo de hierarquia do admin
func (r *resourceAdminUser) GetUserMaxHierarchyLevel(userId string) (int, error) {
	return r.repo.Roles.GetAdminMaxHierarchyLevel(userId)
}

// AssignRoleToAdmin cria AdminRole para admin
func (r *resourceAdminUser) AssignRoleToAdmin(adminRole *models.AdminRole) error {
	return r.repo.Roles.AssignRoleToAdmin(adminRole)
}

// GetAdminRoles busca todos os roles de um admin
func (r *resourceAdminUser) GetAdminRoles(adminId string) ([]models.AdminRole, error) {
	return r.repo.Roles.GetAdminRoles(adminId)
}

func NewAdminUserHandler(repo *repositories.DBconn) IHandlerAdminUser {
	return &resourceAdminUser{repo: repo}
}
