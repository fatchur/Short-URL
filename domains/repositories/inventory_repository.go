package repositories

import (
	"context"
	"short-url/domains/dto"
	"short-url/domains/entities"
	"short-url/domains/values/enums"
)

type InventoryQueryRepositoryInterface interface {
	FindByID(ctx context.Context, id uint) (*entities.Inventory, error)
	FindBySKU(ctx context.Context, sku string) (*entities.Inventory, error)
	FindByCategory(ctx context.Context, category enums.InventoryCategory, pagination dto.Pagination) ([]*entities.Inventory, *dto.PaginationResponse, error)
	FindByDistributor(ctx context.Context, distributorID uint, pagination dto.Pagination) ([]*entities.Inventory, *dto.PaginationResponse, error)
	FindAll(ctx context.Context, pagination dto.Pagination) ([]*entities.Inventory, *dto.PaginationResponse, error)
	FindLowStock(ctx context.Context, pagination dto.Pagination) ([]*entities.Inventory, *dto.PaginationResponse, error)
}

type InventoryCommandRepositoryInterface interface {
	Save(ctx context.Context, inventory *entities.Inventory) error
	Update(ctx context.Context, inventory *entities.Inventory) error
	Delete(ctx context.Context, id uint) error
	UpdateQuantity(ctx context.Context, id uint, quantity int) error
}

type DistributorQueryRepositoryInterface interface {
	FindByID(ctx context.Context, id uint) (*entities.Distributor, error)
	FindByEmail(ctx context.Context, email string) (*entities.Distributor, error)
	FindAll(ctx context.Context) ([]*entities.Distributor, error)
}

type DistributorCommandRepositoryInterface interface {
	Save(ctx context.Context, distributor *entities.Distributor) error
	Update(ctx context.Context, distributor *entities.Distributor) error
	Delete(ctx context.Context, id uint) error
}
