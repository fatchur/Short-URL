package repository

import (
	"context"

	"short-url/domains/entities"
	"short-url/domains/repositories"

	"gorm.io/gorm"
)

type shortUrlCommandRepository struct {
	db *gorm.DB
}

func NewShortUrlCommandRepository(db *gorm.DB) repositories.ShortUrlCommandRepositoryInterface {
	return &shortUrlCommandRepository{
		db: db,
	}
}

func (r *shortUrlCommandRepository) Save(ctx context.Context, shortUrl *entities.ShortUrl) error {
	return r.db.WithContext(ctx).Create(shortUrl).Error
}
