package repository

import (
	"context"

	"short-url/domains/dto"
	"short-url/domains/entities"
	"short-url/domains/repositories"

	"gorm.io/gorm"
)

type shortUrlQueryRepository struct {
	db *gorm.DB
}

func NewShortUrlQueryRepository(db *gorm.DB) repositories.ShortUrlQueryRepositoryInterface {
	return &shortUrlQueryRepository{
		db: db,
	}
}

func (r *shortUrlQueryRepository) FindByID(ctx context.Context, id uint) (*entities.ShortUrl, error) {
	var shortUrl entities.ShortUrl
	err := r.db.WithContext(ctx).First(&shortUrl, id).Error
	if err != nil {
		return nil, err
	}
	return &shortUrl, nil
}

func (r *shortUrlQueryRepository) FindByShortCode(ctx context.Context, shortCode string) (*entities.ShortUrl, error) {
	var shortUrl entities.ShortUrl
	err := r.db.WithContext(ctx).Where("short_code = ? AND is_active = ?", shortCode, true).First(&shortUrl).Error
	if err != nil {
		return nil, err
	}
	return &shortUrl, nil
}

func (r *shortUrlQueryRepository) FindByFilter(ctx context.Context, filter dto.ShortUrlQueryFilter, pagination dto.Pagination) ([]entities.ShortUrl, *dto.PaginationResponse, error) {
	var shortUrls []entities.ShortUrl
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.ShortUrl{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.ExpiredAt != nil {
		query = query.Where("expire_at < ?", *filter.ExpiredAt)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	pagination.SetDefaults()
	offset := pagination.GetOffset()

	if err := query.Offset(offset).Limit(pagination.PageSize).Find(&shortUrls).Error; err != nil {
		return nil, nil, err
	}

	paginationResponse := dto.NewPaginationResponse(pagination.Page, pagination.PageSize, total)

	return shortUrls, paginationResponse, nil
}