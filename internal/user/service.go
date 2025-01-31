package user

import (
	"context"
	"time"

	"github.com/PakornBank/go-backend-example/internal/common/model"
)

//go:generate mockgen -destination=./service_mock.go -package=user github.com/PakornBank/go-backend-example/internal/user Service

// Service defines the methods that a service must implement.
type Service interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

// service is a struct that provides methods to interact with the user service.
type service struct {
	repository  Repository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

// NewService creates a new instance of service with the provided repository and configuration.
func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

// GetUserByID find and return user data with the given id.
func (s *service) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.repository.FindByID(ctx, id)
}
