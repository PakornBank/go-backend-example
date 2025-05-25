package auth

import (
	"context"
	"errors"
	"time"

	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/common/model"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -destination=./service_mock.go -package=auth github.com/PakornBank/go-backend-example/internal/auth Service

// Service defines the methods that a service must implement.
type Service interface {
	Register(ctx context.Context, email, password, fullName string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

// service is a struct that provides methods to interact with the authentication service.
type service struct {
	repository  Repository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

// NewService creates a new instance of service with the provided repository and configuration.
func NewService(repository Repository, config *config.Config) Service {
	return &service{
		repository:  repository,
		jwtSecret:   []byte(config.JWTSecret),
		tokenExpiry: config.TokenExpiryDur,
	}
}

// Register handles the user registration process.
func (s *service) Register(ctx context.Context, email, password, fullName string) (*model.User, error) {
	existingUser, _ := s.repository.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login handles the user login process.
func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.generateToken(user)
}

// generateToken generates a JWT token for the given user.
func (s *service) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
