package service

import (
	"context"

	"short-url/domains/dto"
)

type UserSessionServiceInterface interface {
	CreateSession(ctx context.Context, email, password, deviceInfo, ipAddress string) (*dto.SessionTokenData, error)
}
