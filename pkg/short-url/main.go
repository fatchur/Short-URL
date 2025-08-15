package main

import (
	"context"
	"log"

	"short-url-service/api/controller"
	"short-url-service/api/repository"
	"short-url-service/api/service"
	"short-url-service/router"
	userrepo "user-service/api/repository"

	"short-url/domains/config"
	"short-url/domains/database"
	"short-url/domains/dto"
)

func main() {
	log.Println("Short URL service starting...")

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

	cacheConfig := dto.CacheConfig{
		Host:     cfg.DBHost,
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	db, err := database.DBConnect(ctx, dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	redisClient, err := database.CacheConnect(ctx, cacheConfig)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	commandRepo := repository.NewShortUrlCommandRepository(db)
	queryRepo := repository.NewShortUrlQueryRepository(db)
	redisRepo := repository.NewRedisRepository(redisClient)

	shortUrlService := service.NewShortUrlService(commandRepo, queryRepo, redisRepo)

	shortUrlController := controller.NewShortUrlController(shortUrlService)

	sessionQueryRepo := userrepo.NewUserSessionQueryRepository(db)
	app := router.NewRouter(shortUrlController, sessionQueryRepo)

	log.Println("Starting server on :8080...")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
