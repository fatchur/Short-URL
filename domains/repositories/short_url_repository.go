package repositories

import (
	"short-url/domains/dto"
	"short-url/domains/entities"
)

type ShortUrlCommandRepositoryInterface interface {
	Save(shortUrl *entities.ShortUrl) error
}

type ShortUrlQueryRepositoryInterface interface {
	FindByID(id uint) (*entities.ShortUrl, error)
	FindByShortCode(shortCode string) (*entities.ShortUrl, error)
	FindByFilter(filter dto.ShortUrlQueryFilter, pagination dto.Pagination) ([]entities.ShortUrl, *dto.PaginationResponse, error)
}
