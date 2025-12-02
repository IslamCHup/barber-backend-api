package transport

import (
	"barber-backend-api/internal/models"
	"barber-backend-api/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BarberHandler struct {
	logger  *slog.Logger
	service service.BarberService
}

func NewBarberHandler(logger *slog.Logger, service service.BarberService) *BarberHandler {
	return &BarberHandler{logger: logger, service: service}
}

func (h *BarberHandler) RegisterRoutes(r *gin.Engine) {

	barbers := r.Group("/barbers")
	{
		barbers.GET("/", h.GetAll)
		barbers.GET("/:id", h.GetBarberByID)
		barbers.POST("/", h.AddBarber)
		barbers.PATCH("/:id", h.Update)
		barbers.DELETE("/:id", h.Delete)
	}
}

func (h *BarberHandler) AddBarber(c *gin.Context) {
	method := c.Request.Method
	uri := c.FullPath()

	h.logger.Info("запрос AddBarber начат", "method", method, "uri", uri)

	var req models.BarbersCreateReqDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("ошибка валидации AddBarber", "reason", err.Error(), "body_valid", false, "method", method, "uri", uri)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddBarber(&req); err != nil {
		h.logger.Error("сервисная ошибка AddBarber", "error", err, "method", method, "uri", uri)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("успешный ответ AddBarber", "method", method, "uri", uri, "status_code", http.StatusCreated, "entity_id", nil)
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (h *BarberHandler) Update(c *gin.Context) {
	method := c.Request.Method
	uri := c.FullPath()

	h.logger.Info("запрос Update начат", "method", method, "uri", uri)

	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("некорректный идентификатор в Update", "reason", "некорректный идентификатор", "id", idStr, "body_valid", true, "method", method, "uri", uri)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}
	id := uint(idUint64)

	var barberInput models.BarbersCreateReqDTO
	if err := c.ShouldBindJSON(&barberInput); err != nil {
		h.logger.Warn("ошибка валидации Update", "reason", err.Error(), "id", id, "body_valid", false, "method", method, "uri", uri)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	barber := models.Barber{
		FullName: barberInput.FullName,
	}
	barberReq, err := h.service.Update(id, barber)
	if err != nil {
		h.logger.Error("сервисная ошибка Update", "error", err, "id", id, "method", method, "uri", uri)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("успешный ответ Update", "method", method, "uri", uri, "status_code", http.StatusOK, "entity_id", barberReq.ID)

	barberRes := models.BarberResDTO{
		FullName:  barberReq.FullName,
		AvgRating: barberReq.AvgRating,
	}

	c.JSON(http.StatusOK, barberRes)
}

func (h *BarberHandler) GetAll(c *gin.Context) {
	method := c.Request.Method
	uri := c.FullPath()

	h.logger.Info("запрос GetAll начат", "method", method, "uri", uri)

	barbers, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("сервисная ошибка GetAll", "error", err, "method", method, "uri", uri)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("успешный ответ GetAll", "method", method, "uri", uri, "status_code", http.StatusOK, "count", len(barbers))

	c.JSON(http.StatusOK, barbers)
}

func (h *BarberHandler) GetBarberByID(c *gin.Context) {
	method := c.Request.Method
	uri := c.FullPath()

	h.logger.Info("запрос GetBarberByID начат", "method", method, "uri", uri)

	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("некорректный идентификатор в GetBarberByID", "reason", "некорректный идентификатор", "id", idStr, "method", method, "uri", uri)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}
	id := uint(idUint64)

	barber, err := h.service.GetBarberByID(id)
	if err != nil {
		h.logger.Error("сервисная ошибка GetBarberByID", "error", err, "id", id, "method", method, "uri", uri)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("успешный ответ GetBarberByID", "method", method, "uri", uri, "status_code", http.StatusOK, "entity_id", barber.ID)

	c.JSON(http.StatusOK, barber)
}

func (h *BarberHandler) Delete(c *gin.Context) {
	method := c.Request.Method
	uri := c.FullPath()

	h.logger.Info("запрос Delete начат", "method", method, "uri", uri)

	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("некорректный идентификатор в Delete", "reason", "некорректный идентификатор", "id", idStr, "method", method, "uri", uri)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}
	id := uint(idUint64)

	if err = h.service.Delete(id); err != nil {
		h.logger.Error("сервисная ошибка Delete", "error", err, "id", id, "method", method, "uri", uri)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("успешный ответ Delete", "method", method, "uri", uri, "status_code", http.StatusNoContent, "entity_id", id)

	c.Status(http.StatusNoContent)
}
