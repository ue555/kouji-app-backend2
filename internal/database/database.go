package database

import (
	"fmt"
	"os"
	"strings"
	"time"

	"kouji-app-backend2/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds connection details for MySQL.
type Config struct {
	DSN      string
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Params   string
}

// FromEnv constructs a Config from common MYSQL_* env vars.
func FromEnv() Config {
	return Config{
		DSN:      os.Getenv("MYSQL_DSN"),
		User:     getEnv("MYSQL_USER", "user"),
		Password: getEnv("MYSQL_PASSWORD", "pass1234"),
		Host:     getEnv("MYSQL_HOST", "127.0.0.1"),
		Port:     getEnv("MYSQL_PORT", "3306"),
		Database: getEnv("MYSQL_DATABASE", "kouji_app"),
		Params:   getEnv("MYSQL_PARAMS", "charset=utf8mb4&parseTime=true&loc=Local&collation=utf8mb4_unicode_ci"),
	}
}

// DSNString resolves either the raw DSN or builds one from individual parts.
func (c Config) DSNString() string {
	if strings.TrimSpace(c.DSN) != "" {
		return c.DSN
	}

	params := c.Params
	if params == "" {
		params = "charset=utf8mb4&parseTime=true&loc=Local"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.User, c.Password, c.Host, c.Port, c.Database, params)
}

// Open returns a configured gorm.DB connection.
func Open(cfg Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSNString()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	sqlDB.SetConnMaxIdleTime(15 * time.Minute)

	return db, nil
}

// AutoMigrate runs GORM's automigrate for all models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
