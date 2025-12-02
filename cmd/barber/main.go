package main

import (
	"barber-backend-api/internal/config"
	"barber-backend-api/internal/logging"
	"barber-backend-api/internal/models"
	"barber-backend-api/repository"
	"barber-backend-api/service"
	"barber-backend-api/transport"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	logger := logging.InitLogger()
	
	db := config.SetupDataBase(logger)

	if err := db.AutoMigrate(&models.Appointments{}, &models.Barber{}, &models.Client{}); err != nil {
		logger.Error("ошибка при выполнении автомиграции", "error", err)
		panic(fmt.Sprintf("не удалось выполнить миграции: %v", err))
	}
	

	appointmentsRepo := repository.NewAppointmentsRepository(db)
	clientsRepo := repository.NewClientsRepository(db)
	barberRepo := repository.NewBarbersRepository(logger, db)

	appointmentsService := service.NewAppointmentsService(appointmentsRepo, barberRepo)
	clientsService := service.NewClientsService( clientsRepo)
	barberService := service.NewBarbersService( logger, barberRepo)

	r := gin.Default()

	transport.RegisterRoutes(r, appointmentsService, barberService, clientsService, logger)

	if err := r.Run(); err != nil {
		panic(fmt.Sprintf("ошибка запуска сервера: %v", err))
	}

}
