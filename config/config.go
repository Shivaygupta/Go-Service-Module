package config

import (
	"fmt"
	"net/http"
	"os"

	"kong-go-assignment/errors"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() Config {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		errors.New(http.StatusInternalServerError, "failed to load env file", err.Error())
	}

	cfg := Config{
		Port:       os.Getenv("APP_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}

	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBName == "" {
		debugInfo := fmt.Sprintf("Database configuration is incomplete. Got: host=%q port=%q user=%q name=%q",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName)
		errors.New(http.StatusInternalServerError, "failed to load database configuration", debugInfo)
	}

	return cfg
}
