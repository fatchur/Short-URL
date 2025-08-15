package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"short-url/domains/dto"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DefaultConfig() dto.DBConfig {
	return dto.DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "",
		DBName:   "short_url",
		SSLMode:  "disable",
		Timezone: "UTC",
		LogLevel: "warn",
	}
}

func DBConnect(ctx context.Context, config dto.DBConfig) (*gorm.DB, error) {

	dsn := config.DSN
	if dsn == "" {
		host := defaultIfEmpty(config.Host, "localhost")
		port := defaultIfEmpty(config.Port, "5432")
		user := defaultIfEmpty(config.User, "postgres")
		password := defaultIfEmpty(config.Password, "")
		dbname := defaultIfEmpty(config.DBName, "short_url")
		sslmode := defaultIfEmpty(config.SSLMode, "disable")
		timezone := defaultIfEmpty(config.Timezone, "UTC")

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			host, user, password, dbname, port, sslmode, timezone)
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid database config: %w", err)
	}

	logLevel := parseLogLevel(defaultIfEmpty(config.LogLevel, "warn"))
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logLevel,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func validateConfig(config dto.DBConfig) error {
	if config.DSN == "" {
		if config.Host == "" {
			return fmt.Errorf("host cannot be empty")
		}
		if config.Port == "" {
			return fmt.Errorf("port cannot be empty")
		}
		if config.DBName == "" {
			return fmt.Errorf("database name cannot be empty")
		}
	}
	return nil
}

func defaultIfEmpty(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn", "warning":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}

func MustConnect(ctx context.Context, config dto.DBConfig) *gorm.DB {
	db, err := DBConnect(ctx, config)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	return db
}

func TestConnection(ctx context.Context, config dto.DBConfig) error {
	db, err := DBConnect(ctx, config)
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	defer sqlDB.Close()

	return sqlDB.Ping()
}
