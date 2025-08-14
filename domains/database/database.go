package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	DSN      string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Timezone string
	LogLevel string
}

var (
	currentConfig DBConfig
	configMutex   sync.RWMutex
)

// SetConfig safely sets the database configuration
func SetConfig(cfg DBConfig) {
	configMutex.Lock()
	defer configMutex.Unlock()
	currentConfig = cfg
}

// GetConfig safely gets the current database configuration
func GetConfig() DBConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return currentConfig
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() DBConfig {
	return DBConfig{
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

// DBConnect establishes a connection to the PostgreSQL database
func DBConnect() (*gorm.DB, error) {
	config := GetConfig()

	dsn := config.DSN
	if dsn == "" {
		// Build DSN from individual components
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

	// Validate required fields
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

// validateConfig performs basic validation on the database configuration
func validateConfig(config DBConfig) error {
	if config.DSN == "" {
		// If DSN is not provided, check individual components
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

// defaultIfEmpty returns the fallback value if the given value is empty
func defaultIfEmpty(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

// parseLogLevel converts a string log level to gorm's LogLevel type
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

// MustConnect is a convenience function that panics if connection fails
// Use this only during application startup
func MustConnect() *gorm.DB {
	db, err := DBConnect()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	return db
}

// TestConnection attempts to ping the database to verify connectivity
func TestConnection() error {
	db, err := DBConnect()
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
