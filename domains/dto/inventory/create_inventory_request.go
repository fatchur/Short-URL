package inventory

import "short-url/domains/values/enums"

type CreateInventoryRequest struct {
	DistributorID *uint                    `json:"distributor_id" validate:"omitempty,min=1"`
	Name          string                   `json:"name" validate:"required,min=1,max=255"`
	Description   *string                  `json:"description" validate:"omitempty,max=1000"`
	SKU           string                   `json:"sku" validate:"required,min=1,max=100"`
	CategoryID    *enums.InventoryCategory `json:"category_id" validate:"omitempty"`
	Quantity      int                      `json:"quantity" validate:"required,min=0"`
	MinQuantity   *int                     `json:"min_quantity" validate:"omitempty,min=0"`
	UnitPrice     float64                  `json:"unit_price" validate:"required,min=0"`
}