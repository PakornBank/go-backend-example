package di

import (
	"github.com/PakornBank/go-backend-example/internal/common/health"
	"gorm.io/gorm"
	"log"

	"github.com/PakornBank/go-backend-example/internal/auth"
	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/common/database"
	"github.com/PakornBank/go-backend-example/internal/user"
)

// Container holds the dependencies for the application.
type Container struct {
	UserHandler   user.Handler
	AuthHandler   auth.Handler
	HealthHandler health.Handler
	Config        *config.Config
	db            *gorm.DB
}

// NewContainer creates a new Container with the provided configuration.
func NewContainer(cfg *config.Config) *Container {
	db, err := database.NewDataBase(cfg)
	if err != nil {
		log.Fatal("failed to initialize database: ", err)
	}

	authHandler := auth.NewHandler(auth.NewService(auth.NewRepository(db), cfg))
	userHandler := user.NewHandler(user.NewService(user.NewRepository(db)))
	healthHandler := health.NewHandler(db)

	return &Container{
		AuthHandler:   authHandler,
		UserHandler:   userHandler,
		HealthHandler: healthHandler,
		Config:        cfg,
		db:            db,
	}
}

// GetDB returns the database instance
func (c *Container) GetDB() (*gorm.DB, error) {
	return c.db, nil
}
