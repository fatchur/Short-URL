package service

import (
	"context"
	"short-url/domains/dto/inventory"
	"short-url/domains/entities"
)

type InventoryServiceInterface interface {
	CreateInventory(ctx context.Context, req *inventory.CreateInventoryRequest) (*entities.Inventory, error)
	UpdateInventory(ctx context.Context, req *inventory.UpdateInventoryRequest) (*entities.Inventory, error)
	GetInventoryBySKU(ctx context.Context, sku string) (*entities.Inventory, error)
	GetInventoryList(ctx context.Context) ([]*entities.Inventory, error)
}
