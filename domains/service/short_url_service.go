package service

import (
	"context"

	"short-url/domains/dto"
	"short-url/domains/entities"
)

type ShortUrlServiceInterface interface {
	CreateShortUrl(ctx context.Context, longUrl string, userID uint) (*entities.ShortUrl, error)
	GetByShortCode(ctx context.Context, shortCode string) (*entities.ShortUrl, error)
	GetByFilter(ctx context.Context, filter dto.ShortUrlQueryFilter, pagination dto.Pagination) ([]entities.ShortUrl, *dto.PaginationResponse, error)
	IncrementClickCount(ctx context.Context, shortCode string) error
}