package auth

import (
	"context"

	"github.com/PakornBank/go-backend-example/internal/common/model"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=./repository_mock.go -package=auth github.com/PakornBank/go-backend-example/internal/auth Repository

// Repository defines the methods that a repository must implement.
type Repository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

// repository is a struct that provides methods to interact with the user data in the database.
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new instance of repository with the provided gorm.DB connection.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new user record into the database.
func (r *repository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByEmail retrieves a user from the database by their email address.
func (r *repository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
