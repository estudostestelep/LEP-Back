package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceStaffStock struct {
	repo *repositories.DBconn
}

type IHandlerStaffStock interface {
	// Stock Items
	GetItemById(id string) (*models.StaffStockItem, error)
	ListItems(orgId, projectId string) ([]models.StaffStockItem, error)
	ListItemsBySector(orgId, projectId, sector string) ([]models.StaffStockItem, error)
	ListSectors(orgId, projectId string) ([]models.SectorSummary, error)
	CreateItem(item *models.StaffStockItem) error
	UpdateItem(item *models.StaffStockItem) error
	DeleteItem(id string) error

	// Stock Records
	GetRecordById(id string) (*models.StaffStockRecord, error)
	GetRecordByIdWithItems(id string) (*models.StaffStockRecord, error)
	ListRecords(orgId, projectId string, limit int) ([]models.StaffStockRecord, error)
	ListRecordsBySector(orgId, projectId, sector string, limit int) ([]models.StaffStockRecord, error)
	CreateRecordWithItems(req *models.CreateStockRecordRequest, orgId, projectId, createdById string) (*models.StaffStockRecord, error)
	GenerateShoppingList(recordId string) (*models.ShoppingList, error)
}

func NewStaffStockHandler(repo *repositories.DBconn) IHandlerStaffStock {
	return &resourceStaffStock{repo: repo}
}

// ==================== Stock Items ====================

func (r *resourceStaffStock) GetItemById(id string) (*models.StaffStockItem, error) {
	itemId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffStock.GetItemById(itemId)
}

func (r *resourceStaffStock) ListItems(orgId, projectId string) ([]models.StaffStockItem, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffStock.ListItems(orgUUID, projUUID)
}

func (r *resourceStaffStock) ListItemsBySector(orgId, projectId, sector string) ([]models.StaffStockItem, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffStock.ListItemsBySector(orgUUID, projUUID, sector)
}

func (r *resourceStaffStock) ListSectors(orgId, projectId string) ([]models.SectorSummary, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffStock.ListSectors(orgUUID, projUUID)
}

func (r *resourceStaffStock) CreateItem(item *models.StaffStockItem) error {
	return r.repo.StaffStock.CreateItem(item)
}

func (r *resourceStaffStock) UpdateItem(item *models.StaffStockItem) error {
	return r.repo.StaffStock.UpdateItem(item)
}

func (r *resourceStaffStock) DeleteItem(id string) error {
	itemId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.StaffStock.SoftDeleteItem(itemId)
}

// ==================== Stock Records ====================

func (r *resourceStaffStock) GetRecordById(id string) (*models.StaffStockRecord, error) {
	recordId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffStock.GetRecordById(recordId)
}

func (r *resourceStaffStock) GetRecordByIdWithItems(id string) (*models.StaffStockRecord, error) {
	recordId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffStock.GetRecordByIdWithItems(recordId)
}

func (r *resourceStaffStock) ListRecords(orgId, projectId string, limit int) ([]models.StaffStockRecord, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffStock.ListRecords(orgUUID, projUUID, limit)
}

func (r *resourceStaffStock) ListRecordsBySector(orgId, projectId, sector string, limit int) ([]models.StaffStockRecord, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffStock.ListRecordsBySector(orgUUID, projUUID, sector, limit)
}

func (r *resourceStaffStock) CreateRecordWithItems(req *models.CreateStockRecordRequest, orgId, projectId, createdById string) (*models.StaffStockRecord, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	var createdByUUID *uuid.UUID
	if createdById != "" {
		u, err := uuid.Parse(createdById)
		if err == nil {
			createdByUUID = &u
		}
	}

	record := &models.StaffStockRecord{
		OrganizationId: orgUUID,
		ProjectId:      projUUID,
		RecordDate:     time.Now(),
		Sector:         req.Sector,
		CreatedById:    createdByUUID,
		Notes:          req.Notes,
		CreatedAt:      time.Now(),
	}

	var items []models.StaffStockRecordItem
	for _, itemInput := range req.Items {
		// Get item details for min/max calculation
		item, err := r.repo.StaffStock.GetItemById(itemInput.StockItemId)
		if err != nil {
			continue
		}

		toBuy := 0
		if item.StockMin != nil && item.StockMax != nil {
			if itemInput.CurrentStock < *item.StockMin {
				toBuy = *item.StockMax - itemInput.CurrentStock
			}
		}

		items = append(items, models.StaffStockRecordItem{
			StockItemId:  itemInput.StockItemId,
			CurrentStock: itemInput.CurrentStock,
			ToBuy:        toBuy,
			CreatedAt:    time.Now(),
		})
	}

	err = r.repo.StaffStock.CreateRecordWithItems(record, items)
	if err != nil {
		return nil, err
	}

	return r.GetRecordByIdWithItems(record.Id.String())
}

func (r *resourceStaffStock) GenerateShoppingList(recordId string) (*models.ShoppingList, error) {
	recordUUID, err := uuid.Parse(recordId)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffStock.GenerateShoppingList(recordUUID)
}
