package repository

import (
	"context"

	"short-url/domains/entities"
	"short-url/domains/repositories"
	"short-url/domains/values/enums"

	"gorm.io/gorm"
)

type inventoryQueryRepository struct {
	db *gorm.DB
}

func NewInventoryQueryRepository(db *gorm.DB) repositories.InventoryQueryRepositoryInterface {
	return &inventoryQueryRepository{
		db: db,
	}
}

func (r *inventoryQueryRepository) FindByID(ctx context.Context, id uint) (*entities.Inventory, error) {
	var inventory entities.Inventory
	err := r.db.WithContext(ctx).Preload("Distributor").First(&inventory, id).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryQueryRepository) FindBySKU(ctx context.Context, sku string) (*entities.Inventory, error) {
	var inventory entities.Inventory
	err := r.db.WithContext(ctx).Preload("Distributor").Where("sku = ?", sku).First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryQueryRepository) FindByCategory(ctx context.Context, category enums.InventoryCategory) ([]*entities.Inventory, error) {
	var inventories []*entities.Inventory
	err := r.db.WithContext(ctx).Preload("Distributor").Where("category_id = ?", category).Find(&inventories).Error
	return inventories, err
}

func (r *inventoryQueryRepository) FindByDistributor(ctx context.Context, distributorID uint) ([]*entities.Inventory, error) {
	var inventories []*entities.Inventory
	err := r.db.WithContext(ctx).Preload("Distributor").Where("distributor_id = ?", distributorID).Find(&inventories).Error
	return inventories, err
}

func (r *inventoryQueryRepository) FindAll(ctx context.Context) ([]*entities.Inventory, error) {
	var inventories []*entities.Inventory
	err := r.db.WithContext(ctx).Preload("Distributor").Find(&inventories).Error
	return inventories, err
}

func (r *inventoryQueryRepository) FindLowStock(ctx context.Context) ([]*entities.Inventory, error) {
	var inventories []*entities.Inventory
	err := r.db.WithContext(ctx).Preload("Distributor").Where("quantity <= min_quantity AND min_quantity IS NOT NULL").Find(&inventories).Error
	return inventories, err
}