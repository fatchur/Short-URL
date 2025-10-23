package repository

import (
	"context"
	"testing"
	"time"

	"short-url/domains/entities"
	"short-url/domains/helper"
	"short-url/domains/values/enums"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type InventoryCommandRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *inventoryCommandRepository
	ctx  context.Context
}

func (suite *InventoryCommandRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)
	
	err = db.AutoMigrate(&entities.Inventory{}, &entities.Distributor{})
	suite.Require().NoError(err)
	
	suite.db = db
	suite.repo = &inventoryCommandRepository{db: db}
}

func (suite *InventoryCommandRepositoryTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM inventories")
	suite.db.Exec("DELETE FROM distributors")
}

func (suite *InventoryCommandRepositoryTestSuite) TestSave() {
	inventory := &entities.Inventory{
		DistributorID: helper.UintPtr(1),
		Name:          "Test Product",
		Description:   helper.StringPtr("Test Description"),
		SKU:           "TEST-001",
		CategoryID:    helper.InventoryCategoryPtr(enums.Electronics),
		Quantity:      100,
		MinQuantity:   helper.IntPtr(10),
		UnitPrice:     999.99,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
	}

	err := suite.repo.Save(suite.ctx, inventory)
	
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), inventory.ID)
}

func (suite *InventoryCommandRepositoryTestSuite) TestUpdate() {
	inventory := &entities.Inventory{
		Name:      "Original Product",
		SKU:       "TEST-002",
		Quantity:  50,
		UnitPrice: 500.00,
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err := suite.repo.Save(suite.ctx, inventory)
	suite.Require().NoError(err)

	inventory.Name = "Updated Product"
	inventory.Quantity = 75
	inventory.UpdatedAt = time.Now()

	err = suite.repo.Update(suite.ctx, inventory)
	
	assert.NoError(suite.T(), err)
	
	var updated entities.Inventory
	err = suite.db.First(&updated, inventory.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Product", updated.Name)
	assert.Equal(suite.T(), 75, updated.Quantity)
}

func (suite *InventoryCommandRepositoryTestSuite) TestDelete() {
	inventory := &entities.Inventory{
		Name:      "Test Product",
		SKU:       "TEST-003",
		Quantity:  25,
		UnitPrice: 250.00,
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err := suite.repo.Save(suite.ctx, inventory)
	suite.Require().NoError(err)

	err = suite.repo.Delete(suite.ctx, inventory.ID)
	
	assert.NoError(suite.T(), err)
	
	var count int64
	err = suite.db.Model(&entities.Inventory{}).Where("id = ?", inventory.ID).Count(&count).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *InventoryCommandRepositoryTestSuite) TestUpdateQuantity() {
	inventory := &entities.Inventory{
		Name:      "Test Product",
		SKU:       "TEST-004",
		Quantity:  100,
		UnitPrice: 100.00,
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err := suite.repo.Save(suite.ctx, inventory)
	suite.Require().NoError(err)

	newQuantity := 150
	err = suite.repo.UpdateQuantity(suite.ctx, inventory.ID, newQuantity)
	
	assert.NoError(suite.T(), err)
	
	var updated entities.Inventory
	err = suite.db.First(&updated, inventory.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newQuantity, updated.Quantity)
}

func TestInventoryCommandRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryCommandRepositoryTestSuite))
}