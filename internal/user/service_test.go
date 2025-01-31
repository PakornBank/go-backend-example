package user

import (
	"context"
	"testing"
	"time"

	"github.com/PakornBank/go-backend-example/internal/common/model"
	"github.com/PakornBank/go-backend-example/internal/common/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func setupServiceTest(t *testing.T) (Service, *MockRepository) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockRepository(ctrl)
	userService := &service{
		repository:  mockRepo,
		jwtSecret:   []byte("test-secret"),
		tokenExpiry: time.Hour * 24,
	}
	return userService, mockRepo
}

func TestNewService(t *testing.T) {
	mockRepo := new(MockRepository)
	userService := NewService(mockRepo)

	assert.NotNil(t, userService)
	assert.Equal(t, mockRepo, userService.(*service).repository)
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
			mockFn: func(mr *MockRepository) {
				mr.EXPECT().FindByID(gomock.Any(), mockUser.ID.String()).Return(&mockUser, nil)
			},
			want:    &mockUser,
			wantErr: false,
		},
		{
			name: "user not found",
			id:   mockUser.ID.String(),
			mockFn: func(mr *MockRepository) {
				mr.EXPECT().FindByID(gomock.Any(), mockUser.ID.String()).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errType: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService, mockRepo := setupServiceTest(t)
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
		})
	}
}
