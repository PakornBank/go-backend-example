// Package database provides database connection and configuration functionality.
package database

import (
	"fmt"
	"time"

	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/common/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDataBase initializes a new database connection using the provided configuration.
func NewDataBase(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DBURL()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database SQL DB: %w", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
