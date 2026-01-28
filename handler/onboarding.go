package handler

import (
	"lep/repositories"

	"github.com/google/uuid"
)

// OnboardingHandler handles onboarding status operations
type OnboardingHandler struct {
	repo *repositories.DBconn
}

// IOnboardingHandler interface for onboarding operations
type IOnboardingHandler interface {
	GetOnboardingStatus(orgId, projectId string) (*OnboardingStatus, error)
}

// OnboardingStatus represents the current onboarding progress
type OnboardingStatus struct {
	UsersCount                   int  `json:"users_count"`
	TablesCount                  int  `json:"tables_count"`
	CustomersCount               int  `json:"customers_count"`
	MenusCount                   int  `json:"menus_count"`
	CategoriesCount              int  `json:"categories_count"`
	ProductsCount                int  `json:"products_count"`
	ReservationsCount            int  `json:"reservations_count"`
	OrdersCount                  int  `json:"orders_count"`
	WaitlistCount                int  `json:"waitlist_count"`
	TagsCount                    int  `json:"tags_count"`
	HasNotificationsConfigured   bool `json:"has_notifications_configured"`
	HasDisplaySettingsConfigured bool `json:"has_display_settings_configured"`
}

// NewOnboardingHandler creates a new OnboardingHandler
func NewOnboardingHandler(repo *repositories.DBconn) IOnboardingHandler {
	return &OnboardingHandler{repo: repo}
}

// GetOnboardingStatus returns the current onboarding status for a project
func (h *OnboardingHandler) GetOnboardingStatus(orgId, projectId string) (*OnboardingStatus, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}

	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	status := &OnboardingStatus{}

	// Count users in the project via user_roles
	usersCount, err := h.repo.Roles.CountUsersByOrganization(orgId)
	if err == nil {
		status.UsersCount = usersCount
	}

	// Count tables
	tables, err := h.repo.Tables.ListTablesByProject(orgUUID, projectUUID)
	if err == nil {
		status.TablesCount = len(tables)
	}

	// Count customers
	customers, err := h.repo.Customers.ListCustomers(orgUUID, projectUUID)
	if err == nil {
		status.CustomersCount = len(customers)
	}

	// Count menus
	menus, err := h.repo.Menus.GetMenuList(orgUUID, projectUUID)
	if err == nil {
		status.MenusCount = len(menus)
	}

	// Count categories
	categories, err := h.repo.Categories.GetCategoryList(orgUUID, projectUUID)
	if err == nil {
		status.CategoriesCount = len(categories)
	}

	// Count products
	products, err := h.repo.Products.ListProducts(orgUUID, projectUUID)
	if err == nil {
		status.ProductsCount = len(products)
	}

	// Count reservations
	reservations, err := h.repo.Reservations.ListReservations(orgUUID, projectUUID)
	if err == nil {
		status.ReservationsCount = len(reservations)
	}

	// Count orders
	orders, err := h.repo.Orders.ListOrders(orgId, projectId)
	if err == nil {
		status.OrdersCount = len(orders)
	}

	// Count waitlist entries
	waitlist, err := h.repo.Waitlists.ListWaitlists(orgUUID, projectUUID)
	if err == nil {
		status.WaitlistCount = len(waitlist)
	}

	// Count tags
	tags, err := h.repo.Tags.GetTagList(orgUUID, projectUUID)
	if err == nil {
		status.TagsCount = len(tags)
	}

	// Check if notifications are configured
	notificationConfigs, err := h.repo.Notifications.GetNotificationConfigs(orgUUID, projectUUID)
	if err == nil && len(notificationConfigs) > 0 {
		status.HasNotificationsConfigured = true
	}

	// Check if display settings are configured
	displaySettings, err := h.repo.DisplaySettings.GetSettingsByProject(projectUUID)
	if err == nil && displaySettings != nil {
		status.HasDisplaySettingsConfigured = true
	}

	return status, nil
}
