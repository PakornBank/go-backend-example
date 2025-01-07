package testutil

import (
	"time"

	"github.com/PakornBank/go-backend-example/internal/common/model"
	"github.com/google/uuid"
)

// NewMockUser return a model.User with data.
func NewMockUser() model.User {
	return model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "Test User",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
