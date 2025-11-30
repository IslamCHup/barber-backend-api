package transport

import (
	"barber-backend-api/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	appointments service.AppointmentsService,
	barbers service.BarberService,
	clients service.ClientService,
) {
	// Собираем хендлеры, внедряя зависимости (сервисы)
	appointmentsHandler := NewAppointmentsHandler(appointments)
	barbersHandler := NewBarberHandler(barbers)
	clientsHandler := NewClientsHandler(clients)

	// Каждый хендлер регистрирует маршруты в рамках своей ответственности
	appointmentsHandler.RegisterRoutes(router)
	barbersHandler.RegisterRoutes(router)
	clientsHandler.RegisterRoutes(router)
}
