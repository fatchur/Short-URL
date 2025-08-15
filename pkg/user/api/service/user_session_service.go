package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"short-url/domains/entities"
	"short-url/domains/repositories"
	"short-url/domains/service"
	
	"golang.org/x/crypto/bcrypt"
)

type UserSessionService struct {
	commandRepo repositories.UserSessionCommandRepositoryInterface
	queryRepo   repositories.UserSessionQueryRepositoryInterface
	userQueryRepo repositories.UserQueryRepositoryInterface
}

func NewUserSessionService(
	commandRepo repositories.UserSessionCommandRepositoryInterface,
	queryRepo repositories.UserSessionQueryRepositoryInterface,
	userQueryRepo repositories.UserQueryRepositoryInterface,
) service.UserSessionServiceInterface {
	return &UserSessionService{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
		userQueryRepo: userQueryRepo,
	}
}


func (s *UserSessionService) CreateSession(ctx context.Context, email, password, deviceInfo, ipAddress string) (*entities.UserSession, error) {
	user, err := s.userQueryRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	secretKey, err := generateSecretKey()
	if err != nil {
		return nil, err
	}

	session := &entities.UserSession{
		UserID:       user.ID,
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