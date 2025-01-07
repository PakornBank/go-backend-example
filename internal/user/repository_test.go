package user

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PakornBank/go-backend-example/internal/common/model"
	"github.com/PakornBank/go-backend-example/internal/common/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupRepositoryTest(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, Repository) {
	_, gormDB, sqlMock := testutil.DBMock(t)
	userRepo := &repository{db: gormDB}
	return gormDB, sqlMock, userRepo
}

func TestNewRepository(t *testing.T) {
	_, gormDB, _ := testutil.DBMock(t)
	userRepo := NewRepository(gormDB)
	assert.NotNil(t, userRepo)
	assert.Equal(t, gormDB, userRepo.(*repository).db)
}

func Test_repository_FindByID(t *testing.T) {
	mockUser := testutil.NewMockUser()

	tests := []struct {
		name     string
		id       uuid.UUID
		mockFn   func(sqlmock.Sqlmock)
		wantUser *model.User
		wantErr  bool
		errType  error
	}{
		{
			name: "user found",
			id:   mockUser.ID,
			mockFn: func(sqlMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "created_at", "updated_at"}).
					AddRow(mockUser.ID, mockUser.Email, mockUser.PasswordHash, mockUser.FullName, mockUser.CreatedAt, mockUser.UpdatedAt)
				sqlMock.ExpectQuery(`SELECT .* FROM "users" WHERE id = \$1 (.+) LIMIT \$2`).
					WithArgs(mockUser.ID, 1).
					WillReturnRows(rows)
			},
			wantUser: &mockUser,
			wantErr:  false,
		},
		{
			name: "user not found",
			id:   mockUser.ID,
			mockFn: func(sqlMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "created_at", "updated_at"})
				sqlMock.ExpectQuery(`SELECT .* FROM "users" WHERE id = \$1 (.+) LIMIT \$2`).
					WithArgs(mockUser.ID, 1).
					WillReturnRows(rows)
			},
			wantUser: nil,
			wantErr:  true,
			errType:  gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, sqlMock, userRepo := setupRepositoryTest(t)

			tt.mockFn(sqlMock)
			got, err := userRepo.FindByID(context.Background(), tt.id.String())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser, got)
			}

			assert.NoError(t, sqlMock.ExpectationsWereMet())
		})
	}
}
