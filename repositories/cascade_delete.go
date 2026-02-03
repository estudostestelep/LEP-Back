package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CascadeDeleteRepository handles cascade deletion of related entities
type CascadeDeleteRepository struct {
	db *gorm.DB
}

type ICascadeDeleteRepository interface {
	// Organization cascade operations
	SoftDeleteOrganizationCascade(orgId uuid.UUID) error
	HardDeleteOrganizationCascade(orgId uuid.UUID) error

	// Project cascade operations
	SoftDeleteProjectCascade(projectId uuid.UUID) error
	HardDeleteProjectCascade(projectId uuid.UUID) error
}

func NewCascadeDeleteRepository(db *gorm.DB) ICascadeDeleteRepository {
	return &CascadeDeleteRepository{db: db}
}

// SoftDeleteOrganizationCascade soft deletes an organization and all related data
func (r *CascadeDeleteRepository) SoftDeleteOrganizationCascade(orgId uuid.UUID) error {
	now := time.Now()

	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get all projects of this organization to cascade delete their data
		var projects []models.Project
		if err := tx.Where("organization_id = ? AND deleted_at IS NULL", orgId).Find(&projects).Error; err != nil {
			return err
		}

		// 2. Soft delete all project-related data for each project
		for _, project := range projects {
			if err := r.softDeleteProjectData(tx, project.Id, now); err != nil {
				return err
			}
		}

		// 3. Soft delete all projects
		if err := tx.Model(&models.Project{}).Where("organization_id = ?", orgId).Update("deleted_at", now).Error; err != nil {
			return err
		}

		// 4. Soft delete organization-level data (that doesn't depend on project)
		if err := r.softDeleteOrganizationData(tx, orgId, now); err != nil {
			return err
		}

		// 5. Finally, soft delete the organization itself
		if err := tx.Model(&models.Organization{}).Where("id = ?", orgId).Update("deleted_at", now).Error; err != nil {
			return err
		}

		return nil
	})
}

// HardDeleteOrganizationCascade permanently deletes an organization and all related data
func (r *CascadeDeleteRepository) HardDeleteOrganizationCascade(orgId uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get all projects of this organization
		var projects []models.Project
		if err := tx.Unscoped().Where("organization_id = ?", orgId).Find(&projects).Error; err != nil {
			return err
		}

		// 2. Hard delete all project-related data for each project
		for _, project := range projects {
			if err := r.hardDeleteProjectData(tx, project.Id); err != nil {
				return err
			}
		}

		// 3. Hard delete all projects
		if err := tx.Unscoped().Where("organization_id = ?", orgId).Delete(&models.Project{}).Error; err != nil {
			return err
		}

		// 4. Hard delete organization-level data
		if err := r.hardDeleteOrganizationData(tx, orgId); err != nil {
			return err
		}

		// 5. Finally, hard delete the organization itself
		if err := tx.Unscoped().Delete(&models.Organization{}, "id = ?", orgId).Error; err != nil {
			return err
		}

		return nil
	})
}

// SoftDeleteProjectCascade soft deletes a project and all related data
func (r *CascadeDeleteRepository) SoftDeleteProjectCascade(projectId uuid.UUID) error {
	now := time.Now()

	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Soft delete all project-related data
		if err := r.softDeleteProjectData(tx, projectId, now); err != nil {
			return err
		}

		// 2. Soft delete the project itself
		if err := tx.Model(&models.Project{}).Where("id = ?", projectId).Update("deleted_at", now).Error; err != nil {
			return err
		}

		return nil
	})
}

// HardDeleteProjectCascade permanently deletes a project and all related data
func (r *CascadeDeleteRepository) HardDeleteProjectCascade(projectId uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Hard delete all project-related data
		if err := r.hardDeleteProjectData(tx, projectId); err != nil {
			return err
		}

		// 2. Hard delete the project itself
		if err := tx.Unscoped().Delete(&models.Project{}, "id = ?", projectId).Error; err != nil {
			return err
		}

		return nil
	})
}

// softDeleteProjectData soft deletes all data related to a specific project
func (r *CascadeDeleteRepository) softDeleteProjectData(tx *gorm.DB, projectId uuid.UUID, now time.Time) error {
	// Order deletion
	if err := tx.Model(&models.Order{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Order items (cascade from orders - they have order_id, not project_id directly)
	// Items are deleted when orders are deleted via foreign key

	// Reservations
	if err := tx.Model(&models.Reservation{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Waitlist
	if err := tx.Model(&models.Waitlist{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Tables
	if err := tx.Model(&models.Table{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Products
	if err := tx.Model(&models.Product{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Menus
	if err := tx.Model(&models.Menu{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Categories
	if err := tx.Model(&models.Category{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Subcategories
	if err := tx.Model(&models.Subcategory{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Client Audit Logs (não tem deleted_at - deletar permanentemente)
	if err := tx.Where("project_id = ?", projectId).Delete(&models.ClientAuditLog{}).Error; err != nil {
		return err
	}

	// Customers
	if err := tx.Model(&models.Customer{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Tags
	if err := tx.Model(&models.Tag{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Environments
	if err := tx.Model(&models.Environment{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Blocked Periods
	if err := tx.Model(&models.BlockedPeriod{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Settings
	if err := tx.Model(&models.Settings{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Notification Configs
	if err := tx.Model(&models.NotificationConfig{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Notification Templates
	if err := tx.Model(&models.NotificationTemplate{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Notification Logs
	if err := tx.Model(&models.NotificationLog{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Notification Events (não tem deleted_at)
	if err := tx.Where("project_id = ?", projectId).Delete(&models.NotificationEvent{}).Error; err != nil {
		return err
	}

	// Notification Inbound (não tem deleted_at)
	if err := tx.Where("project_id = ?", projectId).Delete(&models.NotificationInbound{}).Error; err != nil {
		return err
	}

	// Notification Schedule (não tem deleted_at)
	if err := tx.Where("project_id = ?", projectId).Delete(&models.NotificationSchedule{}).Error; err != nil {
		return err
	}

	// Notification Reminder (não tem deleted_at)
	if err := tx.Where("project_id = ?", projectId).Delete(&models.NotificationReminder{}).Error; err != nil {
		return err
	}

	// Response Review Queue (não tem deleted_at)
	if err := tx.Where("project_id = ?", projectId).Delete(&models.ResponseReviewQueue{}).Error; err != nil {
		return err
	}

	// Project Display Settings (não tem deleted_at)
	if err := tx.Where("project_id = ?", projectId).Delete(&models.ProjectDisplaySettings{}).Error; err != nil {
		return err
	}

	// Theme Customization
	if err := tx.Model(&models.ThemeCustomization{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// File References
	if err := tx.Model(&models.FileReference{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Leads
	if err := tx.Model(&models.Lead{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Report Metrics
	if err := tx.Model(&models.ReportMetric{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Audit Logs
	if err := tx.Model(&models.AuditLog{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// User Roles (relationship table)
	if err := tx.Model(&models.ClientRole{}).Where("project_id = ?", projectId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	return nil
}

// hardDeleteProjectData permanently deletes all data related to a specific project
func (r *CascadeDeleteRepository) hardDeleteProjectData(tx *gorm.DB, projectId uuid.UUID) error {
	// Order items first (they reference orders)
	if err := tx.Unscoped().Where("order_id IN (SELECT id FROM orders WHERE project_id = ?)", projectId).Delete(&models.OrderItem{}).Error; err != nil {
		return err
	}

	// Orders
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Order{}).Error; err != nil {
		return err
	}

	// Reservations
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Reservation{}).Error; err != nil {
		return err
	}

	// Waitlist
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Waitlist{}).Error; err != nil {
		return err
	}

	// Tables
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Table{}).Error; err != nil {
		return err
	}

	// Products
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Product{}).Error; err != nil {
		return err
	}

	// Menus
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Menu{}).Error; err != nil {
		return err
	}

	// SubcategoryCategory (junção - deletar primeiro por FK com Category)
	if err := tx.Unscoped().Where("category_id IN (SELECT id FROM categories WHERE project_id = ?)", projectId).Delete(&models.SubcategoryCategory{}).Error; err != nil {
		return err
	}

	// Categories
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Category{}).Error; err != nil {
		return err
	}

	// Subcategories
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Subcategory{}).Error; err != nil {
		return err
	}

	// Client Audit Logs
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.ClientAuditLog{}).Error; err != nil {
		return err
	}

	// Customers
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Customer{}).Error; err != nil {
		return err
	}

	// Tags
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Tag{}).Error; err != nil {
		return err
	}

	// Environments
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Environment{}).Error; err != nil {
		return err
	}

	// Blocked Periods
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.BlockedPeriod{}).Error; err != nil {
		return err
	}

	// Settings
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Settings{}).Error; err != nil {
		return err
	}

	// Notification Configs
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.NotificationConfig{}).Error; err != nil {
		return err
	}

	// Notification Templates
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.NotificationTemplate{}).Error; err != nil {
		return err
	}

	// Notification Logs
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.NotificationLog{}).Error; err != nil {
		return err
	}

	// Notification Events
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.NotificationEvent{}).Error; err != nil {
		return err
	}

	// Notification Inbound
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.NotificationInbound{}).Error; err != nil {
		return err
	}

	// Notification Schedule
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.NotificationSchedule{}).Error; err != nil {
		return err
	}

	// Notification Reminder
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.NotificationReminder{}).Error; err != nil {
		return err
	}

	// Response Review Queue
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.ResponseReviewQueue{}).Error; err != nil {
		return err
	}

	// Project Display Settings
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.ProjectDisplaySettings{}).Error; err != nil {
		return err
	}

	// Theme Customization
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.ThemeCustomization{}).Error; err != nil {
		return err
	}

	// File References
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.FileReference{}).Error; err != nil {
		return err
	}

	// Leads
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.Lead{}).Error; err != nil {
		return err
	}

	// Report Metrics
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.ReportMetric{}).Error; err != nil {
		return err
	}

	// Audit Logs
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.AuditLog{}).Error; err != nil {
		return err
	}

	// User Roles (relationship table) - delete roles linked to this project
	if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(&models.ClientRole{}).Error; err != nil {
		return err
	}

	return nil
}

// softDeleteOrganizationData soft deletes organization-level data (not project-specific)
func (r *CascadeDeleteRepository) softDeleteOrganizationData(tx *gorm.DB, orgId uuid.UUID, now time.Time) error {
	// Organization Package (subscription)
	if err := tx.Model(&models.OrganizationPlan{}).Where("organization_id = ?", orgId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Plan Change Requests
	if err := tx.Model(&models.PlanChangeRequest{}).Where("organization_id = ?", orgId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Roles specific to organization
	if err := tx.Model(&models.Role{}).Where("organization_id = ?", orgId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// User Roles for this organization
	if err := tx.Model(&models.ClientRole{}).Where("organization_id = ?", orgId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	// Sidebar Config
	if err := tx.Model(&models.SidebarConfig{}).Where("organization_id = ?", orgId).Update("deleted_at", now).Error; err != nil {
		return err
	}

	return nil
}

// hardDeleteOrganizationData permanently deletes organization-level data
func (r *CascadeDeleteRepository) hardDeleteOrganizationData(tx *gorm.DB, orgId uuid.UUID) error {
	// Organization Package (subscription)
	if err := tx.Unscoped().Where("organization_id = ?", orgId).Delete(&models.OrganizationPlan{}).Error; err != nil {
		return err
	}

	// Plan Change Requests
	if err := tx.Unscoped().Where("organization_id = ?", orgId).Delete(&models.PlanChangeRequest{}).Error; err != nil {
		return err
	}

	// User Roles for this organization
	if err := tx.Unscoped().Where("organization_id = ?", orgId).Delete(&models.ClientRole{}).Error; err != nil {
		return err
	}

	// Roles specific to organization
	if err := tx.Unscoped().Where("organization_id = ?", orgId).Delete(&models.Role{}).Error; err != nil {
		return err
	}

	// Sidebar Config
	if err := tx.Unscoped().Where("organization_id = ?", orgId).Delete(&models.SidebarConfig{}).Error; err != nil {
		return err
	}

	return nil
}
