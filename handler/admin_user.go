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
	return r.repo.Admins.SoftDeleteAdmin(id)
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

func NewAdminUserHandler(repo *repositories.DBconn) IHandlerAdminUser {
	return &resourceAdminUser{repo: repo}
}
