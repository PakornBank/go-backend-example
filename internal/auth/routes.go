package auth

import (
	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers the auth routes with the provided gin router group.
// It initializes the auth handler with the provided database connection and configuration.
func RegisterRoutes(group *gin.RouterGroup, db *gorm.DB, config *config.Config) {
	authHandler := NewHandler(
		NewService(
			NewRepository(db),
			config,
		),
	)
	auth := group.Group("/auth")
	{
		public := auth.Group("")
		{
			public.POST("/register", authHandler.Register)
			public.POST("/login", authHandler.Login)
		}
	}
}
