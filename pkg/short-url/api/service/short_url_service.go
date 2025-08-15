package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"short-url/domains/dto"
	"short-url/domains/entities"
	"short-url/domains/repositories"
	"short-url/domains/service"
)

type shortUrlService struct {
	commandRepo repositories.ShortUrlCommandRepositoryInterface
	queryRepo   repositories.ShortUrlQueryRepositoryInterface
	redisRepo   repositories.RedisRepositoryInterface
}

func NewShortUrlService(
	commandRepo repositories.ShortUrlCommandRepositoryInterface,
	queryRepo repositories.ShortUrlQueryRepositoryInterface,
	redisRepo repositories.RedisRepositoryInterface,
) service.ShortUrlServiceInterface {
	return &shortUrlService{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
		redisRepo:   redisRepo,
	}
}

func (s *shortUrlService) CreateShortUrl(ctx context.Context, longUrl string, userID uint) (*entities.ShortUrl, error) {
	shortCode := s.generateShortCode()

	shortUrl := &entities.ShortUrl{
		UserID:    userID,
		LongUrl:   longUrl,
		ShortCode: shortCode,
		IsActive:  true,
		CreatedAt: time.Now(),
		CreatedBy: userID,
		UpdatedAt: time.Now(),
		UpdatedBy: userID,
	}

	if err := s.commandRepo.Save(ctx, shortUrl); err != nil {
		return nil, fmt.Errorf("failed to save short url: %w", err)
	}

	return shortUrl, nil
}

func (s *shortUrlService) GetByShortCode(ctx context.Context, shortCode string, userID uint) (*entities.ShortUrl, error) {
	return s.queryRepo.FindByShortCodeAndUserID(ctx, shortCode, userID)
}

func (s *shortUrlService) GetByFilter(ctx context.Context, filter dto.ShortUrlQueryFilter, pagination dto.Pagination) ([]entities.ShortUrl, *dto.PaginationResponse, error) {
	return s.queryRepo.FindByFilter(ctx, filter, pagination)
}

func (s *shortUrlService) IncrementClickCount(ctx context.Context, shortCode string) error {
	key := fmt.Sprintf("click_count:%s", shortCode)
	_, err := s.redisRepo.Increment(ctx, key)
	return err
}

func (s *shortUrlService) generateShortCode() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	encoded := base64.URLEncoding.EncodeToString(bytes)
	return strings.TrimRight(encoded, "=")[:8]
}
