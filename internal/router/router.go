package router

import (
	"github.com/PakornBank/go-backend-example/internal/auth"
	"github.com/PakornBank/go-backend-example/internal/di"
	"github.com/PakornBank/go-backend-example/internal/user"
	"github.com/gin-gonic/gin"
)

// SetupRoutes call functions to register routes on gin router.
func SetupRoutes(router *gin.Engine, container *di.Container) {
	group := router.Group("/api")
	auth.RegisterRoutes(group, container.AuthHandler)
	user.RegisterRoutes(group, container.UserHandler, container.Config)
}
