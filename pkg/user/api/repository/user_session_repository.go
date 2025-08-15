package repository

import (
	"context"
	"time"

	"short-url/domains/entities"
	"short-url/domains/repositories"

	"gorm.io/gorm"
)

type UserSessionCommandRepository struct {
	db *gorm.DB
}

type UserSessionQueryRepository struct {
	db *gorm.DB
}

func NewUserSessionCommandRepository(db *gorm.DB) repositories.UserSessionCommandRepositoryInterface {
	return &UserSessionCommandRepository{db: db}
}

func NewUserSessionQueryRepository(db *gorm.DB) repositories.UserSessionQueryRepositoryInterface {
	return &UserSessionQueryRepository{db: db}
}

func (r *UserSessionCommandRepository) Save(ctx context.Context, session *entities.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *UserSessionCommandRepository) Update(ctx context.Context, session *entities.UserSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *UserSessionCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.UserSession{}, id).Error
}

func (r *UserSessionCommandRepository) DeactivateUserSessions(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&entities.UserSession{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Update("is_active", false).Error
}

func (r *UserSessionQueryRepository) FindByID(ctx context.Context, id uint) (*entities.UserSession, error) {
	var session entities.UserSession
	err := r.db.WithContext(ctx).
		Where("id = ? AND is_active = ? AND expires_at > ?", id, true, time.Now()).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *UserSessionQueryRepository) FindBySessionToken(ctx context.Context, token string) (*entities.UserSession, error) {
	var session entities.UserSession
	err := r.db.WithContext(ctx).
		Where("session_token = ? AND is_active = ? AND expires_at > ?", token, true, time.Now()).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *UserSessionQueryRepository) FindActiveByUserID(ctx context.Context, userID uint) ([]entities.UserSession, error) {
	var sessions []entities.UserSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ? AND expires_at > ?", userID, true, time.Now()).
		Find(&sessions).Error
	return sessions, err
}

func (r *UserSessionQueryRepository) FindByUserIDAndToken(ctx context.Context, userID uint, token string) (*entities.UserSession, error) {
	var session entities.UserSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND session_token = ? AND is_active = ? AND expires_at > ?", userID, token, true, time.Now()).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}