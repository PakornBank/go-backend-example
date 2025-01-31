package auth

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the auth routes with the provided gin router group and handler.
func RegisterRoutes(group *gin.RouterGroup, h Handler) {
	auth := group.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}
