package router

import (
	"inventory-service/api/controller"
	"inventory-service/middleware"
	"short-url/domains/repositories"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(inventoryController *controller.InventoryController, sessionQueryRepo repositories.UserSessionQueryRepositoryInterface) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Inventory Service")
	})

	v1 := app.Group("/api/v1")
	protected := v1.Group("/", middleware.JWTAuth(sessionQueryRepo))
	inventoryController.RegisterRoutes(protected)

	return app
}
