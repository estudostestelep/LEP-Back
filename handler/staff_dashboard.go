package handler

import (
	"encoding/base64"
	"encoding/csv"
	"lep/repositories"
	"lep/repositories/models"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type resourceStaffDashboard struct {
	repo *repositories.DBconn
}

type IHandlerStaffDashboard interface {
	GetDashboardMeta(orgId, projectId string) (*models.DashboardMeta, error)
	GetDashboardRows(orgId, projectId string, filters models.DashboardFilters) ([]models.DashboardRow, error)
	GetDashboardGraphs(orgId, projectId string) (*models.DashboardGraphs, error)
	ImportCSV(req *models.ImportCSVRequest, orgId, projectId, createdById string) (*models.ImportCSVResponse, error)
	ListImportBatches(orgId, projectId string, limit int) ([]models.StaffSalesImportBatch, error)
	GetStaffReportMeta(orgId, projectId string) (*models.StaffReportMeta, error)
	GetStaffReportRows(orgId, projectId string, filters models.StaffReportFilters) ([]models.StaffReportRow, error)
}

func NewStaffDashboardHandler(repo *repositories.DBconn) IHandlerStaffDashboard {
	return &resourceStaffDashboard{repo: repo}
}

func (r *resourceStaffDashboard) GetDashboardMeta(orgId, projectId string) (*models.DashboardMeta, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffDashboard.GetDashboardMeta(orgUUID, projUUID)
}

func (r *resourceStaffDashboard) GetDashboardRows(orgId, projectId string, filters models.DashboardFilters) ([]models.DashboardRow, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffDashboard.GetDashboardRows(orgUUID, projUUID, filters)
}

func (r *resourceStaffDashboard) GetDashboardGraphs(orgId, projectId string) (*models.DashboardGraphs, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffDashboard.GetDashboardGraphs(orgUUID, projUUID)
}

func (r *resourceStaffDashboard) ImportCSV(req *models.ImportCSVRequest, orgId, projectId, createdById string) (*models.ImportCSVResponse, error) {
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

	// Create import batch
	batch := &models.StaffSalesImportBatch{
		OrganizationId: orgUUID,
		ProjectId:      projUUID,
		FileName:       req.FileName,
		Status:         "processing",
		CreatedById:    createdByUUID,
		CreatedAt:      time.Now(),
	}
	if err := r.repo.StaffDashboard.CreateImportBatch(batch); err != nil {
		return nil, err
	}

	// Decode base64 content
	csvContent, err := base64.StdEncoding.DecodeString(req.FileContent)
	if err != nil {
		batch.Status = "failed"
		errMsg := "invalid base64 content"
		batch.ErrorMessage = &errMsg
		r.repo.StaffDashboard.UpdateImportBatch(batch)
		return &models.ImportCSVResponse{
			BatchId: batch.Id,
			Status:  "failed",
			Errors:  []string{errMsg},
		}, nil
	}

	// Parse CSV using streaming (line by line) to handle large files
	reader := csv.NewReader(strings.NewReader(string(csvContent)))

	// Read header first
	header, err := reader.Read()
	if err != nil {
		batch.Status = "failed"
		errMsg := "failed to read CSV header"
		batch.ErrorMessage = &errMsg
		r.repo.StaffDashboard.UpdateImportBatch(batch)
		return &models.ImportCSVResponse{
			BatchId: batch.Id,
			Status:  "failed",
			Errors:  []string{errMsg},
		}, nil
	}

	// Parse header to find column indices
	colIdx := make(map[string]int)
	for i, col := range header {
		colIdx[strings.ToLower(strings.TrimSpace(col))] = i
	}

	var errors []string
	totalRecords := 0
	lineNum := 1 // Start after header

	const batchSize = 500
	salesRecords := make([]models.StaffSalesRecord, 0, batchSize)

	// Process CSV line by line (streaming)
	for {
		row, err := reader.Read()
		if err != nil {
			// End of file or error
			break
		}
		lineNum++

		record := models.StaffSalesRecord{
			OrganizationId: orgUUID,
			ProjectId:      projUUID,
			ImportBatchId:  &batch.Id,
			ImportedAt:     time.Now(),
			CreatedAt:      time.Now(),
		}

		// Parse fields based on CSV structure
		if idx, ok := colIdx["qtd"]; ok && idx < len(row) {
			qty, _ := strconv.Atoi(row[idx])
			record.Quantity = qty
		}
		if idx, ok := colIdx["descricao"]; ok && idx < len(row) {
			record.Description = row[idx]
		}
		if idx, ok := colIdx["grupo"]; ok && idx < len(row) {
			record.GroupName = row[idx]
		}
		if idx, ok := colIdx["nome_funcionario"]; ok && idx < len(row) {
			record.EmployeeName = row[idx]
		}
		if idx, ok := colIdx["nome_modo_venda"]; ok && idx < len(row) {
			record.SaleMode = row[idx]
		}

		// Parse monetary values (Brazilian format: 1.234,56)
		if idx, ok := colIdx["vl_total"]; ok && idx < len(row) {
			record.TotalValue = parseBrazilianFloat(row[idx])
		}
		if idx, ok := colIdx["vl_servico_calculado"]; ok && idx < len(row) {
			record.ServiceFeeCalculated = parseBrazilianFloat(row[idx])
		}
		if idx, ok := colIdx["vl_servico_informado"]; ok && idx < len(row) {
			record.ServiceFeeInformed = parseBrazilianFloat(row[idx])
		}
		if idx, ok := colIdx["vl_preco"]; ok && idx < len(row) {
			record.Price = parseBrazilianFloat(row[idx])
		}
		if idx, ok := colIdx["vl_preco_original"]; ok && idx < len(row) {
			record.OriginalPrice = parseBrazilianFloat(row[idx])
		}
		if idx, ok := colIdx["vl_desconto"]; ok && idx < len(row) {
			record.Discount = parseBrazilianFloat(row[idx])
		}

		// Parse dates
		if idx, ok := colIdx["dt_contabil"]; ok && idx < len(row) {
			date, parseErr := parseBrazilianDate(row[idx])
			if parseErr == nil {
				record.AccountingDate = date
			} else {
				// Only add error if there are few errors (limit to first 100)
				if len(errors) < 100 {
					errors = append(errors, "Linha "+strconv.Itoa(lineNum)+": dt_contabil invalido")
				}
				continue
			}
		}
		if idx, ok := colIdx["dt_hr_lancamento"]; ok && idx < len(row) {
			dateTime, _ := parseBrazilianDateTime(row[idx])
			record.LaunchDateTime = dateTime
		}

		if record.Description == "" || record.GroupName == "" {
			if len(errors) < 100 {
				errors = append(errors, "Linha "+strconv.Itoa(lineNum)+": descricao ou grupo ausente")
			}
			continue
		}

		salesRecords = append(salesRecords, record)

		// Insert in batches of 500 records
		if len(salesRecords) >= batchSize {
			if insertErr := r.repo.StaffDashboard.CreateSalesRecordsBatch(salesRecords); insertErr != nil {
				batch.Status = "failed"
				errMsg := "falha ao inserir registros: " + insertErr.Error()
				batch.ErrorMessage = &errMsg
				batch.RecordsCount = totalRecords
				r.repo.StaffDashboard.UpdateImportBatch(batch)
				return &models.ImportCSVResponse{
					BatchId:      batch.Id,
					RecordsCount: totalRecords,
					Status:       "failed",
					Errors:       append(errors, errMsg),
				}, nil
			}
			totalRecords += len(salesRecords)
			salesRecords = salesRecords[:0] // Reset slice keeping capacity
		}
	}

	// Insert remaining records
	if len(salesRecords) > 0 {
		if insertErr := r.repo.StaffDashboard.CreateSalesRecordsBatch(salesRecords); insertErr != nil {
			batch.Status = "failed"
			errMsg := "falha ao inserir registros finais: " + insertErr.Error()
			batch.ErrorMessage = &errMsg
			batch.RecordsCount = totalRecords
			r.repo.StaffDashboard.UpdateImportBatch(batch)
			return &models.ImportCSVResponse{
				BatchId:      batch.Id,
				RecordsCount: totalRecords,
				Status:       "failed",
				Errors:       append(errors, errMsg),
			}, nil
		}
		totalRecords += len(salesRecords)
	}

	// Check if any records were imported
	if totalRecords == 0 && len(errors) > 0 {
		batch.Status = "failed"
		errMsg := "nenhum registro valido encontrado"
		batch.ErrorMessage = &errMsg
		r.repo.StaffDashboard.UpdateImportBatch(batch)
		return &models.ImportCSVResponse{
			BatchId:      batch.Id,
			RecordsCount: 0,
			Status:       "failed",
			Errors:       errors,
		}, nil
	}

	// Update batch status
	now := time.Now()
	batch.Status = "completed"
	batch.RecordsCount = totalRecords
	batch.CompletedAt = &now
	r.repo.StaffDashboard.UpdateImportBatch(batch)

	return &models.ImportCSVResponse{
		BatchId:      batch.Id,
		RecordsCount: totalRecords,
		Status:       "completed",
		Errors:       errors,
	}, nil
}

func (r *resourceStaffDashboard) ListImportBatches(orgId, projectId string, limit int) ([]models.StaffSalesImportBatch, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffDashboard.ListImportBatches(orgUUID, projUUID, limit)
}

func (r *resourceStaffDashboard) GetStaffReportMeta(orgId, projectId string) (*models.StaffReportMeta, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffDashboard.GetStaffReportMeta(orgUUID, projUUID)
}

func (r *resourceStaffDashboard) GetStaffReportRows(orgId, projectId string, filters models.StaffReportFilters) ([]models.StaffReportRow, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffDashboard.GetStaffReportRows(orgUUID, projUUID, filters)
}

// Helper functions for parsing Brazilian formats

func parseBrazilianFloat(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ".", "")  // Remove thousand separator
	s = strings.ReplaceAll(s, ",", ".") // Replace decimal separator
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func parseBrazilianDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	// Try DD/MM/YYYY
	if t, err := time.Parse("02/01/2006", s); err == nil {
		return t, nil
	}
	// Try YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Time{}, nil
}

func parseBrazilianDateTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	// Try DD/MM/YYYY HH:MM:SS
	if t, err := time.Parse("02/01/2006 15:04:05", s); err == nil {
		return t, nil
	}
	// Try DD/MM/YYYY HH:MM
	if t, err := time.Parse("02/01/2006 15:04", s); err == nil {
		return t, nil
	}
	return time.Time{}, nil
}
