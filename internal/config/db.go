package config

import (
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDataBase(logger *slog.Logger) *gorm.DB {
	
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbConnectPort := os.Getenv("PORT")

	logger.Info("server started addr=:%v env=local", "addr", dbConnectPort, "env", "local")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		logger.Error("ошибка подключения к БД", "error", err)
		panic(err)
	}

	logger.Info("Подключение к БД прошло успешно")
	return db
}
