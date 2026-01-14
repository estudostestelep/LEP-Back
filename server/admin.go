package server

import (
	"lep/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ============================================================================
// DEVELOPER MIGRATION ENDPOINT
// ============================================================================
// Este endpoint é exclusivo para desenvolvedores (pablo@lep.com)
// Permite executar migrações de banco de dados através de um botão no frontend.
//
// COMO USAR:
// 1. Adicione as novas migrações SQL no método ServiceRunDevMigration
// 2. Acesse o frontend logado como pablo@lep.com
// 3. Clique no botão "Run Migration" no header
//
// IMPORTANTE: Sempre adicione migrações de forma idempotente (IF NOT EXISTS)
// para evitar erros em execuções repetidas.
// ============================================================================

type AdminController struct {
	DB *gorm.DB
}

// ServiceResetPasswords reseta as senhas de usuários específicos (endpoint administrativo temporário)
func (a *AdminController) ServiceResetPasswords(c *gin.Context) {
	db := a.DB

	// Gerar hash correto para "senha123"
	password := "senha123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendInternalServerError(c, "Erro ao gerar hash", err)
		return
	}

	// Usuários para resetar
	userEmails := []string{
		"pablo@lep.com",
		"luan@lep.com",
		"eduardo@lep.com",
		"teste@gmail.com",
		"garcom1@gmail.com",
		"gerente1@gmail.com",
	}

	results := make(map[string]string)

	// Atualizar senha de cada usuário
	for _, email := range userEmails {
		result := db.Table("users").
			Where("email = ? AND deleted_at IS NULL", email).
			Update("password", string(hashedPassword))

		if result.Error != nil {
			results[email] = "erro: " + result.Error.Error()
			continue
		}

		if result.RowsAffected == 0 {
			results[email] = "não encontrado"
		} else {
			results[email] = "atualizado"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset completed",
		"results": results,
		"password": password,
	})
}

// ServiceRunDevMigration executa migrações de desenvolvimento
// Este endpoint é protegido e só deve ser acessível por desenvolvedores (pablo@lep.com)
// ============================================================================
// ADICIONE NOVAS MIGRAÇÕES AQUI:
// - Cada migração deve ser idempotente (usar IF NOT EXISTS, IF NOT EXIST, etc.)
// - Comente cada bloco explicando o que a migração faz
// - Mantenha a ordem cronológica
// ============================================================================
func (a *AdminController) ServiceRunDevMigration(c *gin.Context) {
	db := a.DB

	// Resultados de cada migração
	results := make(map[string]string)

	// =========================================================================
	// MIGRAÇÃO: Sistema de Notificações com Agendamento Flexível
	// Data: Janeiro 2026
	// Descrição: Adiciona campos para timing flexível de notificações,
	//            tabela de agendamentos e fila de revisão de respostas
	// =========================================================================

	// 1. Novos campos em settings para timing flexível
	err := db.Exec(`
		ALTER TABLE settings ADD COLUMN IF NOT EXISTS confirmation_hours_before INTEGER DEFAULT 24;
		ALTER TABLE settings ADD COLUMN IF NOT EXISTS reminder_hours_before INTEGER DEFAULT 0;
		ALTER TABLE settings ADD COLUMN IF NOT EXISTS auto_cancel_no_response_hours INTEGER DEFAULT 0;
		ALTER TABLE settings ADD COLUMN IF NOT EXISTS response_processing_mode VARCHAR(20) DEFAULT 'automatic';
	`).Error
	if err != nil {
		results["settings_fields"] = "erro: " + err.Error()
	} else {
		results["settings_fields"] = "sucesso"
	}

	// 2. Novos campos em notification_inbounds para processamento de respostas
	err = db.Exec(`
		ALTER TABLE notification_inbounds ADD COLUMN IF NOT EXISTS reservation_id UUID;
		ALTER TABLE notification_inbounds ADD COLUMN IF NOT EXISTS customer_id UUID;
		ALTER TABLE notification_inbounds ADD COLUMN IF NOT EXISTS response_type VARCHAR(50);
		ALTER TABLE notification_inbounds ADD COLUMN IF NOT EXISTS confidence_score DECIMAL(5,4);
		ALTER TABLE notification_inbounds ADD COLUMN IF NOT EXISTS processing_method VARCHAR(50);
		ALTER TABLE notification_inbounds ADD COLUMN IF NOT EXISTS action_taken VARCHAR(100);
	`).Error
	if err != nil {
		results["notification_inbounds_fields"] = "erro: " + err.Error()
	} else {
		results["notification_inbounds_fields"] = "sucesso"
	}

	// 3. Nova tabela notification_schedules para agendamento flexível
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notification_schedules (
			id UUID PRIMARY KEY,
			organization_id UUID NOT NULL,
			project_id UUID NOT NULL,
			event_type VARCHAR(100) NOT NULL,
			entity_type VARCHAR(50) NOT NULL,
			entity_id UUID NOT NULL,
			scheduled_for TIMESTAMP NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			processed_at TIMESTAMP,
			metadata JSON,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`).Error
	if err != nil {
		results["notification_schedules_table"] = "erro: " + err.Error()
	} else {
		results["notification_schedules_table"] = "sucesso"
	}

	// 4. Nova tabela response_review_queue para fila de revisão
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS response_review_queue (
			id UUID PRIMARY KEY,
			organization_id UUID NOT NULL,
			project_id UUID NOT NULL,
			inbound_id UUID NOT NULL,
			reservation_id UUID NOT NULL,
			customer_id UUID NOT NULL,
			customer_name VARCHAR(255),
			customer_phone VARCHAR(50),
			message_body TEXT,
			suggested_action VARCHAR(20),
			confidence_score DECIMAL(5,4),
			status VARCHAR(20) DEFAULT 'pending_review',
			reviewed_by UUID,
			reviewed_at TIMESTAMP,
			action_taken VARCHAR(50),
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`).Error
	if err != nil {
		results["response_review_queue_table"] = "erro: " + err.Error()
	} else {
		results["response_review_queue_table"] = "sucesso"
	}

	// 5. Índices para performance
	err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_schedules_due ON notification_schedules(scheduled_for, status);
		CREATE INDEX IF NOT EXISTS idx_review_queue_pending ON response_review_queue(organization_id, project_id, status);
	`).Error
	if err != nil {
		results["indexes"] = "erro: " + err.Error()
	} else {
		results["indexes"] = "sucesso"
	}

	// =========================================================================
	// FIM DAS MIGRAÇÕES
	// Adicione novas migrações ACIMA desta linha
	// =========================================================================

	c.JSON(http.StatusOK, gin.H{
		"message": "Dev migration completed",
		"results": results,
	})
}