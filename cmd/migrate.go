package cmd

import (
	"fmt"
	"log"

	"short-url/domains/database"
	"short-url/domains/entities"
)

// MigrateModels contains all the models that need to be migrated
var MigrateModels = []interface{}{
	&entities.User{},
	// Add other entities here as you create them
	// &entities.Institution{},
	// &entities.ShortUrl{},
}

// Migrate runs auto-migration for all registered models
func Migrate() error {
	db, err := database.DBConnect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	defer sqlDB.Close()

	// Test database connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Starting database migration...")

	// Run migrations for all models
	for _, model := range MigrateModels {
		modelName := fmt.Sprintf("%T", model)
		log.Printf("Migrating: %s", modelName)

		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %s: %w", modelName, err)
		}

		log.Printf("Successfully migrated: %s", modelName)
	}

	log.Println("Database migration completed successfully!")
	return nil
}
