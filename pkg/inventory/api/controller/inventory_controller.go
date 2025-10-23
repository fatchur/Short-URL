package controller

import (
	"inventory-service/middleware"
	"short-url/domains/dto"
	"short-url/domains/dto/inventory"
	"short-url/domains/service"

	"github.com/gofiber/fiber/v2"
)

type InventoryController struct {
	service service.InventoryServiceInterface
}

func NewInventoryController(service service.InventoryServiceInterface) *InventoryController {
	return &InventoryController{
		service: service,
	}
}

func (c *InventoryController) CreateInventory(ctx *fiber.Ctx) error {
	var req inventory.CreateInventoryRequest

	if err := ctx.BodyParser(&req); err != nil {
		return middleware.HandleBadRequestError(ctx, err)
	}

	createdInventory, err := c.service.CreateInventory(ctx.Context(), &req)
	if err != nil {
		return middleware.HandleInternalServerError(ctx, err)
	}

	response := dto.NewSuccessResponse(fiber.StatusCreated, "Inventory created successfully", createdInventory)
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (c *InventoryController) UpdateInventory(ctx *fiber.Ctx) error {
	var req inventory.UpdateInventoryRequest

	if err := ctx.BodyParser(&req); err != nil {
		return middleware.HandleBadRequestError(ctx, err)
	}

	updatedInventory, err := c.service.UpdateInventory(ctx.Context(), &req)
	if err != nil {
		return middleware.HandleInternalServerError(ctx, err)
	}

	response := dto.NewSuccessResponse(fiber.StatusOK, "Inventory updated successfully", updatedInventory)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *InventoryController) GetInventoryBySKU(ctx *fiber.Ctx) error {
	sku := ctx.Params("sku")
	if sku == "" {
		return middleware.HandleBadRequestError(ctx, fiber.NewError(fiber.StatusBadRequest, "SKU parameter is required"))
	}

	inventory, err := c.service.GetInventoryBySKU(ctx.Context(), sku)
	if err != nil {
		return middleware.HandleDatabaseError(ctx, err)
	}

	response := dto.NewSuccessResponse(fiber.StatusOK, "Inventory retrieved successfully", inventory)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *InventoryController) GetInventoryList(ctx *fiber.Ctx) error {
	inventories, err := c.service.GetInventoryList(ctx.Context())
	if err != nil {
		return middleware.HandleInternalServerError(ctx, err)
	}

	response := dto.NewSuccessResponse(fiber.StatusOK, "Inventory list retrieved successfully", inventories)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *InventoryController) RegisterRoutes(api fiber.Router) {
	api.Post("/inventory", c.CreateInventory)
	api.Put("/inventory", c.UpdateInventory)
	api.Get("/inventory/:sku", c.GetInventoryBySKU)
	api.Get("/inventories", c.GetInventoryList)
}