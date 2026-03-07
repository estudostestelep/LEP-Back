package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceStaffAvailability struct {
	db *gorm.DB
}

type IStaffAvailabilityRepository interface {
	GetById(id uuid.UUID) (*models.StaffAvailability, error)
	GetByClientAndWeek(orgId, projectId, clientId uuid.UUID, weekStart time.Time) (*models.StaffAvailability, error)
	ListByWeek(orgId, projectId uuid.UUID, weekStart time.Time) ([]models.StaffAvailability, error)
	ListByClient(orgId, projectId, clientId uuid.UUID) ([]models.StaffAvailability, error)
	Create(availability *models.StaffAvailability) error
	Update(availability *models.StaffAvailability) error
	Upsert(availability *models.StaffAvailability) error
	Delete(id uuid.UUID) error
	GetWeekSummary(orgId, projectId uuid.UUID, weekStart time.Time) (*models.WeekAvailabilitySummary, error)
}

func NewStaffAvailabilityRepository(db *gorm.DB) IStaffAvailabilityRepository {
	return &resourceStaffAvailability{db: db}
}

func (r *resourceStaffAvailability) GetById(id uuid.UUID) (*models.StaffAvailability, error) {
	var availability models.StaffAvailability
	err := r.db.Preload("Client").First(&availability, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &availability, nil
}

func (r *resourceStaffAvailability) GetByClientAndWeek(orgId, projectId, clientId uuid.UUID, weekStart time.Time) (*models.StaffAvailability, error) {
	var availability models.StaffAvailability
	err := r.db.Where(
		"organization_id = ? AND project_id = ? AND client_id = ? AND week_start = ?",
		orgId, projectId, clientId, weekStart,
	).First(&availability).Error
	if err != nil {
		return nil, err
	}
	return &availability, nil
}

func (r *resourceStaffAvailability) ListByWeek(orgId, projectId uuid.UUID, weekStart time.Time) ([]models.StaffAvailability, error) {
	var availabilities []models.StaffAvailability
	err := r.db.
		Preload("Client").
		Where("organization_id = ? AND project_id = ? AND week_start = ?", orgId, projectId, weekStart).
		Find(&availabilities).Error
	return availabilities, err
}

func (r *resourceStaffAvailability) ListByClient(orgId, projectId, clientId uuid.UUID) ([]models.StaffAvailability, error) {
	var availabilities []models.StaffAvailability
	err := r.db.
		Where("organization_id = ? AND project_id = ? AND client_id = ?", orgId, projectId, clientId).
		Order("week_start DESC").
		Find(&availabilities).Error
	return availabilities, err
}

func (r *resourceStaffAvailability) Create(availability *models.StaffAvailability) error {
	if availability.Id == uuid.Nil {
		availability.Id = uuid.New()
	}
	return r.db.Create(availability).Error
}

func (r *resourceStaffAvailability) Update(availability *models.StaffAvailability) error {
	return r.db.Save(availability).Error
}

func (r *resourceStaffAvailability) Upsert(availability *models.StaffAvailability) error {
	existing, err := r.GetByClientAndWeek(
		availability.OrganizationId,
		availability.ProjectId,
		availability.ClientId,
		availability.WeekStart,
	)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existing != nil {
		existing.AvailableDays = availability.AvailableDays
		existing.NoAvailability = availability.NoAvailability
		existing.UpdatedAt = time.Now()
		return r.Update(existing)
	}

	return r.Create(availability)
}

func (r *resourceStaffAvailability) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.StaffAvailability{}, "id = ?", id).Error
}

func (r *resourceStaffAvailability) GetWeekSummary(orgId, projectId uuid.UUID, weekStart time.Time) (*models.WeekAvailabilitySummary, error) {
	availabilities, err := r.ListByWeek(orgId, projectId, weekStart)
	if err != nil {
		return nil, err
	}

	summary := &models.WeekAvailabilitySummary{
		WeekStart:      weekStart.Format("2006-01-02"),
		TotalResponses: len(availabilities),
		ByDay:          make(map[string][]models.AvailableClient),
	}

	for _, avail := range availabilities {
		if avail.NoAvailability {
			continue
		}
		clientName := ""
		if avail.Client != nil {
			clientName = avail.Client.Name
		}

		for _, day := range avail.AvailableDays {
			summary.ByDay[day] = append(summary.ByDay[day], models.AvailableClient{
				ClientId:   avail.ClientId,
				ClientName: clientName,
			})
		}
	}

	return summary, nil
}
