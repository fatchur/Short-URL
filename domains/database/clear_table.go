package database

import (
	"context"
	"fmt"
	"log"

	"short-url/domains/dto"
	"short-url/domains/entities"
)

var ClearModels = []interface{}{
	&entities.UrlSafety{},
	&entities.ShortClickDaily{},
	&entities.ShortUrl{},
	&entities.User{},
}

func ClearTables(ctx context.Context, dbConfig dto.DBConfig) error {
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

	log.Println("Starting database table clearing...")

	for _, model := range ClearModels {
		modelName := fmt.Sprintf("%T", model)
		log.Printf("Clearing data from: %s", modelName)

		if err := db.Unscoped().Where("1 = 1").Delete(model).Error; err != nil {
			return fmt.Errorf("failed to clear data from %s: %w", modelName, err)
		}

		log.Printf("Successfully cleared data from: %s", modelName)
	}

	log.Println("Database table clearing completed successfully!")
	return nil
}