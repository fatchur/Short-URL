package middleware

import (
	"strings"

	"short-url/domains/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID       uint   `json:"user_id"`
	Email        string `json:"email"`
	SessionToken string `json:"session_token"`
	jwt.RegisteredClaims
}

const (
	ContextUserID    = "user_id"
	ContextUserEmail = "user_email"
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

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			claims, ok := token.Claims.(*Claims)
			if !ok {
				return nil, jwt.ErrInvalidKey
			}

			session, err := sessionQueryRepo.FindBySessionToken(c.Context(), claims.SessionToken)
			if err != nil {
				return nil, jwt.ErrInvalidKey
			}

			return []byte(session.SecretKey), nil
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