package router

import (
	"testing"

	"github.com/PakornBank/go-backend-example/internal/auth"
	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/di"
	"github.com/PakornBank/go-backend-example/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	ctrl := gomock.NewController(t)
	userHandler := user.NewMockHandler(ctrl)
	authHandler := auth.NewMockHandler(ctrl)
	container := &di.Container{
		UserHandler: userHandler,
		AuthHandler: authHandler,
		Config:      &config.Config{},
	}

	initialRouteCount := len(router.Routes())
	SetupRoutes(router, container)

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
