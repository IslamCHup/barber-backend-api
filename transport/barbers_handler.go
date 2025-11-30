package transport

import (
	"barber-backend-api/internal/models"
	"barber-backend-api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BarberHandler struct {
	service service.BarberService
}

func NewBarberHandler(service service.BarberService) *BarberHandler {
	return &BarberHandler{service: service}
}

func (h *BarberHandler) RegisterRoutes(r *gin.Engine) {
	// Группа маршрутов /barbers
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
	var req models.BarbersCreateReqDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AddBarber(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (h *BarberHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}

	var barberInput models.BarbersCreateReqDTO

	if err := c.ShouldBindJSON(&barberInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	barber := models.Barber{
		FullName: barberInput.FullName,
	}
	barberReq, err := h.service.Update(uint(id), barber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	barberRes := models.BarberResDTO{
		FullName: barberReq.FullName,
		AvgRating: barberReq.AvgRating,
	}

	c.JSON(http.StatusOK, barberRes)
}

func (h *BarberHandler) GetAll(c *gin.Context) {
	barbers, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, barbers)
}

func (h *BarberHandler) GetBarberByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}
	barber, err := h.service.GetBarberByID(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, barber)
}

func (h *BarberHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}
	if err = h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
