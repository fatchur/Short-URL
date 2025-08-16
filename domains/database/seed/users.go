package seed

import (
	"log"
	"time"

	"short-url/domains/entities"
	"golang.org/x/crypto/bcrypt"
)

func StringPtr(s string) *string {
	return &s
}

func UintPtr(u uint) *uint {
	return &u
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}
	return string(hash)
}

var Users = []entities.User{
	{
		InstitutionID: 1,
		Name:          "John Doe",
		Email:         "john@example.com",
		PasswordHash:  hashPassword("password123"),
		PhoneNumber:   StringPtr("081234567890"),
		IsActive:      true,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
		UpdatedBy:     UintPtr(1),
	},
	{
		InstitutionID: 1,
		Name:          "Jane Smith",
		Email:         "jane@example.com",
		PasswordHash:  hashPassword("password456"),
		PhoneNumber:   StringPtr("081987654321"),
		IsActive:      true,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
		UpdatedBy:     UintPtr(1),
	},
	{
		InstitutionID: 2,
		Name:          "Bob Johnson",
		Email:         "bob@example.com",
		PasswordHash:  hashPassword("password789"),
		PhoneNumber:   StringPtr("081555666777"),
		IsActive:      false,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
		UpdatedAt:     time.Now(),
		UpdatedBy:     UintPtr(1),
	},
}
