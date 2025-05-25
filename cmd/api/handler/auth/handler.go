package auth

import (
	"github.com/PakornBank/go-backend-example/cmd/api/model"
	"github.com/PakornBank/go-backend-example/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:generate mockgen -destination=./handler_mock.go -package=auth github.com/PakornBank/go-backend-example/cmd/api/handler/auth Handler

// Handler defines the interface for authentication-related HTTP requests.
type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

// handler handles authentication-related HTTP requests.
type handler struct {
	service auth.Service
}

// NewHandler creates a new instance of handler with the provided service.
func NewHandler(s auth.Service) Handler {

	return &handler{service: s}
}

// Register handles the user registration process.
func (h *handler) Register(c *gin.Context) {
	var input model.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(c.Request.Context(), input.Email, input.Password, input.FullName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := model.User{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusCreated, res)
}

// Login handles the user login process.
func (h *handler) Login(c *gin.Context) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
