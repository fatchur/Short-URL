package controller

import (
	"short-url/domains/dto"
	"short-url/domains/service"
	"user-service/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userSessionService service.UserSessionServiceInterface
}

func NewUserController(userSessionService service.UserSessionServiceInterface) *UserController {
	return &UserController{
		userSessionService: userSessionService,
	}
}

func (c *UserController) CreateSession(ctx *fiber.Ctx) error {
	var req dto.CreateSessionRequest

	if err := ctx.BodyParser(&req); err != nil {
		response := dto.NewErrorResponse(fiber.StatusBadRequest, "Invalid request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		response := dto.NewErrorResponse(fiber.StatusUnauthorized, "User authentication required")
		return ctx.Status(fiber.StatusUnauthorized).JSON(response)
	}

	deviceInfo := ""
	if req.DeviceInfo != "" {
		deviceInfo = req.DeviceInfo
	}

	ipAddress := ""
	if req.IPAddress != "" {
		ipAddress = req.IPAddress
	}

	session, err := c.userSessionService.CreateSession(ctx.Context(), userID, deviceInfo, ipAddress)
	if err != nil {
		response := dto.NewErrorResponse(fiber.StatusInternalServerError, "Failed to create session")
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}

	responseData := dto.CreateSessionResponse{
		SessionToken: session.SessionToken,
		ExpiresAt:    session.ExpiresAt,
	}

	response := dto.NewSuccessResponse(fiber.StatusCreated, "Session created successfully", responseData)
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (c *UserController) RegisterRoutes(api fiber.Router) {
	api.Post("/session", c.CreateSession)
}
