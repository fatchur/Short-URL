package repositories

import (
	"context"
	"short-url/domains/entities"
)

type UserQueryRepositoryInterface interface {
	FindByID(ctx context.Context, id uint) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
}

type UserCommandRepositoryInterface interface {
	Save(ctx context.Context, user *entities.User) error
}