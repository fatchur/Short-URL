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

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

	app.Use(helmet.New(helmet.Config{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "DENY",
		HSTSMaxAge:         31536000,
		// HSTSIncludeSubDomains: true,
		ReferrerPolicy: "strict-origin-when-cross-origin",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 15 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Too many login attempts. Please try again in 15 minutes.",
			})
		},
		SkipFailedRequests: true,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	app.Use(logger.New())

	port := cfg.Port
	if port == "" {
		port = "8081"
	}

	if cfg.Environment == "production" {
		log.Printf("User service starting on port %s with HTTPS", port)
		log.Fatal(app.ListenTLS(":"+port, cfg.TLSCertFile, cfg.TLSKeyFile))
	} else {
		log.Printf("User service starting on port %s (HTTP - Development only)", port)
		log.Printf("WARNING: Using HTTP in development. Use HTTPS in production!")
		log.Fatal(app.Listen(":" + port))
	}
}
