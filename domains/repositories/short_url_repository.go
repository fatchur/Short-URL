package repositories

import (
	"context"

	"short-url/domains/dto"
	"short-url/domains/entities"
)

type ShortUrlCommandRepositoryInterface interface {
	Save(ctx context.Context, shortUrl *entities.ShortUrl) error
}

type ShortUrlQueryRepositoryInterface interface {
	FindByID(ctx context.Context, id uint) (*entities.ShortUrl, error)
	FindByShortCode(ctx context.Context, shortCode string) (*entities.ShortUrl, error)
	FindByFilter(ctx context.Context, filter dto.ShortUrlQueryFilter, pagination dto.Pagination) ([]entities.ShortUrl, *dto.PaginationResponse, error)
}
