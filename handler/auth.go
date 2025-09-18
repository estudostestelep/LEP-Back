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

	_, err := r.repo.LoggedLists.GetLoggedToken(token)
	if err != nil {
		if err := r.repo.LoggedLists.CreateLoggedList(loggedList); err != nil {
			return fmt.Errorf("falha ao criar registro na LoggedLists: %v", err)
		}
	} else {

		if err := r.repo.LoggedLists.DeleteLoggedList(loggedList.Token); err != nil {
			return fmt.Errorf("falha ao deleter registro na LoggedLists: %v", err)
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
		Date:  time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := r.repo.BannedLists.CreateBannedList(bannedList); err != nil {
		return fmt.Errorf("falha ao criar registro na BannedLists: %v", err)
	}

	if err := r.repo.LoggedLists.DeleteLoggedList(token); err != nil {
		return fmt.Errorf("falha ao remover da LoggedLists: %v", err)
	}

	resp, err := r.repo.BannedLists.GetBannedAllList()
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("O ponteiro resp é nulo")
	}

	hoje := time.Now()

	for _, item := range *resp {
		data, err := time.Parse("2006-01-02 15:04:05", item.Date)
		if err != nil {
			continue
		}

		if hoje.After(data) {
			if err := r.repo.BannedLists.DeleteBannedList(item.BannedListId); err != nil {
				continue
			}
		}
	}

	return nil
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
