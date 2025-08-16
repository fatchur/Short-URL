package repositories

import (
	"context"

	"short-url/domains/entities"
)

type UserSessionCommandRepositoryInterface interface {
	Save(ctx context.Context, session *entities.UserSession) error
	Update(ctx context.Context, session *entities.UserSession) error
	Delete(ctx context.Context, id uint) error
	DeactivateUserSessions(ctx context.Context, userID uint) error
}

type UserSessionQueryRepositoryInterface interface {
	FindByID(ctx context.Context, id uint) (*entities.UserSession, error)
	FindBySessionCode(ctx context.Context, code string) (*entities.UserSession, error)
	FindActiveByUserID(ctx context.Context, userID uint) ([]entities.UserSession, error)
	FindByUserIDAndCode(ctx context.Context, userID uint, code string) (*entities.UserSession, error)
}