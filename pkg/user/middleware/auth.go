package middleware

import (
	"context"
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

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			claims, ok := token.Claims.(*Claims)
			if !ok {
				return nil, jwt.ErrInvalidKey
			}

			session, err := sessionQueryRepo.FindBySessionToken(context.Background(), claims.SessionToken)
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

		SetUserIDToContext(c, claims.UserID)
		return c.Next()
	}
}
