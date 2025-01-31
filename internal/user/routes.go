package user

import (
	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/common/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the user routes with the provided gin router group.
// It initializes the auth handler with the provided database connection and configuration.
func RegisterRoutes(group *gin.RouterGroup, h Handler, cfg *config.Config) {
	user := group.Group("/user")
	{
		protected := user.Group("")
		protected.Use(middleware.Auth(cfg.JWTSecret))
		{
			protected.GET("/profile", h.GetProfile)
		}
	}
}
