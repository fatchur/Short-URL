package router

import (
	"short-url-service/api/controller"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(shortUrlController *controller.ShortUrlController) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Short URL Service")
	})

	v1 := app.Group("/api/v1")
	shortUrlController.RegisterRoutes(v1)

	return app
}
