package database

import (
	"fmt"
	"log"

	"short-url/domains/database/seed"
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

	if err := seedUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := seedUserSessions(db); err != nil {
		return fmt.Errorf("failed to seed user sessions: %w", err)
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func seedUsers(db *gorm.DB) error {
	log.Println("Starting user seeding...")

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

func seedUserSessions(db *gorm.DB) error {
	log.Println("Starting user session seeding...")

	sessions := seed.UserSessions

	for _, session := range sessions {
		var existingSession entities.UserSession
		result := db.Where("session_token = ?", session.SessionToken).First(&existingSession)
		if result.Error == nil {
			log.Printf("Session with token %s already exists, skipping", session.SessionToken)
			continue
		}

		if err := db.Create(&session).Error; err != nil {
			log.Printf("Failed to seed session %s: %v", session.SessionToken, err)
			return err
		}
		log.Printf("Successfully seeded session for user ID: %d", session.UserID)
	}

	log.Println("User session seeding completed successfully!")
	return nil
}