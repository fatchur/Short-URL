package router

import (
	"user-service/api/controller"

	"github.com/gofiber/fiber/v2"
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
	userController.RegisterRoutes(user)

	return app
}
