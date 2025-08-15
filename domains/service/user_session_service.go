package service

import (
	"context"

	"short-url/domains/entities"
)

type UserSessionServiceInterface interface {
	CreateSession(ctx context.Context, email, password, deviceInfo, ipAddress string) (*entities.UserSession, error)
}