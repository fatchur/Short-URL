package database

import (
	"context"
	"fmt"
	"log"

	"short-url/domains/dto"
	"short-url/domains/entities"
)

var DropModels = []interface{}{
	&entities.UrlSafety{},
	&entities.ShortClickDaily{},
	&entities.ShortUrl{},
	&entities.UserSession{},
	&entities.User{},
}

func DropTables(ctx context.Context, dbConfig dto.DBConfig) error {
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

	log.Println("Starting database table drop...")

	for _, model := range DropModels {
		modelName := fmt.Sprintf("%T", model)
		log.Printf("Dropping table for: %s", modelName)

		if err := db.Migrator().DropTable(model); err != nil {
			return fmt.Errorf("failed to drop table for %s: %w", modelName, err)
		}

		log.Printf("Successfully dropped table for: %s", modelName)
	}

	log.Println("Database table drop completed successfully!")
	return nil
}