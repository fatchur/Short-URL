package router

import (
	"time"
	"user-service/api/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func NewRouter(userController *controller.UserController) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "User Service API v1.0",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("User Service")
	})

	v1 := app.Group("/api/v1")
	user := v1.Group("/user")

	strictLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 3 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Too many login attempts. Please try again in 3 minutes.",
			})
		},
		SkipFailedRequests: true,
	})

	flexibleLimiter := limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please try again in 1 minute.",
			})
		},
	})

	user.Post("/session", strictLimiter, userController.CreateSession)

	userController.RegisterRoutes(user, flexibleLimiter)

	return app
}
