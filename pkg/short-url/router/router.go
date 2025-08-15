package router

import (
	"short-url-service/api/controller"
	"short-url-service/middleware"
	"short-url/domains/repositories"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(shortUrlController *controller.ShortUrlController, sessionQueryRepo repositories.UserSessionQueryRepositoryInterface) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Short URL Service")
	})

	v1 := app.Group("/api/v1")
	
	protected := v1.Group("/", middleware.JWTAuth(sessionQueryRepo))
	shortUrlController.RegisterRoutes(protected)

	return app
}
