package main

import (
	"log"

	"kouji-app-backend2/internal/config"
	"kouji-app-backend2/internal/database"
)

func main() {
	config.LoadEnv()
	cfg := database.FromEnv()

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	log.Println("migrations applied successfully")
}
