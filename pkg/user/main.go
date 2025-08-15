package main

import (
	"context"
	"log"

	"user-service/api/controller"
	"user-service/api/repository"
	"user-service/api/service"
	"user-service/router"

	"short-url/domains/config"
	"short-url/domains/database"
	"short-url/domains/dto"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
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

	sessionCommandRepo := repository.NewUserSessionCommandRepository(db)
	sessionQueryRepo := repository.NewUserSessionQueryRepository(db)
	userQueryRepo := repository.NewUserQueryRepository(db)

	userSessionService := service.NewUserSessionService(sessionCommandRepo, sessionQueryRepo, userQueryRepo)
	userController := controller.NewUserController(userSessionService)

	app := router.NewRouter(userController)
	app.Use(cors.New())
	app.Use(logger.New())

	port := cfg.Port
	if port == "" {
		port = "8081"
	}

	log.Printf("User service starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
