package handler

import (
	"lep/repositories"
	"lep/repositories/models"
)

type resourcePurchases struct {
	repo *repositories.DBconn
}

type IHandlerPurchases interface {
	GetPurchases(id string) (*models.Purchases, error)
	GetPurchasesByGroup(id string) ([]models.Purchases, error)
	CreatePurchase(purchase *models.Purchases) error
	UpdatePurchase(updatedPurchase *models.Purchases) error
	DeletePurchase(id string) error
}

func (r *resourcePurchases) GetPurchases(id string) (*models.Purchases, error) {
	resp, err := r.repo.Purchases.GetPurchases(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourcePurchases) GetPurchasesByGroup(id string) ([]models.Purchases, error) {
	resp, err := r.repo.Purchases.GetPurchasesByGroup(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourcePurchases) CreatePurchase(purchase *models.Purchases) error {
	purchase.Active = true

	err := r.repo.Purchases.CreatePurchase(purchase)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourcePurchases) UpdatePurchase(updatedPurchase *models.Purchases) error {
	err := r.repo.Purchases.UpdatePurchase(updatedPurchase)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourcePurchases) DeletePurchase(id string) error {
	err := r.repo.Purchases.DeletePurchase(id)
	if err != nil {
		return err
	}
	return nil
}

func NewSourceHandlerPurchases(repo *repositories.DBconn) IHandlerPurchases {
	return &resourcePurchases{repo: repo}
}
