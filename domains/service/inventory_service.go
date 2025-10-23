package service

import (
	"context"
	"short-url/domains/dto"
	"short-url/domains/dto/inventory"
	"short-url/domains/entities"
	"short-url/domains/values/enums"
)

type InventoryServiceInterface interface {
	CreateInventory(ctx context.Context, req *inventory.CreateInventoryRequest) (*entities.Inventory, error)
	UpdateInventory(ctx context.Context, req *inventory.UpdateInventoryRequest) (*entities.Inventory, error)
	GetInventoryBySKU(ctx context.Context, sku string) (*entities.Inventory, error)
	GetInventoryList(ctx context.Context, pagination dto.Pagination) ([]*entities.Inventory, *dto.PaginationResponse, error)
	GetInventoryByCategory(ctx context.Context, category enums.InventoryCategory, pagination dto.Pagination) ([]*entities.Inventory, *dto.PaginationResponse, error)
	GetInventoryByDistributor(ctx context.Context, distributorID uint, pagination dto.Pagination) ([]*entities.Inventory, *dto.PaginationResponse, error)
	GetLowStockInventory(ctx context.Context, pagination dto.Pagination) ([]*entities.Inventory, *dto.PaginationResponse, error)
}
