package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

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
	JWTSecret  string
	Port       string
}

func LoadConfig() *Config {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	envPath := filepath.Join(dir, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: Could not load .env file from %s: %v", envPath, err)
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
		JWTSecret:  getRequiredEnv("JWT_SECRET"),
		Port:       getEnvWithDefault("PORT", "8080"),
	}

	log.Println("Configuration loaded successfully")
	return config
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("Required environment variable %s is not set", key)
	}
	return value
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
