package transport

import (
	"barber-backend-api/internal/models"
	"barber-backend-api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientsHandler struct {
	service service.ClientService
}

func NewClientsHandler(service service.ClientService) *ClientsHandler {
	return &ClientsHandler{service: service}
}

func (h *ClientsHandler) RegisterRoutes(r *gin.Engine) {
	client := r.Group("/clients")
	{
		client.GET("/", h.GetAllClients)
		client.GET("/:id", h.GetClientByID)
		client.POST("/", h.AddClient)
		client.PATCH("/:id", h.Update)
		client.DELETE("/:id", h.Delete)
	}
}

func (h *ClientsHandler) AddClient(c *gin.Context) {
	var req models.ClientCreateReqDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AddClient(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (h *ClientsHandler) GetAllClients(c *gin.Context) {
	clients, err := h.service.GetAllClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clients)
}

func (h *ClientsHandler) GetClientByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}
	client, err := h.service.GetClientByID(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

func (h *ClientsHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}

	var client models.ClientUpdateReqDTO

	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.Update(uint(id), client)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *ClientsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}
	err = h.service.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
