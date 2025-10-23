package repository

import (
	"context"
	"testing"
	"time"

	"short-url/domains/entities"
	"short-url/domains/helper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DistributorQueryRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *distributorQueryRepository
	ctx  context.Context
}

func (suite *DistributorQueryRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)
	
	err = db.AutoMigrate(&entities.Distributor{})
	suite.Require().NoError(err)
	
	suite.db = db
	suite.repo = &distributorQueryRepository{db: db}
}

func (suite *DistributorQueryRepositoryTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM distributors")
}

func (suite *DistributorQueryRepositoryTestSuite) TestFindByID() {
	distributor := &entities.Distributor{
		Name:        "Test Distributor",
		Email:       "test@distributor.com",
		PhoneNumber: helper.StringPtr("123-456-7890"),
		Address:     helper.StringPtr("123 Test Street"),
		CreatedAt:   time.Now(),
		CreatedBy:   1,
		UpdatedAt:   time.Now(),
	}

	err := suite.db.Create(distributor).Error
	suite.Require().NoError(err)

	result, err := suite.repo.FindByID(suite.ctx, distributor.ID)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), distributor.ID, result.ID)
	assert.Equal(suite.T(), distributor.Name, result.Name)
	assert.Equal(suite.T(), distributor.Email, result.Email)
}

func (suite *DistributorQueryRepositoryTestSuite) TestFindByIDNotFound() {
	result, err := suite.repo.FindByID(suite.ctx, 999)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "record not found")
}

func (suite *DistributorQueryRepositoryTestSuite) TestFindByEmail() {
	distributor := &entities.Distributor{
		Name:      "Email Test Distributor",
		Email:     "unique@email.com",
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err := suite.db.Create(distributor).Error
	suite.Require().NoError(err)

	result, err := suite.repo.FindByEmail(suite.ctx, "unique@email.com")
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), distributor.ID, result.ID)
	assert.Equal(suite.T(), "unique@email.com", result.Email)
}

func (suite *DistributorQueryRepositoryTestSuite) TestFindByEmailNotFound() {
	result, err := suite.repo.FindByEmail(suite.ctx, "nonexistent@email.com")
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "record not found")
}

func (suite *DistributorQueryRepositoryTestSuite) TestFindAll() {
	distributor1 := &entities.Distributor{
		Name:      "Distributor 1",
		Email:     "dist1@email.com",
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	distributor2 := &entities.Distributor{
		Name:      "Distributor 2",
		Email:     "dist2@email.com",
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	distributor3 := &entities.Distributor{
		Name:      "Distributor 3",
		Email:     "dist3@email.com",
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err := suite.db.Create(distributor1).Error
	suite.Require().NoError(err)
	err = suite.db.Create(distributor2).Error
	suite.Require().NoError(err)
	err = suite.db.Create(distributor3).Error
	suite.Require().NoError(err)

	results, err := suite.repo.FindAll(suite.ctx)
	
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 3)
	
	names := make([]string, len(results))
	for i, dist := range results {
		names[i] = dist.Name
	}
	assert.Contains(suite.T(), names, "Distributor 1")
	assert.Contains(suite.T(), names, "Distributor 2")
	assert.Contains(suite.T(), names, "Distributor 3")
}

func (suite *DistributorQueryRepositoryTestSuite) TestFindAllEmpty() {
	results, err := suite.repo.FindAll(suite.ctx)
	
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 0)
}

func TestDistributorQueryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DistributorQueryRepositoryTestSuite))
}