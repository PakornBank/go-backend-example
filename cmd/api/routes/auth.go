package routes

import (
	"github.com/PakornBank/go-backend-example/cmd/api/handler/auth"
	"github.com/gin-gonic/gin"
)

// registerAuthRoutes registers the auth routes with the provided gin routes group and handler.
func registerAuthRoutes(r *gin.RouterGroup, h auth.Handler) {
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", h.Register)
		authRoutes.POST("/login", h.Login)
	}
}
