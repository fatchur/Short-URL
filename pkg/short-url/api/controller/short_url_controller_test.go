package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"short-url-service/api/repository"
	"short-url-service/api/service"

	"short-url/domains/config"
	"short-url/domains/database"
	"short-url/domains/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ShortUrlControllerIntegrationTestSuite struct {
	suite.Suite
	app        *fiber.App
	controller *ShortUrlController
	ctx        context.Context
}

func (suite *ShortUrlControllerIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Load test configuration
	cfg := config.LoadConfig()

	dbConfig := dto.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
		Timezone: cfg.DBTimezone,
		LogLevel: cfg.DBLogLevel,
	}

	cacheConfig := dto.CacheConfig{
		Host:     cfg.DBHost,
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	// Clear tables before testing
	err := database.ClearTables(suite.ctx, dbConfig)
	suite.Require().NoError(err)

	// Setup database connections
	db, err := database.DBConnect(suite.ctx, dbConfig)
	suite.Require().NoError(err)

	// Seed tables with test data
	err = database.Seed(db)
	suite.Require().NoError(err)

	redisClient, err := database.CacheConnect(suite.ctx, cacheConfig)
	suite.Require().NoError(err)

	// Setup repositories and service
	commandRepo := repository.NewShortUrlCommandRepository(db)
	queryRepo := repository.NewShortUrlQueryRepository(db)
	redisRepo := repository.NewRedisRepository(redisClient)

	shortUrlService := service.NewShortUrlService(commandRepo, queryRepo, redisRepo)
	suite.controller = NewShortUrlController(shortUrlService)

	// Setup Fiber app
	suite.app = fiber.New()
	v1 := suite.app.Group("/api/v1")
	suite.controller.RegisterRoutes(v1)
}

func (suite *ShortUrlControllerIntegrationTestSuite) TestCreateShortUrl_Success() {
	// Arrange
	requestBody := dto.CreateShortUrlRequest{
		LongUrl: "https://example.com/very-long-url-that-needs-shortening",
		UserID:  1,
	}

	// Act
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode)

	var responseBody dto.CreateShortUrlResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(suite.T(), err)

	assert.NotZero(suite.T(), responseBody.ID)
	assert.Equal(suite.T(), uint(1), responseBody.UserID)
	assert.Equal(suite.T(), "https://example.com/very-long-url-that-needs-shortening", responseBody.LongUrl)
	assert.NotEmpty(suite.T(), responseBody.ShortCode)
	assert.Len(suite.T(), responseBody.ShortCode, 8)
}

func (suite *ShortUrlControllerIntegrationTestSuite) TestCreateShortUrl_InvalidJSON() {
	// Act
	req, _ := http.NewRequest("POST", "/api/v1/url", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 400, resp.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid request body", responseBody["error"])
}

func (suite *ShortUrlControllerIntegrationTestSuite) TestCreateShortUrl_EmptyLongUrl() {
	// Arrange
	requestBody := dto.CreateShortUrlRequest{
		LongUrl: "",
		UserID:  1,
	}

	// Act
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)

	// Assert
	assert.NoError(suite.T(), err)
	// Should still process the request even with empty URL for now
	// This test documents current behavior - validation can be added later
	assert.Equal(suite.T(), 201, resp.StatusCode)
}

func TestShortUrlControllerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ShortUrlControllerIntegrationTestSuite))
}
