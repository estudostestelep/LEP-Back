package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type TwilioService struct {
	AccountSid string
	AuthToken  string
	FromPhone  string
}

type TwilioSMSRequest struct {
	To   string `json:"to"`
	Body string `json:"body"`
}

type TwilioWhatsAppRequest struct {
	To   string `json:"to"`
	Body string `json:"body"`
}

type TwilioResponse struct {
	Sid         string `json:"sid"`
	Status      string `json:"status"`
	ErrorCode   string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type TwilioWebhookStatus struct {
	MessageSid    string `json:"MessageSid"`
	MessageStatus string `json:"MessageStatus"`
	ErrorCode     string `json:"ErrorCode,omitempty"`
}

type TwilioWebhookInbound struct {
	From string `json:"From"`
	To   string `json:"To"`
	Body string `json:"Body"`
	MessageSid string `json:"MessageSid"`
}

func NewTwilioService(accountSid, authToken, fromPhone string) *TwilioService {
	return &TwilioService{
		AccountSid: accountSid,
		AuthToken:  authToken,
		FromPhone:  fromPhone,
	}
}

// SendSMS envia SMS via Twilio
func (t *TwilioService) SendSMS(to, message string) (*TwilioResponse, error) {
	// URL da API do Twilio para SMS
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.AccountSid)

	// Dados do formulário
	data := url.Values{}
	data.Set("From", t.FromPhone)
	data.Set("To", to)
	data.Set("Body", message)

	// Criar requisição
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(t.AccountSid, t.AuthToken)

	// Executar requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var twilioResp TwilioResponse
	err = json.Unmarshal(body, &twilioResp)
	if err != nil {
		return nil, err
	}

	// Verificar se houve erro
	if resp.StatusCode >= 400 {
		return &twilioResp, fmt.Errorf("twilio error: %s", twilioResp.ErrorMessage)
	}

	return &twilioResp, nil
}

// SendWhatsApp envia WhatsApp via Twilio
func (t *TwilioService) SendWhatsApp(to, message, whatsappBusinessNumber string) (*TwilioResponse, error) {
	// URL da API do Twilio para WhatsApp
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.AccountSid)

	// Formatar números para WhatsApp (prefixo whatsapp:)
	whatsappFrom := fmt.Sprintf("whatsapp:%s", whatsappBusinessNumber)
	whatsappTo := fmt.Sprintf("whatsapp:%s", to)

	// Dados do formulário
	data := url.Values{}
	data.Set("From", whatsappFrom)
	data.Set("To", whatsappTo)
	data.Set("Body", message)

	// Criar requisição
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(t.AccountSid, t.AuthToken)

	// Executar requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var twilioResp TwilioResponse
	err = json.Unmarshal(body, &twilioResp)
	if err != nil {
		return nil, err
	}

	// Verificar se houve erro
	if resp.StatusCode >= 400 {
		return &twilioResp, fmt.Errorf("twilio error: %s", twilioResp.ErrorMessage)
	}

	return &twilioResp, nil
}

// ValidateWebhookSignature valida a assinatura do webhook do Twilio
func (t *TwilioService) ValidateWebhookSignature(signature, requestUrl, body string) bool {
	// Obter AuthToken do ambiente ou usar o da instância
	authToken := t.AuthToken
	if authToken == "" {
		authToken = os.Getenv("TWILIO_AUTH_TOKEN")
	}

	if authToken == "" {
		// Sem token de autenticação, não podemos validar
		return false
	}

	// Calcular HMAC-SHA1
	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(requestUrl + body))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// Comparar assinaturas usando função segura contra timing attacks
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}