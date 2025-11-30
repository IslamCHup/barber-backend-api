package transport

import (
	"barber-backend-api/internal/models"
	"barber-backend-api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AppointmentsHandler struct {
	service service.AppointmentsService
}

func NewAppointmentsHandler(service service.AppointmentsService) *AppointmentsHandler {
	return &AppointmentsHandler{service: service}
}

func (h *AppointmentsHandler) RegisterRoutes(r *gin.Engine) {

	appointments := r.Group("/appointments")
	{
		appointments.GET("/", h.GetAllAppointments)
		appointments.GET("/:id", h.GetAppointmentByID)
		appointments.POST("/", h.CreateAppointment)
		appointments.PATCH("/:id", h.Update)
		appointments.DELETE("/:id", h.Delete)
	}
	r.GET("barbers/:barbersID/appointments", h.GetAllAppointmentsByBarberID)
}


func (h *AppointmentsHandler) GetAllAppointments(c *gin.Context) {
	appointments, err := h.service.GetAllAppointments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (h *AppointmentsHandler) GetAppointmentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}

	app, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

func (h *AppointmentsHandler) GetAllAppointmentsByBarberID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}

	appointmentsOfBarber, err := h.service.GetAllAppointmentsByBarberID(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointmentsOfBarber)
}

func (h *AppointmentsHandler) CreateAppointment(c *gin.Context) {
	var req models.AppointmentsCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateAppointment(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (h *AppointmentsHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}

	var req models.AppointmentsUpdateReqDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(uint(id), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *AppointmentsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
