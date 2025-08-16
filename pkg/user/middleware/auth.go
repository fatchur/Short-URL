package middleware

import (
	"strings"

	"short-url/domains/repositories"
	"short-url/domains/helper/jwt"

	"github.com/gofiber/fiber/v2"
)


func GetUserIDFromContext(c *fiber.Ctx) uint {
	userID := c.Locals("userID")
	if userID == nil {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

func SetUserIDToContext(c *fiber.Ctx, userID uint) {
	c.Locals("userID", userID)
}

func ExtractTokenFromHeader(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return ""
	}

	return tokenString
}

func JWTAuth(sessionQueryRepo repositories.UserSessionQueryRepositoryInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := ExtractTokenFromHeader(c)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid authorization header",
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

		SetUserIDToContext(c, claims.UserID)
		return c.Next()
	}
}
