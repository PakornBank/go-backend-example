package router

import (
	"testing"

	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	initialRouteCount := len(router.Routes())
	SetupRoutes(router, &gorm.DB{}, &config.Config{})

	routes := router.Routes()
	assert.Greater(t, len(routes), initialRouteCount)

	hasAPIGroup := false
	for _, route := range routes {
		if route.Path[:4] == "/api" {
			hasAPIGroup = true
			break
		}
	}
	assert.True(t, hasAPIGroup)
}
