package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"short-url/domains/entities"
	"short-url/domains/repositories"
	"short-url/domains/service"
)

type UserSessionService struct {
	commandRepo repositories.UserSessionCommandRepositoryInterface
	queryRepo   repositories.UserSessionQueryRepositoryInterface
}

func NewUserSessionService(
	commandRepo repositories.UserSessionCommandRepositoryInterface,
	queryRepo repositories.UserSessionQueryRepositoryInterface,
) service.UserSessionServiceInterface {
	return &UserSessionService{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func (s *UserSessionService) CreateSession(ctx context.Context, userID uint, deviceInfo, ipAddress string) (*entities.UserSession, error) {
	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	secretKey, err := generateSecretKey()
	if err != nil {
		return nil, err
	}

	session := &entities.UserSession{
		UserID:       userID,
		SessionToken: sessionToken,
		SecretKey:    secretKey,
		DeviceInfo:   &deviceInfo,
		IPAddress:    &ipAddress,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		IsActive:     true,
	}

	err = s.commandRepo.Save(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateSecretKey() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}