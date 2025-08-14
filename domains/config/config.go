package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBTimezone string
	DBLogLevel string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file from domains/config directory
	envPath := filepath.Join("domains", "config", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: Could not load .env file from %s: %v", envPath, err)
	} else {
		log.Printf("Loaded environment variables from %s", envPath)
	}

	config := &Config{
		DBHost:     getRequiredEnv("DB_HOST"),
		DBPort:     getEnvWithDefault("DB_PORT", "5432"),
		DBUser:     getRequiredEnv("DB_USER"),
		DBPassword: getRequiredEnv("DB_PASSWORD"),
		DBName:     getRequiredEnv("DB_NAME"),
		DBSSLMode:  getEnvWithDefault("DB_SSLMODE", "disable"),
		DBTimezone: getEnvWithDefault("DB_TIMEZONE", "UTC"),
		DBLogLevel: getEnvWithDefault("DB_LOG_LEVEL", "warn"),
	}

	log.Println("Configuration loaded successfully")
	return config
}

// getRequiredEnv gets a required environment variable, panics if empty
func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("Required environment variable %s is not set", key)
	}
	return value
}

// getEnvWithDefault gets an environment variable with a fallback default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
