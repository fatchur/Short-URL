package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

const (
	ContextUserID    = "user_id"
	ContextUserEmail = "user_email"
)

func JWTAuth(secretKey string) fiber.Handler {
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

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		c.Locals(ContextUserID, claims.UserID)
		c.Locals(ContextUserEmail, claims.Email)

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

func GetUserEmailFromContext(c *fiber.Ctx) string {
	email, ok := c.Locals(ContextUserEmail).(string)
	if !ok {
		return ""
	}
	return email
}