package routes

import (
	"github.com/PakornBank/go-backend-example/cmd/api/handler/user"
	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/common/middleware"
	"github.com/gin-gonic/gin"
)

// registerUserRoutes registers the user routes with the provided gin routes group and handler.
func registerUserRoutes(r *gin.RouterGroup, h user.Handler, cfg *config.Config) {
	userRoutes := r.Group("/user")
	{
		protected := userRoutes.Group("")
		protected.Use(middleware.Auth(cfg.JWTSecret))
		{
			protected.GET("/profile", h.GetProfile)
		}
	}
}
