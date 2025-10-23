package helper

import "short-url/domains/values/enums"

func StringPtr(s string) *string {
	return &s
}

func UintPtr(u uint) *uint {
	return &u
}

func IntPtr(i int) *int {
	return &i
}

func InventoryCategoryPtr(ic enums.InventoryCategory) *enums.InventoryCategory {
	return &ic
}