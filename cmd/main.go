package main

import (
	"context"
	"flag"
	"log"

	"short-url/domains/config"
	"short-url/domains/database"
	"short-url/domains/dto"
)

func main() {
	var command string
	flag.StringVar(&command, "d", "", "Database command to execute")
	flag.Parse()

	ctx := context.Background()
	cfg := config.LoadConfig()

	dbConfig := dto.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
		Timezone: cfg.DBTimezone,
		LogLevel: cfg.DBLogLevel,
	}

	switch command {
	case "migrate":
		if err := database.Migrate(ctx, dbConfig); err != nil {
			log.Fatal("Migration failed:", err)
		}
	case "seed":
		db, err := database.DBConnect(ctx, dbConfig)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		if err := database.Seed(db); err != nil {
			log.Fatal("Seeding failed:", err)
		}
	case "drop-table":
		if err := database.DropTables(ctx, dbConfig); err != nil {
			log.Fatal("Drop tables failed:", err)
		}
	case "clear-table":
		if err := database.ClearTables(ctx, dbConfig); err != nil {
			log.Fatal("Clear tables failed:", err)
		}
	default:
		log.Fatal("Unknown command. Use: -d migrate, -d seed, -d drop-table, or -d clear-table")
	}
}
