package middleware

import (
	"strings"

	"short-url/domains/helper/jwt"
	"short-url/domains/repositories"

	"github.com/gofiber/fiber/v2"
)

const (
	ContextUserID = "user_id"
)

func JWTAuth(sessionQueryRepo repositories.UserSessionQueryRepositoryInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tempClaims, err := jwt.ParseJWTToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}

		session, err := sessionQueryRepo.FindBySessionCode(c.Context(), tempClaims.SessionCode)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session not found",
			})
		}

		claims, err := jwt.ValidateJWTToken(tokenString, session.SecretKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		c.Locals(ContextUserID, claims.UserID)

		return c.Next()
	}
}

func GetUserIDFromContext(c *fiber.Ctx) uint {
	userID, ok := c.Locals(ContextUserID).(uint)
	if !ok {
		return 0
	}
	return userID
}
