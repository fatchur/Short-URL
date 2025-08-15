package seed

import (
	"time"

	"short-url/domains/entities"
)

func StringPtr(s string) *string {
	return &s
}

func UintPtr(u uint) *uint {
	return &u
}

var Users = []entities.User{
	{
		ID:            1,
		InstitutionID: 1,
		Name:          "John Doe",
		Email:         "john@example.com",
		PhoneNumber:   StringPtr("081234567890"),
		IsActive:      true,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
		UpdatedBy:     UintPtr(1),
	},
	{
		ID:            2,
		InstitutionID: 1,
		Name:          "Jane Smith",
		Email:         "jane@example.com",
		PhoneNumber:   StringPtr("081987654321"),
		IsActive:      true,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
		UpdatedBy:     UintPtr(1),
	},
	{
		ID:            3,
		InstitutionID: 2,
		Name:          "Bob Johnson",
		Email:         "bob@example.com",
		PhoneNumber:   StringPtr("081555666777"),
		IsActive:      false,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
		UpdatedBy:     UintPtr(1),
	},
}
