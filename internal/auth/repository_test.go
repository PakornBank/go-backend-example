package auth

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PakornBank/go-backend-example/internal/common/model"
	"github.com/PakornBank/go-backend-example/internal/common/testutil"
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

func Test_repository_Create(t *testing.T) {
	mockUser := testutil.NewMockUser()

	tests := []struct {
		name    string
		user    *model.User
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
		errType error
	}{
		{
			name: "successful creation",
			user: &model.User{
				Email:        mockUser.Email,
				PasswordHash: mockUser.PasswordHash,
				FullName:     mockUser.FullName,
			},
			mockFn: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
					AddRow(mockUser.ID, mockUser.CreatedAt, mockUser.UpdatedAt)
				sqlMock.ExpectQuery(`INSERT INTO "users"`).
					WithArgs(mockUser.Email, mockUser.PasswordHash, mockUser.FullName).
					WillReturnRows(rows)
				sqlMock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			user: &model.User{
				Email:        mockUser.Email,
				PasswordHash: mockUser.PasswordHash,
				FullName:     mockUser.FullName,
			},
			mockFn: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				sqlMock.ExpectQuery(`INSERT INTO "users"`).
					WithArgs(mockUser.Email, mockUser.PasswordHash, mockUser.FullName).
					WillReturnError(sql.ErrConnDone)
				sqlMock.ExpectRollback()
			},
			wantErr: true,
			errType: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, sqlMock, userRepo := setupRepositoryTest(t)

			tt.mockFn(sqlMock)
			err := userRepo.Create(context.Background(), tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tt.user.ID)
				assert.NotZero(t, tt.user.CreatedAt)
				assert.NotZero(t, tt.user.UpdatedAt)
			}

			assert.NoError(t, sqlMock.ExpectationsWereMet())
		})
	}
}

func Test_repository_FindByEmail(t *testing.T) {
	mockUser := testutil.NewMockUser()

	tests := []struct {
		name     string
		mockFn   func(sqlmock.Sqlmock)
		wantUser *model.User
		wantErr  bool
		errType  error
	}{
		{
			name: "user found",
			mockFn: func(sqlMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "created_at", "updated_at"}).
					AddRow(mockUser.ID, mockUser.Email, mockUser.PasswordHash, mockUser.FullName, mockUser.CreatedAt, mockUser.UpdatedAt)
				sqlMock.ExpectQuery(`SELECT .* FROM "users" WHERE email = \$1 (.+) LIMIT \$2`).
					WithArgs(mockUser.Email, 1).
					WillReturnRows(rows)
			},
			wantUser: &mockUser,
			wantErr:  false,
		},
		{
			name: "user not found",
			mockFn: func(sqlMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "created_at", "updated_at"})
				sqlMock.ExpectQuery(`SELECT .* FROM "users" WHERE email = \$1 (.+) LIMIT \$2`).
					WithArgs(mockUser.Email, 1).
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
			got, err := userRepo.FindByEmail(context.Background(), mockUser.Email)

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
