package router

import (
	"github.com/PakornBank/go-backend-example/internal/auth"
	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes call functions to register routes on gin router.
func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	group := router.Group("/api")
	auth.RegisterRoutes(group, db, cfg)
	user.RegisterRoutes(group, db, cfg)
}
