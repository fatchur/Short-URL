package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"user-service/api/repository"
	"user-service/api/service"
	"user-service/middleware"

	"short-url/domains/config"
	"short-url/domains/database"
	"short-url/domains/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserControllerIntegrationTestSuite struct {
	suite.Suite
	app        *fiber.App
	controller *UserController
	ctx        context.Context
}

func (suite *UserControllerIntegrationTestSuite) SetupSuite() {
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

	err := database.ClearTables(suite.ctx, dbConfig)
	suite.Require().NoError(err)

	db, err := database.DBConnect(suite.ctx, dbConfig)
	suite.Require().NoError(err)

	err = database.Seed(db)
	suite.Require().NoError(err)

	commandRepo := repository.NewUserSessionCommandRepository(db)
	queryRepo := repository.NewUserSessionQueryRepository(db)

	userSessionService := service.NewUserSessionService(commandRepo, queryRepo)
	suite.controller = NewUserController(userSessionService)

	suite.app = fiber.New()

	suite.app.Use(func(c *fiber.Ctx) error {
		middleware.SetUserIDToContext(c, 1)
		return c.Next()
	})

	v1 := suite.app.Group("/api/v1")
	user := v1.Group("/user")
	suite.controller.RegisterRoutes(user)
}

func (suite *UserControllerIntegrationTestSuite) TestCreateSession_Success() {
	requestBody := dto.CreateSessionRequest{
		DeviceInfo: "Test Device",
		IPAddress:  "127.0.0.1",
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/user/session", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req, 10000000)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode)

	var baseResponse dto.BaseResponse
	err = json.NewDecoder(resp.Body).Decode(&baseResponse)
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), baseResponse.Success)
	assert.Equal(suite.T(), 201, baseResponse.Status)
	assert.Equal(suite.T(), "Session created successfully", baseResponse.Message)
	assert.Equal(suite.T(), "v1", baseResponse.APIVersion)
	assert.NotNil(suite.T(), baseResponse.Data)

	dataBytes, _ := json.Marshal(baseResponse.Data)
	var responseData dto.CreateSessionResponse
	err = json.Unmarshal(dataBytes, &responseData)
	assert.NoError(suite.T(), err)

	assert.NotEmpty(suite.T(), responseData.SessionToken)
	assert.True(suite.T(), responseData.ExpiresAt.After(time.Now()))
}

func (suite *UserControllerIntegrationTestSuite) TestCreateSession_InvalidJSON() {
	req, _ := http.NewRequest("POST", "/api/v1/user/session", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

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

func (suite *UserControllerIntegrationTestSuite) TestCreateSession_NoUserContext() {
	app := fiber.New()
	v1 := app.Group("/api/v1")
	user := v1.Group("/user")
	suite.controller.RegisterRoutes(user)

	requestBody := dto.CreateSessionRequest{
		DeviceInfo: "Test Device",
		IPAddress:  "127.0.0.1",
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/user/session", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 401, resp.StatusCode)

	var baseResponse dto.BaseResponse
	err = json.NewDecoder(resp.Body).Decode(&baseResponse)
	assert.NoError(suite.T(), err)

	assert.False(suite.T(), baseResponse.Success)
	assert.Equal(suite.T(), 401, baseResponse.Status)
	assert.Equal(suite.T(), "User authentication required", baseResponse.Message)
	assert.Equal(suite.T(), "v1", baseResponse.APIVersion)
	assert.Nil(suite.T(), baseResponse.Data)
}

func TestUserControllerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerIntegrationTestSuite))
}
