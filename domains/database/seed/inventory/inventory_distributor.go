package inventory

import (
	"time"

	"short-url/domains/entities"
	"short-url/domains/helper"
)

var InventoryDistributors = []entities.Distributor{
	{
		ID:          1,
		Name:        "Tech Solutions Ltd",
		Email:       "contact@techsolutions.com",
		PhoneNumber: helper.StringPtr("021-1234567"),
		Address:     helper.StringPtr("123 Tech Street, Jakarta, Indonesia"),
		CreatedAt:   time.Now(),
		CreatedBy:   1,
		UpdatedAt:   time.Now(),
		UpdatedBy:   helper.UintPtr(1),
	},
	{
		ID:          2,
		Name:        "Global Electronics Corp",
		Email:       "sales@globalelectronics.com",
		PhoneNumber: helper.StringPtr("021-7654321"),
		Address:     helper.StringPtr("456 Electronics Ave, Surabaya, Indonesia"),
		CreatedAt:   time.Now(),
		CreatedBy:   1,
		UpdatedAt:   time.Now(),
		UpdatedBy:   helper.UintPtr(1),
	},
	{
		ID:          3,
		Name:        "Premium Supplies Inc",
		Email:       "info@premiumsupplies.com",
		PhoneNumber: helper.StringPtr("021-9876543"),
		Address:     helper.StringPtr("789 Supply Road, Bandung, Indonesia"),
		CreatedAt:   time.Now(),
		CreatedBy:   1,
		UpdatedAt:   time.Now(),
		UpdatedBy:   helper.UintPtr(1),
	},
}
