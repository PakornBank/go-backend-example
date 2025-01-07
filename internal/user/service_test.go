package user

import (
	"context"
	"testing"
	"time"

	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/common/model"
	"github.com/PakornBank/go-backend-example/internal/common/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockRepository struct {
	mock.Mock
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
	userService := &service{
		repository:  mockRepo,
		jwtSecret:   []byte("test-secret"),
		tokenExpiry: time.Hour * 24,
	}
	return userService, mockRepo
}

func TestNewService(t *testing.T) {
	mockRepo := new(MockRepository)
	cfg := &config.Config{
		JWTSecret:      "test-secret",
		TokenExpiryDur: time.Hour * 24,
	}
	userService := NewService(mockRepo, cfg)

	assert.NotNil(t, userService)
	assert.Equal(t, mockRepo, userService.(*service).repository)
	assert.Equal(t, []byte(cfg.JWTSecret), userService.(*service).jwtSecret)
	assert.Equal(t, cfg.TokenExpiryDur, userService.(*service).tokenExpiry)
}

func Test_service_GetUserByID(t *testing.T) {
	mockUser := testutil.NewMockUser()

	tests := []struct {
		name    string
		id      string
		mockFn  func(*MockRepository)
		want    *model.User
		wantErr bool
		errType error
	}{
		{
			name: "user found",
			id:   mockUser.ID.String(),
			mockFn: func(repo *MockRepository) {
				repo.On("FindByID", mock.Anything, mockUser.ID.String()).Return(&mockUser, nil)
			},
			want:    &mockUser,
			wantErr: false,
		},
		{
			name: "user not found",
			id:   mockUser.ID.String(),
			mockFn: func(repo *MockRepository) {
				repo.On("FindByID", mock.Anything, mockUser.ID.String()).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errType: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService, mockRepo := setupServiceTest()
			tt.mockFn(mockRepo)
			got, err := userService.GetUserByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
