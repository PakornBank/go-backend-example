package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/PakornBank/go-backend-example/cmd/api/model"
	"github.com/PakornBank/go-backend-example/internal/auth"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PakornBank/go-backend-example/internal/common/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupHandlerTest(t *testing.T) (*gin.Engine, *auth.MockService) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockService := auth.NewMockService(ctrl)
	authHandler := &handler{service: mockService}

	router := gin.New()
	group := router.Group("/api")
	{
		group.POST("/register", authHandler.Register)
		group.POST("/login", authHandler.Login)
	}

	return router, mockService
}

func TestNewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := auth.NewMockService(ctrl)
	authHandler := NewHandler(mockService)

	assert.NotNil(t, authHandler)
	assert.Equal(t, mockService, authHandler.(*handler).service)
}

func Test_handler_Register(t *testing.T) {
	user := testutil.NewMockUser()

	tests := []struct {
		name        string
		input       model.RegisterInput
		mockFn      func(*auth.MockService)
		wantCode    int
		errContains string
	}{
		{
			name: "successful registration",
			input: model.RegisterInput{
				Email:    user.Email,
				Password: "password",
				FullName: user.FullName,
			},
			mockFn: func(ms *auth.MockService) {
				ms.EXPECT().Register(gomock.Any(), user.Email, "password", user.FullName).Return(&user, nil)
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "auth_service error",
			input: model.RegisterInput{
				Email:    user.Email,
				Password: "password",
				FullName: user.FullName,
			},
			mockFn: func(ms *auth.MockService) {
				ms.EXPECT().Register(gomock.Any(), user.Email, "password", user.FullName).
					Return(nil, errors.New("auth_service error"))
			},
			wantCode:    http.StatusBadRequest,
			errContains: "auth_service error",
		},
		{
			name: "invalid email",
			input: model.RegisterInput{
				Email:    "",
				Password: "password",
				FullName: user.FullName,
			},
			wantCode:    http.StatusBadRequest,
			errContains: "Error:Field validation for 'Email' failed",
		},
		{
			name: "invalid password",
			input: model.RegisterInput{
				Email:    user.Email,
				Password: "",
				FullName: user.FullName,
			},
			wantCode:    http.StatusBadRequest,
			errContains: "Error:Field validation for 'Password' failed",
		},
		{
			name: "invalid full name",
			input: model.RegisterInput{
				Email:    user.Email,
				Password: "password",
				FullName: "",
			},
			wantCode:    http.StatusBadRequest,
			errContains: "Error:Field validation for 'FullName' failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockService := setupHandlerTest(t)
			if tt.mockFn != nil {
				tt.mockFn(mockService)
			}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &res)
			assert.NoError(t, err)

			assert.NotNil(t, res)

			if tt.wantCode == http.StatusCreated {
				assert.Equal(t, user.ID.String(), res["id"])
				assert.Equal(t, user.FullName, res["full_name"])
				assert.Equal(t, user.Email, res["email"])
				assert.Equal(t, user.CreatedAt.Format(time.RFC3339Nano), res["created_at"])
				assert.Equal(t, user.UpdatedAt.Format(time.RFC3339Nano), res["updated_at"])
				assert.Empty(t, res["password_hash"])
			} else {
				assert.Contains(t, res["error"], tt.errContains)
			}
		})
	}
}

func Test_handler_Login(t *testing.T) {
	const (
		testToken    = "test-token"
		testEmail    = "test@example.com"
		testPassword = "password"
	)

	tests := []struct {
		name        string
		input       model.LoginInput
		mockFn      func(*auth.MockService)
		wantCode    int
		errContains string
	}{
		{
			name: "successful login",
			input: model.LoginInput{
				Email:    testEmail,
				Password: testPassword,
			},
			mockFn: func(ms *auth.MockService) {

				ms.EXPECT().Login(gomock.Any(), testEmail, testPassword).Return(testToken, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "auth_service error",
			input: model.LoginInput{
				Email:    testEmail,
				Password: testPassword,
			},
			mockFn: func(ms *auth.MockService) {
				ms.EXPECT().Login(gomock.Any(), testEmail, testPassword).Return("", errors.New("auth_service error"))
			},
			wantCode:    http.StatusBadRequest,
			errContains: "auth_service error",
		},
		{
			name: "invalid email",
			input: model.LoginInput{
				Email:    "",
				Password: testPassword,
			},
			wantCode:    http.StatusBadRequest,
			errContains: "Error:Field validation for 'Email' failed",
		},
		{
			name: "invalid password",
			input: model.LoginInput{
				Email:    testEmail,
				Password: "",
			},
			wantCode:    http.StatusBadRequest,
			errContains: "Error:Field validation for 'Password' failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockService := setupHandlerTest(t)
			if tt.mockFn != nil {
				tt.mockFn(mockService)
			}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &res)
			assert.NoError(t, err)

			assert.NotNil(t, res)

			if tt.wantCode == http.StatusOK {
				assert.Equal(t, testToken, res["token"])
			} else {
				assert.Contains(t, res["error"], tt.errContains)
			}
		})
	}
}
