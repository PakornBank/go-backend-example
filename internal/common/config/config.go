// Package config provides configuration management for the application.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the configuration values for the application.
type Config struct {
	DBHost         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBPort         string
	ServerPort     string
	JWTSecret      string
	TokenExpiryDur time.Duration
}

// LoadConfig loads the configuration from environment variables and returns a Config struct.
func LoadConfig() (*Config, error) {
	// Only try to load .env file in development (when file exists and is readable)
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Warning: .env file exists but couldn't be loaded: %v\n", err)
		}
	}

	config := &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "go_backend_db"),
		DBPort:         getEnv("DB_PORT", "5432"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		TokenExpiryDur: 24 * time.Hour,
	}

	if config.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable must be set")
	}

	return config, nil
}

// getEnv retrieves the value of the environment variable named by the key.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// DBURL constructs and returns the database connection URL string
func (c *Config) DBURL() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort,
	)
}
