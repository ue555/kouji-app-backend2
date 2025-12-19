package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"kouji-app-backend2/internal/config"
	"kouji-app-backend2/internal/database"
	"kouji-app-backend2/internal/models"
)

func main() {
	var (
		count = flag.Int("n", 10, "number of fake users to generate")
		mode  = flag.String("mode", "db", "output destination: db or json")
	)

	flag.Parse()
	config.LoadEnv()

	users, err := models.GenerateFakeUsers(*count)
	if err != nil {
		log.Fatalf("generate fake users: %v", err)
	}

	switch strings.ToLower(*mode) {
	case "db":
		if err := persist(users); err != nil {
			log.Fatal(err)
		}
	case "json":
		if err := outputJSON(users); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unsupported mode %q", *mode)
	}
}

func outputJSON(users []models.User) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(users)
}

func persist(users []models.User) error {
	if len(users) == 0 {
		fmt.Println("no users to insert")
		return nil
	}

	cfg := database.FromEnv()
	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	if err := db.Create(&users).Error; err != nil {
		return fmt.Errorf("insert users: %w", err)
	}

	log.Printf("inserted %d users\n", len(users))
	log.Println("sample credentials:")
	for i, user := range users {
		if i >= 3 {
			break
		}
		log.Printf("  %s / %s", user.Email, user.PlainPassword)
	}

	return nil
}
