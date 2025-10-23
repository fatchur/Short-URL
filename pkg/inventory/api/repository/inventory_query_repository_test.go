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

type InventoryQueryRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *inventoryQueryRepository
	ctx  context.Context
}

func (suite *InventoryQueryRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)
	
	err = db.AutoMigrate(&entities.Inventory{}, &entities.Distributor{})
	suite.Require().NoError(err)
	
	suite.db = db
	suite.repo = &inventoryQueryRepository{db: db}
}

func (suite *InventoryQueryRepositoryTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM inventories")
	suite.db.Exec("DELETE FROM distributors")
}

func (suite *InventoryQueryRepositoryTestSuite) TestFindByID() {
	inventory := &entities.Inventory{
		Name:      "Test Product",
		SKU:       "TEST-001",
		Quantity:  100,
		UnitPrice: 999.99,
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err := suite.db.Create(inventory).Error
	suite.Require().NoError(err)

	result, err := suite.repo.FindByID(suite.ctx, inventory.ID)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), inventory.ID, result.ID)
	assert.Equal(suite.T(), inventory.Name, result.Name)
	assert.Equal(suite.T(), inventory.SKU, result.SKU)
}

func (suite *InventoryQueryRepositoryTestSuite) TestFindBySKU() {
	inventory := &entities.Inventory{
		Name:      "Test Product",
		SKU:       "UNIQUE-SKU-001",
		Quantity:  50,
		UnitPrice: 500.00,
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err := suite.db.Create(inventory).Error
	suite.Require().NoError(err)

	result, err := suite.repo.FindBySKU(suite.ctx, "UNIQUE-SKU-001")
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), inventory.ID, result.ID)
	assert.Equal(suite.T(), "UNIQUE-SKU-001", result.SKU)
}

func (suite *InventoryQueryRepositoryTestSuite) TestFindByCategory() {
	inventory1 := &entities.Inventory{
		Name:       "Electronics Product",
		SKU:        "ELEC-001",
		CategoryID: helper.InventoryCategoryPtr(enums.Electronics),
		Quantity:   10,
		UnitPrice:  100.00,
		CreatedAt:  time.Now(),
		CreatedBy:  1,
		UpdatedAt:  time.Now(),
	}

	inventory2 := &entities.Inventory{
		Name:       "Clothing Product",
		SKU:        "CLOTH-001",
		CategoryID: helper.InventoryCategoryPtr(enums.Clothing),
		Quantity:   20,
		UnitPrice:  50.00,
		CreatedAt:  time.Now(),
		CreatedBy:  1,
		UpdatedAt:  time.Now(),
	}

	err := suite.db.Create(inventory1).Error
	suite.Require().NoError(err)
	err = suite.db.Create(inventory2).Error
	suite.Require().NoError(err)

	results, err := suite.repo.FindByCategory(suite.ctx, enums.Electronics)
	
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 1)
	assert.Equal(suite.T(), "Electronics Product", results[0].Name)
}

func (suite *InventoryQueryRepositoryTestSuite) TestFindByDistributor() {
	distributor := &entities.Distributor{
		Name:      "Test Distributor",
		Email:     "test@distributor.com",
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err := suite.db.Create(distributor).Error
	suite.Require().NoError(err)

	inventory1 := &entities.Inventory{
		DistributorID: &distributor.ID,
		Name:          "Product 1",
		SKU:           "DIST-001",
		Quantity:      10,
		UnitPrice:     100.00,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
	}

	inventory2 := &entities.Inventory{
		Name:      "Product 2",
		SKU:       "NO-DIST-001",
		Quantity:  20,
		UnitPrice: 200.00,
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err = suite.db.Create(inventory1).Error
	suite.Require().NoError(err)
	err = suite.db.Create(inventory2).Error
	suite.Require().NoError(err)

	results, err := suite.repo.FindByDistributor(suite.ctx, distributor.ID)
	
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 1)
	assert.Equal(suite.T(), "Product 1", results[0].Name)
}

func (suite *InventoryQueryRepositoryTestSuite) TestFindAll() {
	inventory1 := &entities.Inventory{
		Name:      "Product 1",
		SKU:       "ALL-001",
		Quantity:  10,
		UnitPrice: 100.00,
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	inventory2 := &entities.Inventory{
		Name:      "Product 2",
		SKU:       "ALL-002",
		Quantity:  20,
		UnitPrice: 200.00,
		CreatedAt: time.Now(),
		CreatedBy: 1,
		UpdatedAt: time.Now(),
	}

	err := suite.db.Create(inventory1).Error
	suite.Require().NoError(err)
	err = suite.db.Create(inventory2).Error
	suite.Require().NoError(err)

	results, err := suite.repo.FindAll(suite.ctx)
	
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 2)
}

func (suite *InventoryQueryRepositoryTestSuite) TestFindLowStock() {
	inventory1 := &entities.Inventory{
		Name:        "Low Stock Product",
		SKU:         "LOW-001",
		Quantity:    5,
		MinQuantity: helper.IntPtr(10),
		UnitPrice:   100.00,
		CreatedAt:   time.Now(),
		CreatedBy:   1,
		UpdatedAt:   time.Now(),
	}

	inventory2 := &entities.Inventory{
		Name:        "Normal Stock Product",
		SKU:         "NORMAL-001",
		Quantity:    50,
		MinQuantity: helper.IntPtr(10),
		UnitPrice:   200.00,
		CreatedAt:   time.Now(),
		CreatedBy:   1,
		UpdatedAt:   time.Now(),
	}

	err := suite.db.Create(inventory1).Error
	suite.Require().NoError(err)
	err = suite.db.Create(inventory2).Error
	suite.Require().NoError(err)

	results, err := suite.repo.FindLowStock(suite.ctx)
	
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 1)
	assert.Equal(suite.T(), "Low Stock Product", results[0].Name)
}

func TestInventoryQueryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryQueryRepositoryTestSuite))
}