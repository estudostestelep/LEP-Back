package utils

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func NewEmailService(host string, port int, username, password, from string) *EmailService {
	return &EmailService{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

// SendEmail envia email via SMTP
func (e *EmailService) SendEmail(to, subject, body string) (*EmailResponse, error) {
	// Validações básicas
	if strings.TrimSpace(to) == "" {
		return &EmailResponse{
			Status:       "failed",
			ErrorMessage: "recipient email is required",
		}, fmt.Errorf("recipient email is required")
	}

	if strings.TrimSpace(subject) == "" {
		subject = "Notificação LEP System"
	}

	// Configuração SMTP
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	addr := fmt.Sprintf("%s:%d", e.Host, e.Port)

	// Construir mensagem
	msg := e.buildMessage(e.From, to, subject, body)

	// Enviar email
	err := smtp.SendMail(addr, auth, e.From, []string{to}, []byte(msg))
	if err != nil {
		return &EmailResponse{
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, err
	}

	return &EmailResponse{
		Status: "sent",
	}, nil
}

// buildMessage constrói a mensagem de email no formato correto
func (e *EmailService) buildMessage(from, to, subject, body string) string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=utf-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return msg.String()
}

// SendEmailHTML envia email com conteúdo HTML
func (e *EmailService) SendEmailHTML(to, subject, htmlBody string) (*EmailResponse, error) {
	// Envolver o corpo HTML em uma estrutura básica se necessário
	if !strings.Contains(htmlBody, "<html>") {
		htmlBody = fmt.Sprintf(`
		<html>
		<head>
			<meta charset="utf-8">
			<title>%s</title>
		</head>
		<body>
			%s
		</body>
		</html>`, subject, htmlBody)
	}

	return e.SendEmail(to, subject, htmlBody)
}

// TestConnection testa a conexão SMTP
func (e *EmailService) TestConnection() error {
	addr := fmt.Sprintf("%s:%d", e.Host, e.Port)
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)

	// Tenta conectar sem enviar email
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Quit()

	// Testa autenticação
	if err = client.Auth(auth); err != nil {
		return err
	}

	return nil
}