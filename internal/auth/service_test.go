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
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MockRepository struct {
	mock.Mock
}

func (r *MockRepository) Create(ctx context.Context, user *model.User) error {
	args := r.Called(ctx, user)
	return args.Error(0)
}

func (r *MockRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := r.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (r *MockRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	args := r.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func setupServiceTest() (Service, *MockRepository) {
	mockRepo := new(MockRepository)
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
		input       RegisterInput
		mockFn      func(*MockRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "successful registration",
			input: RegisterInput{
				Email:    mockUser.Email,
				Password: "password",
				FullName: mockUser.FullName,
			},
			mockFn: func(repo *MockRepository) {
				repo.On("FindByEmail", mock.Anything, mockUser.Email).Return(nil, gorm.ErrRecordNotFound)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "email already exists",
			input: RegisterInput{
				Email:    mockUser.Email,
				Password: "password",
				FullName: mockUser.FullName,
			},
			mockFn: func(repo *MockRepository) {
				repo.On("FindByEmail", mock.Anything, mockUser.Email).Return(&mockUser, nil)
			},
			wantErr:     true,
			errContains: "email already registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService, mockRepo := setupServiceTest()
			tt.mockFn(mockRepo)
			user, err := authService.Register(context.Background(), tt.input)

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
			mockRepo.AssertExpectations(t)
		})
	}
}

func Test_service_Login(t *testing.T) {
	mockUser := testutil.NewMockUser()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		input       LoginInput
		mockFn      func(*MockRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "successful login",
			input: LoginInput{
				Email:    mockUser.Email,
				Password: "password",
			},
			mockFn: func(repo *MockRepository) {
				mockUser.PasswordHash = string(hashedPassword)
				repo.On("FindByEmail", mock.Anything, mockUser.Email).Return(&mockUser, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid credentials",
			input: LoginInput{
				Email:    mockUser.Email,
				Password: "wrong password",
			},
			mockFn: func(repo *MockRepository) {
				mockUser.PasswordHash = string(hashedPassword)
				repo.On("FindByEmail", mock.Anything, mockUser.Email).Return(&mockUser, nil)
			},
			wantErr:     true,
			errContains: "invalid credentials",
		},
		{
			name: "user not found",
			input: LoginInput{
				Email:    "nonexistent@example.com",
				Password: "password",
			},
			mockFn: func(repo *MockRepository) {
				repo.On("FindByEmail", mock.Anything, "nonexistent@example.com").Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr:     true,
			errContains: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService, mockRepo := setupServiceTest()
			tt.mockFn(mockRepo)
			token, err := authService.Login(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errContains, err.Error())
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	authService, _ := setupServiceTest()
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
