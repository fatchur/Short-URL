package enums

import (
	"encoding/json"
	"strings"
)

type InventoryCategory int

const (
	InvalidInventoryCategory InventoryCategory = 0
	Electronics              InventoryCategory = iota + 1
	Clothing
	Food
	Books
	Furniture
	Automotive
	Health
	Sports
	Toys
	Home
)

func InventoryCategoryFromString(s string) InventoryCategory {
	switch strings.ToLower(s) {
	case "electronics":
		return Electronics
	case "clothing":
		return Clothing
	case "food":
		return Food
	case "books":
		return Books
	case "furniture":
		return Furniture
	case "automotive":
		return Automotive
	case "health":
		return Health
	case "sports":
		return Sports
	case "toys":
		return Toys
	case "home":
		return Home
	default:
		return InvalidInventoryCategory
	}
}

func (ic InventoryCategory) String() string {
	switch ic {
	case Electronics:
		return "electronics"
	case Clothing:
		return "clothing"
	case Food:
		return "food"
	case Books:
		return "books"
	case Furniture:
		return "furniture"
	case Automotive:
		return "automotive"
	case Health:
		return "health"
	case Sports:
		return "sports"
	case Toys:
		return "toys"
	case Home:
		return "home"
	default:
		return "invalid"
	}
}

func (ic InventoryCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(ic.String())
}

func (ic *InventoryCategory) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*ic = InventoryCategoryFromString(s)
	return nil
}
