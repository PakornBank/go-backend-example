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
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "go_backend_db"),
		DBPort:         getEnv("DB_PORT", "5432"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		TokenExpiryDur: 24 * time.Hour,
	}

	if config.JWTSecret == "your-secret-key" {
		return nil, errors.New("jwt secret must be set in environment")
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
