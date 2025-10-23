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

type DistributorCommandRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *distributorCommandRepository
	ctx  context.Context
}

func (suite *DistributorCommandRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	err = db.AutoMigrate(&entities.Distributor{})
	suite.Require().NoError(err)

	suite.db = db
	suite.repo = &distributorCommandRepository{db: db}
}

func (suite *DistributorCommandRepositoryTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM distributors")
}

func (suite *DistributorCommandRepositoryTestSuite) TestSave() {
	distributor := &entities.Distributor{
		Name:        "Test Distributor",
		Email:       "test@distributor.com",
		PhoneNumber: helper.StringPtr("123-456-7890"),
		Address:     helper.StringPtr("123 Test Street"),
		CreatedAt:   time.Now(),
		CreatedBy:   1,
		UpdatedAt:   time.Now(),
	}

	err := suite.repo.Save(suite.ctx, distributor)

	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), distributor.ID)
}

func (suite *DistributorCommandRepositoryTestSuite) TestUpdate() {
	distributor := &entities.Distributor{
		Name:      "Original Distributor",
		Email:     "original@distributor.com",
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err := suite.repo.Save(suite.ctx, distributor)
	suite.Require().NoError(err)

	distributor.Name = "Updated Distributor"
	distributor.Email = "updated@distributor.com"
	distributor.UpdatedAt = time.Now()

	err = suite.repo.Update(suite.ctx, distributor)

	assert.NoError(suite.T(), err)

	var updated entities.Distributor
	err = suite.db.First(&updated, distributor.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Distributor", updated.Name)
	assert.Equal(suite.T(), "updated@distributor.com", updated.Email)
}

func (suite *DistributorCommandRepositoryTestSuite) TestDelete() {
	distributor := &entities.Distributor{
		Name:      "Test Distributor",
		Email:     "delete@distributor.com",
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err := suite.repo.Save(suite.ctx, distributor)
	suite.Require().NoError(err)

	err = suite.repo.Delete(suite.ctx, distributor.ID)

	assert.NoError(suite.T(), err)

	var count int64
	err = suite.db.Model(&entities.Distributor{}).Where("id = ?", distributor.ID).Count(&count).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *DistributorCommandRepositoryTestSuite) TestSaveWithDuplicateEmail() {
	distributor1 := &entities.Distributor{
		Name:      "Distributor 1",
		Email:     "same@email.com",
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	distributor2 := &entities.Distributor{
		Name:      "Distributor 2",
		Email:     "same@email.com",
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err := suite.repo.Save(suite.ctx, distributor1)
	suite.Require().NoError(err)

	err = suite.repo.Save(suite.ctx, distributor2)

	assert.Error(suite.T(), err)
}

func TestDistributorCommandRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DistributorCommandRepositoryTestSuite))
}
