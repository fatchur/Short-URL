package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"short-url-service/api/repository"
	"short-url-service/api/service"
	"short-url-service/middleware"
	userrepo "user-service/api/repository"

	"short-url/domains/config"
	"short-url/domains/database"
	"short-url/domains/dto"
	"short-url/domains/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ShortUrlControllerIntegrationTestSuite struct {
	suite.Suite
	app              *fiber.App
	controller       *ShortUrlController
	ctx              context.Context
	sessionQueryRepo repositories.UserSessionQueryRepositoryInterface
}

func (suite *ShortUrlControllerIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()

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

	err := database.ClearTables(suite.ctx, dbConfig)
	suite.Require().NoError(err)

	db, err := database.DBConnect(suite.ctx, dbConfig)
	suite.Require().NoError(err)

	err = database.Seed(db)
	suite.Require().NoError(err)

	redisClient, err := database.CacheConnect(suite.ctx, cacheConfig)
	suite.Require().NoError(err)

	commandRepo := repository.NewShortUrlCommandRepository(db)
	queryRepo := repository.NewShortUrlQueryRepository(db)
	redisRepo := repository.NewRedisRepository(redisClient)

	shortUrlService := service.NewShortUrlService(commandRepo, queryRepo, redisRepo)
	suite.controller = NewShortUrlController(shortUrlService)

	suite.app = fiber.New()
	v1 := suite.app.Group("/api/v1")

	suite.sessionQueryRepo = userrepo.NewUserSessionQueryRepository(db)
	protected := v1.Group("/", middleware.JWTAuth(suite.sessionQueryRepo))
	suite.controller.RegisterRoutes(protected)
}

func (suite *ShortUrlControllerIntegrationTestSuite) TestCreateShortUrl_Success() {
	requestBody := dto.CreateShortUrlRequest{
		LongUrl: "https://example.com/very-long-url-that-needs-shortening",
	}

	token := suite.generateTestJWT(1, "john@example.com", "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234")

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := suite.app.Test(req, 10000000)

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
	token := suite.generateTestJWT(1, "john@example.com", "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234")
	
	req, _ := http.NewRequest("POST", "/api/v1/url", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := suite.app.Test(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 400, resp.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid request body", responseBody["error"])
}

func (suite *ShortUrlControllerIntegrationTestSuite) TestCreateShortUrl_EmptyLongUrl() {
	requestBody := dto.CreateShortUrlRequest{
		LongUrl: "",
	}

	token := suite.generateTestJWT(1, "john@example.com", "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234")

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := suite.app.Test(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode)
}

func (suite *ShortUrlControllerIntegrationTestSuite) generateTestJWT(userID uint, email, sessionToken string) string {
	claims := &middleware.Claims{
		UserID:       userID,
		Email:        email,
		SessionToken: sessionToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	secretKey := suite.getSessionSecret(sessionToken)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}

func (suite *ShortUrlControllerIntegrationTestSuite) getSessionSecret(sessionToken string) string {
	session, err := suite.sessionQueryRepo.FindBySessionToken(context.Background(), sessionToken)
	if err != nil {
		suite.T().Fatalf("Failed to find session for token %s: %v", sessionToken, err)
	}
	return session.SecretKey
}

func TestShortUrlControllerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ShortUrlControllerIntegrationTestSuite))
}
