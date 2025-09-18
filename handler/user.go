package handler

import (
	"fmt"
	"lep/repositories"
	"lep/repositories/models"

	"errors"

	"golang.org/x/crypto/bcrypt"
)

type resourceUser struct {
	repo *repositories.DBconn
}

type IHandlerUser interface {
	GetUser(id string) (*models.User, error)
	GetUserByGroup(id string) ([]models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(updatedUser *models.User) error
	DeleteUser(id string) error
	GetUserByEmail(email string) (*models.User, error)
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

func (r *resourceUser) CreateUser(user *models.User) error {
	existingUser, _ := r.repo.User.GetUserByEmail(user.Email)

	if existingUser != nil {
		return errors.New("E-mail já cadastrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	err = r.repo.User.CreateUser(user)
	if err != nil {
		return err
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
	err := r.repo.User.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceUser) GetUserByEmail(email string) (*models.User, error) {
	resp, err := r.repo.User.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewSourceHandlerUser(repo *repositories.DBconn) IHandlerUser {
	return &resourceUser{repo: repo}
}
