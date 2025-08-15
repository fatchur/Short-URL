package service

import (
	"context"

	"short-url/domains/entities"
)

type UserSessionServiceInterface interface {
	CreateSession(ctx context.Context, userID uint, deviceInfo, ipAddress string) (*entities.UserSession, error)
}