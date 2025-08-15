package database

import (
	"fmt"
	"log"

	"short-url/domains/database/seed"
	"short-url/domains/entities"

	"gorm.io/gorm"
)


func Seed(db *gorm.DB) error {
	// Test database connection
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Starting database seeding...")

	// Seed users
	if err := seedUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func seedUsers(db *gorm.DB) error {
	log.Println("Starting user seeding...")

	// Load users from seed data
	users := seed.Users

	for _, user := range users {
		var existingUser entities.User
		result := db.Where("email = ?", user.Email).First(&existingUser)
		if result.Error == nil {
			log.Printf("User with email %s already exists, skipping", user.Email)
			continue
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", user.Email, err)
			return err
		}
		log.Printf("Successfully seeded user: %s", user.Name)
	}

	log.Println("User seeding completed successfully!")
	return nil
}