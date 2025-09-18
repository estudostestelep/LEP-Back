package handler

import (
	"errors"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"time"
)

type resourceAuth struct {
	repo *repositories.DBconn
}

type IHandlerAuth interface {
	PostToken(user *models.User, token string) error
	Logout(token string) error
	VerificationToken(token string) (*models.User, error)
}

func (r *resourceAuth) PostToken(user *models.User, token string) error {
	loggedList := &models.LoggedLists{
		Token:     token,
		UserEmail: user.Email,
		UserId:    user.Id,
	}

	if err := r.repo.LoggedLists.CreateLoggedList(loggedList); err != nil {
		if err := r.repo.LoggedLists.DeleteLoggedList(loggedList.Token); err != nil {
			return fmt.Errorf("falha ao remover token existente: %v", err)
		}

		if err := r.repo.LoggedLists.CreateLoggedList(loggedList); err != nil {
			return fmt.Errorf("falha ao criar registro na LoggedLists: %v", err)
		}
	}

	return nil
}

func (r *resourceAuth) Logout(token string) error {
	bannedList := &models.BannedLists{
		Token: token,
	}

	if err := r.repo.BannedLists.CreateBannedList(bannedList); err != nil {
		return fmt.Errorf("falha ao criar registro na BannedLists: %v", err)
	}

	if err := r.repo.LoggedLists.DeleteLoggedList(token); err != nil {
		return fmt.Errorf("falha ao remover da LoggedLists: %v", err)
	}

	r.cleanupExpiredTokens()

	return nil
}

func (r *resourceAuth) cleanupExpiredTokens() {
	cutoffTime := time.Now().AddDate(0, 0, -7)

	resp, err := r.repo.BannedLists.GetBannedAllList()
	if err != nil || resp == nil {
		return
	}

	for _, item := range *resp {
		if item.CreatedAt.Before(cutoffTime) {
			r.repo.BannedLists.DeleteBannedList(item.BannedListId)
		}
	}
}

func (r *resourceAuth) VerificationToken(token string) (*models.User, error) {

	logged, err := r.repo.LoggedLists.GetLoggedToken(token)
	if err != nil {
		return nil, err
	}

	if logged == nil {
		return nil, errors.New("Not found")
	}

	user, err := r.repo.User.GetUserById(logged.UserId.String())
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewAuthHandler(repo *repositories.DBconn) IHandlerAuth {
	return &resourceAuth{repo: repo}
}
