package user

import (
	"fmt"
	"testing"

	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/common/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type routeInfo struct {
	handler string
	method  string
}

const (
	handlerBase = "github.com/PakornBank/go-backend-example/internal/user.Handler."
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		wantRoutes map[string]routeInfo
	}{
		{
			name: "routes registered successfully",
			wantRoutes: map[string]routeInfo{
				"/user/profile": {
					handler: handlerPath("GetProfile"),
					method:  "GET",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			group := router.Group("/")
			_, db, _ := testutil.DBMock(t)
			cfg := &config.Config{
				JWTSecret: "test-secret",
			}

			RegisterRoutes(group, db, cfg)

			routes := router.Routes()
			assert.Equal(t, 1, len(routes))

			for _, route := range routes {
				assert.Equal(t, tt.wantRoutes[route.Path].handler, route.Handler)
				assert.Equal(t, tt.wantRoutes[route.Path].method, route.Method)
			}
		})
	}
}

func handlerPath(name string) string {
	return fmt.Sprintf("%s%s-fm", handlerBase, name)
}
