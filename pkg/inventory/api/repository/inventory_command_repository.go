package repository

import (
	"context"

	"short-url/domains/entities"
	"short-url/domains/repositories"

	"gorm.io/gorm"
)

type inventoryCommandRepository struct {
	db *gorm.DB
}

func NewInventoryCommandRepository(db *gorm.DB) repositories.InventoryCommandRepositoryInterface {
	return &inventoryCommandRepository{
		db: db,
	}
}

func (r *inventoryCommandRepository) Save(ctx context.Context, inventory *entities.Inventory) error {
	return r.db.WithContext(ctx).Create(inventory).Error
}

func (r *inventoryCommandRepository) Update(ctx context.Context, inventory *entities.Inventory) error {
	return r.db.WithContext(ctx).Save(inventory).Error
}

func (r *inventoryCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Inventory{}, id).Error
}

func (r *inventoryCommandRepository) UpdateQuantity(ctx context.Context, id uint, quantity int) error {
	return r.db.WithContext(ctx).Model(&entities.Inventory{}).Where("id = ?", id).Update("quantity", quantity).Error
}