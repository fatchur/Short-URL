package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"short-url-service/api/repository"
	"short-url-service/api/service"
	"short-url-service/middleware"
	userrepo "user-service/api/repository"

	"short-url/domains/config"
	"short-url/domains/database"
	"short-url/domains/dto"
	"short-url/domains/repositories"
	jwthelper "short-url/domains/helper/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ShortUrlControllerIntegrationTestSuite struct {
	suite.Suite
	app              *fiber.App
	controller       *ShortUrlController
	ctx              context.Context
	sessionQueryRepo repositories.UserSessionQueryRepositoryInterface
	userQueryRepo    repositories.UserQueryRepositoryInterface
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

	suite.sessionQueryRepo = userrepo.NewUserSessionQueryRepository(db)
	suite.userQueryRepo = userrepo.NewUserQueryRepository(db)

	suite.app.Get("/url/:shortCode", middleware.JWTAuth(suite.sessionQueryRepo), suite.controller.GetLongUrl)

	v1 := suite.app.Group("/api/v1")
	protected := v1.Group("/", middleware.JWTAuth(suite.sessionQueryRepo))
	suite.controller.RegisterRoutes(protected)
}

func (suite *ShortUrlControllerIntegrationTestSuite) TestCreateShortUrl_Success() {
	requestBody := dto.CreateShortUrlRequest{
		LongUrl: "https://example.com/very-long-url-that-needs-shortening",
	}

	userID := suite.getFirstUser()
	token := suite.generateTestJWT(userID, "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234")

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := suite.app.Test(req, 10000000)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode)

	var baseResponse dto.BaseResponse
	err = json.NewDecoder(resp.Body).Decode(&baseResponse)
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), baseResponse.Success)
	assert.Equal(suite.T(), 201, baseResponse.Status)
	assert.Equal(suite.T(), "Short URL created successfully", baseResponse.Message)
	assert.Equal(suite.T(), "v1", baseResponse.APIVersion)
	assert.NotNil(suite.T(), baseResponse.Data)

	dataBytes, _ := json.Marshal(baseResponse.Data)
	var responseData dto.CreateShortUrlResponse
	err = json.Unmarshal(dataBytes, &responseData)
	assert.NoError(suite.T(), err)

	assert.NotZero(suite.T(), responseData.ID)
	assert.Equal(suite.T(), userID, responseData.UserID)
	assert.Equal(suite.T(), "https://example.com/very-long-url-that-needs-shortening", responseData.LongUrl)
	assert.NotEmpty(suite.T(), responseData.ShortCode)
	assert.Len(suite.T(), responseData.ShortCode, 8)
}

func (suite *ShortUrlControllerIntegrationTestSuite) TestCreateShortUrl_InvalidJSON() {
	userID := suite.getFirstUser()
	token := suite.generateTestJWT(userID, "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234")

	req, _ := http.NewRequest("POST", "/api/v1/url", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := suite.app.Test(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 400, resp.StatusCode)

	var baseResponse dto.BaseResponse
	err = json.NewDecoder(resp.Body).Decode(&baseResponse)
	assert.NoError(suite.T(), err)

	assert.False(suite.T(), baseResponse.Success)
	assert.Equal(suite.T(), 400, baseResponse.Status)
	assert.Equal(suite.T(), "Invalid request body", baseResponse.Message)
	assert.Equal(suite.T(), "v1", baseResponse.APIVersion)
	assert.Nil(suite.T(), baseResponse.Data)
}

func (suite *ShortUrlControllerIntegrationTestSuite) TestCreateShortUrl_EmptyLongUrl() {
	requestBody := dto.CreateShortUrlRequest{
		LongUrl: "",
	}

	userID := suite.getFirstUser()
	token := suite.generateTestJWT(userID, "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234")

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := suite.app.Test(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode)

	var baseResponse dto.BaseResponse
	err = json.NewDecoder(resp.Body).Decode(&baseResponse)
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), baseResponse.Success)
	assert.Equal(suite.T(), 201, baseResponse.Status)
	assert.Equal(suite.T(), "Short URL created successfully", baseResponse.Message)
	assert.Equal(suite.T(), "v1", baseResponse.APIVersion)
}

func (suite *ShortUrlControllerIntegrationTestSuite) TestCreateAndGetShortUrl_Integration() {
	userID := suite.getFirstUser()
	token := suite.generateTestJWT(userID, "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234")

	createRequest := dto.CreateShortUrlRequest{
		LongUrl: "https://integration-test.example.com/very-long-url-for-testing",
	}

	createBody, _ := json.Marshal(createRequest)
	createReq, _ := http.NewRequest("POST", "/api/v1/url", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)

	createResp, err := suite.app.Test(createReq, 100000000)
	suite.Require().NoError(err)
	suite.Require().Equal(201, createResp.StatusCode)

	var createBaseResponse dto.BaseResponse
	err = json.NewDecoder(createResp.Body).Decode(&createBaseResponse)
	suite.Require().NoError(err)
	suite.Require().True(createBaseResponse.Success)

	createDataBytes, _ := json.Marshal(createBaseResponse.Data)
	var createResponseData dto.CreateShortUrlResponse
	err = json.Unmarshal(createDataBytes, &createResponseData)
	suite.Require().NoError(err)
	shortCode := createResponseData.ShortCode

	suite.T().Logf("Generated short code: %s", shortCode)

	// get via json
	getJSONReq, _ := http.NewRequest("GET", "/url/"+shortCode, nil)
	getJSONReq.Header.Set("Authorization", "Bearer "+token)
	getJSONReq.Header.Set("Accept", "application/json")

	getJSONResp, err := suite.app.Test(getJSONReq, 100000000)
	suite.Require().NoError(err)
	suite.Require().Equal(200, getJSONResp.StatusCode)

	var getBaseResponse dto.BaseResponse
	err = json.NewDecoder(getJSONResp.Body).Decode(&getBaseResponse)
	suite.Require().NoError(err)

	assert.True(suite.T(), getBaseResponse.Success)
	assert.Equal(suite.T(), 200, getBaseResponse.Status)
	assert.Equal(suite.T(), "Short URL retrieved successfully", getBaseResponse.Message)
	assert.Equal(suite.T(), "v1", getBaseResponse.APIVersion)
	assert.NotNil(suite.T(), getBaseResponse.Data)

	getDataMap, ok := getBaseResponse.Data.(map[string]interface{})
	suite.Require().True(ok)

	assert.Equal(suite.T(), shortCode, getDataMap["short_code"])
	assert.Equal(suite.T(), "https://integration-test.example.com/very-long-url-for-testing", getDataMap["long_url"])
	assert.Equal(suite.T(), float64(userID), getDataMap["user_id"])

	// get via redirect
	getRedirectReq, _ := http.NewRequest("GET", "/url/"+shortCode, nil)
	getRedirectReq.Header.Set("Authorization", "Bearer "+token)

	getRedirectResp, err := suite.app.Test(getRedirectReq, 100000000)
	suite.Require().NoError(err)
	suite.Require().Equal(302, getRedirectResp.StatusCode)

	location := getRedirectResp.Header.Get("Location")
	assert.Equal(suite.T(), "https://integration-test.example.com/very-long-url-for-testing", location)
}

func (suite *ShortUrlControllerIntegrationTestSuite) generateTestJWT(userID uint, sessionCode string) string {
	secretKey := suite.getSessionSecret(sessionCode)
	tokenString, _, _ := jwthelper.GenerateJWTToken(userID, sessionCode, secretKey)
	return tokenString
}

func (suite *ShortUrlControllerIntegrationTestSuite) getSessionSecret(sessionCode string) string {
	session, err := suite.sessionQueryRepo.FindBySessionCode(context.Background(), sessionCode)
	if err != nil {
		suite.T().Fatalf("Failed to find session for code %s: %v", sessionCode, err)
	}
	return session.SecretKey
}

func (suite *ShortUrlControllerIntegrationTestSuite) getFirstUser() uint {
	user, err := suite.userQueryRepo.FindByEmail(context.Background(), "john@example.com")
	if err != nil {
		suite.T().Fatalf("Failed to find first user: %v", err)
	}
	return user.ID
}

func TestShortUrlControllerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ShortUrlControllerIntegrationTestSuite))
}
