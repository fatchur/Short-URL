package repository

import (
	"context"

	"short-url/domains/dto"
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

func (r *inventoryQueryRepository) FindByCategory(ctx context.Context, category enums.InventoryCategory, pagination dto.Pagination) ([]*entities.Inventory, *dto.PaginationResponse, error) {
	pagination.SetDefaults()
	
	var inventories []*entities.Inventory
	var total int64
	
	query := r.db.WithContext(ctx).Model(&entities.Inventory{}).Where("category_id = ?", category)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}
	
	err := query.Preload("Distributor").
		Offset(pagination.GetOffset()).
		Limit(pagination.PageSize).
		Find(&inventories).Error
	
	if err != nil {
		return nil, nil, err
	}
	
	paginationResponse := dto.NewPaginationResponse(pagination.Page, pagination.PageSize, total)
	return inventories, paginationResponse, nil
}

func (r *inventoryQueryRepository) FindByDistributor(ctx context.Context, distributorID uint, pagination dto.Pagination) ([]*entities.Inventory, *dto.PaginationResponse, error) {
	pagination.SetDefaults()
	
	var inventories []*entities.Inventory
	var total int64
	
	query := r.db.WithContext(ctx).Model(&entities.Inventory{}).Where("distributor_id = ?", distributorID)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}
	
	err := query.Preload("Distributor").
		Offset(pagination.GetOffset()).
		Limit(pagination.PageSize).
		Find(&inventories).Error
	
	if err != nil {
		return nil, nil, err
	}
	
	paginationResponse := dto.NewPaginationResponse(pagination.Page, pagination.PageSize, total)
	return inventories, paginationResponse, nil
}

func (r *inventoryQueryRepository) FindAll(ctx context.Context, pagination dto.Pagination) ([]*entities.Inventory, *dto.PaginationResponse, error) {
	pagination.SetDefaults()
	
	var inventories []*entities.Inventory
	var total int64
	
	query := r.db.WithContext(ctx).Model(&entities.Inventory{})
	
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}
	
	err := query.Preload("Distributor").
		Offset(pagination.GetOffset()).
		Limit(pagination.PageSize).
		Find(&inventories).Error
	
	if err != nil {
		return nil, nil, err
	}
	
	paginationResponse := dto.NewPaginationResponse(pagination.Page, pagination.PageSize, total)
	return inventories, paginationResponse, nil
}

func (r *inventoryQueryRepository) FindLowStock(ctx context.Context, pagination dto.Pagination) ([]*entities.Inventory, *dto.PaginationResponse, error) {
	pagination.SetDefaults()
	
	var inventories []*entities.Inventory
	var total int64
	
	query := r.db.WithContext(ctx).Model(&entities.Inventory{}).Where("quantity <= min_quantity AND min_quantity IS NOT NULL")
	
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}
	
	err := query.Preload("Distributor").
		Offset(pagination.GetOffset()).
		Limit(pagination.PageSize).
		Find(&inventories).Error
	
	if err != nil {
		return nil, nil, err
	}
	
	paginationResponse := dto.NewPaginationResponse(pagination.Page, pagination.PageSize, total)
	return inventories, paginationResponse, nil
}