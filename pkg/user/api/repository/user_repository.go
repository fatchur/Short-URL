package repository

import (
	"context"
	"short-url/domains/entities"
	"short-url/domains/repositories"

	"gorm.io/gorm"
)

type UserQueryRepository struct {
	db *gorm.DB
}

type UserCommandRepository struct {
	db *gorm.DB
}

func NewUserQueryRepository(db *gorm.DB) repositories.UserQueryRepositoryInterface {
	return &UserQueryRepository{db: db}
}

func NewUserCommandRepository(db *gorm.DB) repositories.UserCommandRepositoryInterface {
	return &UserCommandRepository{db: db}
}

func (r *UserQueryRepository) FindByID(ctx context.Context, id uint) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserQueryRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("email = ? AND is_active = ?", email, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserCommandRepository) Save(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}