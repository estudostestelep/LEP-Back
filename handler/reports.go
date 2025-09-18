package handler

import (
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type ReportsHandler struct {
	repo *repositories.DBconn
}

type IReportsHandler interface {
	GetOccupancyReport(orgId, projectId string, startDate, endDate time.Time) (*OccupancyReport, error)
	GetReservationReport(orgId, projectId string, startDate, endDate time.Time) (*ReservationReport, error)
	GetWaitlistReport(orgId, projectId string, startDate, endDate time.Time) (*WaitlistReport, error)
	GetLeadReport(orgId, projectId string, startDate, endDate time.Time) (*LeadReport, error)
	ExportReportToCSV(reportType, orgId, projectId string, startDate, endDate time.Time) ([]byte, error)
	GenerateDailyMetrics(orgId, projectId uuid.UUID, date time.Time) error
}

// Estruturas de relatórios
type OccupancyReport struct {
	Period       string                 `json:"period"`
	TotalTables  int                    `json:"total_tables"`
	DailyMetrics []DailyOccupancyMetric `json:"daily_metrics"`
	Summary      OccupancySummary       `json:"summary"`
}

type DailyOccupancyMetric struct {
	Date              string  `json:"date"`
	TablesOccupied    int     `json:"tables_occupied"`
	OccupancyRate     float64 `json:"occupancy_rate"`
	TotalReservations int     `json:"total_reservations"`
}

type OccupancySummary struct {
	AverageOccupancy float64 `json:"average_occupancy"`
	PeakOccupancy    float64 `json:"peak_occupancy"`
	BestDay          string  `json:"best_day"`
	WorstDay         string  `json:"worst_day"`
}

type ReservationReport struct {
	Period            string                   `json:"period"`
	TotalReservations int                      `json:"total_reservations"`
	DailyMetrics      []DailyReservationMetric `json:"daily_metrics"`
	StatusBreakdown   map[string]int           `json:"status_breakdown"`
	Summary           ReservationSummary       `json:"summary"`
}

type DailyReservationMetric struct {
	Date          string `json:"date"`
	Reservations  int    `json:"reservations"`
	Cancellations int    `json:"cancellations"`
	NoShows       int    `json:"no_shows"`
}

type ReservationSummary struct {
	AverageDaily     float64 `json:"average_daily"`
	CancellationRate float64 `json:"cancellation_rate"`
	NoShowRate       float64 `json:"no_show_rate"`
}

type WaitlistReport struct {
	Period         string                `json:"period"`
	TotalWaitlist  int                   `json:"total_waitlist"`
	DailyMetrics   []DailyWaitlistMetric `json:"daily_metrics"`
	ConversionRate float64               `json:"conversion_rate"`
}

type DailyWaitlistMetric struct {
	Date        string  `json:"date"`
	Added       int     `json:"added"`
	Seated      int     `json:"seated"`
	Left        int     `json:"left"`
	AvgWaitTime float64 `json:"avg_wait_time"`
}

type LeadReport struct {
	Period          string         `json:"period"`
	TotalLeads      int            `json:"total_leads"`
	SourceBreakdown map[string]int `json:"source_breakdown"`
	StatusBreakdown map[string]int `json:"status_breakdown"`
	ConversionRate  float64        `json:"conversion_rate"`
}

func NewReportsHandler(repo *repositories.DBconn) IReportsHandler {
	return &ReportsHandler{repo: repo}
}

// GetOccupancyReport - Relatório de ocupação de mesas
func (r *ReportsHandler) GetOccupancyReport(orgId, projectId string, startDate, endDate time.Time) (*OccupancyReport, error) {
	orgUUID, _ := uuid.Parse(orgId)
	projectUUID, _ := uuid.Parse(projectId)

	// Buscar todas as mesas do projeto
	tables, err := r.repo.Tables.GetTablesByProject(orgUUID, projectUUID)
	if err != nil {
		return nil, err
	}
	totalTables := len(tables)

	// Buscar reservas no período
	reservations, err := r.repo.Reservations.GetReservationsByProject(orgUUID, projectUUID)
	if err != nil {
		return nil, err
	}

	// Filtrar reservas no período
	var periodReservations []models.Reservation
	for _, res := range reservations {
		if res.Datetime.After(startDate) && res.Datetime.Before(endDate) && res.Status == "confirmed" {
			periodReservations = append(periodReservations, res)
		}
	}

	// Calcular métricas diárias
	dailyMetrics := []DailyOccupancyMetric{}
	var totalOccupancy float64
	bestDay := ""
	worstDay := ""
	var bestRate, worstRate float64 = 0, 100

	// Iterar por cada dia do período
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dayStart := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
		dayEnd := dayStart.Add(24 * time.Hour)

		reservationsThisDay := 0
		for _, res := range periodReservations {
			if res.Datetime.After(dayStart) && res.Datetime.Before(dayEnd) {
				reservationsThisDay++
			}
		}

		// Assumir que cada reserva ocupa a mesa por 2h, em um dia de 12h úteis = 6 turnos
		maxPossibleReservations := totalTables * 6
		occupancyRate := float64(reservationsThisDay) / float64(maxPossibleReservations) * 100

		dailyMetrics = append(dailyMetrics, DailyOccupancyMetric{
			Date:              d.Format("2006-01-02"),
			TablesOccupied:    reservationsThisDay,
			OccupancyRate:     occupancyRate,
			TotalReservations: reservationsThisDay,
		})

		totalOccupancy += occupancyRate

		// Rastrear melhor e pior dia
		if occupancyRate > bestRate {
			bestRate = occupancyRate
			bestDay = d.Format("2006-01-02")
		}
		if occupancyRate < worstRate {
			worstRate = occupancyRate
			worstDay = d.Format("2006-01-02")
		}
	}

	avgOccupancy := totalOccupancy / float64(len(dailyMetrics))

	return &OccupancyReport{
		Period:       fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		TotalTables:  totalTables,
		DailyMetrics: dailyMetrics,
		Summary: OccupancySummary{
			AverageOccupancy: avgOccupancy,
			PeakOccupancy:    bestRate,
			BestDay:          bestDay,
			WorstDay:         worstDay,
		},
	}, nil
}

// GetReservationReport - Relatório de reservas
func (r *ReportsHandler) GetReservationReport(orgId, projectId string, startDate, endDate time.Time) (*ReservationReport, error) {
	orgUUID, _ := uuid.Parse(orgId)
	projectUUID, _ := uuid.Parse(projectId)

	// Buscar reservas no período
	reservations, err := r.repo.Reservations.GetReservationsByProject(orgUUID, projectUUID)
	if err != nil {
		return nil, err
	}

	// Filtrar e contar por status
	statusBreakdown := make(map[string]int)
	var periodReservations []models.Reservation

	for _, res := range reservations {
		if res.Datetime.After(startDate) && res.Datetime.Before(endDate) {
			periodReservations = append(periodReservations, res)
			statusBreakdown[res.Status]++
		}
	}

	// Calcular métricas
	totalReservations := len(periodReservations)
	cancelled := statusBreakdown["cancelled"]
	// noShows seria um status específico - por agora assumir 0
	noShows := 0

	cancellationRate := 0.0
	if totalReservations > 0 {
		cancellationRate = float64(cancelled) / float64(totalReservations) * 100
	}

	// Calcular média diária
	days := int(endDate.Sub(startDate).Hours() / 24)
	if days == 0 {
		days = 1
	}
	averageDaily := float64(totalReservations) / float64(days)

	return &ReservationReport{
		Period:            fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		TotalReservations: totalReservations,
		DailyMetrics:      []DailyReservationMetric{}, // Seria calculado como no exemplo anterior
		StatusBreakdown:   statusBreakdown,
		Summary: ReservationSummary{
			AverageDaily:     averageDaily,
			CancellationRate: cancellationRate,
			NoShowRate:       float64(noShows) / float64(totalReservations) * 100,
		},
	}, nil
}

// GetWaitlistReport - Relatório de fila de espera
func (r *ReportsHandler) GetWaitlistReport(orgId, projectId string, startDate, endDate time.Time) (*WaitlistReport, error) {
	orgUUID, _ := uuid.Parse(orgId)
	projectUUID, _ := uuid.Parse(projectId)

	// Buscar waitlist no período
	waitlist, err := r.repo.Waitlists.GetWaitlistByProject(orgUUID, projectUUID)
	if err != nil {
		return nil, err
	}

	// Filtrar por período
	var periodWaitlist []models.Waitlist
	seated := 0

	for _, wait := range waitlist {
		if wait.CreatedAt.After(startDate) && wait.CreatedAt.Before(endDate) {
			periodWaitlist = append(periodWaitlist, wait)
			if wait.Status == "seated" {
				seated++
			}
		}
	}

	totalWaitlist := len(periodWaitlist)
	conversionRate := 0.0
	if totalWaitlist > 0 {
		conversionRate = float64(seated) / float64(totalWaitlist) * 100
	}

	return &WaitlistReport{
		Period:         fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		TotalWaitlist:  totalWaitlist,
		DailyMetrics:   []DailyWaitlistMetric{}, // Seria calculado detalhadamente
		ConversionRate: conversionRate,
	}, nil
}

// GetLeadReport - Relatório de leads (CRM básico)
func (r *ReportsHandler) GetLeadReport(orgId, projectId string, startDate, endDate time.Time) (*LeadReport, error) {
	// Implementação simplificada - assumindo que o repository de leads existirá

	// Por agora, relatório mock
	return &LeadReport{
		Period:     fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		TotalLeads: 0,
		SourceBreakdown: map[string]int{
			"waitlist":    0,
			"reservation": 0,
			"walk_in":     0,
		},
		StatusBreakdown: map[string]int{
			"new":       0,
			"contacted": 0,
			"converted": 0,
			"lost":      0,
		},
		ConversionRate: 0.0,
	}, nil
}

// ExportReportToCSV - Exporta relatório para CSV
func (r *ReportsHandler) ExportReportToCSV(reportType, orgId, projectId string, startDate, endDate time.Time) ([]byte, error) {
	// Implementação básica - na prática usaria uma biblioteca CSV

	csvContent := "Type,Date,Value\n"

	switch reportType {
	case "occupancy":
		report, err := r.GetOccupancyReport(orgId, projectId, startDate, endDate)
		if err != nil {
			return nil, err
		}

		for _, metric := range report.DailyMetrics {
			csvContent += fmt.Sprintf("occupancy,%s,%.2f\n", metric.Date, metric.OccupancyRate)
		}

	case "reservations":
		report, err := r.GetReservationReport(orgId, projectId, startDate, endDate)
		if err != nil {
			return nil, err
		}

		csvContent += fmt.Sprintf("reservations,%s,%d\n", report.Period, report.TotalReservations)

	default:
		return nil, fmt.Errorf("unsupported report type: %s", reportType)
	}

	return []byte(csvContent), nil
}

// GenerateDailyMetrics - Gera métricas diárias para armazenamento
func (r *ReportsHandler) GenerateDailyMetrics(orgId, projectId uuid.UUID, date time.Time) error {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	// Contar reservas do dia
	reservations, err := r.repo.Reservations.GetReservationsByProject(orgId, projectId)
	if err != nil {
		return err
	}

	reservationCount := 0
	for _, res := range reservations {
		if res.Datetime.After(dayStart) && res.Datetime.Before(dayEnd) && res.Status == "confirmed" {
			reservationCount++
		}
	}

	// Criar métrica de reservas
	reservationMetric := &models.ReportMetric{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		MetricType:     "reservation_count",
		MetricDate:     dayStart,
		Value:          float64(reservationCount),
		CreatedAt:      time.Now(),
	}

	// Salvar métrica - assumindo que o repository existirá
	// return r.repo.ReportMetrics.CreateMetric(reservationMetric)

	// Por agora apenas log
	fmt.Printf("Would save metric: %+v\n", reservationMetric)

	return nil
}
