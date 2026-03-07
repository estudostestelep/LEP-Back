package server

import (
	"fmt"
	"lep/handler"
	"lep/repositories/models"
	"lep/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceStaffSchedule struct {
	handler *handler.Handlers
}

type IServerStaffSchedule interface {
	ServiceGetById(c *gin.Context)
	ServiceListByWeek(c *gin.Context)
	ServiceListByDateRange(c *gin.Context)
	ServiceCreate(c *gin.Context)
	ServiceCreateBatch(c *gin.Context)
	ServiceUpdate(c *gin.Context)
	ServiceDelete(c *gin.Context)
	ServiceGetWeekSummary(c *gin.Context)
	ServiceSendEmails(c *gin.Context)
}

func NewStaffScheduleServer(handlers *handler.Handlers) IServerStaffSchedule {
	return &ResourceStaffSchedule{handler: handlers}
}

func (r *ResourceStaffSchedule) ServiceGetById(c *gin.Context) {
	id := c.Param("id")

	resp, err := r.handler.HandlerStaffSchedule.GetById(id)
	if err != nil {
		utils.SendNotFoundError(c, "Schedule")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffSchedule) ServiceListByWeek(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	weekStart := c.Param("weekStart")

	resp, err := r.handler.HandlerStaffSchedule.ListByWeek(orgId, projectId, weekStart)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing schedules", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffSchedule) ServiceListByDateRange(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	resp, err := r.handler.HandlerStaffSchedule.ListByDateRange(orgId, projectId, startDate, endDate)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing schedules", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffSchedule) ServiceCreate(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var req models.StaffSchedule
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	orgUUID, _ := uuid.Parse(orgId)
	projUUID, _ := uuid.Parse(projectId)
	req.OrganizationId = orgUUID
	req.ProjectId = projUUID

	err := r.handler.HandlerStaffSchedule.Create(&req)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating schedule", err)
		return
	}

	utils.SendCreatedSuccess(c, "Schedule created successfully", req)
}

func (r *ResourceStaffSchedule) ServiceCreateBatch(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var schedules []models.StaffSchedule
	if err := c.BindJSON(&schedules); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	orgUUID, _ := uuid.Parse(orgId)
	projUUID, _ := uuid.Parse(projectId)

	for i := range schedules {
		schedules[i].OrganizationId = orgUUID
		schedules[i].ProjectId = projUUID
	}

	err := r.handler.HandlerStaffSchedule.CreateBatch(schedules)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating schedules", err)
		return
	}

	utils.SendCreatedSuccess(c, "Schedules created successfully", gin.H{"count": len(schedules)})
}

func (r *ResourceStaffSchedule) ServiceUpdate(c *gin.Context) {
	id := c.Param("id")

	existing, err := r.handler.HandlerStaffSchedule.GetById(id)
	if err != nil {
		utils.SendNotFoundError(c, "Schedule")
		return
	}

	var req models.StaffSchedule
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	existing.ClientId = req.ClientId
	existing.ScheduleDate = req.ScheduleDate
	existing.Shift = req.Shift
	existing.Status = req.Status
	existing.SlotNumber = req.SlotNumber

	err = r.handler.HandlerStaffSchedule.Update(existing)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating schedule", err)
		return
	}

	utils.SendOKSuccess(c, "Schedule updated successfully", existing)
}

func (r *ResourceStaffSchedule) ServiceDelete(c *gin.Context) {
	id := c.Param("id")

	err := r.handler.HandlerStaffSchedule.Delete(id)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting schedule", err)
		return
	}

	utils.SendOKSuccess(c, "Schedule deleted successfully", nil)
}

func (r *ResourceStaffSchedule) ServiceGetWeekSummary(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	weekStart := c.Param("weekStart")

	resp, err := r.handler.HandlerStaffSchedule.GetWeekSummary(orgId, projectId, weekStart)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting summary", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ServiceSendEmails envia emails para todos os funcionários escalados na semana
func (r *ResourceStaffSchedule) ServiceSendEmails(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var req models.SendScheduleEmailsRequest
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Buscar escalas da semana (com Client preloaded)
	schedules, err := r.handler.HandlerStaffSchedule.ListByWeek(orgId, projectId, req.WeekStart)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing schedules", err)
		return
	}

	if len(schedules) == 0 {
		utils.SendOKSuccess(c, "Nenhuma escala encontrada para a semana", gin.H{"sent": 0})
		return
	}

	// Buscar projeto para configurações SMTP
	project, err := r.handler.HandlerProject.GetProjectById(projectId)
	if err != nil || project == nil {
		utils.SendNotFoundError(c, "Project")
		return
	}

	// Agrupar escalas por cliente
	type clientSchedules struct {
		name      string
		email     string
		schedules []models.StaffSchedule
	}
	clientMap := make(map[string]*clientSchedules)

	for _, sc := range schedules {
		if sc.Client == nil || sc.Client.Email == "" {
			continue
		}
		clientId := sc.ClientId.String()
		if _, ok := clientMap[clientId]; !ok {
			clientMap[clientId] = &clientSchedules{
				name:  sc.Client.Name,
				email: sc.Client.Email,
			}
		}
		clientMap[clientId].schedules = append(clientMap[clientId].schedules, sc)
	}

	if len(clientMap) == 0 {
		utils.SendOKSuccess(c, "Nenhum funcionário com email encontrado na escala", gin.H{"sent": 0})
		return
	}

	// Configurar serviço de notificação
	notif := utils.NewNotificationService()

	// Enviar emails e coletar IDs para marcar como enviado
	var sentIds []string
	var errors []string

	for _, cs := range clientMap {
		body := buildScheduleEmailBody(cs.name, cs.schedules, req.WeekStart)
		subject := fmt.Sprintf("Sua escala - semana de %s", formatWeekStartBR(req.WeekStart))

		result, err := notif.SendNotification(utils.NotificationRequest{
			Channel:   "email",
			Recipient: cs.email,
			Subject:   subject,
			Message:   body,
		}, project)

		if err != nil || result.Status == "failed" {
			errMsg := cs.email
			if result != nil && result.ErrorMessage != "" {
				errMsg += ": " + result.ErrorMessage
			} else if err != nil {
				errMsg += ": " + err.Error()
			}
			errors = append(errors, errMsg)
			continue
		}

		for _, sc := range cs.schedules {
			sentIds = append(sentIds, sc.Id.String())
		}
	}

	// Marcar escalas como email enviado
	if len(sentIds) > 0 {
		if err := r.handler.HandlerStaffSchedule.MarkEmailSent(sentIds); err != nil {
			utils.SendInternalServerError(c, "Error marking emails as sent", err)
			return
		}
	}

	response := gin.H{
		"sent":   len(clientMap) - len(errors),
		"failed": len(errors),
	}
	if len(errors) > 0 {
		response["errors"] = errors
	}

	utils.SendOKSuccess(c, "Emails processados", response)
}

// buildScheduleEmailBody monta o corpo do email com os dias escalados
func buildScheduleEmailBody(name string, schedules []models.StaffSchedule, weekStart string) string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Olá, %s!", name))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Você está escalado(a) para os seguintes dias na semana de %s:", formatWeekStartBR(weekStart)))
	lines = append(lines, "")

	for _, sc := range schedules {
		dayLabel := sc.ScheduleDate.Format("02/01/2006") + " (" + weekdayPT(sc.ScheduleDate.Weekday()) + ")"
		shiftLabel := shiftLabelPT(sc.Shift)
		lines = append(lines, fmt.Sprintf("• %s - %s", dayLabel, shiftLabel))
	}

	lines = append(lines, "")
	lines = append(lines, "Até lá!")

	return strings.Join(lines, "\n")
}

func shiftLabelPT(shift string) string {
	switch shift {
	case "almoco":
		return "Almoço"
	case "noite":
		return "Noite"
	default:
		return shift
	}
}

func weekdayPT(w time.Weekday) string {
	switch w {
	case time.Sunday:
		return "Domingo"
	case time.Monday:
		return "Segunda"
	case time.Tuesday:
		return "Terça"
	case time.Wednesday:
		return "Quarta"
	case time.Thursday:
		return "Quinta"
	case time.Friday:
		return "Sexta"
	case time.Saturday:
		return "Sábado"
	default:
		return w.String()
	}
}

func formatWeekStartBR(weekStart string) string {
	// weekStart is "2026-01-13" → "13/01/2026"
	if len(weekStart) == 10 {
		return weekStart[8:10] + "/" + weekStart[5:7] + "/" + weekStart[0:4]
	}
	return weekStart
}
