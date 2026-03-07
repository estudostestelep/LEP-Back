package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceStaffDashboard struct {
	db *gorm.DB
}

type IStaffDashboardRepository interface {
	// Sales Records
	CreateSalesRecord(record *models.StaffSalesRecord) error
	CreateSalesRecordsBatch(records []models.StaffSalesRecord) error
	ListSalesRecords(orgId, projectId uuid.UUID, filters models.DashboardFilters) ([]models.StaffSalesRecord, error)
	GetDashboardMeta(orgId, projectId uuid.UUID) (*models.DashboardMeta, error)
	GetDashboardRows(orgId, projectId uuid.UUID, filters models.DashboardFilters) ([]models.DashboardRow, error)
	GetDashboardGraphs(orgId, projectId uuid.UUID) (*models.DashboardGraphs, error)
	DeleteSalesRecordsByBatch(batchId uuid.UUID) error

	// Import Batches
	CreateImportBatch(batch *models.StaffSalesImportBatch) error
	UpdateImportBatch(batch *models.StaffSalesImportBatch) error
	GetImportBatchById(id uuid.UUID) (*models.StaffSalesImportBatch, error)
	ListImportBatches(orgId, projectId uuid.UUID, limit int) ([]models.StaffSalesImportBatch, error)

	// Staff Reports
	GetStaffReportMeta(orgId, projectId uuid.UUID) (*models.StaffReportMeta, error)
	GetStaffReportRows(orgId, projectId uuid.UUID, filters models.StaffReportFilters) ([]models.StaffReportRow, error)
}

func NewStaffDashboardRepository(db *gorm.DB) IStaffDashboardRepository {
	return &resourceStaffDashboard{db: db}
}

// ==================== Sales Records ====================

func (r *resourceStaffDashboard) CreateSalesRecord(record *models.StaffSalesRecord) error {
	if record.Id == uuid.Nil {
		record.Id = uuid.New()
	}
	return r.db.Create(record).Error
}

func (r *resourceStaffDashboard) CreateSalesRecordsBatch(records []models.StaffSalesRecord) error {
	for i := range records {
		if records[i].Id == uuid.Nil {
			records[i].Id = uuid.New()
		}
	}
	return r.db.CreateInBatches(&records, 500).Error
}

func (r *resourceStaffDashboard) ListSalesRecords(orgId, projectId uuid.UUID, filters models.DashboardFilters) ([]models.StaffSalesRecord, error) {
	var records []models.StaffSalesRecord
	query := r.db.Where("organization_id = ? AND project_id = ?", orgId, projectId)

	if filters.StartDate != nil {
		startDate, _ := time.Parse("2006-01-02", *filters.StartDate)
		query = query.Where("accounting_date >= ?", startDate)
	}
	if filters.EndDate != nil {
		endDate, _ := time.Parse("2006-01-02", *filters.EndDate)
		query = query.Where("accounting_date <= ?", endDate)
	}
	if len(filters.Groups) > 0 {
		query = query.Where("group_name IN ?", filters.Groups)
	}
	if len(filters.Descriptions) > 0 {
		query = query.Where("description IN ?", filters.Descriptions)
	}
	if len(filters.Employees) > 0 {
		query = query.Where("employee_name IN ?", filters.Employees)
	}

	err := query.Order("accounting_date DESC").Find(&records).Error
	return records, err
}

func (r *resourceStaffDashboard) GetDashboardMeta(orgId, projectId uuid.UUID) (*models.DashboardMeta, error) {
	meta := &models.DashboardMeta{}

	// Groups
	var groups []string
	r.db.Model(&models.StaffSalesRecord{}).
		Distinct("group_name").
		Where("organization_id = ? AND project_id = ?", orgId, projectId).
		Pluck("group_name", &groups)
	meta.Groups = groups

	// Descriptions
	var descriptions []string
	r.db.Model(&models.StaffSalesRecord{}).
		Distinct("description").
		Where("organization_id = ? AND project_id = ?", orgId, projectId).
		Order("description ASC").
		Pluck("description", &descriptions)
	meta.Descriptions = descriptions

	// Employees
	var employees []string
	r.db.Model(&models.StaffSalesRecord{}).
		Distinct("employee_name").
		Where("organization_id = ? AND project_id = ?", orgId, projectId).
		Order("employee_name ASC").
		Pluck("employee_name", &employees)
	meta.Employees = employees

	// Sale Modes
	var saleModes []string
	r.db.Model(&models.StaffSalesRecord{}).
		Distinct("sale_mode").
		Where("organization_id = ? AND project_id = ?", orgId, projectId).
		Pluck("sale_mode", &saleModes)
	meta.SaleModes = saleModes

	// Date range
	var minDate, maxDate *time.Time
	r.db.Model(&models.StaffSalesRecord{}).
		Select("MIN(accounting_date)").
		Where("organization_id = ? AND project_id = ?", orgId, projectId).
		Scan(&minDate)
	r.db.Model(&models.StaffSalesRecord{}).
		Select("MAX(accounting_date)").
		Where("organization_id = ? AND project_id = ?", orgId, projectId).
		Scan(&maxDate)
	meta.MinDate = minDate
	meta.MaxDate = maxDate

	return meta, nil
}

func (r *resourceStaffDashboard) GetDashboardRows(orgId, projectId uuid.UUID, filters models.DashboardFilters) ([]models.DashboardRow, error) {
	var rows []models.DashboardRow

	query := r.db.Model(&models.StaffSalesRecord{}).
		Select(`
			TO_CHAR(accounting_date, 'DD/MM/YYYY') as accounting_date,
			TO_CHAR(accounting_date, 'Day') as weekday,
			group_name,
			description,
			employee_name,
			SUM(qtd) as quantity,
			SUM(vl_servico_info) as service_fee_informed,
			SUM(vl_servico_calc) as service_fee_calculated,
			SUM(vl_total) as total_value
		`).
		Where("organization_id = ? AND project_id = ?", orgId, projectId)

	if filters.StartDate != nil {
		startDate, _ := time.Parse("2006-01-02", *filters.StartDate)
		query = query.Where("accounting_date >= ?", startDate)
	}
	if filters.EndDate != nil {
		endDate, _ := time.Parse("2006-01-02", *filters.EndDate)
		query = query.Where("accounting_date <= ?", endDate)
	}
	if len(filters.Groups) > 0 {
		query = query.Where("group_name IN ?", filters.Groups)
	}
	if len(filters.Descriptions) > 0 {
		query = query.Where("description IN ?", filters.Descriptions)
	}

	err := query.
		Group("accounting_date, group_name, description, employee_name").
		Order("accounting_date DESC, group_name ASC").
		Scan(&rows).Error

	return rows, err
}

func (r *resourceStaffDashboard) GetDashboardGraphs(orgId, projectId uuid.UUID) (*models.DashboardGraphs, error) {
	graphs := &models.DashboardGraphs{
		ByEmployee: []models.EmployeeGraphPoint{},
		ByGroup:    []models.GroupGraphPoint{},
		ByDate:     []models.DateGraphPoint{},
		ByWeekday:  []models.WeekdayGraphPoint{},
	}

	// Por Funcionário (top 10)
	var employeeData []struct {
		Employee string
		Total    float64
	}
	r.db.Model(&models.StaffSalesRecord{}).
		Select("employee_name as employee, SUM(vl_total) as total").
		Where("organization_id = ? AND project_id = ? AND employee_name IS NOT NULL AND employee_name != ''",
			orgId, projectId).
		Group("employee_name").
		Order("total DESC").
		Limit(10).
		Scan(&employeeData)

	for _, e := range employeeData {
		graphs.ByEmployee = append(graphs.ByEmployee, models.EmployeeGraphPoint{
			Employee: e.Employee,
			Total:    e.Total,
		})
	}

	// Por Grupo (top 10)
	var groupData []struct {
		Group string
		Total float64
	}
	r.db.Model(&models.StaffSalesRecord{}).
		Select("group_name as group, SUM(vl_total) as total").
		Where("organization_id = ? AND project_id = ? AND group_name IS NOT NULL AND group_name != ''",
			orgId, projectId).
		Group("group_name").
		Order("total DESC").
		Limit(10).
		Scan(&groupData)

	for _, g := range groupData {
		graphs.ByGroup = append(graphs.ByGroup, models.GroupGraphPoint{
			Group: g.Group,
			Total: g.Total,
		})
	}

	// Últimos 30 dias
	var dateData []struct {
		Date  string
		Total float64
	}
	r.db.Model(&models.StaffSalesRecord{}).
		Select("TO_CHAR(accounting_date, 'DD/MM') as date, SUM(vl_total) as total").
		Where("organization_id = ? AND project_id = ? AND accounting_date >= ?",
			orgId, projectId, time.Now().AddDate(0, 0, -30)).
		Group("accounting_date").
		Order("accounting_date ASC").
		Scan(&dateData)

	for _, d := range dateData {
		graphs.ByDate = append(graphs.ByDate, models.DateGraphPoint{
			Date:  d.Date,
			Total: d.Total,
		})
	}

	// Por Dia da Semana
	weekdayNames := []string{"Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado"}
	var weekdayData []struct {
		Dow   int
		Total float64
	}
	r.db.Model(&models.StaffSalesRecord{}).
		Select("EXTRACT(DOW FROM accounting_date)::int as dow, SUM(vl_total) as total").
		Where("organization_id = ? AND project_id = ?", orgId, projectId).
		Group("EXTRACT(DOW FROM accounting_date)").
		Order("dow ASC").
		Scan(&weekdayData)

	for _, wd := range weekdayData {
		weekday := "Dia " + string(rune('0'+wd.Dow))
		if wd.Dow >= 0 && wd.Dow < len(weekdayNames) {
			weekday = weekdayNames[wd.Dow]
		}
		graphs.ByWeekday = append(graphs.ByWeekday, models.WeekdayGraphPoint{
			Weekday: weekday,
			Total:   wd.Total,
		})
	}

	return graphs, nil
}

func (r *resourceStaffDashboard) DeleteSalesRecordsByBatch(batchId uuid.UUID) error {
	return r.db.Delete(&models.StaffSalesRecord{}, "import_batch_id = ?", batchId).Error
}

// ==================== Import Batches ====================

func (r *resourceStaffDashboard) CreateImportBatch(batch *models.StaffSalesImportBatch) error {
	if batch.Id == uuid.Nil {
		batch.Id = uuid.New()
	}
	return r.db.Create(batch).Error
}

func (r *resourceStaffDashboard) UpdateImportBatch(batch *models.StaffSalesImportBatch) error {
	return r.db.Save(batch).Error
}

func (r *resourceStaffDashboard) GetImportBatchById(id uuid.UUID) (*models.StaffSalesImportBatch, error) {
	var batch models.StaffSalesImportBatch
	err := r.db.First(&batch, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *resourceStaffDashboard) ListImportBatches(orgId, projectId uuid.UUID, limit int) ([]models.StaffSalesImportBatch, error) {
	var batches []models.StaffSalesImportBatch
	query := r.db.
		Where("organization_id = ? AND project_id = ?", orgId, projectId).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&batches).Error
	return batches, err
}

// ==================== Staff Reports ====================

func (r *resourceStaffDashboard) GetStaffReportMeta(orgId, projectId uuid.UUID) (*models.StaffReportMeta, error) {
	meta := &models.StaffReportMeta{}

	// Years
	var years []int
	r.db.Model(&models.StaffAttendance{}).
		Distinct("EXTRACT(YEAR FROM attendance_date)::int").
		Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", orgId, projectId).
		Order("EXTRACT(YEAR FROM attendance_date)::int DESC").
		Pluck("EXTRACT(YEAR FROM attendance_date)::int", &years)
	meta.Years = years

	// Employees (from clients with attendance)
	var employees []string
	r.db.Model(&models.StaffAttendance{}).
		Distinct("clients.name").
		Joins("JOIN clients ON clients.id = staff_attendances.client_id").
		Where("staff_attendances.organization_id = ? AND staff_attendances.project_id = ? AND staff_attendances.deleted_at IS NULL",
			orgId, projectId).
		Order("clients.name ASC").
		Pluck("clients.name", &employees)
	meta.Employees = employees

	// Sectors (placeholder - could come from role or custom field)
	meta.Sectors = []string{}

	return meta, nil
}

func (r *resourceStaffDashboard) GetStaffReportRows(orgId, projectId uuid.UUID, filters models.StaffReportFilters) ([]models.StaffReportRow, error) {
	var rows []models.StaffReportRow

	query := r.db.Model(&models.StaffAttendance{}).
		Select(`
			TO_CHAR(staff_attendances.attendance_date, 'DD/MM/YYYY') as date,
			TO_CHAR(staff_attendances.attendance_date, 'Day') as weekday,
			staff_attendances.shift,
			clients.name as client_name,
			'' as sector,
			0 as daily_rate,
			COALESCE(
				(SELECT SUM(staff_consumption_products.unit_cost * staff_attendance_consumptions.quantity)
				 FROM staff_attendance_consumptions
				 JOIN staff_consumption_products ON staff_consumption_products.id = staff_attendance_consumptions.product_id
				 WHERE staff_attendance_consumptions.attendance_id = staff_attendances.id), 0
			) as consumption,
			0 as commission,
			0 as total
		`).
		Joins("JOIN clients ON clients.id = staff_attendances.client_id").
		Where("staff_attendances.organization_id = ? AND staff_attendances.project_id = ? AND staff_attendances.deleted_at IS NULL",
			orgId, projectId).
		Where("EXTRACT(YEAR FROM staff_attendances.attendance_date) = ?", filters.Year)

	if filters.StartDate != nil {
		startDate, _ := time.Parse("2006-01-02", *filters.StartDate)
		query = query.Where("staff_attendances.attendance_date >= ?", startDate)
	}
	if filters.EndDate != nil {
		endDate, _ := time.Parse("2006-01-02", *filters.EndDate)
		query = query.Where("staff_attendances.attendance_date <= ?", endDate)
	}
	if len(filters.Employees) > 0 {
		query = query.Where("clients.name IN ?", filters.Employees)
	}

	err := query.Order("staff_attendances.attendance_date DESC").Scan(&rows).Error
	return rows, err
}
