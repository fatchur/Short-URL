package database

import (
	"context"
	"fmt"
	"log"

	"short-url/domains/dto"
	"short-url/domains/entities"
)

var MigrateModels = []interface{}{
	&entities.User{},
	&entities.UserSession{},
	&entities.ShortUrl{},
	&entities.ShortClickDaily{},
	&entities.UrlSafety{},
	&entities.Distributor{},
	&entities.Inventory{},
}

func Migrate(ctx context.Context, dbConfig dto.DBConfig) error {
	db, err := DBConnect(ctx, dbConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Starting database migration...")

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
