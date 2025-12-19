package main

import (
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"kouji-app-backend2/internal/config"
	"kouji-app-backend2/internal/database"
	userhandler "kouji-app-backend2/internal/user"
)

func main() {
	config.LoadEnv()

	cfg := database.FromEnv()
	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler := userhandler.NewHandler(db)
	handler.Register(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("server listening on %s", addr)
	e.Logger.Fatal(e.Start(addr))
}
