package main

import (
	"flag"
	"log"
)

func main() {
	var command string
	flag.StringVar(&command, "d", "", "Database command to execute")
	flag.Parse()

	switch command {
	case "migrate":
		if err := Migrate(); err != nil {
			log.Fatal("Migration failed:", err)
		}
	default:
		log.Fatal("Unknown command. Use: -d migrate")
	}
}
