package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceStaffStock struct {
	db *gorm.DB
}

type IStaffStockRepository interface {
	// Stock Items CRUD
	GetItemById(id uuid.UUID) (*models.StaffStockItem, error)
	ListItems(orgId, projectId uuid.UUID) ([]models.StaffStockItem, error)
	ListItemsBySector(orgId, projectId uuid.UUID, sector string) ([]models.StaffStockItem, error)
	ListSectors(orgId, projectId uuid.UUID) ([]models.SectorSummary, error)
	CreateItem(item *models.StaffStockItem) error
	UpdateItem(item *models.StaffStockItem) error
	SoftDeleteItem(id uuid.UUID) error

	// Stock Records
	GetRecordById(id uuid.UUID) (*models.StaffStockRecord, error)
	GetRecordByIdWithItems(id uuid.UUID) (*models.StaffStockRecord, error)
	ListRecords(orgId, projectId uuid.UUID, limit int) ([]models.StaffStockRecord, error)
	ListRecordsBySector(orgId, projectId uuid.UUID, sector string, limit int) ([]models.StaffStockRecord, error)
	CreateRecord(record *models.StaffStockRecord) error
	CreateRecordItem(item *models.StaffStockRecordItem) error
	CreateRecordWithItems(record *models.StaffStockRecord, items []models.StaffStockRecordItem) error
	UpdateRecordPdfUrl(id uuid.UUID, pdfUrl string) error
	MarkRecordEmailSent(id uuid.UUID) error

	// Queries
	GetLastRecordForItem(orgId, projectId, itemId uuid.UUID) (*models.StaffStockRecordItem, error)
	GenerateShoppingList(recordId uuid.UUID) (*models.ShoppingList, error)
}

func NewStaffStockRepository(db *gorm.DB) IStaffStockRepository {
	return &resourceStaffStock{db: db}
}

// ==================== Stock Items ====================

func (r *resourceStaffStock) GetItemById(id uuid.UUID) (*models.StaffStockItem, error) {
	var item models.StaffStockItem
	err := r.db.First(&item, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *resourceStaffStock) ListItems(orgId, projectId uuid.UUID) ([]models.StaffStockItem, error) {
	var items []models.StaffStockItem
	err := r.db.
		Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL AND active = true", orgId, projectId).
		Order("sector ASC, category ASC, name ASC").
		Find(&items).Error
	return items, err
}

func (r *resourceStaffStock) ListItemsBySector(orgId, projectId uuid.UUID, sector string) ([]models.StaffStockItem, error) {
	var items []models.StaffStockItem
	err := r.db.
		Where("organization_id = ? AND project_id = ? AND sector = ? AND deleted_at IS NULL AND active = true", orgId, projectId, sector).
		Order("category ASC, name ASC").
		Find(&items).Error
	return items, err
}

func (r *resourceStaffStock) ListSectors(orgId, projectId uuid.UUID) ([]models.SectorSummary, error) {
	var sectors []models.SectorSummary
	err := r.db.Model(&models.StaffStockItem{}).
		Select("sector, COUNT(*) as item_count").
		Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL AND active = true", orgId, projectId).
		Group("sector").
		Order("sector ASC").
		Scan(&sectors).Error
	return sectors, err
}

func (r *resourceStaffStock) CreateItem(item *models.StaffStockItem) error {
	if item.Id == uuid.Nil {
		item.Id = uuid.New()
	}
	return r.db.Create(item).Error
}

func (r *resourceStaffStock) UpdateItem(item *models.StaffStockItem) error {
	return r.db.Save(item).Error
}

func (r *resourceStaffStock) SoftDeleteItem(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.StaffStockItem{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// ==================== Stock Records ====================

func (r *resourceStaffStock) GetRecordById(id uuid.UUID) (*models.StaffStockRecord, error) {
	var record models.StaffStockRecord
	err := r.db.First(&record, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *resourceStaffStock) GetRecordByIdWithItems(id uuid.UUID) (*models.StaffStockRecord, error) {
	var record models.StaffStockRecord
	err := r.db.
		Preload("CreatedBy").
		Preload("Items").
		Preload("Items.StockItem").
		First(&record, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *resourceStaffStock) ListRecords(orgId, projectId uuid.UUID, limit int) ([]models.StaffStockRecord, error) {
	var records []models.StaffStockRecord
	query := r.db.
		Preload("CreatedBy").
		Where("organization_id = ? AND project_id = ?", orgId, projectId).
		Order("record_date DESC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&records).Error
	return records, err
}

func (r *resourceStaffStock) ListRecordsBySector(orgId, projectId uuid.UUID, sector string, limit int) ([]models.StaffStockRecord, error) {
	var records []models.StaffStockRecord
	query := r.db.
		Preload("CreatedBy").
		Where("organization_id = ? AND project_id = ? AND sector = ?", orgId, projectId, sector).
		Order("record_date DESC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&records).Error
	return records, err
}

func (r *resourceStaffStock) CreateRecord(record *models.StaffStockRecord) error {
	if record.Id == uuid.Nil {
		record.Id = uuid.New()
	}
	return r.db.Create(record).Error
}

func (r *resourceStaffStock) CreateRecordItem(item *models.StaffStockRecordItem) error {
	if item.Id == uuid.Nil {
		item.Id = uuid.New()
	}
	return r.db.Create(item).Error
}

func (r *resourceStaffStock) CreateRecordWithItems(record *models.StaffStockRecord, items []models.StaffStockRecordItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if record.Id == uuid.Nil {
			record.Id = uuid.New()
		}

		if err := tx.Create(record).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].Id = uuid.New()
			items[i].StockRecordId = record.Id
		}

		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *resourceStaffStock) UpdateRecordPdfUrl(id uuid.UUID, pdfUrl string) error {
	return r.db.Model(&models.StaffStockRecord{}).Where("id = ?", id).Update("pdf_url", pdfUrl).Error
}

func (r *resourceStaffStock) MarkRecordEmailSent(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.StaffStockRecord{}).Where("id = ?", id).Update("email_sent_at", now).Error
}

// ==================== Queries ====================

func (r *resourceStaffStock) GetLastRecordForItem(orgId, projectId, itemId uuid.UUID) (*models.StaffStockRecordItem, error) {
	var recordItem models.StaffStockRecordItem
	err := r.db.
		Joins("JOIN staff_stock_records ON staff_stock_records.id = staff_stock_record_items.stock_record_id").
		Where("staff_stock_records.organization_id = ? AND staff_stock_records.project_id = ? AND staff_stock_record_items.stock_item_id = ?",
			orgId, projectId, itemId).
		Order("staff_stock_records.record_date DESC, staff_stock_records.created_at DESC").
		First(&recordItem).Error
	if err != nil {
		return nil, err
	}
	return &recordItem, nil
}

func (r *resourceStaffStock) GenerateShoppingList(recordId uuid.UUID) (*models.ShoppingList, error) {
	record, err := r.GetRecordByIdWithItems(recordId)
	if err != nil {
		return nil, err
	}

	shoppingList := &models.ShoppingList{
		RecordId:   record.Id,
		RecordDate: record.RecordDate,
		Sector:     record.Sector,
		ExtraNotes: record.Notes,
		Items:      []models.ShoppingListItem{},
	}

	for _, item := range record.Items {
		if item.ToBuy > 0 && item.StockItem != nil {
			stockItem := item.StockItem
			shoppingList.Items = append(shoppingList.Items, models.ShoppingListItem{
				ItemId:       stockItem.Id,
				ItemName:     stockItem.Name,
				Category:     stockItem.Category,
				CurrentStock: item.CurrentStock,
				StockMin:     derefIntOrZero(stockItem.StockMin),
				StockMax:     derefIntOrZero(stockItem.StockMax),
				ToBuy:        item.ToBuy,
				WhereToBuy:   stockItem.WhereToBuy,
				Notes:        stockItem.Notes,
			})
		}
	}

	return shoppingList, nil
}

func derefIntOrZero(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}
