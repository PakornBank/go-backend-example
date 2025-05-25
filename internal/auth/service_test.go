package auth

import (
	"context"
	"testing"
	"time"

	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/common/model"
	"github.com/PakornBank/go-backend-example/internal/common/testutil"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type registerInput struct {
	Email    string
	Password string
	FullName string
}

type loginInput struct {
	Email    string
	Password string
}

func setupServiceTest(t *testing.T) (Service, *MockRepository) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockRepository(ctrl)
	authService := &service{
		repository:  mockRepo,
		jwtSecret:   []byte("test-secret"),
		tokenExpiry: time.Hour * 24,
	}
	return authService, mockRepo
}

func TestNewService(t *testing.T) {
	mockRepo := new(MockRepository)
	cfg := &config.Config{
		JWTSecret:      "test-secret",
		TokenExpiryDur: time.Hour * 24,
	}
	authService := NewService(mockRepo, cfg)

	assert.NotNil(t, authService)
	assert.Equal(t, mockRepo, authService.(*service).repository)
	assert.Equal(t, []byte(cfg.JWTSecret), authService.(*service).jwtSecret)
	assert.Equal(t, cfg.TokenExpiryDur, authService.(*service).tokenExpiry)
}

func Test_service_Register(t *testing.T) {
	mockUser := testutil.NewMockUser()

	tests := []struct {
		name        string
		input       registerInput
		mockFn      func(*MockRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "successful registration",
			input: registerInput{
				Email:    mockUser.Email,
				Password: "password",
				FullName: mockUser.FullName,
			},
			mockFn: func(mr *MockRepository) {
				mr.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(&model.User{})).Return(nil)
				mr.EXPECT().FindByEmail(gomock.Any(), mockUser.Email).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: false,
		},
		{
			name: "email already exists",
			input: registerInput{
				Email:    mockUser.Email,
				Password: "password",
				FullName: mockUser.FullName,
			},
			mockFn: func(mr *MockRepository) {
				mr.EXPECT().FindByEmail(gomock.Any(), mockUser.Email).Return(&mockUser, nil)
			},
			wantErr:     true,
			errContains: "email already registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService, mockRepo := setupServiceTest(t)
			tt.mockFn(mockRepo)
			user, err := authService.Register(context.Background(), tt.input.Email, tt.input.Password, tt.input.FullName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errContains, err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.input.Email, user.Email)
				assert.Equal(t, tt.input.FullName, user.FullName)
			}
		})
	}
}

func Test_service_Login(t *testing.T) {
	mockUser := testutil.NewMockUser()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		input       loginInput
		mockFn      func(*MockRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "successful login",
			input: loginInput{
				Email:    mockUser.Email,
				Password: "password",
			},
			mockFn: func(mr *MockRepository) {
				mockUser.PasswordHash = string(hashedPassword)
				mr.EXPECT().FindByEmail(gomock.Any(), mockUser.Email).Return(&mockUser, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid credentials",
			input: loginInput{
				Email:    mockUser.Email,
				Password: "wrong password",
			},
			mockFn: func(mr *MockRepository) {
				mockUser.PasswordHash = string(hashedPassword)
				mr.EXPECT().FindByEmail(gomock.Any(), mockUser.Email).Return(&mockUser, nil)
			},
			wantErr:     true,
			errContains: "invalid credentials",
		},
		{
			name: "user not found",
			input: loginInput{
				Email:    "nonexistent@example.com",
				Password: "password",
			},
			mockFn: func(mr *MockRepository) {
				mr.EXPECT().FindByEmail(gomock.Any(), "nonexistent@example.com").Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr:     true,
			errContains: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService, mockRepo := setupServiceTest(t)
			tt.mockFn(mockRepo)
			token, err := authService.Login(context.Background(), tt.input.Email, tt.input.Password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errContains, err.Error())
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	authService, _ := setupServiceTest(t)
	mockUser := testutil.NewMockUser()

	token, err := authService.(*service).generateToken(&mockUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(_ *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, mockUser.ID.String(), claims["user_id"])
	assert.Equal(t, mockUser.Email, claims["email"])
}
