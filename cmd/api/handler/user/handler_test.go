package user

import (
	"encoding/json"
	"errors"
	"github.com/PakornBank/go-backend-example/internal/user"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PakornBank/go-backend-example/internal/common/model"
	"github.com/PakornBank/go-backend-example/internal/common/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupHandlerTest(ctrl *gomock.Controller, middleware gin.HandlerFunc) (*gin.Engine, *user.MockService) {
	gin.SetMode(gin.TestMode)

	mockService := user.NewMockService(ctrl)
	userHandler := &handler{service: mockService}

	router := gin.New()
	group := router.Group("/api")
	if middleware != nil {
		group.Use(middleware)
	}
	{
		group.GET("/profile", userHandler.GetProfile)
	}

	return router, mockService
}

func TestNewHandler(t *testing.T) {
	mockService := new(user.MockService)
	userHandler := NewHandler(mockService)

	assert.NotNil(t, userHandler)
	assert.Equal(t, mockService, userHandler.(*handler).service)
}

func Test_handler_GetProfile(t *testing.T) {
	mockUser := testutil.NewMockUser()

	tests := []struct {
		name        string
		middleware  gin.HandlerFunc
		mockFn      func(*user.MockService)
		wantCode    int
		errContains string
	}{
		{
			name: "successful profile retrieval",
			middleware: func(c *gin.Context) {
				c.Set("user_id", mockUser.ID.String())
			},
			mockFn: func(ms *user.MockService) {
				ms.EXPECT().GetUserByID(gomock.Any(), mockUser.ID.String()).Return(&mockUser, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "user_service error",
			middleware: func(c *gin.Context) {
				c.Set("user_id", mockUser.ID.String())
			},
			mockFn: func(ms *user.MockService) {
				ms.EXPECT().GetUserByID(gomock.Any(), mockUser.ID.String()).Return(nil, errors.New("user_service error"))
			},
			wantCode:    http.StatusNotFound,
			errContains: "user_service error",
		},
		{
			name:        "no user_id input context",
			wantCode:    http.StatusUnauthorized,
			errContains: "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			router, mockService := setupHandlerTest(ctrl, tt.middleware)
			if tt.mockFn != nil {
				tt.mockFn(mockService)
			}

			req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantCode == http.StatusOK {
				var res model.User
				err := json.Unmarshal(w.Body.Bytes(), &res)
				assert.NoError(t, err)

				assert.Equal(t, mockUser.ID, res.ID)
				assert.Equal(t, mockUser.Email, res.Email)
				assert.Equal(t, mockUser.FullName, res.FullName)
				assert.Equal(t, mockUser.CreatedAt.Format(time.RFC3339Nano), res.CreatedAt.Format(time.RFC3339Nano))
				assert.Equal(t, mockUser.UpdatedAt.Format(time.RFC3339Nano), res.UpdatedAt.Format(time.RFC3339Nano))
				assert.Empty(t, res.PasswordHash)
			} else {
				var res map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &res)
				assert.NoError(t, err)

				assert.Contains(t, res["error"], tt.errContains)
			}
		})
	}
}
