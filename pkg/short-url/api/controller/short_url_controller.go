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
		response := dto.NewErrorResponse(fiber.StatusBadRequest, "Invalid request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		response := dto.NewErrorResponse(fiber.StatusUnauthorized, "User authentication required")
		return ctx.Status(fiber.StatusUnauthorized).JSON(response)
	}

	shortUrl, err := c.service.CreateShortUrl(ctx.Context(), req.LongUrl, userID)
	if err != nil {
		response := dto.NewErrorResponse(fiber.StatusInternalServerError, "Failed to create short URL")
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}

	responseData := dto.CreateShortUrlResponse{
		ID:        shortUrl.ID,
		ShortCode: shortUrl.ShortCode,
		LongUrl:   shortUrl.LongUrl,
		UserID:    shortUrl.UserID,
	}

	response := dto.NewSuccessResponse(fiber.StatusCreated, "Short URL created successfully", responseData)
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (c *ShortUrlController) GetLongUrl(ctx *fiber.Ctx) error {
	shortCode := ctx.Params("shortCode")
	if shortCode == "" {
		response := dto.NewErrorResponse(fiber.StatusBadRequest, "Short code is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		response := dto.NewErrorResponse(fiber.StatusUnauthorized, "User authentication required")
		return ctx.Status(fiber.StatusUnauthorized).JSON(response)
	}

	shortUrl, err := c.service.GetByShortCode(ctx.Context(), shortCode, userID)
	if err != nil {
		response := dto.NewErrorResponse(fiber.StatusNotFound, "Short URL not found or access denied")
		return ctx.Status(fiber.StatusNotFound).JSON(response)
	}

	acceptHeader := ctx.Get("Accept")
	if acceptHeader == "application/json" {
		responseData := map[string]interface{}{
			"short_code": shortUrl.ShortCode,
			"long_url":   shortUrl.LongUrl,
			"user_id":    shortUrl.UserID,
		}

		response := dto.NewSuccessResponse(fiber.StatusOK, "Short URL retrieved successfully", responseData)
		return ctx.Status(fiber.StatusOK).JSON(response)
	}

	return ctx.Redirect(shortUrl.LongUrl, fiber.StatusFound)
}

func (c *ShortUrlController) PublicRedirect(ctx *fiber.Ctx) error {
	shortCode := ctx.Params("shortCode")
	if shortCode == "" {
		response := dto.NewErrorResponse(fiber.StatusBadRequest, "Short code is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	shortUrl, err := c.service.GetByShortCodePublic(ctx.Context(), shortCode)
	if err != nil {
		response := dto.NewErrorResponse(fiber.StatusNotFound, "Short URL not found")
		return ctx.Status(fiber.StatusNotFound).JSON(response)
	}

	return ctx.Redirect(shortUrl.LongUrl, fiber.StatusFound)
}

func (c *ShortUrlController) RegisterRoutes(api fiber.Router) {
	api.Post("/url", c.CreateShortUrl)
}
