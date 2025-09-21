package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	Code      string    `json:"code,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
}

type SuccessResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// SendError envia uma resposta de erro padronizada
func SendError(c *gin.Context, statusCode int, message string, err error) {
	response := ErrorResponse{
		Error:     http.StatusText(statusCode),
		Message:   message,
		Timestamp: time.Now(),
		Path:      c.Request.URL.Path,
	}

	if err != nil {
		response.Details = err.Error()
	}

	c.JSON(statusCode, response)
}

// SendValidationError envia erro de validação específico
func SendValidationError(c *gin.Context, message string, validationErr error) {
	response := ErrorResponse{
		Error:     "Validation Error",
		Message:   message,
		Details:   validationErr.Error(),
		Code:      "VALIDATION_FAILED",
		Timestamp: time.Now(),
		Path:      c.Request.URL.Path,
	}

	c.JSON(http.StatusBadRequest, response)
}

// SendBadRequestError envia erro 400 padronizado
func SendBadRequestError(c *gin.Context, message string, err error) {
	SendError(c, http.StatusBadRequest, message, err)
}

// SendNotFoundError envia erro 404 padronizado
func SendNotFoundError(c *gin.Context, resource string) {
	message := resource + " not found"
	SendError(c, http.StatusNotFound, message, nil)
}

// SendInternalServerError envia erro 500 padronizado
func SendInternalServerError(c *gin.Context, message string, err error) {
	SendError(c, http.StatusInternalServerError, message, err)
}

// SendForbiddenError envia erro 403 padronizado
func SendForbiddenError(c *gin.Context, message string) {
	SendError(c, http.StatusForbidden, message, nil)
}

// SendUnauthorizedError envia erro 401 padronizado
func SendUnauthorizedError(c *gin.Context, message string) {
	SendError(c, http.StatusUnauthorized, message, nil)
}

// SendSuccess envia resposta de sucesso padronizada
func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	response := SuccessResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}

	c.JSON(statusCode, response)
}

// SendCreatedSuccess envia resposta 201 padronizada
func SendCreatedSuccess(c *gin.Context, message string, data interface{}) {
	SendSuccess(c, http.StatusCreated, message, data)
}

// SendOKSuccess envia resposta 200 padronizada
func SendOKSuccess(c *gin.Context, message string, data interface{}) {
	SendSuccess(c, http.StatusOK, message, data)
}