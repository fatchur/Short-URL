package controller

import (
	"inventory-service/middleware"
	"short-url/domains/dto"
	"short-url/domains/dto/inventory"
	"short-url/domains/helper/validation"
	"short-url/domains/service"
	"short-url/domains/values/enums"
	"strconv"

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

func (c *InventoryController) parsePagination(ctx *fiber.Ctx) dto.Pagination {
	pagination := dto.Pagination{
		Page:     1,
		PageSize: 10,
	}

	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			pagination.Page = page
		}
	}

	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			pagination.PageSize = pageSize
		}
	}

	return pagination
}

func (c *InventoryController) CreateInventory(ctx *fiber.Ctx) error {
	var req inventory.CreateInventoryRequest

	if err := ctx.BodyParser(&req); err != nil {
		return middleware.HandleBadRequestError(ctx, err)
	}

	if validationErrors := validation.ValidateStruct(&req); len(validationErrors) > 0 {
		return middleware.HandleValidationError(ctx, validationErrors)
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

	if validationErrors := validation.ValidateStruct(&req); len(validationErrors) > 0 {
		return middleware.HandleValidationError(ctx, validationErrors)
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
	pagination := c.parsePagination(ctx)
	
	inventories, paginationResponse, err := c.service.GetInventoryList(ctx.Context(), pagination)
	if err != nil {
		return middleware.HandleInternalServerError(ctx, err)
	}

	responseData := map[string]interface{}{
		"inventories": inventories,
		"pagination":  paginationResponse,
	}

	response := dto.NewSuccessResponse(fiber.StatusOK, "Inventory list retrieved successfully", responseData)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *InventoryController) GetInventoryByCategory(ctx *fiber.Ctx) error {
	categoryStr := ctx.Params("category")
	if categoryStr == "" {
		return middleware.HandleBadRequestError(ctx, fiber.NewError(fiber.StatusBadRequest, "Category parameter is required"))
	}

	category := enums.InventoryCategoryFromString(categoryStr)
	if category == enums.InvalidInventoryCategory {
		return middleware.HandleBadRequestError(ctx, fiber.NewError(fiber.StatusBadRequest, "Invalid category"))
	}

	pagination := c.parsePagination(ctx)
	
	inventories, paginationResponse, err := c.service.GetInventoryByCategory(ctx.Context(), category, pagination)
	if err != nil {
		return middleware.HandleInternalServerError(ctx, err)
	}

	responseData := map[string]interface{}{
		"inventories": inventories,
		"pagination":  paginationResponse,
	}

	response := dto.NewSuccessResponse(fiber.StatusOK, "Inventory by category retrieved successfully", responseData)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *InventoryController) GetInventoryByDistributor(ctx *fiber.Ctx) error {
	distributorIDStr := ctx.Params("distributor_id")
	if distributorIDStr == "" {
		return middleware.HandleBadRequestError(ctx, fiber.NewError(fiber.StatusBadRequest, "Distributor ID parameter is required"))
	}

	distributorID, err := strconv.ParseUint(distributorIDStr, 10, 32)
	if err != nil {
		return middleware.HandleBadRequestError(ctx, fiber.NewError(fiber.StatusBadRequest, "Invalid distributor ID"))
	}

	pagination := c.parsePagination(ctx)
	
	inventories, paginationResponse, err := c.service.GetInventoryByDistributor(ctx.Context(), uint(distributorID), pagination)
	if err != nil {
		return middleware.HandleInternalServerError(ctx, err)
	}

	responseData := map[string]interface{}{
		"inventories": inventories,
		"pagination":  paginationResponse,
	}

	response := dto.NewSuccessResponse(fiber.StatusOK, "Inventory by distributor retrieved successfully", responseData)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *InventoryController) GetLowStockInventory(ctx *fiber.Ctx) error {
	pagination := c.parsePagination(ctx)
	
	inventories, paginationResponse, err := c.service.GetLowStockInventory(ctx.Context(), pagination)
	if err != nil {
		return middleware.HandleInternalServerError(ctx, err)
	}

	responseData := map[string]interface{}{
		"inventories": inventories,
		"pagination":  paginationResponse,
	}

	response := dto.NewSuccessResponse(fiber.StatusOK, "Low stock inventory retrieved successfully", responseData)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *InventoryController) RegisterRoutes(api fiber.Router) {
	api.Post("/inventory", c.CreateInventory)
	api.Put("/inventory", c.UpdateInventory)
	api.Get("/inventory/:sku", c.GetInventoryBySKU)
	api.Get("/inventories", c.GetInventoryList)
	api.Get("/inventories/category/:category", c.GetInventoryByCategory)
	api.Get("/inventories/distributor/:distributor_id", c.GetInventoryByDistributor)
	api.Get("/inventories/low-stock", c.GetLowStockInventory)
}