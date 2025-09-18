package utils

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

// CreateDefaultNotificationTemplates - Cria templates padrão para um projeto
func CreateDefaultNotificationTemplates(orgId, projectId uuid.UUID) []models.NotificationTemplate {
	templates := []models.NotificationTemplate{}

	// === SMS TEMPLATES ===

	// Reserva criada - SMS
	templates = append(templates, models.NotificationTemplate{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           "Reserva Criada - SMS",
		Channel:        "sms",
		Subject:        "",
		Body:           "Olá {{nome}}! Sua reserva foi confirmada para {{data_hora}} na mesa {{mesa}} para {{pessoas}} pessoas. Restaurante LEP.",
		Variables:      []string{"nome", "data_hora", "mesa", "pessoas"},
		Active:         true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	// Reserva atualizada - SMS
	templates = append(templates, models.NotificationTemplate{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           "Reserva Atualizada - SMS",
		Channel:        "sms",
		Subject:        "",
		Body:           "{{nome}}, sua reserva foi atualizada para {{data_hora}} na mesa {{mesa}} para {{pessoas}} pessoas. Restaurante LEP.",
		Variables:      []string{"nome", "data_hora", "mesa", "pessoas"},
		Active:         true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	// Reserva cancelada - SMS
	templates = append(templates, models.NotificationTemplate{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           "Reserva Cancelada - SMS",
		Channel:        "sms",
		Subject:        "",
		Body:           "{{nome}}, sua reserva para {{data_hora}} foi cancelada. Em caso de dúvidas, entre em contato. Restaurante LEP.",
		Variables:      []string{"nome", "data_hora"},
		Active:         true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	// Mesa disponível - SMS
	templates = append(templates, models.NotificationTemplate{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           "Mesa Disponível - SMS",
		Channel:        "sms",
		Subject:        "",
		Body:           "{{nome}}, uma mesa está disponível! Mesa {{mesa}} livre agora. Você tem 15 minutos para confirmar. Restaurante LEP.",
		Variables:      []string{"nome", "mesa"},
		Active:         true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	// Confirmação 24h - SMS
	templates = append(templates, models.NotificationTemplate{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           "Confirmação 24h - SMS",
		Channel:        "sms",
		Subject:        "",
		Body:           "{{nome}}, lembramos que sua reserva é amanhã {{data_hora}} na mesa {{mesa}}. Confirme por favor. Restaurante LEP.",
		Variables:      []string{"nome", "data_hora", "mesa"},
		Active:         true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	// === EMAIL TEMPLATES ===

	// Reserva criada - Email
	templates = append(templates, models.NotificationTemplate{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           "Reserva Criada - Email",
		Channel:        "email",
		Subject:        "Reserva Confirmada - Restaurante LEP",
		Body: `<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #2c3e50;">Reserva Confirmada!</h2>
			<p>Olá <strong>{{nome}}</strong>,</p>
			<p>Sua reserva foi confirmada com sucesso!</p>
			<div style="background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0;">
				<h3>Detalhes da Reserva:</h3>
				<p><strong>Data e Hora:</strong> {{data_hora}}</p>
				<p><strong>Mesa:</strong> {{mesa}}</p>
				<p><strong>Pessoas:</strong> {{pessoas}}</p>
			</div>
			<p>Aguardamos você!</p>
			<p>Atenciosamente,<br><strong>Restaurante LEP</strong></p>
		</div>`,
		Variables: []string{"nome", "data_hora", "mesa", "pessoas"},
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// Reserva atualizada - Email
	templates = append(templates, models.NotificationTemplate{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           "Reserva Atualizada - Email",
		Channel:        "email",
		Subject:        "Reserva Atualizada - Restaurante LEP",
		Body: `<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #f39c12;">Reserva Atualizada</h2>
			<p>Olá <strong>{{nome}}</strong>,</p>
			<p>Sua reserva foi atualizada.</p>
			<div style="background-color: #fff3cd; padding: 20px; border-radius: 5px; margin: 20px 0;">
				<h3>Novos Detalhes:</h3>
				<p><strong>Data e Hora:</strong> {{data_hora}}</p>
				<p><strong>Mesa:</strong> {{mesa}}</p>
				<p><strong>Pessoas:</strong> {{pessoas}}</p>
			</div>
			<p>Aguardamos você!</p>
			<p>Atenciosamente,<br><strong>Restaurante LEP</strong></p>
		</div>`,
		Variables: []string{"nome", "data_hora", "mesa", "pessoas"},
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// Reserva cancelada - Email
	templates = append(templates, models.NotificationTemplate{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           "Reserva Cancelada - Email",
		Channel:        "email",
		Subject:        "Reserva Cancelada - Restaurante LEP",
		Body: `<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #e74c3c;">Reserva Cancelada</h2>
			<p>Olá <strong>{{nome}}</strong>,</p>
			<p>Sua reserva para <strong>{{data_hora}}</strong> foi cancelada.</p>
			<p>Em caso de dúvidas, entre em contato conosco.</p>
			<p>Esperamos vê-lo em breve!</p>
			<p>Atenciosamente,<br><strong>Restaurante LEP</strong></p>
		</div>`,
		Variables: []string{"nome", "data_hora"},
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// Mesa disponível - Email
	templates = append(templates, models.NotificationTemplate{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           "Mesa Disponível - Email",
		Channel:        "email",
		Subject:        "Mesa Disponível Agora! - Restaurante LEP",
		Body: `<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #27ae60;">Mesa Disponível!</h2>
			<p>Olá <strong>{{nome}}</strong>,</p>
			<p>Ótima notícia! Uma mesa está disponível agora.</p>
			<div style="background-color: #d4edda; padding: 20px; border-radius: 5px; margin: 20px 0;">
				<h3>Mesa {{mesa}} - Disponível Agora!</h3>
				<p><strong>Atenção:</strong> Você tem 15 minutos para confirmar sua presença.</p>
			</div>
			<p>Venha já!</p>
			<p>Atenciosamente,<br><strong>Restaurante LEP</strong></p>
		</div>`,
		Variables: []string{"nome", "mesa"},
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// Confirmação 24h - Email
	templates = append(templates, models.NotificationTemplate{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           "Confirmação 24h - Email",
		Channel:        "email",
		Subject:        "Lembrete: Sua reserva é amanhã - Restaurante LEP",
		Body: `<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #3498db;">Lembrete de Reserva</h2>
			<p>Olá <strong>{{nome}}</strong>,</p>
			<p>Lembramos que sua reserva é amanhã!</p>
			<div style="background-color: #e3f2fd; padding: 20px; border-radius: 5px; margin: 20px 0;">
				<h3>Detalhes da Reserva:</h3>
				<p><strong>Data e Hora:</strong> {{data_hora}}</p>
				<p><strong>Mesa:</strong> {{mesa}}</p>
			</div>
			<p>Por favor, confirme sua presença.</p>
			<p>Aguardamos você!</p>
			<p>Atenciosamente,<br><strong>Restaurante LEP</strong></p>
		</div>`,
		Variables: []string{"nome", "data_hora", "mesa"},
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// === WHATSAPP TEMPLATES ===

	// Reserva criada - WhatsApp
	templates = append(templates, models.NotificationTemplate{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           "Reserva Criada - WhatsApp",
		Channel:        "whatsapp",
		Subject:        "",
		Body:           "🎉 *Reserva Confirmada!*\n\nOlá *{{nome}}*!\n\nSua reserva foi confirmada:\n📅 {{data_hora}}\n🪑 Mesa {{mesa}}\n👥 {{pessoas}} pessoas\n\nRestaurante LEP 🍽️",
		Variables:      []string{"nome", "data_hora", "mesa", "pessoas"},
		Active:         true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	// Mesa disponível - WhatsApp
	templates = append(templates, models.NotificationTemplate{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Name:           "Mesa Disponível - WhatsApp",
		Channel:        "whatsapp",
		Subject:        "",
		Body:           "🚨 *Mesa Disponível!*\n\n{{nome}}, temos uma mesa livre!\n\n🪑 Mesa {{mesa}} disponível AGORA\n⏰ Você tem 15 minutos para confirmar\n\nVenha já! 🏃‍♂️\n\nRestaurante LEP 🍽️",
		Variables:      []string{"nome", "mesa"},
		Active:         true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	return templates
}

// CreateDefaultNotificationConfigs - Cria configurações padrão para um projeto
func CreateDefaultNotificationConfigs(orgId, projectId uuid.UUID) []models.NotificationConfig {
	configs := []models.NotificationConfig{}

	// Configuração para criação de reserva
	configs = append(configs, models.NotificationConfig{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		EventType:      "reservation_create",
		Enabled:        true,
		Channels:       []string{"sms", "email"},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	// Configuração para atualização de reserva
	configs = append(configs, models.NotificationConfig{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		EventType:      "reservation_update",
		Enabled:        true,
		Channels:       []string{"sms", "email"},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	// Configuração para cancelamento de reserva
	configs = append(configs, models.NotificationConfig{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		EventType:      "reservation_cancel",
		Enabled:        true,
		Channels:       []string{"sms", "email"},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	// Configuração para mesa disponível
	configs = append(configs, models.NotificationConfig{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		EventType:      "table_available",
		Enabled:        true,
		Channels:       []string{"sms", "whatsapp"},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	// Configuração para confirmação 24h
	configs = append(configs, models.NotificationConfig{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		EventType:      "confirmation_24h",
		Enabled:        true,
		Channels:       []string{"sms", "email"},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	return configs
}