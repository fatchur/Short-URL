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

func (s *shortUrlService) GetByShortCodePublic(ctx context.Context, shortCode string) (*entities.ShortUrl, error) {
	if s.redisRepo != nil {
		cachedUrl, err := s.redisRepo.Get(ctx, fmt.Sprintf("short_url:%s", shortCode))
		if err == nil && cachedUrl != "" {
			shortUrl, err := s.queryRepo.FindByShortCode(ctx, shortCode)
			if err == nil {
				return shortUrl, nil
			}
		}
	}

	shortUrl, err := s.queryRepo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	if s.redisRepo != nil {
		s.redisRepo.Set(ctx, fmt.Sprintf("short_url:%s", shortCode), shortUrl.LongUrl, time.Hour*24)
	}

	return shortUrl, nil
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
