package models

import (
	"time"

	"github.com/google/uuid"
)

// AccessLog - Log de acessos/login do usuário
type AccessLog struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserId    uuid.UUID `gorm:"not null;index" json:"user_id"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Location  string    `json:"location,omitempty"` // Estimativa baseada no IP
	LoginAt   time.Time `json:"login_at"`
	CreatedAt time.Time `json:"created_at"`
}

// AccessLogResponse - DTO para resposta da API
type AccessLogResponse struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Location  string    `json:"location,omitempty"`
	LoginAt   time.Time `json:"login_at"`
}

// AccessLogPaginatedResponse - Resposta paginada de logs de acesso
type AccessLogPaginatedResponse struct {
	Data       []AccessLogResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	TotalPages int                 `json:"total_pages"`
}
