package controller

import (
	"short-url/domains/dto"
	"short-url/domains/service"

	"github.com/gofiber/fiber/v2"
)

type ShortUrlController struct {
	service service.ShortUrlServiceInterface
}

func NewShortUrlController(service service.ShortUrlServiceInterface) *ShortUrlController {
	return &ShortUrlController{
		service: service,
	}
}

func (c *ShortUrlController) CreateShortUrl(ctx *fiber.Ctx) error {
	var req dto.CreateShortUrlRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	shortUrl, err := c.service.CreateShortUrl(ctx.Context(), req.LongUrl, req.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create short URL",
		})
	}

	response := dto.CreateShortUrlResponse{
		ID:        shortUrl.ID,
		ShortCode: shortUrl.ShortCode,
		LongUrl:   shortUrl.LongUrl,
		UserID:    shortUrl.UserID,
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (c *ShortUrlController) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Post("/url", c.CreateShortUrl)
}
