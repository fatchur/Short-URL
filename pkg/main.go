package main

import (
	"context"
	"log"

	"short-url/domains/config"
	"short-url/domains/database"
	"short-url/domains/dto"

	// User service imports
	userController "user-service/api/controller"
	userRepo "user-service/api/repository"
	userService "user-service/api/service"

	// Short URL service imports
	shortUrlController "short-url-service/api/controller"
	shortUrlRepo "short-url-service/api/repository"
	shortUrlService "short-url-service/api/service"
	shortUrlMiddleware "short-url-service/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	// Initialize repositories
	// User repositories
	userSessionCommandRepo := userRepo.NewUserSessionCommandRepository(db)
	userSessionQueryRepo := userRepo.NewUserSessionQueryRepository(db)
	userQueryRepo := userRepo.NewUserQueryRepository(db)

	// Short URL repositories
	shortUrlCommandRepo := shortUrlRepo.NewShortUrlCommandRepository(db)
	shortUrlQueryRepo := shortUrlRepo.NewShortUrlQueryRepository(db)
	redisRepo := shortUrlRepo.NewRedisRepository(redisClient)

	// Initialize services
	userSessionService := userService.NewUserSessionService(userSessionCommandRepo, userSessionQueryRepo, userQueryRepo)
	shortUrlSvc := shortUrlService.NewShortUrlService(shortUrlCommandRepo, shortUrlQueryRepo, redisRepo)

	userCtrl := userController.NewUserController(userSessionService)
	shortUrlCtrl := shortUrlController.NewShortUrlController(shortUrlSvc)

	app := fiber.New(fiber.Config{
		AppName: "Short URL Monolith v1.0",
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(helmet.New(helmet.Config{
		XSSProtection:             "1; mode=block",
		ContentTypeNosniff:        "nosniff",
		XFrameOptions:             "DENY",
		ReferrerPolicy:            "no-referrer",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "cross-origin",
		OriginAgentCluster:        "?1",
		XDNSPrefetchControl:       "off",
		XDownloadOptions:          "noopen",
		// XPermittedCrossDomainPolicies: "none",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Rate limiting
	strictLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: cfg.RateLimitDuration,
		// Message:    "Too many requests, please try again later",
	})

	flexibleLimiter := limiter.New(limiter.Config{
		Max:        100,
		Expiration: cfg.RateLimitDuration,
		// Message:    "Too many requests, please try again later",
	})

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "short-url-monolith",
			"version": "1.0.0",
		})
	})

	// API routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// User service routes
	user := v1.Group("/user")
	user.Post("/session", strictLimiter, userCtrl.CreateSession)
	userCtrl.RegisterRoutes(user, flexibleLimiter)

	// Short URL service routes
	url := v1.Group("/url")
	url.Post("/", flexibleLimiter, shortUrlMiddleware.JWTAuth(userSessionQueryRepo), shortUrlCtrl.CreateShortUrl)
	url.Get("/:shortCode", flexibleLimiter, shortUrlMiddleware.JWTAuth(userSessionQueryRepo), shortUrlCtrl.GetLongUrl)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	// Direct redirect routes (no auth required for public access) - MUST be absolutely last
	app.Get("/url/:shortCode", shortUrlCtrl.PublicRedirect) // Temporary: keep old route
	app.Get("/:shortCode", shortUrlCtrl.PublicRedirect)

	log.Printf("Monolith server starting on port %s", port)
	log.Printf("Health check available at: http://localhost:%s/health", port)
	log.Printf("User API available at: http://localhost:%s/api/v1/user", port)
	log.Printf("Short URL API available at: http://localhost:%s/api/v1/url", port)
	log.Printf("Short URL redirect available at: http://localhost:%s/:shortCode", port)

	if cfg.Environment == "production" {
		log.Fatal(app.ListenTLS(":"+port, cfg.TLSCertFile, cfg.TLSKeyFile))
	} else {
		log.Fatal(app.Listen(":" + port))
	}
}
