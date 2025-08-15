package controller

import (
	"short-url/domains/dto"
	"short-url/domains/service"
	"short-url-service/middleware"

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

	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User authentication required",
		})
	}

	shortUrl, err := c.service.CreateShortUrl(ctx.Context(), req.LongUrl, userID)
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

func (c *ShortUrlController) GetLongUrl(ctx *fiber.Ctx) error {
	shortCode := ctx.Params("shortCode")
	if shortCode == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User authentication required",
		})
	}

	shortUrl, err := c.service.GetByShortCode(ctx.Context(), shortCode, userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short URL not found or access denied",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"short_code": shortUrl.ShortCode,
		"long_url":   shortUrl.LongUrl,
		"user_id":    shortUrl.UserID,
	})
}

func (c *ShortUrlController) RegisterRoutes(api fiber.Router) {
	api.Post("/url", c.CreateShortUrl)
}
