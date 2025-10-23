package database

import (
	"fmt"
	"log"

	"short-url/domains/database/seed"
	"short-url/domains/database/seed/inventory"
	"short-url/domains/entities"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Starting database seeding...")

	userIDs, err := seedUsers(db)
	if err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := seedUserSessions(db, userIDs); err != nil {
		return fmt.Errorf("failed to seed user sessions: %w", err)
	}

	distributorIDs, err := seedDistributors(db)
	if err != nil {
		return fmt.Errorf("failed to seed distributors: %w", err)
	}

	if err := seedInventories(db, distributorIDs); err != nil {
		return fmt.Errorf("failed to seed inventories: %w", err)
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func seedUsers(db *gorm.DB) ([]uint, error) {
	log.Println("Starting user seeding...")

	users := seed.Users
	var userIDs []uint

	for _, user := range users {
		var existingUser entities.User
		result := db.Where("email = ?", user.Email).First(&existingUser)
		if result.Error == nil {
			log.Printf("User with email %s already exists, using existing ID", user.Email)
			userIDs = append(userIDs, existingUser.ID)
			continue
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", user.Email, err)
			return nil, err
		}
		log.Printf("Successfully seeded user: %s", user.Name)
		userIDs = append(userIDs, user.ID)
	}

	log.Println("User seeding completed successfully!")
	return userIDs, nil
}

func seedUserSessions(db *gorm.DB, userIDs []uint) error {
	log.Println("Starting user session seeding...")

	sessions := seed.UserSessions

	for i, session := range sessions {
		if i >= len(userIDs) {
			log.Printf("Not enough user IDs for session %d, skipping", i)
			continue
		}

		session.UserID = userIDs[i]

		var existingSession entities.UserSession
		result := db.Where("session_code = ?", session.SessionCode).First(&existingSession)
		if result.Error == nil {
			log.Printf("Session with code %s already exists, skipping", session.SessionCode)
			continue
		}

		if err := db.Create(&session).Error; err != nil {
			log.Printf("Failed to seed session %s: %v", session.SessionCode, err)
			return err
		}
		log.Printf("Successfully seeded session for user ID: %d", session.UserID)
	}

	log.Println("User session seeding completed successfully!")
	return nil
}

func seedDistributors(db *gorm.DB) ([]uint, error) {
	log.Println("Starting distributor seeding...")

	distributors := inventory.InventoryDistributors
	var distributorIDs []uint

	for _, distributor := range distributors {
		var existingDistributor entities.Distributor
		result := db.Where("email = ?", distributor.Email).First(&existingDistributor)
		if result.Error == nil {
			log.Printf("Distributor with email %s already exists, using existing ID", distributor.Email)
			distributorIDs = append(distributorIDs, existingDistributor.ID)
			continue
		}

		if err := db.Create(&distributor).Error; err != nil {
			log.Printf("Failed to seed distributor %s: %v", distributor.Email, err)
			return nil, err
		}
		log.Printf("Successfully seeded distributor: %s", distributor.Name)
		distributorIDs = append(distributorIDs, distributor.ID)
	}

	log.Println("Distributor seeding completed successfully!")
	return distributorIDs, nil
}

func seedInventories(db *gorm.DB, distributorIDs []uint) error {
	log.Println("Starting inventory seeding...")

	inventories := inventory.Inventories

	for _, inv := range inventories {
		var existingInventory entities.Inventory
		result := db.Where("sku = ?", inv.SKU).First(&existingInventory)
		if result.Error == nil {
			log.Printf("Inventory with SKU %s already exists, skipping", inv.SKU)
			continue
		}

		if err := db.Create(&inv).Error; err != nil {
			log.Printf("Failed to seed inventory %s: %v", inv.SKU, err)
			return err
		}
		log.Printf("Successfully seeded inventory: %s", inv.Name)
	}

	log.Println("Inventory seeding completed successfully!")
	return nil
}
