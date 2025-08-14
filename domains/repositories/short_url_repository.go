package repositories

import (
	"short-url/domains/entities"

	"gorm.io/gorm"
)

type ShortUrlRepository struct {
	db *gorm.DB
}

func NewShortUrlRepository(db *gorm.DB) *ShortUrlRepository {
	return &ShortUrlRepository{db: db}
}

func (r *ShortUrlRepository) Save(shortUrl *entities.ShortUrl) error {
	return r.db.Create(shortUrl).Error
}
