package main

import (
	"context"
	"log"

	"inventory-service/api/controller"
	"inventory-service/api/repository"
	"inventory-service/api/service"
	"inventory-service/router"
	userrepo "user-service/api/repository"

	"short-url/domains/config"
	"short-url/domains/database"
	"short-url/domains/dto"
)

func main() {
	log.Println("Inventory service starting...")

	ctx := context.Background()
	cfg := config.LoadConfig()

	dbConfig := dto.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
		Timezone: cfg.DBTimezone,
		LogLevel: cfg.DBLogLevel,
	}

	db, err := database.DBConnect(ctx, dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	inventoryCommandRepo := repository.NewInventoryCommandRepository(db)
	inventoryQueryRepo := repository.NewInventoryQueryRepository(db)

	inventoryService := service.NewInventoryService(inventoryCommandRepo, inventoryQueryRepo)

	inventoryController := controller.NewInventoryController(inventoryService)

	sessionQueryRepo := userrepo.NewUserSessionQueryRepository(db)
	app := router.NewRouter(inventoryController, sessionQueryRepo)

	log.Println("Starting server on :8081...")
	if err := app.Listen(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
