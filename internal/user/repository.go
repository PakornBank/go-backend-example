package user

import (
	"context"
	"time"

	"github.com/PakornBank/go-backend-example/internal/common/model"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=./repository_mock.go -package=user github.com/PakornBank/go-backend-example/internal/user Repository

// Repository defines the methods that a repository must implement.
type Repository interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
}

// repository is a struct that provides methods to interact with the user data in the database.
type repository struct {
	db      *gorm.DB
	timeout time.Duration
}

// NewRepository creates a new instance of repository with the provided gorm.DB connection.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db, timeout: 5 * time.Second}
}

// FindByID retrieves a user from the database by their ID.
func (r *repository) FindByID(ctx context.Context, id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var user model.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
