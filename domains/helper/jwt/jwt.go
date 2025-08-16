package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID      uint   `json:"user_id"`
	SessionCode string `json:"session_code"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(userID uint, sessionCode, secretKey string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := JWTClaims{
		UserID:      userID,
		SessionCode: sessionCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

func ValidateJWTToken(tokenString, secretKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ParseJWTToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	if err != nil {
		if claims, ok := token.Claims.(*JWTClaims); ok {
			return claims, nil
		}
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token format")
}
