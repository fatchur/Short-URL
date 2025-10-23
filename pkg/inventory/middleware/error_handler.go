package middleware

import (
	"errors"
	"short-url/domains/dto"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		switch code {
		case fiber.StatusUnauthorized:
			return HandleUnauthorizedError(c, err)
		case fiber.StatusBadRequest:
			return HandleBadRequestError(c, err)
		case fiber.StatusNotFound:
			return HandleNotFoundError(c, err)
		case fiber.StatusForbidden:
			return HandleForbiddenError(c, err)
		default:
			return HandleInternalServerError(c, err)
		}
	}
}

func HandleUnauthorizedError(c *fiber.Ctx, err error) error {
	response := dto.NewErrorResponse(fiber.StatusUnauthorized, "Unauthorized access. Please provide valid authentication credentials.")
	return c.Status(fiber.StatusUnauthorized).JSON(response)
}

func HandleBadRequestError(c *fiber.Ctx, err error) error {
	message := "Bad request. Please check your input data."

	if err != nil {
		message = err.Error()
	}

	response := dto.NewErrorResponse(fiber.StatusBadRequest, message)
	return c.Status(fiber.StatusBadRequest).JSON(response)
}

func HandleNotFoundError(c *fiber.Ctx, err error) error {
	message := "Resource not found."

	if errors.Is(err, gorm.ErrRecordNotFound) {
		message = "Inventory item not found."
	}

	response := dto.NewErrorResponse(fiber.StatusNotFound, message)
	return c.Status(fiber.StatusNotFound).JSON(response)
}

func HandleForbiddenError(c *fiber.Ctx, err error) error {
	response := dto.NewErrorResponse(fiber.StatusForbidden, "Access forbidden. You don't have permission to access this resource.")
	return c.Status(fiber.StatusForbidden).JSON(response)
}

func HandleInternalServerError(c *fiber.Ctx, err error) error {
	message := "Internal server error. Please try again later."

	if err != nil {
	}

	response := dto.NewErrorResponse(fiber.StatusInternalServerError, message)
	return c.Status(fiber.StatusInternalServerError).JSON(response)
}

func HandleValidationError(c *fiber.Ctx, err error) error {
	message := "Validation failed. Please check your input data."

	if err != nil {
		message = err.Error()
	}

	response := dto.NewErrorResponse(fiber.StatusUnprocessableEntity, message)
	return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
}

func HandleDatabaseError(c *fiber.Ctx, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return HandleNotFoundError(c, err)
	}

	return HandleInternalServerError(c, err)
}
