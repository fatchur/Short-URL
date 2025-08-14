package main

import (
	"flag"
	"log"

	"short-url/domains/config"
	"short-url/domains/database"
)

func main() {
	var command string
	flag.StringVar(&command, "d", "", "Database command to execute")
	flag.Parse()

	// Load configuration and setup database connection
	cfg := config.LoadConfig()
	dbConfig := database.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
		Timezone: cfg.DBTimezone,
		LogLevel: cfg.DBLogLevel,
	}
	database.SetConfig(dbConfig)

	db, err := database.DBConnect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	switch command {
	case "migrate":
		if err := Migrate(); err != nil {
			log.Fatal("Migration failed:", err)
		}
	case "seed":
		if err := Seed(db); err != nil {
			log.Fatal("Seeding failed:", err)
		}
	default:
		log.Fatal("Unknown command. Use: -d migrate or -d seed")
	}
}
