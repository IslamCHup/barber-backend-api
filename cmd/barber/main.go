package main

import (
	"barber-backend-api/internal/config"
	"barber-backend-api/internal/models"
	"barber-backend-api/repository"
	"barber-backend-api/service"
	"barber-backend-api/transport"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.SetupDataBase()

	if err := db.AutoMigrate(&models.Appointments{}, &models.Barber{}, &models.Client{}); err != nil{
		panic(fmt.Sprintf("не удалось выполнить миграции: %v", err))
	}

	appointmentsRepo := repository.NewAppointmentsRepository(db)
	clientsRepo := repository.NewClientsRepository(db)
	barberRepo := repository.NewBarbersRepository(db)

	appointmentsService := service.NewAppointmentsService(appointmentsRepo, barberRepo)
	clientsService := service.NewClientsService(clientsRepo)
	barberService := service.NewBarbersService(barberRepo)
	
	r := gin.Default()

	transport.RegisterRoutes(r, appointmentsService, barberService, clientsService)

	if err := r.Run(); err != nil{
		panic(fmt.Sprintf("ошибка запуска сервера: %v", err))
	}

}