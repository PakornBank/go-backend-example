package user

import (
	"context"

	"github.com/PakornBank/go-backend-example/internal/common/model"
	"gorm.io/gorm"
)

// Repository defines the methods that a repository must implement.
type Repository interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
}

// repository is a struct that provides methods to interact with the user data in the database.
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new instance of repository with the provided gorm.DB connection.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// FindByID retrieves a user from the database by their ID.
func (r *repository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
