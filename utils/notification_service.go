package utils

import (
	"fmt"
	"lep/repositories/models"
	"strings"
	"time"
)

type NotificationService struct {
	twilioService *TwilioService
	emailService  *EmailService
}

type NotificationRequest struct {
	Channel   string            `json:"channel"` // "sms", "email", "whatsapp"
	Recipient string            `json:"recipient"`
	Subject   string            `json:"subject,omitempty"`
	Message   string            `json:"message"`
	Variables map[string]string `json:"variables,omitempty"` // Para templates
}

type NotificationResult struct {
	Status       string `json:"status"` // "sent", "failed"
	ExternalId   string `json:"external_id,omitempty"` // MessageSid, etc.
	ErrorMessage string `json:"error_message,omitempty"`
	Channel      string `json:"channel"`
	Recipient    string `json:"recipient"`
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// ConfigureTwilio configura o service Twilio
func (n *NotificationService) ConfigureTwilio(accountSid, authToken, fromPhone string) {
	n.twilioService = NewTwilioService(accountSid, authToken, fromPhone)
}

// ConfigureEmail configura o service de email
func (n *NotificationService) ConfigureEmail(host string, port int, username, password, from string) {
	n.emailService = NewEmailService(host, port, username, password, from)
}

// SendNotification envia notificação pelo canal especificado
func (n *NotificationService) SendNotification(req NotificationRequest, project *models.Project) (*NotificationResult, error) {
	// Processar template se houver variáveis
	message := n.processTemplate(req.Message, req.Variables)
	subject := n.processTemplate(req.Subject, req.Variables)

	switch strings.ToLower(req.Channel) {
	case "sms":
		return n.sendSMS(req.Recipient, message, project)
	case "email":
		return n.sendEmail(req.Recipient, subject, message, project)
	case "whatsapp":
		return n.sendWhatsApp(req.Recipient, message, project)
	default:
		return &NotificationResult{
			Status:       "failed",
			ErrorMessage: fmt.Sprintf("unsupported channel: %s", req.Channel),
			Channel:      req.Channel,
			Recipient:    req.Recipient,
		}, fmt.Errorf("unsupported channel: %s", req.Channel)
	}
}

// sendSMS envia SMS via Twilio
func (n *NotificationService) sendSMS(to, message string, project *models.Project) (*NotificationResult, error) {
	if n.twilioService == nil {
		// Configurar Twilio com dados do projeto
		if project.TwilioAccountSid == nil || project.TwilioAuthToken == nil || project.TwilioPhone == nil {
			return &NotificationResult{
				Status:       "failed",
				ErrorMessage: "Twilio not configured for this project",
				Channel:      "sms",
				Recipient:    to,
			}, fmt.Errorf("twilio not configured")
		}
		n.ConfigureTwilio(*project.TwilioAccountSid, *project.TwilioAuthToken, *project.TwilioPhone)
	}

	resp, err := n.twilioService.SendSMS(to, message)
	if err != nil {
		return &NotificationResult{
			Status:       "failed",
			ErrorMessage: err.Error(),
			Channel:      "sms",
			Recipient:    to,
		}, err
	}

	return &NotificationResult{
		Status:     "sent",
		ExternalId: resp.Sid,
		Channel:    "sms",
		Recipient:  to,
	}, nil
}

// sendWhatsApp envia WhatsApp via Twilio
func (n *NotificationService) sendWhatsApp(to, message string, project *models.Project) (*NotificationResult, error) {
	if n.twilioService == nil {
		// Configurar Twilio com dados do projeto
		if project.TwilioAccountSid == nil || project.TwilioAuthToken == nil || project.WhatsappBusinessNumber == nil {
			return &NotificationResult{
				Status:       "failed",
				ErrorMessage: "WhatsApp not configured for this project",
				Channel:      "whatsapp",
				Recipient:    to,
			}, fmt.Errorf("whatsapp not configured")
		}
		n.ConfigureTwilio(*project.TwilioAccountSid, *project.TwilioAuthToken, *project.WhatsappBusinessNumber)
	}

	resp, err := n.twilioService.SendWhatsApp(to, message, *project.WhatsappBusinessNumber)
	if err != nil {
		return &NotificationResult{
			Status:       "failed",
			ErrorMessage: err.Error(),
			Channel:      "whatsapp",
			Recipient:    to,
		}, err
	}

	return &NotificationResult{
		Status:     "sent",
		ExternalId: resp.Sid,
		Channel:    "whatsapp",
		Recipient:  to,
	}, nil
}

// sendEmail envia email via SMTP
func (n *NotificationService) sendEmail(to, subject, message string, project *models.Project) (*NotificationResult, error) {
	if n.emailService == nil {
		// Configurar email com dados do projeto
		if project.SmtpHost == nil || project.SmtpPort == nil || project.SmtpUsername == nil || project.SmtpPassword == nil {
			return &NotificationResult{
				Status:       "failed",
				ErrorMessage: "Email not configured for this project",
				Channel:      "email",
				Recipient:    to,
			}, fmt.Errorf("email not configured")
		}

		from := to // fallback
		if project.SmtpFrom != nil {
			from = *project.SmtpFrom
		}

		n.ConfigureEmail(*project.SmtpHost, *project.SmtpPort, *project.SmtpUsername, *project.SmtpPassword, from)
	}

	resp, err := n.emailService.SendEmailHTML(to, subject, message)
	if err != nil {
		return &NotificationResult{
			Status:       "failed",
			ErrorMessage: err.Error(),
			Channel:      "email",
			Recipient:    to,
		}, err
	}

	return &NotificationResult{
		Status:    resp.Status,
		Channel:   "email",
		Recipient: to,
	}, nil
}

// processTemplate processa template substituindo variáveis
func (n *NotificationService) processTemplate(template string, variables map[string]string) string {
	if variables == nil {
		return template
	}

	result := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// Adicionar variáveis padrão
	result = strings.ReplaceAll(result, "{{data}}", time.Now().Format("02/01/2006"))
	result = strings.ReplaceAll(result, "{{hora}}", time.Now().Format("15:04"))

	return result
}

// ValidateChannel verifica se o canal é suportado
func (n *NotificationService) ValidateChannel(channel string) bool {
	supportedChannels := []string{"sms", "email", "whatsapp"}
	for _, supported := range supportedChannels {
		if strings.ToLower(channel) == supported {
			return true
		}
	}
	return false
}

// GetSupportedChannels retorna lista de canais suportados
func (n *NotificationService) GetSupportedChannels() []string {
	return []string{"sms", "email", "whatsapp"}
}