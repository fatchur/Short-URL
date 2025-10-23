package inventory

import (
	"time"

	"short-url/domains/entities"
	"short-url/domains/helper"
	"short-url/domains/values/enums"
)

var Inventories = []entities.Inventory{
	{
		DistributorID: helper.UintPtr(1),
		Name:          "MacBook Pro 14-inch",
		Description:   helper.StringPtr("Apple MacBook Pro with M2 chip, 16GB RAM, 512GB SSD"),
		SKU:           "MBP-14-M2-16-512",
		CategoryID:    helper.InventoryCategoryPtr(enums.Electronics),
		Quantity:      25,
		MinQuantity:   helper.IntPtr(5),
		UnitPrice:     29999000.00,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
		UpdatedBy:     helper.UintPtr(1),
	},
	{
		DistributorID: helper.UintPtr(1),
		Name:          "iPhone 15 Pro",
		Description:   helper.StringPtr("Apple iPhone 15 Pro with A17 Pro chip, 128GB storage"),
		SKU:           "IP15P-128-TI",
		CategoryID:    helper.InventoryCategoryPtr(enums.Electronics),
		Quantity:      50,
		MinQuantity:   helper.IntPtr(10),
		UnitPrice:     17999000.00,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
		UpdatedBy:     helper.UintPtr(1),
	},
	{
		DistributorID: helper.UintPtr(2),
		Name:          "Nike Air Force 1",
		Description:   helper.StringPtr("Classic white Nike Air Force 1 sneakers"),
		SKU:           "NAF1-WHT-42",
		CategoryID:    helper.InventoryCategoryPtr(enums.Clothing),
		Quantity:      100,
		MinQuantity:   helper.IntPtr(20),
		UnitPrice:     1500000.00,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
		UpdatedBy:     helper.UintPtr(1),
	},
	{
		DistributorID: helper.UintPtr(3),
		Name:          "Organic Coffee Beans",
		Description:   helper.StringPtr("Premium organic Arabica coffee beans from Aceh"),
		SKU:           "OCB-ACEH-1KG",
		CategoryID:    helper.InventoryCategoryPtr(enums.Food),
		Quantity:      200,
		MinQuantity:   helper.IntPtr(50),
		UnitPrice:     150000.00,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
		UpdatedBy:     helper.UintPtr(1),
	},
	{
		DistributorID: helper.UintPtr(2),
		Name:          "Programming Books Set",
		Description:   helper.StringPtr("Collection of modern programming books including Go, React, and System Design"),
		SKU:           "PBS-PROG-SET",
		CategoryID:    helper.InventoryCategoryPtr(enums.Books),
		Quantity:      30,
		MinQuantity:   helper.IntPtr(5),
		UnitPrice:     750000.00,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
		UpdatedBy:     helper.UintPtr(1),
	},
}
