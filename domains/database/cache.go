package database

import (
	"context"
	"fmt"

	"short-url/domains/dto"

	"github.com/redis/go-redis/v9"
)


func DefaultCacheConfig() dto.CacheConfig {
	return dto.CacheConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}
}

func CacheConnect(ctx context.Context, config dto.CacheConfig) (*redis.Client, error) {
	host := defaultIfEmpty(config.Host, "localhost")
	port := defaultIfEmpty(config.Port, "6379")
	addr := fmt.Sprintf("%s:%s", host, port)

	if err := validateCacheConfig(config); err != nil {
		return nil, fmt.Errorf("invalid cache config: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.DB,
	})

	return rdb, nil
}

func validateCacheConfig(config dto.CacheConfig) error {
	if config.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if config.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	return nil
}

func MustCacheConnect(ctx context.Context, config dto.CacheConfig) *redis.Client {
	rdb, err := CacheConnect(ctx, config)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to cache: %v", err))
	}
	return rdb
}

func TestCacheConnection(ctx context.Context, config dto.CacheConfig) error {
	rdb, err := CacheConnect(ctx, config)
	if err != nil {
		return err
	}
	defer rdb.Close()

	return rdb.Ping(ctx).Err()
}
