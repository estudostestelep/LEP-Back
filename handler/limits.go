package handler

import (
	"fmt"
	"lep/repositories"
	"time"

	"github.com/google/uuid"
)

// LimitType define os tipos de limites verificáveis
type LimitType string

const (
	LimitUsers             LimitType = "users"
	LimitTables            LimitType = "tables"
	LimitProducts          LimitType = "products"
	LimitReservationsDay   LimitType = "reservations_per_day"
	LimitAuditLogs         LimitType = "audit_logs_limit"
	LimitAuditLogsRetention LimitType = "audit_logs_retention"
)

// UsageData representa o uso atual de recursos
type UsageData struct {
	TablesCount       int `json:"tables_count"`
	UsersCount        int `json:"users_count"`
	ProductsCount     int `json:"products_count"`
	ReservationsToday int `json:"reservations_today"`
}

// LimitsData representa os limites do plano
type LimitsData struct {
	MaxTables              int  `json:"max_tables"`
	MaxUsers               int  `json:"max_users"`
	MaxProducts            int  `json:"max_products"`
	MaxReservationsPerDay  int  `json:"max_reservations_per_day"`
	MaxAuditLogs           int  `json:"max_audit_logs"`            // -1 = ilimitado, 0 = desabilitado
	AuditLogsRetentionDays int  `json:"audit_logs_retention_days"` // Dias de retenção dos logs
	NotificationsEnabled   bool `json:"notifications_enabled"`
	ReportsEnabled         bool `json:"reports_enabled"`
	ReservationsEnabled    bool `json:"reservations_enabled"`
	WaitlistEnabled        bool `json:"waitlist_enabled"`
	AuditLogsEnabled       bool `json:"audit_logs_enabled"`
}

// UsageLimitsResponse representa a resposta completa de uso e limites
type UsageLimitsResponse struct {
	PackageCode string     `json:"package_code"`
	PackageName string     `json:"package_name"`
	Usage       UsageData  `json:"usage"`
	Limits      LimitsData `json:"limits"`
}

// ILimitHandler interface para verificação de limites
type ILimitHandler interface {
	CheckLimit(orgId, projectId string, limitType LimitType) (canCreate bool, current, limit int, err error)
	GetUsageAndLimits(orgId, projectId string) (*UsageLimitsResponse, error)
	HasModule(orgId, moduleCode string) (bool, error)
}

// LimitHandler implementa ILimitHandler
type LimitHandler struct {
	packageRepo     repositories.IPackageRepository
	tableRepo       repositories.ITableRepository
	userOrgRepo     repositories.IUserOrganizationRepository
	productRepo     repositories.IProductRepository
	reservationRepo repositories.IReservationRepository
	moduleRepo      repositories.IModuleRepository
}

// NewLimitHandler cria uma nova instância de LimitHandler
func NewLimitHandler(
	packageRepo repositories.IPackageRepository,
	tableRepo repositories.ITableRepository,
	userOrgRepo repositories.IUserOrganizationRepository,
	productRepo repositories.IProductRepository,
	reservationRepo repositories.IReservationRepository,
	moduleRepo repositories.IModuleRepository,
) *LimitHandler {
	return &LimitHandler{
		packageRepo:     packageRepo,
		tableRepo:       tableRepo,
		userOrgRepo:     userOrgRepo,
		productRepo:     productRepo,
		reservationRepo: reservationRepo,
		moduleRepo:      moduleRepo,
	}
}

// CheckLimit verifica se a organização pode criar mais um recurso do tipo especificado
// Retorna: canCreate (pode criar), current (quantidade atual), limit (limite máximo), error
func (h *LimitHandler) CheckLimit(orgId, projectId string, limitType LimitType) (bool, int, int, error) {
	// 1. Buscar assinatura da organização
	orgPackage, err := h.packageRepo.GetOrganizationPackage(orgId)
	if err != nil {
		// Sem assinatura = sem plano, não pode criar
		return false, 0, 0, fmt.Errorf("organização sem plano ativo")
	}

	// 2. Buscar limites do pacote
	limits, err := h.packageRepo.GetPackageLimits(orgPackage.PackageId.String())
	if err != nil {
		// Sem limites definidos = ilimitado
		return true, 0, -1, nil
	}

	// 3. Encontrar o limite específico
	var maxLimit int = -1 // -1 = ilimitado
	for _, l := range limits {
		if l.LimitType == string(limitType) {
			maxLimit = l.LimitValue
			break
		}
	}

	// Se -1 = ilimitado, permitir
	if maxLimit == -1 {
		return true, 0, -1, nil
	}

	// Se 0 = funcionalidade desabilitada
	if maxLimit == 0 {
		return false, 0, 0, nil
	}

	// 4. Contar recursos atuais
	currentCount, err := h.countResource(orgId, projectId, limitType)
	if err != nil {
		return false, 0, maxLimit, err
	}

	return currentCount < maxLimit, currentCount, maxLimit, nil
}

// countResource conta os recursos atuais por tipo
func (h *LimitHandler) countResource(orgId, projectId string, limitType LimitType) (int, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return 0, fmt.Errorf("ID de organização inválido: %w", err)
	}

	projUUID, _ := uuid.Parse(projectId) // projectId pode ser vazio para alguns tipos

	switch limitType {
	case LimitTables:
		tables, err := h.tableRepo.ListTables(orgUUID, projUUID)
		if err != nil {
			return 0, err
		}
		return len(tables), nil

	case LimitUsers:
		users, err := h.userOrgRepo.ListByOrganization(orgId)
		if err != nil {
			return 0, err
		}
		return len(users), nil

	case LimitProducts:
		products, err := h.productRepo.ListProducts(orgUUID, projUUID)
		if err != nil {
			return 0, err
		}
		return len(products), nil

	case LimitReservationsDay:
		// Contar reservas do dia atual
		reservations, err := h.reservationRepo.GetReservationsByProject(orgUUID, projUUID)
		if err != nil {
			return 0, err
		}
		// Filtrar apenas reservas de hoje (Datetime é string no formato "YYYY-MM-DD...")
		today := time.Now().Format("2006-01-02")
		count := 0
		for _, r := range reservations {
			if len(r.Datetime) >= 10 && r.Datetime[:10] == today {
				count++
			}
		}
		return count, nil

	default:
		return 0, fmt.Errorf("tipo de limite desconhecido: %s", limitType)
	}
}

// GetUsageAndLimits retorna o uso atual e os limites da organização
func (h *LimitHandler) GetUsageAndLimits(orgId, projectId string) (*UsageLimitsResponse, error) {
	response := &UsageLimitsResponse{
		Usage:  UsageData{},
		Limits: LimitsData{},
	}

	// Buscar assinatura
	orgPackage, err := h.packageRepo.GetOrganizationPackage(orgId)
	if err != nil {
		return nil, fmt.Errorf("organização sem plano ativo: %w", err)
	}

	// Preencher informações do pacote
	if orgPackage.Package != nil {
		response.PackageCode = orgPackage.Package.CodeName
		response.PackageName = orgPackage.Package.DisplayName
	}

	// Buscar limites do pacote
	limits, _ := h.packageRepo.GetPackageLimits(orgPackage.PackageId.String())

	// Mapear limites (default = -1 ilimitado, exceto audit logs que default = 0 desabilitado)
	response.Limits.MaxTables = -1
	response.Limits.MaxUsers = -1
	response.Limits.MaxProducts = -1
	response.Limits.MaxReservationsPerDay = -1
	response.Limits.MaxAuditLogs = 0
	response.Limits.AuditLogsRetentionDays = 0

	for _, l := range limits {
		switch l.LimitType {
		case "users":
			response.Limits.MaxUsers = l.LimitValue
		case "tables":
			response.Limits.MaxTables = l.LimitValue
		case "products":
			response.Limits.MaxProducts = l.LimitValue
		case "reservations_per_day":
			response.Limits.MaxReservationsPerDay = l.LimitValue
		case "audit_logs_limit":
			response.Limits.MaxAuditLogs = l.LimitValue
		case "audit_logs_retention":
			response.Limits.AuditLogsRetentionDays = l.LimitValue
		}
	}

	// Verificar módulos habilitados
	modules, _ := h.packageRepo.GetPackageModules(orgPackage.PackageId.String())
	for _, m := range modules {
		switch m.CodeName {
		case "client_notifications":
			response.Limits.NotificationsEnabled = true
		case "client_reports":
			response.Limits.ReportsEnabled = true
		case "client_reservations":
			response.Limits.ReservationsEnabled = true
		case "client_waitlist":
			response.Limits.WaitlistEnabled = true
		case "client_audit_logs":
			response.Limits.AuditLogsEnabled = true
		}
	}

	// Contar uso atual
	orgUUID, _ := uuid.Parse(orgId)
	projUUID, _ := uuid.Parse(projectId)

	// Contar mesas
	if tables, err := h.tableRepo.ListTables(orgUUID, projUUID); err == nil {
		response.Usage.TablesCount = len(tables)
	}

	// Contar usuários
	if users, err := h.userOrgRepo.ListByOrganization(orgId); err == nil {
		response.Usage.UsersCount = len(users)
	}

	// Contar produtos
	if products, err := h.productRepo.ListProducts(orgUUID, projUUID); err == nil {
		response.Usage.ProductsCount = len(products)
	}

	// Contar reservas do dia
	if reservations, err := h.reservationRepo.GetReservationsByProject(orgUUID, projUUID); err == nil {
		today := time.Now().Format("2006-01-02")
		count := 0
		for _, r := range reservations {
			// Datetime está em formato string, comparar apenas a data
			if len(r.Datetime) >= 10 && r.Datetime[:10] == today {
				count++
			}
		}
		response.Usage.ReservationsToday = count
	}

	return response, nil
}

// HasModule verifica se a organização tem acesso a um módulo específico
func (h *LimitHandler) HasModule(orgId, moduleCode string) (bool, error) {
	// Buscar assinatura
	orgPackage, err := h.packageRepo.GetOrganizationPackage(orgId)
	if err != nil {
		return false, fmt.Errorf("organização sem plano ativo: %w", err)
	}

	// Buscar módulos do pacote
	modules, err := h.packageRepo.GetPackageModules(orgPackage.PackageId.String())
	if err != nil {
		return false, err
	}

	// Verificar se o módulo está na lista
	for _, m := range modules {
		if m.CodeName == moduleCode && m.Active {
			return true, nil
		}
	}

	return false, nil
}
