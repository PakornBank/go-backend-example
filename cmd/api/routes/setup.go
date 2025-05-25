package routes

import (
	"github.com/PakornBank/go-backend-example/cmd/api/di"
	"github.com/gin-gonic/gin"
)

// SetupRoutes call functions to register routes on gin routes.
func SetupRoutes(router *gin.Engine, container *di.Container) {
	router.GET("/health", container.HealthHandler.Check)

	group := router.Group("/api")
	registerAuthRoutes(group, container.AuthHandler)
	registerUserRoutes(group, container.UserHandler, container.Config)
}
