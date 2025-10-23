package service

import (
	"context"
	"time"

	"short-url/domains/dto/inventory"
	"short-url/domains/entities"
	"short-url/domains/repositories"
	"short-url/domains/service"
)

type inventoryService struct {
	commandRepo repositories.InventoryCommandRepositoryInterface
	queryRepo   repositories.InventoryQueryRepositoryInterface
}

func NewInventoryService(
	commandRepo repositories.InventoryCommandRepositoryInterface,
	queryRepo repositories.InventoryQueryRepositoryInterface,
) service.InventoryServiceInterface {
	return &inventoryService{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func (s *inventoryService) CreateInventory(ctx context.Context, req *inventory.CreateInventoryRequest) (*entities.Inventory, error) {
	inventoryEntity := &entities.Inventory{
		DistributorID: req.DistributorID,
		Name:          req.Name,
		Description:   req.Description,
		SKU:           req.SKU,
		CategoryID:    req.CategoryID,
		Quantity:      req.Quantity,
		MinQuantity:   req.MinQuantity,
		UnitPrice:     req.UnitPrice,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
	}

	if err := s.commandRepo.Save(ctx, inventoryEntity); err != nil {
		return nil, err
	}
	return inventoryEntity, nil
}

func (s *inventoryService) UpdateInventory(ctx context.Context, req *inventory.UpdateInventoryRequest) (*entities.Inventory, error) {
	inventoryEntity := &entities.Inventory{
		ID:            req.ID,
		DistributorID: req.DistributorID,
		Name:          req.Name,
		Description:   req.Description,
		SKU:           req.SKU,
		CategoryID:    req.CategoryID,
		Quantity:      req.Quantity,
		MinQuantity:   req.MinQuantity,
		UnitPrice:     req.UnitPrice,
		UpdatedAt:     time.Now(),
		UpdatedBy:     &[]uint{1}[0],
	}

	if err := s.commandRepo.Update(ctx, inventoryEntity); err != nil {
		return nil, err
	}
	return inventoryEntity, nil
}

func (s *inventoryService) GetInventoryBySKU(ctx context.Context, sku string) (*entities.Inventory, error) {
	return s.queryRepo.FindBySKU(ctx, sku)
}

func (s *inventoryService) GetInventoryList(ctx context.Context) ([]*entities.Inventory, error) {
	return s.queryRepo.FindAll(ctx)
}