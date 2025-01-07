package user

import (
	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/common/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers the user routes with the provided gin router group.
// It initializes the auth handler with the provided database connection and configuration.
func RegisterRoutes(group *gin.RouterGroup, db *gorm.DB, config *config.Config) {
	userHandler := NewHandler(
		NewService(
			NewRepository(db),
			config,
		),
	)
	user := group.Group("/user")
	{
		protected := user.Group("")
		protected.Use(middleware.Auth(config.JWTSecret))
		{
			protected.GET("/profile", userHandler.GetProfile)
		}
	}
}
