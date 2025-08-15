package controller

import (
	"short-url/domains/dto"
	"short-url/domains/service"

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

	if req.Email == "" || req.Password == "" {
		response := dto.NewErrorResponse(fiber.StatusBadRequest, "Email and password are required")
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	session, err := c.userSessionService.CreateSession(ctx.Context(), req.Email, req.Password, req.DeviceInfo, req.IPAddress)
	if err != nil {
		response := dto.NewErrorResponse(fiber.StatusUnauthorized, "Invalid credentials")
		return ctx.Status(fiber.StatusUnauthorized).JSON(response)
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
