package entities

import (
	"short-url/domains/values/enums"
	"time"
)

type Inventory struct {
	ID            uint                     `json:"-" gorm:"primaryKey;autoIncrement"`
	DistributorID *uint                    `json:"distributor_id" gorm:"index;foreignKey:DistributorID;references:ID"`
	Name          string                   `json:"name" gorm:"type:varchar(255);not null"`
	Description   *string                  `json:"description" gorm:"type:text"`
	SKU           string                   `json:"sku" gorm:"type:varchar(100);not null;uniqueIndex"`
	CategoryID    *enums.InventoryCategory `json:"category_id" gorm:"index"`
	Quantity      int                      `json:"quantity" gorm:"not null;default:0"`
	MinQuantity   *int                     `json:"min_quantity"`
	UnitPrice     float64                  `json:"unit_price" gorm:"type:decimal(10,2);not null;default:0"`
	CreatedAt     time.Time                `json:"-" gorm:"autoCreateTime"`
	CreatedBy     uint                     `json:"-" gorm:"index"`
	UpdatedAt     time.Time                `json:"-" gorm:"autoUpdateTime"`
	UpdatedBy     *uint                    `json:"-" gorm:"index"`
	Distributor   *Distributor             `json:"distributor,omitempty" gorm:"foreignKey:DistributorID"`
}

func (Inventory) TableName() string {
	return "inventories"
}
