package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

// handler handles authentication-related HTTP requests.
type handler struct {
	service Service
}

// NewHandler creates a new instance of handler with the provided service.
func NewHandler(s Service) Handler {

	return &handler{service: s}
}

// Register handles the user registration process.
func (h *handler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login handles the user login process.
func (h *handler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
